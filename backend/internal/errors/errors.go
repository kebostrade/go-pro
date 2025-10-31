// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package errors provides custom error types and error handling utilities.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error types for the application.
var (
	ErrNotFound           = errors.New("resource not found")
	ErrValidation         = errors.New("validation failed")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInternalServer     = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrConflict           = errors.New("resource conflict")
	ErrServiceUnavailable = errors.New("service unavailable")
)

// APIError represents an error in API operations.
type APIError struct {
	Type       string                 `json:"type"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Cause      error                  `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}

	return e.Message
}

// Unwrap returns the underlying error.
func (e *APIError) Unwrap() error {
	return e.Cause
}

// NewAPIError creates a new API error.
func NewAPIError(errType, message string, statusCode int, cause error) *APIError {
	return &APIError{
		Type:       errType,
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
	}
}

// Predefined API errors.
func NewNotFoundError(message string) *APIError {
	return &APIError{
		Type:       "NOT_FOUND",
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewValidationError(message string, cause error) *APIError {
	return &APIError{
		Type:       "VALIDATION_ERROR",
		Code:       "VALIDATION_ERROR",
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Cause:      cause,
	}
}

func NewInternalError(message string, cause error) *APIError {
	return &APIError{
		Type:       "INTERNAL_ERROR",
		Code:       "INTERNAL_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Cause:      cause,
	}
}

func NewBadRequestError(message string) *APIError {
	return &APIError{
		Type:       "BAD_REQUEST",
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewConflictError(message string) *APIError {
	return &APIError{
		Type:       "CONFLICT",
		Code:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewUnauthorizedError(message string) *APIError {
	return &APIError{
		Type:       "UNAUTHORIZED",
		Code:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *APIError {
	return &APIError{
		Type:       "FORBIDDEN",
		Code:       "FORBIDDEN",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewInternalServerError(message string) *APIError {
	return &APIError{
		Type:       "INTERNAL_SERVER_ERROR",
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// IsAPIError checks if an error is an APIError.
func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}

	return nil, false
}
