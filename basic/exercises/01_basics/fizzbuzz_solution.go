//go:build ignore

package main

import (
	"fmt"
	"strings"
)

// Solution: FizzBuzz

func FizzBuzzSolution(n int) string {
	if n%15 == 0 {
		return "FizzBuzz"
	} else if n%3 == 0 {
		return "Fizz"
	} else if n%5 == 0 {
		return "Buzz"
	}
	return fmt.Sprintf("%d", n)
}

func main() {
	fmt.Println("FizzBuzz Solution")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()

	for i := 1; i <= 100; i++ {
		result := FizzBuzzSolution(i)
		fmt.Printf("%3d: %s\n", i, result)
	}
}
