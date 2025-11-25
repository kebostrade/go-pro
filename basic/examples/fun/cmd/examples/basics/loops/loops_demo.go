package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Loops & Control Flow Demo")

	demo1BasicForLoop()
	demo2WhileStyleLoop()
	demo3InfiniteLoop()
	demo4RangeLoop()
	demo5NestedLoops()
	demo6ControlStatements()
}

func demo1BasicForLoop() {
	utils.PrintSubHeader("1. Basic For Loop")

	fmt.Println("Count from 1 to 5:")
	for i := 1; i <= 5; i++ {
		fmt.Printf("  %d\n", i)
	}

	fmt.Println("\nCount down from 5 to 1:")
	for i := 5; i >= 1; i-- {
		fmt.Printf("  %d\n", i)
	}

	fmt.Println("\nEven numbers from 0 to 10:")
	for i := 0; i <= 10; i += 2 {
		fmt.Printf("  %d\n", i)
	}
}

func demo2WhileStyleLoop() {
	utils.PrintSubHeader("2. While-Style Loop")

	// Go doesn't have 'while', use 'for' instead
	count := 1
	fmt.Println("While-style loop (count to 5):")
	for count <= 5 {
		fmt.Printf("  Count: %d\n", count)
		count++
	}
}

func demo3InfiniteLoop() {
	utils.PrintSubHeader("3. Infinite Loop with Break")

	count := 0
	fmt.Println("Infinite loop with break condition:")
	for {
		count++
		fmt.Printf("  Iteration: %d\n", count)

		if count >= 5 {
			fmt.Println("  Breaking out of loop")
			break
		}
	}
}

func demo4RangeLoop() {
	utils.PrintSubHeader("4. Range Loop")

	// Range over slice
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Println("Range over slice:")
	for index, value := range numbers {
		fmt.Printf("  Index: %d, Value: %d\n", index, value)
	}

	// Range with only value
	fmt.Println("\nRange with only values:")
	for _, value := range numbers {
		fmt.Printf("  %d\n", value)
	}

	// Range over string
	text := "Hello"
	fmt.Println("\nRange over string:")
	for index, char := range text {
		fmt.Printf("  Index: %d, Char: %c (rune: %d)\n", index, char, char)
	}

	// Range over map
	scores := map[string]int{
		"Alice": 95,
		"Bob":   87,
		"Carol": 92,
	}
	fmt.Println("\nRange over map:")
	for name, score := range scores {
		fmt.Printf("  %s: %d\n", name, score)
	}
}

func demo5NestedLoops() {
	utils.PrintSubHeader("5. Nested Loops")

	// Multiplication table
	fmt.Println("Multiplication table (3x3):")
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("%4d", i*j)
		}
		fmt.Println()
	}

	// Pattern printing
	fmt.Println("\nTriangle pattern:")
	rows := 5
	for i := 1; i <= rows; i++ {
		// Print spaces
		for j := 1; j <= rows-i; j++ {
			fmt.Print(" ")
		}
		// Print stars
		for k := 1; k <= 2*i-1; k++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}

func demo6ControlStatements() {
	utils.PrintSubHeader("6. Control Statements (break, continue)")

	// Continue - skip to next iteration
	fmt.Println("Print odd numbers (1-10) using continue:")
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			continue // Skip even numbers
		}
		fmt.Printf("  %d\n", i)
	}

	// Break - exit loop
	fmt.Println("\nFind first number divisible by 7:")
	for i := 1; i <= 100; i++ {
		if i%7 == 0 {
			fmt.Printf("  Found: %d\n", i)
			break
		}
	}

	// Labeled break (break out of nested loop)
	fmt.Println("\nLabeled break example:")
outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("  i=%d, j=%d\n", i, j)
			if i == 2 && j == 2 {
				fmt.Println("  Breaking out of both loops")
				break outer
			}
		}
	}

	// Labeled continue
	fmt.Println("\nLabeled continue example:")
outerLoop:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if j == 2 {
				continue outerLoop
			}
			fmt.Printf("  i=%d, j=%d\n", i, j)
		}
	}
}
