// Package config loads runtime settings from environment variables.
package config

import (
	"os"
	"time"
)

const (
	defaultPort        = "8080"
	defaultHolidaysURL = "https://api.boostr.cl/holidays.json"
	defaultHTTPTimeout = 10 * time.Second
)

// Settings holds the runtime configuration of the application.
type Settings struct {
	Port        string
	HolidaysURL string
	HTTPTimeout time.Duration
}

// Load reads the settings from the PORT, HOLIDAYS_API_URL and HTTP_TIMEOUT environment
// variables, falling back to defaults for missing or invalid values.
func Load() Settings {
	return Settings{
		Port:        stringFromEnv("PORT", defaultPort),
		HolidaysURL: stringFromEnv("HOLIDAYS_API_URL", defaultHolidaysURL),
		HTTPTimeout: durationFromEnv("HTTP_TIMEOUT", defaultHTTPTimeout),
	}
}

func stringFromEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	parsed, err := time.ParseDuration(os.Getenv(key))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
