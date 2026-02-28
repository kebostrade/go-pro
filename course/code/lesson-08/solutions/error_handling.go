package exercises

import (
	"errors"
	"fmt"
)

// Exercise 1: Basic error
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Exercise 3: ValidateUser
func ValidateUser(email, username, password string) error {
	if username == "" {
		return ErrEmptyUsername
	}
	
	// Simple email validation
	if email == "" || email[len(email)-4:] != ".com" && email[len(email)-3:] != ".ua" {
		return ErrInvalidEmail
	}
	
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	
	return nil
}

// Exercise 4: Error wrapping
func ReadConfig(filename string) ([]byte, error) {
	// Simulate error
	err := errors.New("file not found")
	return nil, fmt.Errorf("failed to read config: %w", err)
}

// Exercise 5: Error unwrapping
func UnwrapError(err error) string {
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return "no wrapped error"
	}
	return unwrapped.Error()
}

// Exercise 6: Panic and recover
func SafeDivide(a, b float64) (result float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic recovered: division by zero")
			result = 0
		}
	}()
	
	if b == 0 {
		panic("division by zero")
	}
	return a / b, nil
}

// Exercise 7: Error types assertion
func IsValidationError(err error) bool {
	var ve ValidationError
	return errors.As(err, &ve)
}

// Exercise 8: Multiple error handling
func ValidateForm(name, email, phone string) error {
	var errs []error
	
	if name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	if email == "" {
		errs = append(errs, errors.New("email is required"))
	}
	if phone == "" {
		errs = append(errs, errors.New("phone is required"))
	}
	
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// Exercise 9: Custom error with context
func NewStackError(msg string) error {
	return StackError{
		Message: msg,
		Stack:   "stack trace here", // In real code, use runtime.Callers
	}
}

// Exercise 10: Error handling best practices
func ProcessData(data string) (string, error) {
	if data == "" {
		return "", errors.New("data cannot be empty")
	}
	
	// Process data...
	// If error occurs:
	// return "", fmt.Errorf("failed to process data: %w", err)
	
	return "processed: " + data, nil
}
