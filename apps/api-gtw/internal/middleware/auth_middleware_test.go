package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authpb "github.com/buemura/minibank/api-gtw/proto/auth/v1"
	"github.com/buemura/minibank/api-gtw/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuth_NoHeader(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	mw := NewAuthMiddleware(authClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	mw.Authenticate()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestAuth_InvalidFormat_NoBearerPrefix(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	mw := NewAuthMiddleware(authClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Token abc123")

	mw.Authenticate()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization format")
}

func TestAuth_InvalidFormat_TooManyParts(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	mw := NewAuthMiddleware(authClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer token extra")

	mw.Authenticate()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization format")
}

func TestAuth_InvalidToken(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	mw := NewAuthMiddleware(authClient)

	authClient.On("ValidateToken", mock.Anything, mock.Anything).Return(&authpb.ValidateTokenResponse{
		Valid: false,
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	mw.Authenticate()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
	authClient.AssertExpectations(t)
}

func TestAuth_gRPCError(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	mw := NewAuthMiddleware(authClient)

	authClient.On("ValidateToken", mock.Anything, mock.Anything).Return(nil, errors.New("gRPC error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer some-token")

	mw.Authenticate()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	authClient.AssertExpectations(t)
}

func TestAuth_ValidToken(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	mw := NewAuthMiddleware(authClient)

	authClient.On("ValidateToken", mock.Anything, mock.Anything).Return(&authpb.ValidateTokenResponse{
		Valid:  true,
		UserId: "user-123",
		Email:  "test@example.com",
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	nextCalled := false
	c.Set("_gin_handler_index", 0)

	router := gin.New()
	router.Use(mw.Authenticate())
	router.GET("/", func(c *gin.Context) {
		nextCalled = true
		assert.Equal(t, "user-123", c.GetString("user_id"))
		assert.Equal(t, "test@example.com", c.GetString("email"))
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	authClient.AssertExpectations(t)
}
