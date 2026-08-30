package config

import (
	"fmt"
	"time"
)

type AppConfig struct {
	HTTPAddr        string
	ProcessingTime  time.Duration
	ShutdownTimeout time.Duration
}

const (
	envVarAppAddress         = "HTTP_ADDR"
	envVarAppProcessingTime  = "TASK_PROCESSING_TIME"
	envVarAppShutdownTimeout = "SHUTDOWN_TIMEOUT"
)

const (
	defaultAppAddress         = ":8000"
	defaultAppProcessingTime  = 2 * time.Second
	defaultAppShutdownTimeout = 10 * time.Second
)

func loadAppConfig() (AppConfig, error) {
	appCfg := AppConfig{
		HTTPAddr:        envString(envVarAppAddress, defaultAppAddress),
		ProcessingTime:  envDuration(envVarAppProcessingTime, defaultAppProcessingTime),
		ShutdownTimeout: envDuration(envVarAppShutdownTimeout, defaultAppShutdownTimeout),
	}

	if appCfg.ProcessingTime <= 0 {
		return AppConfig{}, fmt.Errorf("TASK_PROCESSING_TIME must be positive, got %s", appCfg.ProcessingTime)
	}
	if appCfg.ShutdownTimeout <= 0 {
		return AppConfig{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be positive, got %s", appCfg.ShutdownTimeout)
	}
	return appCfg, nil
}
