package exercises

import (
	"errors"
	"fmt"
)

// Suppress unused import warning - fmt is used in exercise solutions
var _ = fmt.Sprint

// Exercise 2: Functions Practice
// Complete the following functions to practice with Go function patterns

// SimpleGreeting creates a greeting message
// Parameter: name (string)
// Returns: greeting message (string)
// Format: "Hello, [name]! Welcome to Go programming."
func SimpleGreeting(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome to Go programming.", name)
}

// Calculator performs basic arithmetic operations
// Parameters: a, b (both int), operation (string)
// Returns: result (int), error (error)
// Supported operations: "add", "subtract", "multiply", "divide"
// Return error for unsupported operations or division by zero
func Calculator(a, b int, operation string) (int, error) {
	switch operation {
	case "add":
		return a + b, nil
	case "subtract":
		return a - b, nil
	case "multiply":
		return a * b, nil
	case "divide":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", operation)
	}
}

// MultipleReturns demonstrates functions with multiple return values
// Parameters: x, y (both float64)
// Returns: sum, difference, product, quotient (all float64)
func MultipleReturns(x, y float64) (float64, float64, float64, float64) {
	sum := x + y
	difference := x - y
	product := x * y
	var quotient float64
	if y != 0 {
		quotient = x / y
	}
	return sum, difference, product, quotient
}

// NamedReturns uses named return values to calculate rectangle properties
// Parameters: length, width (both float64)
// Returns: area, perimeter (both float64) - use named returns
func NamedReturns(length, width float64) (area, perimeter float64) {
	area = length * width
	perimeter = 2 * (length + width)
	return // naked return
}

// VariadicSum calculates the sum of variable number of integers
// Parameters: numbers (...int) - variadic parameter
// Returns: sum of all numbers (int)
func VariadicSum(numbers ...int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// VariadicAverage calculates the average of variable number of float64 values
// Parameters: values (...float64) - variadic parameter
// Returns: average (float64), count (int)
// Return 0.0, 0 if no values provided
func VariadicAverage(values ...float64) (float64, int) {
	if len(values) == 0 {
		return 0.0, 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values)), len(values)
}

// StringJoiner joins strings with a separator
// Parameters: separator (string), strings (...string) - variadic parameter
// Returns: joined string (string)
// Example: StringJoiner("-", "a", "b", "c") should return "a-b-c"
func StringJoiner(separator string, strings ...string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}
	return result
}

// FunctionAsParameter demonstrates using functions as parameters
// Parameters: a, b (both int), operation (func(int, int) int)
// Returns: result of applying the operation function to a and b (int)
func FunctionAsParameter(a, b int, operation func(int, int) int) int {
	return operation(a, b)
}

// ReturnFunction returns a function that adds a fixed value
// Parameter: addValue (int)
// Returns: a function that takes an int and returns an int
// The returned function should add addValue to its parameter
func ReturnFunction(addValue int) func(int) int {
	return func(x int) int {
		return x + addValue
	}
}

// Closure demonstrates closures by creating a counter
// Returns: a function that increments and returns a counter value
// Each call to the returned function should increment the counter
func Closure() func() int {
	counter := 0
	return func() int {
		counter++
		return counter
	}
}

// ErrorHandling demonstrates proper error handling
// Parameters: dividend, divisor (both float64)
// Returns: result (float64), error (error)
// Return error if divisor is zero with message "division by zero"
func ErrorHandling(dividend, divisor float64) (float64, error) {
	if divisor == 0 {
		return 0.0, errors.New("division by zero")
	}
	return dividend / divisor, nil
}

// MultipleErrorReturns demonstrates multiple operations with error handling
// Parameters: a, b (both int)
// Returns: sum, product (both int), error (error)
// Return error if either a or b is negative with message "negative numbers not allowed"
func MultipleErrorReturns(a, b int) (int, int, error) {
	if a < 0 || b < 0 {
		return 0, 0, errors.New("negative numbers not allowed")
	}
	return a + b, a * b, nil
}

// RecursiveFactorial calculates factorial using recursion
// Parameter: n (int)
// Returns: factorial of n (int)
// Return 1 for n <= 1
func RecursiveFactorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * RecursiveFactorial(n-1)
}

// RecursiveFibonacci calculates Fibonacci number using recursion
// Parameter: n (int)
// Returns: nth Fibonacci number (int)
// Fibonacci sequence: 0, 1, 1, 2, 3, 5, 8, 13, ...
func RecursiveFibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return RecursiveFibonacci(n-1) + RecursiveFibonacci(n-2)
}

// Helper functions for testing (you can use these in your implementations)

// Add function for testing FunctionAsParameter
func Add(a, b int) int {
	return a + b
}

// Multiply function for testing FunctionAsParameter
func Multiply(a, b int) int {
	return a * b
}

// HigherOrderFunction demonstrates a function that takes and returns functions
// Parameter: transform (func(int) int) - a function that transforms an int
// Returns: a function that applies transform twice to its input
func HigherOrderFunction(transform func(int) int) func(int) int {
	return func(x int) int {
		return transform(transform(x))
	}
}
