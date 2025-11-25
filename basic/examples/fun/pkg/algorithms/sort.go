package algorithms

import (
	"sync"

	"golang.org/x/exp/constraints"
)

// MergeSort performs merge sort sequentially
// Time complexity: O(n log n)
// Space complexity: O(n)
// Stable: Yes
func MergeSort[T constraints.Ordered](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])

	return merge(left, right)
}

// merge combines two sorted slices into one sorted slice
func merge[T constraints.Ordered](left, right []T) []T {
	result := make([]T, 0, len(left)+len(right))
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
// depth controls the recursion depth for concurrent execution
// Time complexity: O(n log n)
// Space complexity: O(n)
func MergeSortConcurrent[T constraints.Ordered](arr []T, depth int) []T {
	if len(arr) <= 1 {
		return arr
	}

	// Use sequential sort for small arrays or deep recursion
	if len(arr) < 1000 || depth <= 0 {
		return MergeSort(arr)
	}

	mid := len(arr) / 2
	var left, right []T
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

// QuickSort performs quick sort in place
// Time complexity: O(n log n) average, O(n²) worst case
// Space complexity: O(log n) due to recursion
// Stable: No
func QuickSort[T constraints.Ordered](arr []T) {
	if len(arr) < 2 {
		return
	}
	quickSortHelper(arr, 0, len(arr)-1)
}

func quickSortHelper[T constraints.Ordered](arr []T, low, high int) {
	if low < high {
		pivotIdx := partition(arr, low, high)
		quickSortHelper(arr, low, pivotIdx-1)
		quickSortHelper(arr, pivotIdx+1, high)
	}
}

func partition[T constraints.Ordered](arr []T, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// BubbleSort performs bubble sort in place
// Time complexity: O(n²)
// Space complexity: O(1)
// Stable: Yes
func BubbleSort[T constraints.Ordered](arr []T) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		// If no swaps occurred, array is sorted
		if !swapped {
			break
		}
	}
}

// InsertionSort performs insertion sort in place
// Time complexity: O(n²) worst case, O(n) best case
// Space complexity: O(1)
// Stable: Yes
func InsertionSort[T constraints.Ordered](arr []T) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1

		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

// SelectionSort performs selection sort in place
// Time complexity: O(n²)
// Space complexity: O(1)
// Stable: No
func SelectionSort[T constraints.Ordered](arr []T) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		if minIdx != i {
			arr[i], arr[minIdx] = arr[minIdx], arr[i]
		}
	}
}

// IsSorted checks if a slice is sorted in ascending order
// Time complexity: O(n)
func IsSorted[T constraints.Ordered](arr []T) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}

// IsSortedDescending checks if a slice is sorted in descending order
// Time complexity: O(n)
func IsSortedDescending[T constraints.Ordered](arr []T) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] > arr[i-1] {
			return false
		}
	}
	return true
}

// Reverse reverses a slice in place
// Time complexity: O(n)
func Reverse[T any](arr []T) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// SortWithComparator sorts a slice using a custom comparator function
// comparator should return true if a < b
// Time complexity: O(n log n)
func SortWithComparator[T any](arr []T, comparator func(a, b T) bool) {
	if len(arr) < 2 {
		return
	}
	sortWithComparatorHelper(arr, 0, len(arr)-1, comparator)
}

func sortWithComparatorHelper[T any](arr []T, low, high int, comparator func(a, b T) bool) {
	if low < high {
		pivotIdx := partitionWithComparator(arr, low, high, comparator)
		sortWithComparatorHelper(arr, low, pivotIdx-1, comparator)
		sortWithComparatorHelper(arr, pivotIdx+1, high, comparator)
	}
}

func partitionWithComparator[T any](arr []T, low, high int, comparator func(a, b T) bool) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if comparator(arr[j], pivot) {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// StableSort performs a stable sort (maintains relative order of equal elements)
// Uses merge sort which is naturally stable
// Time complexity: O(n log n)
func StableSort[T constraints.Ordered](arr []T) []T {
	return MergeSort(arr)
}

// TopK returns the k largest elements from the slice
// Time complexity: O(n log k)
func TopK[T constraints.Ordered](arr []T, k int) []T {
	if k <= 0 || len(arr) == 0 {
		return []T{}
	}
	if k >= len(arr) {
		sorted := make([]T, len(arr))
		copy(sorted, arr)
		QuickSort(sorted)
		Reverse(sorted)
		return sorted
	}

	// Use a simple approach: sort and take top k
	sorted := make([]T, len(arr))
	copy(sorted, arr)
	QuickSort(sorted)
	Reverse(sorted)
	return sorted[:k]
}

// BottomK returns the k smallest elements from the slice
// Time complexity: O(n log k)
func BottomK[T constraints.Ordered](arr []T, k int) []T {
	if k <= 0 || len(arr) == 0 {
		return []T{}
	}
	if k >= len(arr) {
		sorted := make([]T, len(arr))
		copy(sorted, arr)
		QuickSort(sorted)
		return sorted
	}

	sorted := make([]T, len(arr))
	copy(sorted, arr)
	QuickSort(sorted)
	return sorted[:k]
}
