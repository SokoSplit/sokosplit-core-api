package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware validates a JWT bearer token and sets "user_id" in the
// Gin context. Tokens are issued at login/GitHub-OAuth-callback time
// (not implemented in this skeleton — wire up your OAuth flow and mint a
// JWT with a "sub" claim set to the user's UUID).
func AuthMiddleware() gin.HandlerFunc {
	secret := []byte(os.Getenv("JWT_SECRET"))
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		sub, _ := claims["sub"].(string)
		userID, err := uuid.Parse(sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid subject claim"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
