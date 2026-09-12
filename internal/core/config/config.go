package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr            string
	HTTPShutdownTimeout time.Duration
	DatabaseHost        string
	DatabasePort        string
	DatabaseUser        string
	DatabasePassword    string
	DatabaseName        string
	MinIOEndpoint       string
	MinIOAccessKey      string
	MinIOSecretKey      string
	MinIOBucket         string
	MinIOUseSSL         bool
}

func Load() (Config, error) {
	_ = godotenv.Load()
	c := Config{
		HTTPAddr:            env("HTTP_ADDR", ":8080"),
		HTTPShutdownTimeout: durationEnv("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		DatabaseHost:        env("POSTGRES_HOST", "localhost"), DatabasePort: env("POSTGRES_PORT", "5432"),
		DatabaseUser:     env("POSTGRES_USER", "postgres"),
		DatabasePassword: env("POSTGRES_PASSWORD", "postgres"),
		DatabaseName:     env("POSTGRES_DB", "genealogy"),
		MinIOEndpoint:    env("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:   env("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:   env("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:      env("MINIO_BUCKET", "genealogy"),
		MinIOUseSSL:      boolEnv("MINIO_USE_SSL", false),
	}
	if c.DatabaseUser == "" || c.DatabaseName == "" {
		return Config{}, fmt.Errorf("POSTGRES_USER and POSTGRES_DB are required")
	}
	return c, nil
}

func boolEnv(key string, fallback bool) bool {
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

func (c Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return fallback
}
