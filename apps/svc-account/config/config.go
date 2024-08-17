package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	GRPC_PORT       string
	DATABASE_URL    string
	DATABASE_HOST   string
	DATABASE_PORT   string
	DATABASE_USER   string
	DATABASE_PASS   string
	DATABASE_DBNAME string
	REDIS_URL       string
	REDIS_PASSWORD  string = ""
)

func LoadEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		GRPC_PORT = os.Getenv("GRPC_PORT")
		DATABASE_URL = os.Getenv("DATABASE_URL")
		DATABASE_HOST = os.Getenv("DATABASE_HOST")
		DATABASE_PORT = os.Getenv("DATABASE_PORT")
		DATABASE_USER = os.Getenv("DATABASE_USER")
		DATABASE_PASS = os.Getenv("DATABASE_PASS")
		DATABASE_DBNAME = os.Getenv("DATABASE_DBNAME")
		REDIS_URL = os.Getenv("REDIS_URL")
		REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	} else {
		GRPC_PORT = viper.GetString("GRPC_PORT")
		DATABASE_URL = viper.GetString("DATABASE_URL")
		DATABASE_HOST = viper.GetString("DATABASE_HOST")
		DATABASE_PORT = viper.GetString("DATABASE_PORT")
		DATABASE_USER = viper.GetString("DATABASE_USER")
		DATABASE_PASS = viper.GetString("DATABASE_PASS")
		DATABASE_DBNAME = viper.GetString("DATABASE_DBNAME")
		REDIS_URL = viper.GetString("REDIS_URL")
		REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
	}

	if len(GRPC_PORT) == 0 ||
		len(DATABASE_URL) == 0 ||
		len(DATABASE_HOST) == 0 ||
		len(DATABASE_PORT) == 0 ||
		len(DATABASE_USER) == 0 ||
		len(DATABASE_PASS) == 0 ||
		len(DATABASE_DBNAME) == 0 ||
		len(REDIS_URL) == 0 {
		panic("Missing environment variables")
	}
}
