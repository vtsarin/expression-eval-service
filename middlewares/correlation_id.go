package middlewares

import (
	"go-service/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/**
 * BindCorrelationId middleware ensures that every request has a correlation ID
 * It follows these rules:
 * 1. Use the correlation ID from the request header if present
 * 2. If not present in the header, generate a new UUID
 * 3. Set the correlation ID in the context and response header
 */
func BindCorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader(constants.X_CORRELATION_ID)
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		c.Set("correlationId", correlationID)

		// Set the correlation ID in the response header. Only if the response is not nil
		if c.Writer != nil {
			c.Writer.Header().Set(constants.X_CORRELATION_ID, correlationID)
		}

		c.Next()
	}
}
