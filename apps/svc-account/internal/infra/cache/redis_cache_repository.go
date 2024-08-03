package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
	rdb *redis.Client
}

func NewRedisCacheRepository() *RedisCacheRepository {
	return &RedisCacheRepository{
		rdb: RDB,
	}
}

func (r *RedisCacheRepository) Get(key string) (string, error) {
	value, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (r *RedisCacheRepository) Set(key string, value string, expiration time.Duration) error {
	err := r.rdb.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
