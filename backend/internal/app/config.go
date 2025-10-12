package app

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration for the backend.
type Config struct {
	HTTPPort          string
	DBDSN             string
	AuthSecret        string
	EvaluationTimeout time.Duration
	RadiusKm          float64
}

// LoadConfig reads environment variables into Config with defaults applied.
func LoadConfig() (Config, error) {
	cfg := Config{
		HTTPPort:          getEnv("BACKEND_HTTP_PORT", "8080"),
		DBDSN:             os.Getenv("BACKEND_DB_DSN"),
		AuthSecret:        getEnv("BACKEND_AUTH_SECRET", "dev-secret"),
		EvaluationTimeout: 3 * time.Second,
		RadiusKm:          5,
	}

	if v := os.Getenv("BACKEND_EVALUATION_TIMEOUT"); v != "" {
		dur, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("parse BACKEND_EVALUATION_TIMEOUT: %w", err)
		}
		cfg.EvaluationTimeout = dur
	}

	if v := os.Getenv("BACKEND_RADIUS_KM"); v != "" {
		km, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return Config{}, fmt.Errorf("parse BACKEND_RADIUS_KM: %w", err)
		}
		cfg.RadiusKm = km
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
