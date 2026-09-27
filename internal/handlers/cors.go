package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware allows the dashboard (running on a different domain, e.g.
// Vercel) to call this API from the browser. Set ALLOWED_ORIGIN to your
// dashboard's URL in production; falls back to "*" (allow any origin) if
// unset, which is fine for early development.
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}