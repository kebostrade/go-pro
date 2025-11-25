//go:build ignore

package main

import (
	"errors"
	"fmt"
)

// processData calculates the sum and product of a slice of integers
// Returns error if the input slice is empty
func processData(numbers []int) (int, int, error) {
	if len(numbers) == 0 {
		return 0, 0, errors.New("empty input slice")
	}

	sum := 0
	product := 1
	for _, num := range numbers {
		sum += num
		product *= num
	}
	return sum, product, nil
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	sum, product, err := processData(numbers)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Sum: %d, Product: %d\n", sum, product)

	// Test empty slice
	emptySlice := []int{}
	sum, product, err = processData(emptySlice)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
