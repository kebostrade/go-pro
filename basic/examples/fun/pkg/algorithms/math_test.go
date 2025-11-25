package algorithms

import (
	"math"
	"testing"
)

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"negative number", -5, false},
		{"zero", 0, false},
		{"one", 1, false},
		{"two", 2, true},
		{"three", 3, true},
		{"four", 4, false},
		{"five", 5, true},
		{"seventeen", 17, true},
		{"hundred", 100, false},
		{"large prime", 97, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPrime(tt.n)
			if got != tt.want {
				t.Errorf("IsPrime(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestGeneratePrimes(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{"less than 2", 1, []int{}},
		{"up to 10", 10, []int{2, 3, 5, 7}},
		{"up to 20", 20, []int{2, 3, 5, 7, 11, 13, 17, 19}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GeneratePrimes(tt.n)
			if !intSliceEqual(got, tt.want) {
				t.Errorf("GeneratePrimes(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"negative", -1, 0},
		{"zero", 0, 0},
		{"one", 1, 1},
		{"two", 2, 1},
		{"three", 3, 2},
		{"five", 5, 5},
		{"ten", 10, 55},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fibonacci(tt.n)
			if got != tt.want {
				t.Errorf("Fibonacci(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestFibonacciSequence(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{"zero", 0, []int{}},
		{"one", 1, []int{0}},
		{"five", 5, []int{0, 1, 1, 2, 3}},
		{"ten", 10, []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FibonacciSequence(tt.n)
			if !intSliceEqual(got, tt.want) {
				t.Errorf("FibonacciSequence(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 12, 8, 4},
		{"one negative", -12, 8, 4},
		{"both negative", -12, -8, 4},
		{"coprime", 17, 19, 1},
		{"same number", 7, 7, 7},
		{"with zero", 12, 0, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GCD(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 12, 8, 24},
		{"coprime", 3, 7, 21},
		{"with zero", 12, 0, 0},
		{"same number", 5, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LCM(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("LCM(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"negative", -1, 0},
		{"zero", 0, 1},
		{"one", 1, 1},
		{"five", 5, 120},
		{"ten", 10, 3628800},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Factorial(tt.n)
			if got != tt.want {
				t.Errorf("Factorial(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		name     string
		base     int
		exponent int
		want     int
	}{
		{"zero exponent", 5, 0, 1},
		{"negative exponent", 2, -1, 0},
		{"positive", 2, 10, 1024},
		{"base zero", 0, 5, 0},
		{"base one", 1, 100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Power(tt.base, tt.exponent)
			if got != tt.want {
				t.Errorf("Power(%d, %d) = %d, want %d", tt.base, tt.exponent, got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"positive numbers", []int{1, 2, 3, 4, 5}, 15},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"mixed", []int{1, -2, 3, -4}, -2},
		{"empty", []int{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.nums)
			if got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want float64
	}{
		{"positive numbers", []int{1, 2, 3, 4, 5}, 3.0},
		{"single number", []int{10}, 10.0},
		{"empty", []int{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Average(tt.nums)
			if got != tt.want {
				t.Errorf("Average() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want float64
	}{
		{"odd count", []int{3, 1, 2}, 2.0},
		{"even count", []int{1, 2, 3, 4}, 2.5},
		{"single", []int{5}, 5.0},
		{"empty", []int{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy since Median sorts in place
			numsCopy := make([]int, len(tt.nums))
			copy(numsCopy, tt.nums)
			got := Median(numsCopy)
			if got != tt.want {
				t.Errorf("Median() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestMode(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"single mode", []int{1, 2, 2, 3}, 2},
		{"all same", []int{5, 5, 5}, 5},
		{"empty", []int{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mode(tt.nums)
			if got != tt.want {
				t.Errorf("Mode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestStandardDeviation(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want float64
	}{
		{"normal distribution", []int{2, 4, 4, 4, 5, 5, 7, 9}, 2.0},
		{"all same", []int{5, 5, 5, 5}, 0.0},
		{"empty", []int{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StandardDeviation(tt.nums)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("StandardDeviation() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestIsPerfectSquare(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"negative", -4, false},
		{"zero", 0, true},
		{"one", 1, true},
		{"perfect square", 16, true},
		{"not perfect", 15, false},
		{"large perfect", 144, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPerfectSquare(tt.n)
			if got != tt.want {
				t.Errorf("IsPerfectSquare(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"negative", -2, false},
		{"zero", 0, false},
		{"one", 1, true},
		{"power of two", 16, true},
		{"not power", 15, false},
		{"large power", 1024, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPowerOfTwo(tt.n)
			if got != tt.want {
				t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"negative", -5, 1},
		{"zero", 0, 1},
		{"already power", 16, 16},
		{"not power", 15, 16},
		{"large number", 100, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextPowerOfTwo(tt.n)
			if got != tt.want {
				t.Errorf("NextPowerOfTwo(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

// Helper function
func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Benchmarks
func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(97)
	}
}

func BenchmarkGeneratePrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GeneratePrimes(1000)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(30)
	}
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorial(20)
	}
}

func BenchmarkPower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Power(2, 10)
	}
}
