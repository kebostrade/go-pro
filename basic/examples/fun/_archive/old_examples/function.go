//go:build ignore

package main

import "fmt"

// Basic function with parameters and return value
func add(x, y int) int {
	return x + y
}

// Function with multiple return values
func divide(x, y float64) (float64, error) {
	if y == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return x / y, nil
}

// Function with named return values
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // naked return
}

// Variadic function
func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

// Function type declaration
type Operation func(int, int) int

// Higher-order function that returns a function
func getOperation(op string) Operation {
	switch op {
	case "+":
		return func(x, y int) int { return x + y }
	case "-":
		return func(x, y int) int { return x - y }
	default:
		return func(x, y int) int { return 0 }
	}
}

func main() {
	fmt.Println("Basic function:", add(5, 3))

	if result, err := divide(10, 2); err == nil {
		fmt.Println("Division result:", result)
	}

	x, y := split(17)
	fmt.Println("Split result:", x, y)

	fmt.Println("Variadic sum:", sum(1, 2, 3, 4, 5))

	plus := getOperation("+")
	fmt.Println("Operation result:", plus(10, 5))
}
