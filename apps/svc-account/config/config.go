package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	GRPC_PORT      string
	DATABASE_URL   string
	REDIS_URL      string
	REDIS_PASSWORD string = ""
)

func LoadEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		GRPC_PORT = os.Getenv("GRPC_PORT")
		DATABASE_URL = os.Getenv("DATABASE_URL")
		REDIS_URL = os.Getenv("REDIS_URL")
		REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	} else {
		GRPC_PORT = viper.GetString("GRPC_PORT")
		DATABASE_URL = viper.GetString("DATABASE_URL")
		REDIS_URL = viper.GetString("REDIS_URL")
		REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
	}

	if len(GRPC_PORT) == 0 || len(DATABASE_URL) == 0 || len(REDIS_URL) == 0 {
		panic("Missing environment variables")
	}
}
