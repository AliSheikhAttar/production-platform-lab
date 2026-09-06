package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alisheikhattar/production-platform-lab/application/api-gateway/internal/config"
	"github.com/alisheikhattar/production-platform-lab/application/api-gateway/internal/grpcclient"
	"github.com/alisheikhattar/production-platform-lab/application/api-gateway/internal/httpapi"
)

func Run(logger *slog.Logger) int {
	cfg, err := config.Load()
	if err != nil {
		logger.Error(
			"configuration failed",
			slog.String("error", err.Error()),
		)
		return 1
	}

	taskClient, err := grpcclient.NewTaskClient(cfg.TaskServiceGRPCAddress)
	if err != nil {
		logger.Error(
			"grpc client creation failed",
			slog.String("error", err.Error()),
		)
		return 1
	}
	defer taskClient.Close()

	handler := httpapi.NewHandler(logger, taskClient.Client)

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info(
			"api gateway starting",
			slog.String("address", cfg.HTTPAddress),
			slog.String("task_service_grpc", cfg.TaskServiceGRPCAddress),
		)

		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"http server failed",
				slog.String("error", err.Error()),
			)
			return 1
		}

	case <-signalCtx.Done():
		logger.Info(
			"shutdown signal received",
			slog.String("signal", signalCtx.Err().Error()),
		)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.ShutdownTimeout,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"graceful shutdown failed",
			slog.String("error", err.Error()),
		)
		return 1
	}

	logger.Info("api gateway stopped")

	return 0
}

func Main() {
	os.Exit(Run(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		),
	))
}
