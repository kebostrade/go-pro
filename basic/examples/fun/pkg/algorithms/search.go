package algorithms

import "golang.org/x/exp/constraints"

// BinarySearch performs binary search on a sorted slice using iteration
// Returns the index of the target element, or -1 if not found
// Time complexity: O(log n)
// Space complexity: O(1)
func BinarySearch[T constraints.Ordered](arr []T, target T) int {
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
// Returns the index of the target element, or -1 if not found
// Time complexity: O(log n)
// Space complexity: O(log n) due to recursion stack
func BinarySearchRecursive[T constraints.Ordered](arr []T, target T, left, right int) int {
	// Bounds validation
	if left < 0 || right >= len(arr) || left > right {
		return -1
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

// FindFirstOccurrence finds the first occurrence of target in a sorted array
// Returns the index of the first occurrence, or -1 if not found
// Time complexity: O(log n)
func FindFirstOccurrence[T constraints.Ordered](arr []T, target T) int {
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

// FindLastOccurrence finds the last occurrence of target in a sorted array
// Returns the index of the last occurrence, or -1 if not found
// Time complexity: O(log n)
func FindLastOccurrence[T constraints.Ordered](arr []T, target T) int {
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

// CountOccurrences counts the number of occurrences of target in a sorted array
// Time complexity: O(log n)
func CountOccurrences[T constraints.Ordered](arr []T, target T) int {
	first := FindFirstOccurrence(arr, target)
	if first == -1 {
		return 0
	}

	last := FindLastOccurrence(arr, target)
	return last - first + 1
}

// LinearSearch performs a linear search on an unsorted slice
// Returns the index of the target element, or -1 if not found
// Time complexity: O(n)
// Space complexity: O(1)
func LinearSearch[T comparable](arr []T, target T) int {
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1
}

// LinearSearchAll finds all occurrences of target in a slice
// Returns a slice of indices where the target was found
// Time complexity: O(n)
func LinearSearchAll[T comparable](arr []T, target T) []int {
	indices := make([]int, 0)
	for i, v := range arr {
		if v == target {
			indices = append(indices, i)
		}
	}
	return indices
}

// SearchWithPredicate finds the first element that satisfies the predicate
// Returns the index and true if found, -1 and false otherwise
// Time complexity: O(n)
func SearchWithPredicate[T any](arr []T, predicate func(T) bool) (int, bool) {
	for i, v := range arr {
		if predicate(v) {
			return i, true
		}
	}
	return -1, false
}

// SearchAllWithPredicate finds all elements that satisfy the predicate
// Returns a slice of indices
// Time complexity: O(n)
func SearchAllWithPredicate[T any](arr []T, predicate func(T) bool) []int {
	indices := make([]int, 0)
	for i, v := range arr {
		if predicate(v) {
			indices = append(indices, i)
		}
	}
	return indices
}

// FindMin finds the minimum element in a slice
// Returns the minimum element and its index
// Time complexity: O(n)
func FindMin[T constraints.Ordered](arr []T) (T, int) {
	if len(arr) == 0 {
		var zero T
		return zero, -1
	}

	minVal := arr[0]
	minIdx := 0

	for i := 1; i < len(arr); i++ {
		if arr[i] < minVal {
			minVal = arr[i]
			minIdx = i
		}
	}

	return minVal, minIdx
}

// FindMax finds the maximum element in a slice
// Returns the maximum element and its index
// Time complexity: O(n)
func FindMax[T constraints.Ordered](arr []T) (T, int) {
	if len(arr) == 0 {
		var zero T
		return zero, -1
	}

	maxVal := arr[0]
	maxIdx := 0

	for i := 1; i < len(arr); i++ {
		if arr[i] > maxVal {
			maxVal = arr[i]
			maxIdx = i
		}
	}

	return maxVal, maxIdx
}

// FindMinMax finds both minimum and maximum elements in a single pass
// Returns (min, minIdx, max, maxIdx)
// Time complexity: O(n)
func FindMinMax[T constraints.Ordered](arr []T) (T, int, T, int) {
	if len(arr) == 0 {
		var zero T
		return zero, -1, zero, -1
	}

	minVal, maxVal := arr[0], arr[0]
	minIdx, maxIdx := 0, 0

	for i := 1; i < len(arr); i++ {
		if arr[i] < minVal {
			minVal = arr[i]
			minIdx = i
		}
		if arr[i] > maxVal {
			maxVal = arr[i]
			maxIdx = i
		}
	}

	return minVal, minIdx, maxVal, maxIdx
}

// Contains checks if a slice contains a specific element
// Time complexity: O(n)
func Contains[T comparable](arr []T, target T) bool {
	return LinearSearch(arr, target) != -1
}

// ContainsAll checks if a slice contains all elements from another slice
// Time complexity: O(n * m) where n is len(arr) and m is len(targets)
func ContainsAll[T comparable](arr []T, targets []T) bool {
	for _, target := range targets {
		if !Contains(arr, target) {
			return false
		}
	}
	return true
}

// ContainsAny checks if a slice contains any element from another slice
// Time complexity: O(n * m) where n is len(arr) and m is len(targets)
func ContainsAny[T comparable](arr []T, targets []T) bool {
	for _, target := range targets {
		if Contains(arr, target) {
			return true
		}
	}
	return false
}
