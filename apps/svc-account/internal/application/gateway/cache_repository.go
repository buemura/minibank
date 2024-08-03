package gateway

import "time"

type CacheRepository interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
}
