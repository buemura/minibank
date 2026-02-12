package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/buemura/minibank/api-gtw/internal/mocks"
	"github.com/buemura/minibank/api-gtw/internal/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimit_Allowed(t *testing.T) {
	mockLimiter := new(mocks.MockRateLimiter)
	mw := NewRateLimitMiddleware(mockLimiter)

	resetAt := time.Now().Add(60 * time.Second)
	mockLimiter.On("Allow", mock.Anything, "ratelimit:api:192.0.2.1", int64(100), 1*time.Minute).
		Return(&ratelimit.RateLimitResult{
			Allowed:   true,
			Limit:     100,
			Remaining: 99,
			ResetAt:   resetAt,
		}, nil)

	nextCalled := false
	router := gin.New()
	router.Use(mw.RateLimit(RateLimitConfig{
		Limit:      100,
		Window:     1 * time.Minute,
		RouteGroup: "api",
	}))
	router.GET("/", func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1:1234"
	router.ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "100", w.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "99", w.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"))
	mockLimiter.AssertExpectations(t)
}

func TestRateLimit_Exceeded(t *testing.T) {
	mockLimiter := new(mocks.MockRateLimiter)
	mw := NewRateLimitMiddleware(mockLimiter)

	resetAt := time.Now().Add(30 * time.Second)
	mockLimiter.On("Allow", mock.Anything, mock.Anything, int64(10), 1*time.Minute).
		Return(&ratelimit.RateLimitResult{
			Allowed:   false,
			Limit:     10,
			Remaining: 0,
			ResetAt:   resetAt,
		}, nil)

	nextCalled := false
	router := gin.New()
	router.Use(mw.RateLimit(RateLimitConfig{
		Limit:      10,
		Window:     1 * time.Minute,
		RouteGroup: "auth",
	}))
	router.GET("/", func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "RATE_LIMITED")
	assert.Equal(t, "10", w.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", w.Header().Get("X-RateLimit-Remaining"))
	mockLimiter.AssertExpectations(t)
}

func TestRateLimit_UsesUserIDWhenAvailable(t *testing.T) {
	mockLimiter := new(mocks.MockRateLimiter)
	mw := NewRateLimitMiddleware(mockLimiter)

	resetAt := time.Now().Add(60 * time.Second)
	mockLimiter.On("Allow", mock.Anything, "ratelimit:api:user-123", int64(100), 1*time.Minute).
		Return(&ratelimit.RateLimitResult{
			Allowed:   true,
			Limit:     100,
			Remaining: 99,
			ResetAt:   resetAt,
		}, nil)

	nextCalled := false
	router := gin.New()
	// Simulate auth middleware setting user_id
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user-123")
		c.Next()
	})
	router.Use(mw.RateLimit(RateLimitConfig{
		Limit:      100,
		Window:     1 * time.Minute,
		RouteGroup: "api",
	}))
	router.GET("/", func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	mockLimiter.AssertExpectations(t)
}

func TestRateLimit_UsesIPWhenNoUserID(t *testing.T) {
	mockLimiter := new(mocks.MockRateLimiter)
	mw := NewRateLimitMiddleware(mockLimiter)

	resetAt := time.Now().Add(60 * time.Second)
	mockLimiter.On("Allow", mock.Anything, "ratelimit:auth:192.0.2.1", int64(20), 1*time.Minute).
		Return(&ratelimit.RateLimitResult{
			Allowed:   true,
			Limit:     20,
			Remaining: 19,
			ResetAt:   resetAt,
		}, nil)

	router := gin.New()
	router.Use(mw.RateLimit(RateLimitConfig{
		Limit:      20,
		Window:     1 * time.Minute,
		RouteGroup: "auth",
	}))
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1:5678"
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockLimiter.AssertExpectations(t)
}

func TestRateLimit_RedisError_FailsOpen(t *testing.T) {
	mockLimiter := new(mocks.MockRateLimiter)
	mw := NewRateLimitMiddleware(mockLimiter)

	mockLimiter.On("Allow", mock.Anything, mock.Anything, int64(100), 1*time.Minute).
		Return(nil, errors.New("redis connection refused"))

	nextCalled := false
	router := gin.New()
	router.Use(mw.RateLimit(RateLimitConfig{
		Limit:      100,
		Window:     1 * time.Minute,
		RouteGroup: "api",
	}))
	router.GET("/", func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	mockLimiter.AssertExpectations(t)
}
