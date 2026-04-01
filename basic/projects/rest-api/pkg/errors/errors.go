package errors

import "errors"

// Sentinel errors for the application
var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrValidation is returned when input validation fails
	ErrValidation = errors.New("validation error")

	// ErrInternal is returned for internal server errors
	ErrInternal = errors.New("internal server error")
)
