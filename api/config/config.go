package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	// DatabaseURL is the PostgreSQL connection string (required).
	DatabaseURL string

	// Port is the HTTP listen port for the Go API (default: 8080).
	Port string

	// R2AccountID is the Cloudflare account ID used to derive the R2 endpoint.
	R2AccountID string

	// R2AccessKeyID is the R2 / S3-compatible access key ID.
	R2AccessKeyID string

	// R2SecretAccessKey is the R2 / S3-compatible secret access key.
	R2SecretAccessKey string

	// R2Bucket is the name of the Cloudflare R2 bucket for photo uploads.
	R2Bucket string

	// JWTSecret is the HMAC signing key for JWT tokens (required for Phase 8 auth).
	JWTSecret string
}

// Load reads environment variables (with optional .env file) and returns a
// populated Config. Returns an error if any required variable is missing.
func Load() (*Config, error) {
	// Load .env if present; ignore the error so production (real env vars) works fine.
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		Port:              os.Getenv("API_PORT"),
		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2Bucket:          os.Getenv("R2_BUCKET"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL is required")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
