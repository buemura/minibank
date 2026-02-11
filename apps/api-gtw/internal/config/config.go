package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort               string
	AuthServiceAddr        string
	AccountServiceAddr     string
	TransactionServiceAddr string
	RedisAddr              string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		HTTPPort:               getEnv("HTTP_PORT", "8080"),
		AuthServiceAddr:        getEnv("AUTH_SERVICE_ADDR", "localhost:50051"),
		AccountServiceAddr:     getEnv("ACCOUNT_SERVICE_ADDR", "localhost:50052"),
		TransactionServiceAddr: getEnv("TRANSACTION_SERVICE_ADDR", "localhost:50053"),
		RedisAddr:              getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
