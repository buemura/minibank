package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	GRPC_PORT             string
	DATABASE_URL          string
	BROKER_URL            string
	REDIS_URL             string
	REDIS_PASSWORD        string = ""
	GRPC_HOST_SVC_ACCOUNT string
)

func LoadEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		GRPC_PORT = os.Getenv("GRPC_PORT")
		DATABASE_URL = os.Getenv("DATABASE_URL")
		BROKER_URL = os.Getenv("BROKER_URL")
		REDIS_URL = os.Getenv("REDIS_URL")
		REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
		GRPC_HOST_SVC_ACCOUNT = os.Getenv("GRPC_HOST_SVC_ACCOUNT")
	} else {
		GRPC_PORT = viper.GetString("GRPC_PORT")
		DATABASE_URL = viper.GetString("DATABASE_URL")
		BROKER_URL = viper.GetString("BROKER_URL")
		REDIS_URL = viper.GetString("REDIS_URL")
		REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
		GRPC_HOST_SVC_ACCOUNT = viper.GetString("GRPC_HOST_SVC_ACCOUNT")
	}

	if len(GRPC_PORT) == 0 ||
		len(DATABASE_URL) == 0 ||
		len(BROKER_URL) == 0 ||
		len(REDIS_URL) == 0 ||
		len(GRPC_HOST_SVC_ACCOUNT) == 0 {
		panic("Missing environment variables")
	}
}
