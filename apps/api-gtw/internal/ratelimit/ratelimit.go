package ratelimit

import (
	"context"
	"time"
)

// RateLimitResult holds the outcome of a rate limit check.
type RateLimitResult struct {
	Allowed   bool
	Limit     int64
	Remaining int64
	ResetAt   time.Time
}

// RateLimiter checks whether a request identified by key is within the allowed rate.
type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int64, window time.Duration) (*RateLimitResult, error)
}
