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
	S3Endpoint          string
	S3AccessKey         string
	S3SecretKey         string
	S3Bucket            string
	S3UseSSL            bool
	JWTSecret           string
	JWTIssuer           string
	JWTAccessTTL        time.Duration
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
		S3Endpoint:       env("S3_ENDPOINT", "localhost:8333"),
		S3AccessKey:      env("S3_ACCESS_KEY", "s3admin"),
		S3SecretKey:      env("S3_SECRET_KEY", "s3secret"),
		S3Bucket:         env("S3_BUCKET", "genealogy"),
		S3UseSSL:         boolEnv("S3_USE_SSL", false),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTIssuer:        env("JWT_ISSUER", "genealogy-tree"),
		JWTAccessTTL:     durationEnv("JWT_ACCESS_TTL", 15*time.Minute),
	}
	if c.DatabaseUser == "" || c.DatabaseName == "" {
		return Config{}, fmt.Errorf("POSTGRES_USER and POSTGRES_DB are required")
	}
	if c.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
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
