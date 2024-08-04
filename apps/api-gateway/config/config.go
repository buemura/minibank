package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	PORT                      string
	GRPC_HOST_SVC_ACCOUNT     string
	GRPC_HOST_SVC_TRANSACTION string
)

func LoadEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		PORT = os.Getenv("PORT")
		GRPC_HOST_SVC_ACCOUNT = os.Getenv("GRPC_HOST_SVC_ACCOUNT")
		GRPC_HOST_SVC_TRANSACTION = os.Getenv("GRPC_HOST_SVC_TRANSACTION")
	} else {
		PORT = viper.GetString("PORT")
		GRPC_HOST_SVC_ACCOUNT = viper.GetString("GRPC_HOST_SVC_ACCOUNT")
		GRPC_HOST_SVC_TRANSACTION = viper.GetString("GRPC_HOST_SVC_TRANSACTION")
	}

	if len(PORT) == 0 ||
		len(GRPC_HOST_SVC_ACCOUNT) == 0 ||
		len(GRPC_HOST_SVC_TRANSACTION) == 0 {
		panic("Missing environment variables")
	}
}
