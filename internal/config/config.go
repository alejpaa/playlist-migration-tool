package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port                  string
	GoogleCredentialsFile string
	TokenFile             string
	Environment           string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:                  getEnv("PORT", "8080"),
		GoogleCredentialsFile: getEnv("GOOGLE_CREDENTIALS_FILE", "client_secret_332431762901-dthq67hje7hcldkt4edg2n6dlbujsuck.apps.googleusercontent.com.json"),
		TokenFile:             getEnv("TOKEN_FILE", "token.json"),
		Environment:           getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
