package tests

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/ITheCorgi/analyticevents/internal/app"
	"github.com/ITheCorgi/analyticevents/internal/config"
	"github.com/ITheCorgi/analyticevents/internal/consts"
	"github.com/ITheCorgi/analyticevents/internal/repository"
	"github.com/ITheCorgi/analyticevents/internal/usecase"
	"github.com/ITheCorgi/analyticevents/migrations"
	"github.com/ory/dockertest/v3"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chmigrate"
	"go.uber.org/zap"
)

var (
	configFile string         // nolint:gochecknoglobals
	configApp  *config.Config // nolint:gochecknoglobals
	db         *ch.DB
)

func init() {
	flag.StringVar(&configFile, consts.ConfigFileName, "config.yaml", "--config ./config.yaml")
}

func TestMain(m *testing.M) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configApp, err = config.New(configFile)
	if err != nil {
		log.Fatalf("failed to get config instance, configFile: %s\nerr: %v\n", configFile, err)
	}

	resource, pool, err := prepareClickhouseDockerForTest()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}

		err := resource.Close()
		if err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}()

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer zapLogger.Sync()

	eventRepo := repository.New(db)
	service := usecase.New(zapLogger, eventRepo)

	srv, err := app.ConfigureServer(configApp, service)
	if err != nil {
		log.Fatalf(err.Error())
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			zapLogger.Fatal(err.Error())
		}
	}()

	defer func() {
		err = srv.Shutdown(ctx)
		if err != nil {
			zapLogger.Fatal(err.Error())
		}
	}()

	os.Exit(m.Run())
}

func prepareClickhouseDockerForTest() (*dockertest.Resource, *dockertest.Pool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to Docker: %s", err)
	}

	resource, err := pool.BuildAndRun("analytic-ch-test-image", "./Dockerfile", []string{})
	if err != nil && !strings.EqualFold(err.Error()[2:], "container already exists") {
		return nil, nil, fmt.Errorf("could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db, err = repository.ProvideChConnection(context.Background(), configApp)
		if err != nil {
			fmt.Errorf(err.Error())
			return err
		}

		migrator := chmigrate.NewMigrator(db, migrations.Migrations)

		err = migrator.Init(context.Background())
		if err != nil {
			fmt.Errorf(err.Error())
			return err
		}

		group, err := migrator.Migrate(context.Background())
		if err != nil {
			migrator.Rollback(context.Background())
			fmt.Errorf(err.Error())
			return err
		}

		if group.IsZero() {
			fmt.Println("there are no new migrations to run")
			return nil
		}

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("c not connect to docker: %s", err)
	}

	return resource, pool, nil
}
