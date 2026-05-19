package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string
	DatabaseURL     string
	AppURL          string
	CookieSecure    bool
	ShutdownTimeout time.Duration
	MigrationsDir   string
}

func Load() (Config, error) {
	cfg := Config{
		Addr:            env("SYNC_SERVER_ADDR", ":8080"),
		DatabaseURL:     env("DATABASE_URL", "postgres://miniclaude:miniclaude@localhost:5432/miniclaude_sync?sslmode=disable"),
		AppURL:          env("SYNC_SERVER_APP_URL", "http://localhost:8080"),
		CookieSecure:    envBool("SYNC_SERVER_COOKIE_SECURE", false),
		ShutdownTimeout: envDuration("SYNC_SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		MigrationsDir:   env("SYNC_SERVER_MIGRATIONS_DIR", "migrations"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.AppURL == "" {
		return Config{}, fmt.Errorf("SYNC_SERVER_APP_URL is required")
	}

	return cfg, nil
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
