package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host            string
	GRPCPort        int
	ShutdownTimeout time.Duration
	DatabaseURL     string
}

func Load() (Config, error) {
	cfg := Config{
		Host:            envOrDefault("TASK_SERVICE_HOST", "127.0.0.1"),
		GRPCPort:        9091,
		ShutdownTimeout: 10 * time.Second,
	}

	if value := os.Getenv("TASK_SERVICE_GRPC_PORT"); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("invalid TASK_SERVICE_GRPC_PORT: %w", err)
		}

		if port < 1 || port > 65535 {
			return Config{}, fmt.Errorf(
				"TASK_SERVICE_GRPC_PORT must be between 1 and 65535",
			)
		}

		cfg.GRPCPort = port
	}

	databaseURL := os.Getenv("TASK_SERVICE_DATABASE_URL")

	if databaseURL == "" {
		return Config{},
			fmt.Errorf("TASK_SERVICE_DATABASE_URL is required")
	}

	cfg.DatabaseURL = databaseURL

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
