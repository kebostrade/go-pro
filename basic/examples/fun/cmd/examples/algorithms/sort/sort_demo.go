package main

import (
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/algorithms"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Sort Algorithms Demo")

	demo1MergeSort()
	demo2QuickSort()
	demo3CompareAlgorithms()
	demo4CustomSort()
}

func demo1MergeSort() {
	utils.PrintSubHeader("1. Merge Sort")

	arr := []int{38, 27, 43, 3, 9, 82, 10}
	fmt.Printf("Original: %v\n", arr)

	sorted := algorithms.MergeSort(arr)
	fmt.Printf("Sorted:   %v\n", sorted)
	fmt.Printf("Is sorted: %v\n", algorithms.IsSorted(sorted))

	// Concurrent merge sort
	fmt.Println("\nConcurrent merge sort on larger array:")
	largeArr := utils.GenerateRandomInts(1000, 1, 1000)

	start := time.Now()
	sorted = algorithms.MergeSortConcurrent(largeArr, 4)
	duration := time.Since(start)

	fmt.Printf("Sorted %d elements in %v\n", len(sorted), duration)
	fmt.Printf("Is sorted: %v\n", algorithms.IsSorted(sorted))
}

func demo2QuickSort() {
	utils.PrintSubHeader("2. Quick Sort")

	arr := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Original: %v\n", arr)

	algorithms.QuickSort(arr)
	fmt.Printf("Sorted:   %v\n", arr)
}

func demo3CompareAlgorithms() {
	utils.PrintSubHeader("3. Algorithm Comparison")

	sizes := []int{100, 500, 1000}

	for _, size := range sizes {
		fmt.Printf("\nArray size: %d\n", size)

		// Generate random array
		original := utils.GenerateRandomInts(size, 1, 1000)

		// Merge Sort
		arr1 := make([]int, len(original))
		copy(arr1, original)
		start := time.Now()
		algorithms.MergeSort(arr1)
		fmt.Printf("  Merge Sort:     %v\n", time.Since(start))

		// Quick Sort
		arr2 := make([]int, len(original))
		copy(arr2, original)
		start = time.Now()
		algorithms.QuickSort(arr2)
		fmt.Printf("  Quick Sort:     %v\n", time.Since(start))

		// Bubble Sort (only for small arrays)
		if size <= 500 {
			arr3 := make([]int, len(original))
			copy(arr3, original)
			start = time.Now()
			algorithms.BubbleSort(arr3)
			fmt.Printf("  Bubble Sort:    %v\n", time.Since(start))
		}
	}
}

func demo4CustomSort() {
	utils.PrintSubHeader("4. Custom Sort with Comparator")

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

	fmt.Println("Original:", people)

	// Sort by age
	algorithms.SortWithComparator(people, func(a, b Person) bool {
		return a.Age < b.Age
	})

	fmt.Println("Sorted by age:", people)

	// Sort by name
	algorithms.SortWithComparator(people, func(a, b Person) bool {
		return a.Name < b.Name
	})

	fmt.Println("Sorted by name:", people)
}
