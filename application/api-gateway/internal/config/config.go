package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddress            string
	TaskServiceGRPCAddress string
	ShutdownTimeout        time.Duration
}

func Load() (Config, error) {
	httpAddress := os.Getenv("GATEWAY_HTTP_ADDRESS")
	if httpAddress == "" {
		httpAddress = ":8080"
	}

	taskServiceAddress := os.Getenv("TASK_SERVICE_GRPC_ADDRESS")
	if taskServiceAddress == "" {
		return Config{}, fmt.Errorf("TASK_SERVICE_GRPC_ADDRESS is required")
	}

	return Config{
		HTTPAddress:            httpAddress,
		TaskServiceGRPCAddress: taskServiceAddress,
		ShutdownTimeout:        10 * time.Second,
	}, nil
}
