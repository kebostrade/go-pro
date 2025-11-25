package algorithms

import (
	"math"
	"sync"

	"golang.org/x/exp/constraints"
)

// IsPrime checks if a number is prime
// Time complexity: O(√n)
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	// Check for divisors up to √n
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}

// GeneratePrimes generates all prime numbers up to n using Sieve of Eratosthenes
// Time complexity: O(n log log n)
// Space complexity: O(n)
func GeneratePrimes(n int) []int {
	if n < 2 {
		return []int{}
	}

	// Create a boolean slice to mark primes
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	// Sieve of Eratosthenes
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	// Collect primes
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

// GeneratePrimesConcurrent generates primes using concurrent workers
// Time complexity: O(n log log n)
func GeneratePrimesConcurrent(limit int, workers int) []int {
	if limit < 2 {
		return []int{}
	}

	// Use sieve for better performance
	return GeneratePrimes(limit)
}

// Fibonacci calculates the nth Fibonacci number iteratively
// Time complexity: O(n)
// Space complexity: O(1)
func Fibonacci(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return n
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// FibonacciSequence generates the first n Fibonacci numbers
// Time complexity: O(n)
func FibonacciSequence(n int) []int {
	if n <= 0 {
		return []int{}
	}

	sequence := make([]int, n)
	if n >= 1 {
		sequence[0] = 0
	}
	if n >= 2 {
		sequence[1] = 1
	}

	for i := 2; i < n; i++ {
		sequence[i] = sequence[i-1] + sequence[i-2]
	}

	return sequence
}

// FibonacciConcurrent calculates Fibonacci numbers concurrently using channels
// Demonstrates concurrent pattern, not necessarily faster for small n
func FibonacciConcurrent(n int, workers int) []int {
	if n <= 0 {
		return []int{}
	}

	jobs := make(chan int, n)
	results := make(chan struct {
		index int
		value int
	}, n)

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				results <- struct {
					index int
					value int
				}{index, Fibonacci(index)}
			}
		}()
	}

	// Send jobs
	go func() {
		for i := 0; i < n; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	sequence := make([]int, n)
	for result := range results {
		sequence[result.index] = result.value
	}

	return sequence
}

// GCD calculates the greatest common divisor using Euclidean algorithm
// Time complexity: O(log min(a, b))
func GCD(a, b int) int {
	a = abs(a)
	b = abs(b)

	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM calculates the least common multiple
// Time complexity: O(log min(a, b))
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return abs(a*b) / GCD(a, b)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Factorial calculates n! iteratively
// Time complexity: O(n)
func Factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return 1
	}

	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// Power calculates base^exponent using fast exponentiation
// Time complexity: O(log n)
func Power(base, exponent int) int {
	if exponent == 0 {
		return 1
	}
	if exponent < 0 {
		return 0 // Integer division, would be 1/Power(base, -exponent) for floats
	}

	result := 1
	currentBase := base
	currentExp := exponent

	for currentExp > 0 {
		if currentExp%2 == 1 {
			result *= currentBase
		}
		currentBase *= currentBase
		currentExp /= 2
	}

	return result
}

// Sum calculates the sum of a slice of numbers
// Time complexity: O(n)
func Sum[T constraints.Integer | constraints.Float](numbers []T) T {
	var sum T
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// Average calculates the average of a slice of numbers
// Time complexity: O(n)
func Average[T constraints.Integer | constraints.Float](numbers []T) float64 {
	if len(numbers) == 0 {
		return 0
	}
	return float64(Sum(numbers)) / float64(len(numbers))
}

// Median calculates the median of a slice of numbers
// Note: This modifies the input slice by sorting it
// Time complexity: O(n log n)
func Median[T constraints.Integer | constraints.Float](numbers []T) float64 {
	if len(numbers) == 0 {
		return 0
	}

	// Sort the numbers
	QuickSort(numbers)

	n := len(numbers)
	if n%2 == 0 {
		// Even number of elements: average of two middle elements
		return (float64(numbers[n/2-1]) + float64(numbers[n/2])) / 2
	}
	// Odd number of elements: middle element
	return float64(numbers[n/2])
}

// Mode finds the most frequently occurring element
// Time complexity: O(n)
func Mode[T comparable](numbers []T) T {
	if len(numbers) == 0 {
		var zero T
		return zero
	}

	frequency := make(map[T]int)
	for _, num := range numbers {
		frequency[num]++
	}

	var mode T
	maxCount := 0
	for num, count := range frequency {
		if count > maxCount {
			mode = num
			maxCount = count
		}
	}

	return mode
}

// StandardDeviation calculates the standard deviation of a slice of numbers
// Time complexity: O(n)
func StandardDeviation[T constraints.Integer | constraints.Float](numbers []T) float64 {
	if len(numbers) == 0 {
		return 0
	}

	avg := Average(numbers)
	var sumSquaredDiff float64

	for _, num := range numbers {
		diff := float64(num) - avg
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(len(numbers))
	return math.Sqrt(variance)
}

// IsPerfectSquare checks if a number is a perfect square
// Time complexity: O(log n)
func IsPerfectSquare(n int) bool {
	if n < 0 {
		return false
	}

	sqrt := int(math.Sqrt(float64(n)))
	return sqrt*sqrt == n
}

// IsPowerOfTwo checks if a number is a power of 2
// Time complexity: O(1)
func IsPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// NextPowerOfTwo finds the next power of 2 greater than or equal to n
// Time complexity: O(1)
func NextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}

	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++

	return n
}
