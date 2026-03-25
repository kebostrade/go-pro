package exercises

import (
	"errors"
	"fmt"
)

// Exercise 1: Basic error
// Create a function that divides two numbers and returns error if divisor is zero
func Divide(a, b float64) (float64, error) {
	// TODO: Return error if b == 0
	return 0, nil
}

// Exercise 2: Custom error type
// Define a custom error type "ValidationError" with Field and Message
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	// TODO: Return formatted error message
	return ""
}

// Exercise 3: Sentinel errors
// Define sentinel errors for different validation cases
var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrEmptyUsername   = errors.New("username cannot be empty")
)

// ValidateUser validates user input and returns ValidationError
func ValidateUser(email, username, password string) error {
	// TODO: Check each validation and return appropriate error
	// - Empty username -> ErrEmptyUsername
	// - Invalid email -> ErrInvalidEmail  
	// - Short password -> ErrPasswordTooShort
	return nil
}

// Exercise 4: Error wrapping
// Create a function that wraps errors with context
func ReadConfig(filename string) ([]byte, error) {
	// TODO: Simulate reading file and wrap any error with fmt.Errorf
	// Return error like: fmt.Errorf("failed to read config: %w", originalErr)
	return nil, nil
}

// Exercise 5: Error unwrapping
// Demonstrate error unwrapping
func UnwrapError(err error) string {
	// TODO: Use errors.Unwrap to get underlying error
	// Return the unwrapped error message or "no wrapped error"
	return ""
}

// Exercise 6: Panic and recover
// Implement a safe division that recovers from panic
func SafeDivide(a, b float64) (result float64, err error) {
	// TODO: Use defer with recover to catch panic
	// Return error if b is zero
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic recovered: division by zero")
			result = 0
		}
	}()
	
	// TODO: Perform division (may panic)
	return 0, nil
}

// Exercise 7: Error types assertion
// Check if error is of specific type
func IsValidationError(err error) bool {
	// TODO: Use errors.As to check if error is ValidationError
	return false
}

// Exercise 8: Multiple error handling
// Collect multiple errors
func ValidateForm(name, email, phone string) error {
	// TODO: Collect multiple errors using errors.Join
	// Return all errors together
	return nil
}

// Exercise 9: Custom error with context
// Create an error that includes stack trace info
type StackError struct {
	Message string
	Stack   string
}

func NewStackError(msg string) error {
	// TODO: Create error with stack information
	return nil
}

func (e StackError) Error() string {
	return e.Message
}

// Exercise 10: Error handling best practices
// Rewrite function with proper error handling
func ProcessData(data string) (string, error) {
	// TODO: Demonstrate proper error handling patterns:
	// - Check errors early
	// - Wrap errors with context
	// - Return specific errors
	return "", nil
}
