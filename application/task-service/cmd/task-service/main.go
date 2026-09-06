package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	taskv1 "github.com/alisheikhattar/production-platform-lab/application/task-service/gen/task/v1"
	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/config"
	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/grpcapi"
	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/postgresrepo"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	dbCtx, dbCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer dbCancel()

	pool, err := pgxpool.New(
		dbCtx,
		cfg.DatabaseURL,
	)
	if err != nil {
		logger.Error(
			"database pool creation failed",
			slog.String("error", err.Error()),
		)
		return 1
	}
	defer pool.Close()

	if err := pool.Ping(dbCtx); err != nil {
		logger.Error(
			"database connection failed",
			slog.String("error", err.Error()),
		)
		return 1
	}

	repository := postgresrepo.NewTaskRepository(pool)

	grpcAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort)

	grpcServer := grpc.NewServer()
	taskServer := grpcapi.NewServer(repository)
	taskv1.RegisterTaskServiceServer(grpcServer, taskServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Error(
			"grpc listener failed",
			slog.String("error", err.Error()),
		)
		return 1
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
			"grpc server starting",
			slog.String("address", grpcAddr),
		)

		serverErrors <- grpcServer.Serve(listener)
	}()

	select {
	case err := <-serverErrors:
		if err != nil {
			logger.Error(
				"grpc server failed",
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

	grpcServer.GracefulStop()

	logger.Info("task service stopped")

	return 0
}
