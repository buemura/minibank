package mocks

import (
	"context"
	"time"

	"github.com/buemura/minibank/api-gtw/internal/ratelimit"
	"github.com/stretchr/testify/mock"
)

type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (*ratelimit.RateLimitResult, error) {
	args := m.Called(ctx, key, limit, window)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ratelimit.RateLimitResult), args.Error(1)
}
