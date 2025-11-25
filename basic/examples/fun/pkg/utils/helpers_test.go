package utils

import (
	"testing"
	"time"
)

func TestMin(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{-1, -5, -5},
	}

	for _, tt := range tests {
		got := Min(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Min(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{1, 2, 2},
		{5, 3, 5},
		{-1, -5, -1},
	}

	for _, tt := range tests {
		got := Max(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Max(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		got := Abs(tt.n)
		if got != tt.want {
			t.Errorf("Abs(%d) = %d; want %d", tt.n, got, tt.want)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		value int
		want  bool
	}{
		{"contains", []int{1, 2, 3}, 2, true},
		{"not contains", []int{1, 2, 3}, 4, false},
		{"empty", []int{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.slice, tt.value)
			if got != tt.want {
				t.Errorf("Contains() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6}
	result := Filter(slice, func(n int) bool {
		return n%2 == 0
	})

	want := []int{2, 4, 6}
	if !sliceEqual(result, want) {
		t.Errorf("Filter(even) = %v; want %v", result, want)
	}
}

func TestMap(t *testing.T) {
	slice := []int{1, 2, 3}
	result := Map(slice, func(n int) int {
		return n * 2
	})

	want := []int{2, 4, 6}
	if !sliceEqual(result, want) {
		t.Errorf("Map(*2) = %v; want %v", result, want)
	}
}

func TestReduce(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Reduce(slice, 0, func(acc, n int) int {
		return acc + n
	})

	want := 15
	if result != want {
		t.Errorf("Reduce(sum) = %d; want %d", result, want)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"normal", []int{1, 2, 3}, []int{3, 2, 1}},
		{"single", []int{1}, []int{1}},
		{"empty", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.in))
			copy(slice, tt.in)
			Reverse(slice)
			if !sliceEqual(slice, tt.want) {
				t.Errorf("Reverse() = %v; want %v", slice, tt.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	original := []int{1, 2, 3}
	cloned := Clone(original)

	cloned[0] = 999

	if original[0] == 999 {
		t.Error("Clone() should create independent copy")
	}
}

func TestUnique(t *testing.T) {
	slice := []int{1, 2, 2, 3, 3, 3, 4}
	result := Unique(slice)

	// Result should contain 1, 2, 3, 4 (order may vary)
	if len(result) != 4 {
		t.Errorf("Unique() length = %d; want 4", len(result))
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name string
		slice []int
		size  int
		want  int // number of chunks
	}{
		{"even division", []int{1, 2, 3, 4, 5, 6}, 2, 3},
		{"uneven division", []int{1, 2, 3, 4, 5}, 2, 3},
		{"size larger", []int{1, 2}, 5, 1},
		{"size zero", []int{1, 2}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Chunk(tt.slice, tt.size)
			if result == nil && tt.want == 0 {
				return
			}
			if len(result) != tt.want {
				t.Errorf("Chunk() chunks = %d; want %d", len(result), tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{500 * time.Nanosecond, "500ns"},
		{1500 * time.Microsecond, "1.50ms"},
		{2500 * time.Millisecond, "2.50s"},
	}

	for _, tt := range tests {
		got := FormatDuration(tt.duration)
		if got != tt.want {
			t.Errorf("FormatDuration(%v) = %s; want %s", tt.duration, got, tt.want)
		}
	}
}

func TestRetry(t *testing.T) {
	t.Run("success on first attempt", func(t *testing.T) {
		attempts := 0
		err := Retry(3, 10*time.Millisecond, func() error {
			attempts++
			return nil
		})

		if err != nil || attempts != 1 {
			t.Errorf("Retry() err = %v, attempts = %d; want nil, 1", err, attempts)
		}
	})

	t.Run("success on retry", func(t *testing.T) {
		attempts := 0
		err := Retry(3, 10*time.Millisecond, func() error {
			attempts++
			if attempts < 2 {
				return &testError{}
			}
			return nil
		})

		if err != nil || attempts != 2 {
			t.Errorf("Retry() err = %v, attempts = %d; want nil, 2", err, attempts)
		}
	})

	t.Run("fail all attempts", func(t *testing.T) {
		attempts := 0
		err := Retry(3, 10*time.Millisecond, func() error {
			attempts++
			return &testError{}
		})

		if err == nil || attempts != 3 {
			t.Errorf("Retry() err = %v, attempts = %d; want error, 3", err, attempts)
		}
	})
}

func TestGenerateRandomInts(t *testing.T) {
	result := GenerateRandomInts(10, 1, 100)

	if len(result) != 10 {
		t.Errorf("GenerateRandomInts() length = %d; want 10", len(result))
	}

	for _, v := range result {
		if v < 1 || v > 100 {
			t.Errorf("GenerateRandomInts() value %d out of range [1, 100]", v)
		}
	}
}

// Helper functions
func sliceEqual(a, b []int) bool {
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

type testError struct{}

func (e *testError) Error() string {
	return "test error"
}

// Benchmarks
func BenchmarkFilter(b *testing.B) {
	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Filter(slice, func(n int) bool {
			return n%2 == 0
		})
	}
}

func BenchmarkMap(b *testing.B) {
	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Map(slice, func(n int) int {
			return n * 2
		})
	}
}
