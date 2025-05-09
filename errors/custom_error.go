package errors

import (
	"fmt"
)

type CustomError struct {
	StatusCode    int
	ErrorCode     string
	Message       string
	Description   string
	InternalError error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("StatusCode: %d, ErrorCode: %s, Message: %s, Description: %s",
		e.StatusCode, e.ErrorCode, e.Message, e.Description)
}

func NewCustomError(statusCode int, errorCode string, message, description string, internalError error) *CustomError {
	return &CustomError{
		StatusCode:    statusCode,
		ErrorCode:     errorCode,
		Message:       message,
		Description:   description,
		InternalError: internalError,
	}
}
