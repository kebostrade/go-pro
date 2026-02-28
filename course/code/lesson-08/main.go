package main

import (
	"errors"
	"fmt"
	"lesson-08/exercises"
)

func main() {
	fmt.Println("=== Lesson 08: Error Handling Patterns ===")
	fmt.Println()

	// Exercise 1: Basic error
	fmt.Println("1. Basic Error:")
	result, err := exercises.Divide(10, 2)
	fmt.Printf("   Divide(10, 2) = %.2f, err = %v\n", result, err)
	
	result, err = exercises.Divide(10, 0)
	fmt.Printf("   Divide(10, 0) = %.2f, err = %v\n", result, err)
	fmt.Println()

	// Exercise 2 & 3: Custom and sentinel errors
	fmt.Println("2 & 3. Custom & Sentinel Errors:")
	err = exercises.ValidateUser("test@example.com", "john", "password123")
	fmt.Printf("   ValidateUser(valid) = %v\n", err)
	
	err = exercises.ValidateUser("", "john", "password123")
	fmt.Printf("   ValidateUser(empty email) = %v\n", err)
	
	err = exercises.ValidateUser("test@example.com", "john", "short")
	fmt.Printf("   ValidateUser(short password) = %v\n", err)
	fmt.Println()

	// Exercise 4: Error wrapping
	fmt.Println("4. Error Wrapping:")
	err = exercises.ReadConfig("config.json")
	fmt.Printf("   ReadConfig error: %v\n", err)
	fmt.Println()

	// Exercise 5: Error unwrapping
	fmt.Println("5. Error Unwrapping:")
	wrapped := fmt.Errorf("wrapped: %w", exercises.ErrInvalidEmail)
	unwrapped := exercises.UnwrapError(wrapped)
	fmt.Printf("   Unwrapped: %s\n", unwrapped)
	fmt.Println()

	// Exercise 6: Panic and recover
	fmt.Println("6. Panic & Recover:")
	result, err = exercises.SafeDivide(10, 2)
	fmt.Printf("   SafeDivide(10, 2) = %.2f, err = %v\n", result, err)
	
	result, err = exercises.SafeDivide(10, 0)
	fmt.Printf("   SafeDivide(10, 0) = %.2f, err = %v\n", result, err)
	fmt.Println()

	// Exercise 7: Error type assertion
	fmt.Println("7. Error Type Assertion:")
	customErr := exercises.ValidationError{Field: "email", Message: "invalid"}
	fmt.Printf("   IsValidationError(custom): %v\n", exercises.IsValidationError(customErr))
	fmt.Printf("   IsValidationError(regular): %v\n", exercises.IsValidationError(errors.New("regular")))
	fmt.Println()

	// Exercise 8: Multiple errors
	fmt.Println("8. Multiple Errors:")
	err = exercises.ValidateForm("", "", "")
	fmt.Printf("   ValidateForm(empty): %v\n", err)
	
	err = exercises.ValidateForm("John", "john@example.com", "1234567890")
	fmt.Printf("   ValidateForm(valid): %v\n", err)
	fmt.Println()

	// Exercise 9: Custom error with context
	fmt.Println("9. Custom Error with Context:")
	err = exercises.NewStackError("something went wrong")
	fmt.Printf("   StackError: %v\n", err)
	fmt.Println()

	// Exercise 10: Best practices
	fmt.Println("10. Error Handling Best Practices:")
	result, err = exercises.ProcessData("hello")
	fmt.Printf("   ProcessData(valid): %s, err = %v\n", result, err)
	
	result, err = exercises.ProcessData("")
	fmt.Printf("   ProcessData(empty): %s, err = %v\n", result, err)
	fmt.Println()

	fmt.Println("=== All exercises completed! ===")
}
