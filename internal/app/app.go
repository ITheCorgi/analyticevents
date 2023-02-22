package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ITheCorgi/analyticevents/internal/config"
	"github.com/ITheCorgi/analyticevents/internal/controllers"
	"github.com/ITheCorgi/analyticevents/internal/repository"
	"github.com/ITheCorgi/analyticevents/internal/usecase"
	"github.com/ITheCorgi/analyticevents/pkg/server"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Run executes service operation
func Run(ctx context.Context, cancelFunc context.CancelFunc, configApp *config.Config, zapLogger *zap.Logger) {
	defer zapLogger.Sync()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	chConn, err := repository.ProvideChConnection(ctx, configApp)
	if err != nil {
		zapLogger.Fatal("failed to get clickhouse connection", zap.Error(err))
	}
	defer chConn.Close()

	eventRepo := repository.New(chConn)

	service := usecase.New(zapLogger, eventRepo)

	srv, err := ConfigureServer(configApp, service)
	if err != nil {
		zapLogger.Fatal(err.Error())
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			zapLogger.Fatal(err.Error())
		}
	}()

	sig := <-stopChan
	zapLogger.Info("start to gracefully shutdown due to signal", zap.String("sig", sig.String()))

	err = srv.Shutdown(ctx)
	if err != nil {
		zapLogger.Fatal(err.Error())
	}

	cancelFunc()

	os.Exit(0)
}

func ConfigureServer(cfg *config.Config, writer controllers.EventWriter) (*http.Server, error) {
	swagger, err := server.GetSwagger()
	if err != nil {
		return nil, err
	}

	swagger.Servers = nil
	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))

	server.RegisterHandlers(r, controllers.New(writer))

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
	}

	return s, nil
}
