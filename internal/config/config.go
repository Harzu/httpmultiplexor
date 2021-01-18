package config

import (
	"os"
	"strconv"
)

const (
	envPortKey                = "PORT"
	envHandleRateLimitKey     = "HANDLE_RATE_LIMIT"
	envClientRequestRateLimit = "CLIENT_REQUEST_RATE_LIMIT"
)

type Config struct {
	Port                   string
	HandleRateLimit        int
	ClientRequestRateLimit int
}

func InitConfig() *Config {
	cfg := &Config{
		Port:                   "8080",
		HandleRateLimit:        100,
		ClientRequestRateLimit: 4,
	}

	port := os.Getenv(envPortKey)
	if _, err := strconv.Atoi(port); err == nil {
		cfg.Port = port
	}

	handleRateLimit, err := strconv.Atoi(os.Getenv(envHandleRateLimitKey))
	if err == nil {
		cfg.HandleRateLimit = handleRateLimit
	}

	clientRequestRateLimit, err := strconv.Atoi(os.Getenv(envClientRequestRateLimit))
	if err == nil {
		cfg.ClientRequestRateLimit = clientRequestRateLimit
	}

	return cfg
}
