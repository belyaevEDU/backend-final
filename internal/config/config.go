package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	ProcessingTime  time.Duration
	ShutdownTimeout time.Duration
}

const (
	envVarAddress         = "HTTP_ADDR"
	envVarProcessing      = "TASK_PROCESSING_TIME"
	envVarShutdownTimeout = "SHUTDOWN_TIMEOUT"
)

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:        envString(envVarAddress, ":8000"),
		ProcessingTime:  envDuration(envVarProcessing, 2*time.Second),
		ShutdownTimeout: envDuration(envVarShutdownTimeout, 10*time.Second),
	}

	if cfg.ProcessingTime <= 0 {
		return Config{}, fmt.Errorf("TASK_PROCESSING_TIME must be positive, got %s", cfg.ProcessingTime)
	}
	if cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be positive, got %s", cfg.ShutdownTimeout)
	}
	return cfg, nil
}

func envString(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		d, err := time.ParseDuration(strings.TrimSpace(v))
		if err != nil {
			return def
		}
		return d
	}
	return def
}
