package middleware

import (
	"net/http"
	"strings"

	authpb "github.com/buemura/minibank/api-gateway/proto/auth/v1"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authClient authpb.AuthServiceClient
}

func NewAuthMiddleware(authClient authpb.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		token := parts[1]

		resp, err := m.authClient.ValidateToken(c.Request.Context(), &authpb.ValidateTokenRequest{
			AccessToken: token,
		})

		if err != nil || !resp.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("user_id", resp.UserId)
		c.Set("email", resp.Email)

		c.Next()
	}
}
