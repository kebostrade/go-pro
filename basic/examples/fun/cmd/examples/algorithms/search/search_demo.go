package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/algorithms"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Search Algorithms Demo")

	demo1BinarySearch()
	demo2LinearSearch()
	demo3FindMinMax()
	demo4SearchWithPredicate()
}

func demo1BinarySearch() {
	utils.PrintSubHeader("1. Binary Search")

	// Sorted array
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	fmt.Printf("Sorted array: %v\n", arr)

	// Search for existing elements
	targets := []int{7, 15, 1, 19}
	for _, target := range targets {
		index := algorithms.BinarySearch(arr, target)
		if index != -1 {
			fmt.Printf("Found %d at index %d\n", target, index)
		}
	}

	// Search for non-existing element
	target := 10
	index := algorithms.BinarySearch(arr, target)
	if index == -1 {
		fmt.Printf("%d not found in array\n", target)
	}

	// Recursive binary search
	fmt.Println("\nUsing recursive binary search:")
	target = 13
	index = algorithms.BinarySearchRecursive(arr, target, 0, len(arr)-1)
	fmt.Printf("Found %d at index %d\n", target, index)
}

func demo2LinearSearch() {
	utils.PrintSubHeader("2. Linear Search")

	// Unsorted array
	arr := []string{"apple", "banana", "cherry", "date", "elderberry"}
	fmt.Printf("Array: %v\n", arr)

	// Search for element
	target := "cherry"
	index := algorithms.LinearSearch(arr, target)
	if index != -1 {
		fmt.Printf("Found '%s' at index %d\n", target, index)
	}

	// Search for all occurrences
	arr2 := []int{1, 3, 5, 3, 7, 3, 9}
	fmt.Printf("\nArray with duplicates: %v\n", arr2)
	indices := algorithms.LinearSearchAll(arr2, 3)
	fmt.Printf("Found 3 at indices: %v\n", indices)
}

func demo3FindMinMax() {
	utils.PrintSubHeader("3. Find Min/Max")

	arr := []int{42, 17, 93, 8, 56, 31, 74}
	fmt.Printf("Array: %v\n", arr)

	// Find minimum
	minVal, minIdx := algorithms.FindMin(arr)
	fmt.Printf("Minimum: %d at index %d\n", minVal, minIdx)

	// Find maximum
	maxVal, maxIdx := algorithms.FindMax(arr)
	fmt.Printf("Maximum: %d at index %d\n", maxVal, maxIdx)

	// Find both in one pass
	minV, minI, maxV, maxI := algorithms.FindMinMax(arr)
	fmt.Printf("Min: %d (index %d), Max: %d (index %d)\n", minV, minI, maxV, maxI)
}

func demo4SearchWithPredicate() {
	utils.PrintSubHeader("4. Search with Predicate")

	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
		{Name: "Diana", Age: 28},
	}

	fmt.Println("People:", people)

	// Find first person over 30
	index, found := algorithms.SearchWithPredicate(people, func(p Person) bool {
		return p.Age > 30
	})

	if found {
		fmt.Printf("\nFirst person over 30: %s (age %d) at index %d\n",
			people[index].Name, people[index].Age, index)
	}

	// Find all people under 30
	indices := algorithms.SearchAllWithPredicate(people, func(p Person) bool {
		return p.Age < 30
	})

	fmt.Println("\nPeople under 30:")
	for _, idx := range indices {
		fmt.Printf("  %s (age %d)\n", people[idx].Name, people[idx].Age)
	}
}
