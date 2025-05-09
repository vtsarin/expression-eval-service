package errors

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var errorMetrics *ErrorMetrics

// InitializeErrorMetrics initializes the error metrics
func InitializeErrorMetrics(logger *zap.Logger) {
	errorMetrics = NewErrorMetrics(logger)
}

// HandleError handles errors and sends appropriate responses
func HandleError(c *gin.Context, status int, err error) {
	// Type assertion to check if the error is a CustomError
	if customErr, ok := err.(*CustomError); ok {
		// Record error metrics
		if errorMetrics != nil {
			errorMetrics.RecordError(customErr.ErrorCode)
		}

		errorResponse := map[string]interface{}{
			"statusCode":  customErr.StatusCode,
			"message":     customErr.Message,
			"description": customErr.Description,
			"errorCode":   customErr.ErrorCode,
		}
		c.JSON(customErr.StatusCode, errorResponse)
		c.Abort()
		return
	}

	// Handle generic errors
	errorResponse := map[string]interface{}{
		"statusCode": status,
		"message":    err.Error(),
	}
	c.JSON(status, errorResponse)
	c.Abort()
}

// GetError returns a standardized error response
func GetError(err CustomError) map[string]interface{} {
	errorResponse := map[string]interface{}{
		"statusCode":  err.StatusCode,
		"message":     err.Message,
		"description": err.Description,
		"errorCode":   err.ErrorCode,
	}
	return errorResponse
}

// GetErrorMetrics returns the current error metrics
func GetErrorMetrics() map[string]int64 {
	if errorMetrics == nil {
		return make(map[string]int64)
	}
	return errorMetrics.GetErrorCounts()
}

// ResetErrorMetrics resets the error metrics
func ResetErrorMetrics() {
	if errorMetrics != nil {
		errorMetrics.Reset()
	}
}
