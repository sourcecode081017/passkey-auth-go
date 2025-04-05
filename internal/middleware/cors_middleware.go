package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a middleware handler that adds CORS headers to each response
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from any origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific HTTP methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")

		// Allow specific headers
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Allow credentials (cookies, etc.)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Continue to the next middleware/handler
		c.Next()
	}
}
