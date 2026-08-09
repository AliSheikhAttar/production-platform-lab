package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/config"
	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/httpapi"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	cfg, err := config.Load()
	if err != nil {
		logger.Error(
			"configuration failed",
			slog.String("error", err.Error()),
		)
		return 1
	}

	handler := httpapi.NewHandler(logger)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	server := &http.Server{
		Addr:              addr,
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
			"task service starting",
			slog.String("address", addr),
		)

		serverErrors <- server.ListenAndServe()
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

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"graceful shutdown failed",
			slog.String("error", err.Error()),
		)
		return 1
	}

	logger.Info("task service stopped")

	return 0
}
