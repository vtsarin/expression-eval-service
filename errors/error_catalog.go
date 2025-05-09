package errors

// Predefined errors catalog
var (
	/*
		400
	*/
	ErrInvalidRequestBody = NewCustomError(400, "E4007700", "Request body is invalid", "The request is invalid", nil)

	/**
	401
	*/
	ErrUnauthorized = NewCustomError(401, "E4017700", "Unauthorized", "Unauthorized", nil)

	/*
		404
	*/
	ErrUserNotFound = NewCustomError(404, "E4047700", "Resource not found", "Requested resource not found", nil)

	/*
		500
	*/
	ErrInternalError = NewCustomError(500, "E5007700", "Internal error", "Something went wrong", nil)
)
