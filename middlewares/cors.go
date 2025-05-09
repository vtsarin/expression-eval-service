package middlewares

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates a middleware that handles CORS
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

				if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(204)
					return
				}
				break
			}
		}
		c.Next()
	}
}
