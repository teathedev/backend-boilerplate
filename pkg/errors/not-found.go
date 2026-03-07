package errors

import "github.com/teathedev/backend-boilerplate/internal/ent"

func NewNotFoundError(module string, message string) *CustomError {
	return &CustomError{
		Module:  module,
		Message: message,
		Params:  nil,
		Status:  404,
	}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	if ent.IsNotFound(err) {
		return true
	}

	// Type assert the error to *CustomError
	customErr, ok := err.(*CustomError)
	if !ok {
		return false
	}

	// Check if the status is 404
	return customErr.Status == 404
}
