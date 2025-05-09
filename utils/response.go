package utils

import (
	"go-service/errors"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	StatusCode  int         `json:"statusCode"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	ErrorCode   string      `json:"errorCode,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

// SendResponse sends a standardized JSON response
func SendResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		StatusCode: status,
		Message:    message,
		Data:       data,
	})
}

// SendError sends a standardized error response
func SendError(c *gin.Context, err error) {
	// Use the error handler for consistent error responses
	errors.HandleError(c, 500, err)
}

// SendSuccess sends a success response with data
func SendSuccess(c *gin.Context, message string, data interface{}) {
	SendResponse(c, 200, message, data)
}

// SendNoContent sends a no content response
func SendNoContent(c *gin.Context) {
	c.Status(204)
}
