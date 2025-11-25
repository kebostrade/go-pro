//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task: Implement merge sort algorithm with both sequential and concurrent versions.
// Merge sort is a divide-and-conquer algorithm with O(n log n) time complexity.

// MergeSortSequential performs merge sort sequentially
func MergeSortSequential(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := MergeSortSequential(arr[:mid])
	right := MergeSortSequential(arr[mid:])

	return merge(left, right)
}

// merge combines two sorted arrays into one sorted array
func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Append remaining elements
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

// MergeSortConcurrent performs merge sort using goroutines
func MergeSortConcurrent(arr []int, depth int) []int {
	if len(arr) <= 1 {
		return arr
	}

	// Use sequential sort for small arrays or deep recursion
	if len(arr) < 1000 || depth <= 0 {
		return MergeSortSequential(arr)
	}

	mid := len(arr) / 2
	var left, right []int
	var wg sync.WaitGroup

	wg.Add(2)

	// Sort left half concurrently
	go func() {
		defer wg.Done()
		left = MergeSortConcurrent(arr[:mid], depth-1)
	}()

	// Sort right half concurrently
	go func() {
		defer wg.Done()
		right = MergeSortConcurrent(arr[mid:], depth-1)
	}()

	wg.Wait()

	return merge(left, right)
}

// IsSorted checks if an array is sorted
func IsSorted(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}

// generateRandomArray creates an array of random integers
func generateRandomArray(size int, max int) []int {
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(max)
	}
	return arr
}

// copyArray creates a copy of an array
func copyArray(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	return result
}

func main() {
	fmt.Println("Merge Sort Algorithm Demo")
	fmt.Println("=".repeat(60))

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Demo 1: Basic Merge Sort
	fmt.Println("\n1. Basic Sequential Merge Sort")
	fmt.Println("-".repeat(60))

	arr1 := []int{38, 27, 43, 3, 9, 82, 10}
	fmt.Printf("Original array: %v\n", arr1)

	sorted1 := MergeSortSequential(copyArray(arr1))
	fmt.Printf("Sorted array:   %v\n", sorted1)
	fmt.Printf("Is sorted: %v\n", IsSorted(sorted1))

	// Demo 2: Larger Array
	fmt.Println("\n2. Sorting Larger Array")
	fmt.Println("-".repeat(60))

	arr2 := generateRandomArray(20, 100)
	fmt.Printf("Original (20 elements): %v\n", arr2)

	sorted2 := MergeSortSequential(copyArray(arr2))
	fmt.Printf("Sorted:                 %v\n", sorted2)

	// Demo 3: Performance Comparison
	fmt.Println("\n3. Performance Comparison (Sequential vs Concurrent)")
	fmt.Println("-".repeat(60))

	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		arr := generateRandomArray(size, 10000)

		// Sequential sort
		arrCopy1 := copyArray(arr)
		start := time.Now()
		MergeSortSequential(arrCopy1)
		seqDuration := time.Since(start)

		// Concurrent sort
		arrCopy2 := copyArray(arr)
		start = time.Now()
		MergeSortConcurrent(arrCopy2, 4) // depth of 4 for concurrency
		concDuration := time.Since(start)

		fmt.Printf("\nArray size: %d elements\n", size)
		fmt.Printf("  Sequential: %v\n", seqDuration)
		fmt.Printf("  Concurrent: %v\n", concDuration)
		if seqDuration > concDuration {
			speedup := float64(seqDuration) / float64(concDuration)
			fmt.Printf("  Speedup: %.2fx faster\n", speedup)
		}
	}

	// Demo 4: Edge Cases
	fmt.Println("\n4. Edge Cases")
	fmt.Println("-".repeat(60))

	// Empty array
	empty := []int{}
	fmt.Printf("Empty array: %v -> %v\n", empty, MergeSortSequential(empty))

	// Single element
	single := []int{42}
	fmt.Printf("Single element: %v -> %v\n", single, MergeSortSequential(single))

	// Already sorted
	sorted := []int{1, 2, 3, 4, 5}
	fmt.Printf("Already sorted: %v -> %v\n", sorted, MergeSortSequential(copyArray(sorted)))

	// Reverse sorted
	reverse := []int{5, 4, 3, 2, 1}
	fmt.Printf("Reverse sorted: %v -> %v\n", reverse, MergeSortSequential(copyArray(reverse)))

	// Duplicates
	duplicates := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	fmt.Printf("With duplicates: %v -> %v\n", duplicates, MergeSortSequential(copyArray(duplicates)))

	// Demo 5: Stability Test
	fmt.Println("\n5. Merge Sort is Stable")
	fmt.Println("-".repeat(60))
	fmt.Println("Merge sort maintains the relative order of equal elements")

	type Item struct {
		Value int
		Index int
	}

	items := []Item{
		{3, 1}, {1, 2}, {3, 3}, {2, 4}, {1, 5},
	}

	fmt.Println("Original items (Value, OriginalIndex):")
	for _, item := range items {
		fmt.Printf("  (%d, %d)\n", item.Value, item.Index)
	}

	// Extract values for sorting
	values := make([]int, len(items))
	for i, item := range items {
		values[i] = item.Value
	}

	sortedValues := MergeSortSequential(values)
	fmt.Println("\nSorted values:", sortedValues)
	fmt.Println("(Equal values maintain their relative order)")

	fmt.Println("\nMerge Sort demo completed!")
}

// Helper function to repeat strings
func (s string) repeat(count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
