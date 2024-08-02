package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	GRPC_PORT    string
	DATABASE_URL string
)

func LoadEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		GRPC_PORT = os.Getenv("GRPC_PORT")
		DATABASE_URL = os.Getenv("DATABASE_URL")
	} else {
		GRPC_PORT = viper.GetString("GRPC_PORT")
		DATABASE_URL = viper.GetString("DATABASE_URL")
	}

	if len(GRPC_PORT) == 0 || len(DATABASE_URL) == 0 {
		panic("Missing environment variables")
	}
}
