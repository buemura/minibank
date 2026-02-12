package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/buemura/minibank/api-gtw/internal/cache"
	authpb "github.com/buemura/minibank/api-gtw/proto/auth/v1"
	"github.com/gin-gonic/gin"
)

const tokenCacheTTL = 5 * time.Minute

type AuthMiddleware struct {
	authClient authpb.AuthServiceClient
	cache      cache.Cache
}

func NewAuthMiddleware(authClient authpb.AuthServiceClient, cache cache.Cache) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient, cache: cache}
}

type cachedToken struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHENTICATED", "message": "Authorization header required"}})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHENTICATED", "message": "Invalid authorization format"}})
			return
		}

		token := parts[1]
		cacheKey := "auth:token:" + token

		// Try cache first
		if cached, err := m.cache.Get(cacheKey); err == nil && cached != "" {
			var data cachedToken
			if json.Unmarshal([]byte(cached), &data) == nil {
				c.Set("user_id", data.UserID)
				c.Set("email", data.Email)
				c.Next()
				return
			}
		}

		// Cache miss — call gRPC
		resp, err := m.authClient.ValidateToken(c.Request.Context(), &authpb.ValidateTokenRequest{
			AccessToken: token,
		})

		if err != nil || !resp.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHENTICATED", "message": "Invalid or expired token"}})
			return
		}

		// Cache valid result
		tokenJSON, _ := json.Marshal(cachedToken{
			UserID: resp.UserId,
			Email:  resp.Email,
		})
		_ = m.cache.Set(cacheKey, string(tokenJSON), tokenCacheTTL)

		c.Set("user_id", resp.UserId)
		c.Set("email", resp.Email)

		c.Next()
	}
}
