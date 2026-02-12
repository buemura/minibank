package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/buemura/minibank/api-gtw/internal/ratelimit"
	"github.com/gin-gonic/gin"
)

// RateLimitConfig holds the configuration for a rate limit middleware instance.
type RateLimitConfig struct {
	Limit      int64
	Window     time.Duration
	RouteGroup string
}

type RateLimitMiddleware struct {
	limiter ratelimit.RateLimiter
}

func NewRateLimitMiddleware(limiter ratelimit.RateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{limiter: limiter}
}

// RateLimit returns a Gin handler that enforces the given rate limit config.
// It uses user_id from context (set by auth middleware) if available,
// otherwise falls back to client IP.
func (m *RateLimitMiddleware) RateLimit(cfg RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.GetString("user_id")
		if identifier == "" {
			identifier = c.ClientIP()
		}

		key := fmt.Sprintf("ratelimit:%s:%s", cfg.RouteGroup, identifier)

		result, err := m.limiter.Allow(c.Request.Context(), key, cfg.Limit, cfg.Window)
		if err != nil {
			// Fail open — if Redis is down, allow the request through.
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

		if !result.Allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "Too many requests. Please try again later.",
				},
			})
			return
		}

		c.Next()
	}
}
