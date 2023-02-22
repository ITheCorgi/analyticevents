package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ITheCorgi/analyticevents/internal/app"
	"github.com/ITheCorgi/analyticevents/internal/config"
	"github.com/ITheCorgi/analyticevents/internal/consts"
	"github.com/ITheCorgi/analyticevents/internal/repository"
	"github.com/ITheCorgi/analyticevents/migrations"
	"github.com/spf13/cobra"
	"github.com/uptrace/go-clickhouse/chmigrate"
	"go.uber.org/zap"
)

var (
	configFile string         // nolint:gochecknoglobals
	fileName   string         // nolint:gochecknoglobals
	configApp  *config.Config // nolint:gochecknoglobals
	zapLogger  *zap.Logger    // nolint:gochecknoglobals

	rootCmd = &cobra.Command{ // nolint:gochecknoglobals
		Use:              "root",
		Short:            "analytic",
		PersistentPreRun: preRun,
		Run:              runApp,
	}

	migrateCmd = &cobra.Command{
		Use: "migrate [command]",
	}
	migrateCmdUp = &cobra.Command{
		Use:              consts.MigrateActionUp,
		PersistentPreRun: preRun,
		Run:              runMigration(consts.MigrateActionUp),
	}
	migrateCreateFile = &cobra.Command{
		Use:              consts.MigrateActionCreate,
		PersistentPreRun: preRun,
		Run:              runMigration(consts.MigrateActionCreate),
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, consts.ConfigFileName, "c", "config.yaml", "--config ./config.yaml")
	migrateCmd.AddCommand(migrateCmdUp)

	migrateCreateFile.Flags().StringVarP(&fileName, consts.MigrateFlagName, "n", "", "migrate --name add_new_analytic_table")
	migrateCreateFile.MarkFlagRequired(consts.MigrateFlagName)
	migrateCmd.AddCommand(migrateCreateFile)

	rootCmd.AddCommand(migrateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func preRun(_ *cobra.Command, _ []string) {
	if configApp == nil {
		var err error
		configApp, err = config.New(configFile)
		if err != nil {
			log.Fatalf("failed to get config instance, configFile: %s\nerr: %v\n", configFile, err)
		}

		log.Println("config object parsed")
	}

	if zapLogger == nil {
		err := getZapLogger()
		if err != nil {
			log.Fatalln("failed to get zap logger")
		}
	}

	return
}

func runApp(_ *cobra.Command, _ []string) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	app.Run(ctx, cancelFunc, configApp, zapLogger)
}

func getZapLogger() error {
	var err error

	switch configApp.App.Environment {
	case consts.ProdEnvironment:
		zapLogger, err = zap.NewProduction()
		if err != nil {
			return err
		}

		zapLogger.Info("production zap logger started")

	default:
		zapLogger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}

		zapLogger.Info("dev zap logger started")
	}

	return nil
}

func runMigration(action string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		pCtx := context.Background()

		conn, err := repository.ProvideChConnection(pCtx, configApp)
		if err != nil {
			zapLogger.Fatal(err.Error())
		}

		migrator := chmigrate.NewMigrator(conn, migrations.Migrations)

		err = migrator.Init(context.Background())
		if err != nil {
			zapLogger.Fatal(err.Error())
		}

		switch action {
		case consts.MigrateActionUp:
			mCtx, cancel := context.WithTimeout(pCtx, time.Minute)
			defer cancel()

			group, err := migrator.Migrate(mCtx)
			if err != nil {
				migrator.Rollback(mCtx)
				zapLogger.Fatal(err.Error())
			}

			if group.IsZero() {
				zapLogger.Info("there are no new migrations to run")
			}

			zapLogger.Info("migration up done")
			return

		case consts.MigrateActionCreate:
			zapLogger.Info("started job for creating migration file")

			files, err := migrator.CreateSQLMigrations(context.Background(), fileName)
			if err != nil {
				zapLogger.Fatal(err.Error())
			}

			for i := range files {
				zapLogger.Info("migration file created", zap.String("path", fmt.Sprintf("%s/%s", files[i].Path, files[i].Name)))
			}

		default:
			zapLogger.Fatal("unknown migrate command")
		}
	}
}
