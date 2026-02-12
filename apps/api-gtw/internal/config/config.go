package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort               string
	AuthServiceAddr        string
	AccountServiceAddr     string
	TransactionServiceAddr string
	RedisAddr              string
	OTLPEndpoint           string
	AllowedOrigins         string

	// Rate limiting
	RateLimitAuthLimit  int64
	RateLimitAuthWindow time.Duration
	RateLimitAPILimit   int64
	RateLimitAPIWindow  time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		HTTPPort:               getEnv("HTTP_PORT", "8080"),
		AuthServiceAddr:        getEnv("AUTH_SERVICE_ADDR", "localhost:50051"),
		AccountServiceAddr:     getEnv("ACCOUNT_SERVICE_ADDR", "localhost:50052"),
		TransactionServiceAddr: getEnv("TRANSACTION_SERVICE_ADDR", "localhost:50053"),
		RedisAddr:              getEnv("REDIS_ADDR", "localhost:6379"),
		OTLPEndpoint:           getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		AllowedOrigins:         getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),

		RateLimitAuthLimit:  getEnvInt64("RATE_LIMIT_AUTH_LIMIT", 20),
		RateLimitAuthWindow: getEnvDuration("RATE_LIMIT_AUTH_WINDOW", 1*time.Minute),
		RateLimitAPILimit:   getEnvInt64("RATE_LIMIT_API_LIMIT", 100),
		RateLimitAPIWindow:  getEnvDuration("RATE_LIMIT_API_WINDOW", 1*time.Minute),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}
