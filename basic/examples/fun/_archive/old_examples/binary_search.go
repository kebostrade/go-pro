//go:build ignore

package main

import (
	"fmt"
	"sort"
	"strings"
)

// Task: Implement binary search algorithm with both iterative and recursive approaches.
// Binary search works on sorted arrays and has O(log n) time complexity.

// BinarySearchIterative performs binary search using iteration
func BinarySearchIterative(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		}

		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1 // Not found
}

// BinarySearchRecursive performs binary search using recursion
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if left > right {
		return -1 // Not found
	}

	mid := left + (right-left)/2

	if arr[mid] == target {
		return mid
	}

	if arr[mid] < target {
		return BinarySearchRecursive(arr, target, mid+1, right)
	}

	return BinarySearchRecursive(arr, target, left, mid-1)
}

// FindFirstOccurrence finds the first occurrence of target in sorted array
func FindFirstOccurrence(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			result = mid
			right = mid - 1 // Continue searching in left half
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}

// FindLastOccurrence finds the last occurrence of target in sorted array
func FindLastOccurrence(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			result = mid
			left = mid + 1 // Continue searching in right half
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}

// CountOccurrences counts how many times target appears in sorted array
func CountOccurrences(arr []int, target int) int {
	first := FindFirstOccurrence(arr, target)
	if first == -1 {
		return 0
	}

	last := FindLastOccurrence(arr, target)
	return last - first + 1
}

// FindInsertPosition finds the position where target should be inserted
func FindInsertPosition(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return left
}

// SearchInRotatedArray searches in a rotated sorted array
func SearchInRotatedArray(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		}

		// Determine which half is sorted
		if arr[left] <= arr[mid] {
			// Left half is sorted
			if target >= arr[left] && target < arr[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			// Right half is sorted
			if target > arr[mid] && target <= arr[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return -1
}

func main() {
	fmt.Println("Binary Search Algorithms Demo")
	fmt.Println(strings.Repeat("=", 60))

	// Demo 1: Basic Binary Search
	fmt.Println("\n1. Basic Binary Search (Iterative)")
	fmt.Println(strings.Repeat("-", 60))

	arr := []int{2, 5, 8, 12, 16, 23, 38, 45, 56, 67, 78}
	fmt.Printf("Array: %v\n", arr)

	targets := []int{23, 50, 2, 78}
	for _, target := range targets {
		index := BinarySearchIterative(arr, target)
		if index != -1 {
			fmt.Printf("Found %d at index %d\n", target, index)
		} else {
			fmt.Printf("%d not found in array\n", target)
		}
	}

	// Demo 2: Recursive Binary Search
	fmt.Println("\n2. Binary Search (Recursive)")
	fmt.Println(strings.Repeat("-", 60))

	target := 16
	index := BinarySearchRecursive(arr, target, 0, len(arr)-1)
	if index != -1 {
		fmt.Printf("Found %d at index %d (recursive)\n", target, index)
	}

	// Demo 3: Find First and Last Occurrence
	fmt.Println("\n3. Find First and Last Occurrence")
	fmt.Println(strings.Repeat("-", 60))

	arrWithDuplicates := []int{1, 2, 2, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7}
	fmt.Printf("Array: %v\n", arrWithDuplicates)

	searchValue := 6
	first := FindFirstOccurrence(arrWithDuplicates, searchValue)
	last := FindLastOccurrence(arrWithDuplicates, searchValue)
	count := CountOccurrences(arrWithDuplicates, searchValue)

	fmt.Printf("Value %d:\n", searchValue)
	fmt.Printf("  First occurrence at index: %d\n", first)
	fmt.Printf("  Last occurrence at index: %d\n", last)
	fmt.Printf("  Total occurrences: %d\n", count)

	// Demo 4: Find Insert Position
	fmt.Println("\n4. Find Insert Position")
	fmt.Println(strings.Repeat("-", 60))

	sortedArr := []int{1, 3, 5, 7, 9, 11}
	fmt.Printf("Array: %v\n", sortedArr)

	insertValues := []int{4, 0, 12, 7}
	for _, val := range insertValues {
		pos := FindInsertPosition(sortedArr, val)
		fmt.Printf("Insert %d at position %d\n", val, pos)
	}

	// Demo 5: Search in Rotated Array
	fmt.Println("\n5. Search in Rotated Sorted Array")
	fmt.Println(strings.Repeat("-", 60))

	rotatedArr := []int{15, 18, 2, 3, 6, 12}
	fmt.Printf("Rotated array: %v\n", rotatedArr)

	rotatedTargets := []int{6, 18, 2, 20}
	for _, target := range rotatedTargets {
		index := SearchInRotatedArray(rotatedArr, target)
		if index != -1 {
			fmt.Printf("Found %d at index %d\n", target, index)
		} else {
			fmt.Printf("%d not found in array\n", target)
		}
	}

	// Demo 6: Performance Comparison
	fmt.Println("\n6. Performance Comparison (Large Array)")
	fmt.Println(strings.Repeat("-", 60))

	// Create a large sorted array
	largeArr := make([]int, 1000000)
	for i := range largeArr {
		largeArr[i] = i * 2
	}

	searchTarget := 999998
	fmt.Printf("Searching for %d in array of %d elements\n", searchTarget, len(largeArr))

	// Binary search
	index = BinarySearchIterative(largeArr, searchTarget)
	if index != -1 {
		fmt.Printf("Binary Search: Found at index %d\n", index)
	}

	// Using Go's built-in sort.Search
	index = sort.Search(len(largeArr), func(i int) bool {
		return largeArr[i] >= searchTarget
	})
	if index < len(largeArr) && largeArr[index] == searchTarget {
		fmt.Printf("sort.Search: Found at index %d\n", index)
	}

	fmt.Println("\nBinary Search demo completed!")
}
