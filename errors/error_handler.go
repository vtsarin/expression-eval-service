package errors

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, status int, err error) {
	// Type assertion to check if the error is a CustomError
	if customErr, ok := err.(*CustomError); ok {
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

func GetError(err CustomError) map[string]interface{} {
	errorResponse := map[string]interface{}{
		"statusCode":  err.StatusCode,
		"message":     err.Message,
		"description": err.Description,
		"errorCode":   err.ErrorCode,
	}
	return errorResponse
}
