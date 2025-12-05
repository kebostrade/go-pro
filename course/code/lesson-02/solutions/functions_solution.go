package solutions

import (
	"errors"
	"fmt"
	"strings"
)

// SimpleGreeting creates a greeting message
func SimpleGreeting(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome to Go programming.", name)
}

// Calculator performs basic arithmetic operations
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
func NamedReturns(length, width float64) (area, perimeter float64) {
	area = length * width
	perimeter = 2 * (length + width)
	return
}

// VariadicSum calculates the sum of variable number of integers
func VariadicSum(numbers ...int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// VariadicAverage calculates the average of variable number of float64 values
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
func FunctionAsParameter(a, b int, operation func(int, int) int) int {
	return operation(a, b)
}

// ReturnFunction returns a function that adds a fixed value
func ReturnFunction(addValue int) func(int) int {
	return func(x int) int {
		return x + addValue
	}
}

// Closure demonstrates closures by creating a counter
func Closure() func() int {
	counter := 0
	return func() int {
		counter++
		return counter
	}
}

// ErrorHandling demonstrates proper error handling
func ErrorHandling(dividend, divisor float64) (float64, error) {
	if divisor == 0 {
		return 0.0, errors.New("division by zero")
	}
	return dividend / divisor, nil
}

// MultipleErrorReturns demonstrates multiple operations with error handling
func MultipleErrorReturns(a, b int) (int, int, error) {
	if a < 0 || b < 0 {
		return 0, 0, errors.New("negative numbers not allowed")
	}
	return a + b, a * b, nil
}

// RecursiveFactorial calculates factorial using recursion
func RecursiveFactorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * RecursiveFactorial(n-1)
}

// RecursiveFibonacci calculates Fibonacci number using recursion
func RecursiveFibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return RecursiveFibonacci(n-1) + RecursiveFibonacci(n-2)
}

// HigherOrderFunction demonstrates a function that takes and returns functions
func HigherOrderFunction(transform func(int) int) func(int) int {
	return func(x int) int {
		return transform(transform(x))
	}
}

// Helper functions for testing
func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}
