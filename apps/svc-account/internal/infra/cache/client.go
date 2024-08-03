package cache

import (
	"github.com/buemura/minibank/svc-account/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Connect() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_URL,
		Password: config.REDIS_PASSWORD,
		DB:       0, // use default DB
	})
}
