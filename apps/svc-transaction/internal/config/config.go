package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort           string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	AccountServiceAddr string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		GRPCPort:           getEnv("GRPC_PORT", "50053"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5434"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "transaction_db"),
		AccountServiceAddr: getEnv("ACCOUNT_SERVICE_ADDR", "localhost:50052"),
	}
}

func (c *Config) DatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
