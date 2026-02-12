package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// luaScript atomically increments the counter and sets expiry if new.
// KEYS[1] = rate limit key
// ARGV[1] = window duration in seconds
// Returns: [current_count, ttl_in_seconds]
const luaScript = `
local key = KEYS[1]
local window = tonumber(ARGV[1])
local current = redis.call("INCR", key)
if current == 1 then
    redis.call("EXPIRE", key, window)
end
local ttl = redis.call("TTL", key)
return {current, ttl}
`

type RedisRateLimiter struct {
	rdb    *redis.Client
	script *redis.Script
}

func NewRedisRateLimiter(addr, password string) *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisRateLimiter{
		rdb:    client,
		script: redis.NewScript(luaScript),
	}
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (*RateLimitResult, error) {
	windowSecs := int64(window.Seconds())

	result, err := r.script.Run(ctx, r.rdb, []string{key}, windowSecs).Int64Slice()
	if err != nil {
		return nil, err
	}

	current := result[0]
	ttl := result[1]

	remaining := limit - current
	if remaining < 0 {
		remaining = 0
	}

	resetAt := time.Now().Add(time.Duration(ttl) * time.Second)

	return &RateLimitResult{
		Allowed:   current <= limit,
		Limit:     limit,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}
