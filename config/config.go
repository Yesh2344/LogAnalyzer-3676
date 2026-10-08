package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	LogDir       string // directory containing log files
	OutputFormat string // "json" or "text"
}

// Load builds a Config from environment variables and optional overrides.
// Overrides (logDir, format) are applied only if non‑empty.
func Load(logDirOverride, formatOverride string) (*Config, error) {
	cfg := &Config{
		LogDir:       getEnv("LOG_DIR", "./logs"),
		OutputFormat: strings.ToLower(getEnv("OUTPUT_FORMAT", "text")),
	}

	if logDirOverride != "" {
		cfg.LogDir = logDirOverride
	}
	if formatOverride != "" {
		cfg.OutputFormat = strings.ToLower(formatOverride)
	}

	if cfg.LogDir == "" {
		return nil, errors.New("log directory cannot be empty")
	}
	if cfg.OutputFormat != "json" && cfg.OutputFormat != "text" {
		return nil, errors.New("output format must be 'json' or 'text'")
	}
	return cfg, nil
}

// getEnv reads an environment variable or returns fallback.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}