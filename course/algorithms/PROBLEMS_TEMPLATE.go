// Package main provides a template for algorithm problem solutions
// Each problem follows a consistent structure for learning and practice
package main

import (
	"fmt"
	"time"
)

// Problem represents a coding problem
type Problem struct {
	ID          int
	Title       string
	Difficulty  string
	Category    string
	Description string
	Examples    []Example
	Constraints []string
}

// Example represents a test case example
type Example struct {
	Input    string
	Output   string
	Explanation string
}

// Solution tracks solution attempts
type Solution struct {
	ProblemID  int
	Date       time.Time
	Time       time.Duration
	Attempts   int
	Optimized  bool
	Notes      string
	Patterns   []string
}

// ============================================
// PROBLEM: [Problem Name]
// ============================================

/*
PROBLEM DESCRIPTION:
[Describe the problem here]

EXAMPLES:
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: Because nums[0] + nums[1] == 9, we return [0, 1].

CONSTRAINTS:
- 2 <= nums.length <= 10^4
- -10^9 <= nums[i] <= 10^9
- Only one valid answer exists

PATTERN: [Hash Map / Two Pointers / Sliding Window / etc.]
TIME: O(n)
SPACE: O(n)
*/

// Solution1: [Approach Name]
func problemName_solution1(input []int) []int {
	// Approach: [Describe the approach]
	// Time Complexity: O(?)
	// Space Complexity: O(?)

	result := []int{}
	// Implementation here
	return result
}

// Solution2: [Alternative Approach - Optimized]
func problemName_solution2(input []int) []int {
	// Approach: [Describe optimized approach]
	// Time Complexity: O(?)
	// Space Complexity: O(?)

	result := []int{}
	// Optimized implementation here
	return result
}

// ============================================
// UTILITY FUNCTIONS
// ============================================

// ListNode for linked list problems
type ListNode struct {
	Val  int
	Next *ListNode
}

// TreeNode for tree problems
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// PrintSlice prints a slice for debugging
func PrintSlice[T any](s []T) {
	fmt.Printf("%v\n", s)
}

// Min returns minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abs returns absolute value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ============================================
// COMMON PATTERNS
// ============================================

// SlidingWindow template
func slidingWindow(arr []int, k int) int {
	// Initialize window
	left := 0
	result := 0

	// Expand window
	for right := 0; right < len(arr); right++ {
		// Add arr[right] to window

		// Shrink window if needed
		for /* invalid condition */ left <= right {
			// Remove arr[left] from window
			left++
		}

		// Update result
		result = Max(result, right-left+1)
	}

	return result
}

// TwoPointers template
func twoPointers(arr []int, target int) []int {
	left, right := 0, len(arr)-1

	for left < right {
		sum := arr[left] + arr[right]

		if sum == target {
			return []int{left, right}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}

	return []int{-1, -1}
}

// BinarySearch template
func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}

// FastSlowPointers template (cycle detection)
func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow, fast := head, head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			return true
		}
	}

	return false
}

// BFS template
func bfs(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	result := [][]int{}
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue)
		level := []int{}

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		result = append(result, level)
	}

	return result
}

// DFS template
func dfs(root *TreeNode) []int {
	result := []int{}

	var traverse func(node *TreeNode)
	traverse = func(node *TreeNode) {
		if node == nil {
			return
		}

		// Pre-order
		result = append(result, node.Val)
		traverse(node.Left)
		traverse(node.Right)
		// Post-order would append after recursive calls
	}

	traverse(root)
	return result
}

// Backtracking template
func backtrack(candidates []int, target int) [][]int {
	result := [][]int{}

	var backtrackHelper func(start int, current []int, remaining int)
	backtrackHelper = func(start int, current []int, remaining int) {
		if remaining == 0 {
			// Make a copy of current
			temp := make([]int, len(current))
			copy(temp, current)
			result = append(result, temp)
			return
		}

		for i := start; i < len(candidates); i++ {
			if candidates[i] > remaining {
				continue
			}

			current = append(current, candidates[i])
			backtrackHelper(i, current, remaining-candidates[i])
			current = current[:len(current)-1] // backtrack
		}
	}

	backtrackHelper(0, []int{}, target)
	return result
}

// ============================================
// TESTING UTILITIES
// ============================================

// TestCase represents a test case
type TestCase struct {
	Name     string
	Input    interface{}
	Expected interface{}
}

// AssertEqual checks if two values are equal
func AssertEqual[T comparable](got, expected T) bool {
	return got == expected
}

// RunTests runs a batch of test cases
func RunTests(tests []TestCase, solutionFn func(interface{}) interface{}) {
	passed := 0
	failed := 0

	for _, test := range tests {
		start := time.Now()
		result := solutionFn(test.Input)
		duration := time.Since(start)

		if AssertEqual(result, test.Expected) {
			fmt.Printf("✅ PASS: %s (%v)\n", test.Name, duration)
			passed++
		} else {
			fmt.Printf("❌ FAIL: %s - Expected: %v, Got: %v\n", test.Name, test.Expected, result)
			failed++
		}
	}

	fmt.Printf("\n📊 Results: %d passed, %d failed\n", passed, failed)
}

// Benchmark measures execution time
func Benchmark(name string, fn func(), iterations int) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	duration := time.Since(start)
	avg := duration / time.Duration(iterations)
	fmt.Printf("⏱️  %s: %v total, %v avg (%d iterations)\n", name, duration, avg, iterations)
}

// ============================================
// MAIN - Example Usage
// ============================================

func main() {
	fmt.Println("=== Algorithm Practice Template ===\n")

	// Example: Running test cases
	tests := []TestCase{
		{Name: "Example 1", Input: []int{2, 7, 11, 15}, Expected: []int{0, 1}},
		{Name: "Example 2", Input: []int{3, 2, 4}, Expected: []int{1, 2}},
		{Name: "Example 3", Input: []int{3, 3}, Expected: []int{0, 1}},
	}

	fmt.Println("Running tests...")
	// RunTests(tests, twoSum) // Uncomment when implementing

	// Example: Benchmarking
	fmt.Println("\nBenchmarking...")
	Benchmark("Binary Search", func() {
		binarySearch([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 7)
	}, 10000)
}
