package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Functions Demo")

	demo1BasicFunctions()
	demo2MultipleReturns()
	demo3NamedReturns()
	demo4VariadicFunctions()
	demo5HigherOrderFunctions()
	demo6Closures()
}

func demo1BasicFunctions() {
	utils.PrintSubHeader("1. Basic Functions")

	result := addNumbers(5, 3)
	fmt.Printf("add(5, 3) = %d\n", result)

	greeting := greet("Alice")
	fmt.Println(greeting)
}

func addNumbers(x, y int) int {
	return x + y
}

func greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func demo2MultipleReturns() {
	utils.PrintSubHeader("2. Multiple Return Values")

	quotient, remainder := divideWithRemainder(17, 5)
	fmt.Printf("17 ÷ 5 = %d remainder %d\n", quotient, remainder)

	// Error handling pattern
	result, err := safeDivide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 ÷ 2 = %.2f\n", result)
	}

	result, err = safeDivide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func divideWithRemainder(a, b int) (int, int) {
	return a / b, a % b
}

func safeDivide(x, y float64) (float64, error) {
	if y == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return x / y, nil
}

func demo3NamedReturns() {
	utils.PrintSubHeader("3. Named Return Values")

	x, y := splitSum(17)
	fmt.Printf("split(17) = %d, %d\n", x, y)

	min, max := findMinMax([]int{3, 7, 2, 9, 1, 5})
	fmt.Printf("min = %d, max = %d\n", min, max)
}

func splitSum(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // naked return
}

func findMinMax(numbers []int) (min, max int) {
	if len(numbers) == 0 {
		return
	}

	min, max = numbers[0], numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}
	return
}

func demo4VariadicFunctions() {
	utils.PrintSubHeader("4. Variadic Functions")

	total := sumAll(1, 2, 3, 4, 5)
	fmt.Printf("sum(1, 2, 3, 4, 5) = %d\n", total)

	total = sumAll(10, 20, 30)
	fmt.Printf("sum(10, 20, 30) = %d\n", total)

	// Passing a slice
	numbers := []int{5, 10, 15, 20}
	total = sumAll(numbers...)
	fmt.Printf("sum(5, 10, 15, 20) = %d\n", total)

	// Multiple types
	printAll("Numbers:", 1, 2, 3, 4, 5)
}

func sumAll(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

func printAll(label string, values ...interface{}) {
	fmt.Printf("%s ", label)
	for i, val := range values {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(val)
	}
	fmt.Println()
}

func demo5HigherOrderFunctions() {
	utils.PrintSubHeader("5. Higher-Order Functions")

	// Function as a type
	type Operation func(int, int) int

	// Functions that return functions
	add := getOperation("+")
	subtract := getOperation("-")
	multiply := getOperation("*")

	fmt.Printf("10 + 5 = %d\n", add(10, 5))
	fmt.Printf("10 - 5 = %d\n", subtract(10, 5))
	fmt.Printf("10 * 5 = %d\n", multiply(10, 5))

	// Function that takes a function
	numbers := []int{1, 2, 3, 4, 5}
	doubled := applyToAll(numbers, func(x int) int { return x * 2 })
	fmt.Printf("Doubled: %v\n", doubled)
}

func getOperation(op string) func(int, int) int {
	switch op {
	case "+":
		return func(x, y int) int { return x + y }
	case "-":
		return func(x, y int) int { return x - y }
	case "*":
		return func(x, y int) int { return x * y }
	default:
		return func(x, y int) int { return 0 }
	}
}

func applyToAll(numbers []int, fn func(int) int) []int {
	result := make([]int, len(numbers))
	for i, num := range numbers {
		result[i] = fn(num)
	}
	return result
}

func demo6Closures() {
	utils.PrintSubHeader("6. Closures")

	// Closure captures variables from outer scope
	counter := makeCounter()

	fmt.Println("Counter:")
	fmt.Printf("  Call 1: %d\n", counter())
	fmt.Printf("  Call 2: %d\n", counter())
	fmt.Printf("  Call 3: %d\n", counter())

	// Another counter instance
	counter2 := makeCounter()
	fmt.Printf("  New counter: %d\n", counter2())

	// Closure with parameters
	multiplier := makeMultiplier(3)
	fmt.Printf("\nMultiplier by 3:\n")
	fmt.Printf("  5 * 3 = %d\n", multiplier(5))
	fmt.Printf("  10 * 3 = %d\n", multiplier(10))
}

func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}
