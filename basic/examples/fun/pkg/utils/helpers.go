package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// RepeatString repeats a string n times
func RepeatString(s string, count int) string {
	return strings.Repeat(s, count)
}

// PrintSeparator prints a separator line with the given character
func PrintSeparator(char string, length int) {
	fmt.Println(strings.Repeat(char, length))
}

// PrintHeader prints a formatted header
func PrintHeader(title string) {
	length := 60
	fmt.Println()
	PrintSeparator("=", length)
	fmt.Println(title)
	PrintSeparator("=", length)
}

// PrintSubHeader prints a formatted sub-header
func PrintSubHeader(title string) {
	length := 60
	fmt.Println()
	fmt.Println(title)
	PrintSeparator("-", length)
}

// NewRandSource creates a new random source with current time seed
func NewRandSource() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GenerateRandomInts generates n random integers between min and max
func GenerateRandomInts(n, min, max int) []int {
	r := NewRandSource()
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = r.Intn(max-min+1) + min
	}
	return result
}

// GenerateRandomInt64s generates n random int64 values between min and max
func GenerateRandomInt64s(n int, min, max int64) []int64 {
	r := NewRandSource()
	result := make([]int64, n)
	for i := 0; i < n; i++ {
		result[i] = r.Int63n(max-min+1) + min
	}
	return result
}

// GenerateRandomFloats generates n random float64 values between min and max
func GenerateRandomFloats(n int, min, max float64) []float64 {
	r := NewRandSource()
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = min + r.Float64()*(max-min)
	}
	return result
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of an integer
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Contains checks if a slice contains a value
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Filter filters a slice based on a predicate function
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map applies a function to each element of a slice
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Reduce reduces a slice to a single value using a reducer function
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// Reverse reverses a slice in place
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Clone creates a deep copy of a slice
func Clone[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}

// Unique returns a slice with duplicate elements removed
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Chunk splits a slice into chunks of the specified size
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	chunks := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// MeasureTime measures the execution time of a function
func MeasureTime(name string, fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	fmt.Printf("%s took %s\n", name, FormatDuration(duration))
	return duration
}

// Retry retries a function up to maxAttempts times with a delay between attempts
func Retry(maxAttempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i < maxAttempts-1 {
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxAttempts, err)
}

// Debounce creates a debounced version of a function
func Debounce(delay time.Duration, fn func()) func() {
	var timer *time.Timer
	return func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(delay, fn)
	}
}

// Throttle creates a throttled version of a function
func Throttle(interval time.Duration, fn func()) func() {
	var lastRun time.Time
	return func() {
		if time.Since(lastRun) >= interval {
			fn()
			lastRun = time.Now()
		}
	}
}
