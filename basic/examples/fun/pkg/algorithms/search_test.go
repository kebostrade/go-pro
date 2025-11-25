package algorithms

import (
	"testing"
)

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"found at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
		{"found at middle", []int{1, 2, 3, 4, 5}, 3, 2},
		{"found at end", []int{1, 2, 3, 4, 5}, 5, 4},
		{"not found below range", []int{1, 2, 3, 4, 5}, 0, -1},
		{"not found above range", []int{1, 2, 3, 4, 5}, 6, -1},
		{"not found in range", []int{1, 2, 4, 5}, 3, -1},
		{"empty array", []int{}, 5, -1},
		{"single element found", []int{5}, 5, 0},
		{"single element not found", []int{5}, 3, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("BinarySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinarySearchRecursive(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"found at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
		{"found at middle", []int{1, 2, 3, 4, 5}, 3, 2},
		{"found at end", []int{1, 2, 3, 4, 5}, 5, 4},
		{"not found", []int{1, 2, 3, 4, 5}, 6, -1},
		{"empty array", []int{}, 5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearchRecursive(tt.arr, tt.target, 0, len(tt.arr)-1)
			if got != tt.want {
				t.Errorf("BinarySearchRecursive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindFirstOccurrence(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"single occurrence", []int{1, 2, 3, 4, 5}, 3, 2},
		{"multiple occurrences", []int{1, 2, 2, 2, 3}, 2, 1},
		{"at beginning", []int{1, 1, 1, 2, 3}, 1, 0},
		{"not found", []int{1, 2, 3, 4, 5}, 6, -1},
		{"all same", []int{5, 5, 5, 5}, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindFirstOccurrence(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("FindFirstOccurrence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindLastOccurrence(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"single occurrence", []int{1, 2, 3, 4, 5}, 3, 2},
		{"multiple occurrences", []int{1, 2, 2, 2, 3}, 2, 3},
		{"at end", []int{1, 2, 3, 3, 3}, 3, 4},
		{"not found", []int{1, 2, 3, 4, 5}, 6, -1},
		{"all same", []int{5, 5, 5, 5}, 5, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindLastOccurrence(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("FindLastOccurrence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountOccurrences(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"single occurrence", []int{1, 2, 3, 4, 5}, 3, 1},
		{"multiple occurrences", []int{1, 2, 2, 2, 3}, 2, 3},
		{"not found", []int{1, 2, 3, 4, 5}, 6, 0},
		{"all same", []int{5, 5, 5, 5}, 5, 4},
		{"empty array", []int{}, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountOccurrences(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("CountOccurrences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinearSearch(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"found at beginning", []int{5, 2, 3, 1, 4}, 5, 0},
		{"found at middle", []int{5, 2, 3, 1, 4}, 3, 2},
		{"found at end", []int{5, 2, 3, 1, 4}, 4, 4},
		{"not found", []int{5, 2, 3, 1, 4}, 6, -1},
		{"empty array", []int{}, 5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LinearSearch(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("LinearSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinearSearchAll(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   []int
	}{
		{"single occurrence", []int{1, 2, 3, 4, 5}, 3, []int{2}},
		{"multiple occurrences", []int{1, 2, 2, 2, 3}, 2, []int{1, 2, 3}},
		{"not found", []int{1, 2, 3, 4, 5}, 6, []int{}},
		{"all same", []int{5, 5, 5, 5}, 5, []int{0, 1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LinearSearchAll(tt.arr, tt.target)
			if len(got) != len(tt.want) {
				t.Errorf("LinearSearchAll() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("LinearSearchAll()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSearchWithPredicate(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6}

	t.Run("find even number", func(t *testing.T) {
		idx, found := SearchWithPredicate(arr, func(n int) bool { return n%2 == 0 })
		if !found || arr[idx]%2 != 0 {
			t.Errorf("SearchWithPredicate() = %v, %v, want even number", idx, found)
		}
	})

	t.Run("not found", func(t *testing.T) {
		idx, found := SearchWithPredicate(arr, func(n int) bool { return n > 10 })
		if found {
			t.Errorf("SearchWithPredicate() = %v, %v, want not found", idx, found)
		}
	})
}

func TestFindMin(t *testing.T) {
	tests := []struct {
		name    string
		arr     []int
		wantVal int
		wantIdx int
	}{
		{"normal array", []int{5, 2, 8, 1, 9}, 1, 3},
		{"min at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
		{"min at end", []int{5, 4, 3, 2, 1}, 1, 4},
		{"single element", []int{42}, 42, 0},
		{"negative numbers", []int{-1, -5, -3}, -5, 1},
		{"empty array", []int{}, 0, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotIdx := FindMin(tt.arr)
			if gotVal != tt.wantVal || gotIdx != tt.wantIdx {
				t.Errorf("FindMin() = (%v, %v), want (%v, %v)", gotVal, gotIdx, tt.wantVal, tt.wantIdx)
			}
		})
	}
}

func TestFindMax(t *testing.T) {
	tests := []struct {
		name    string
		arr     []int
		wantVal int
		wantIdx int
	}{
		{"normal array", []int{5, 2, 8, 1, 9}, 9, 4},
		{"max at beginning", []int{9, 2, 3, 4, 5}, 9, 0},
		{"max at end", []int{1, 2, 3, 4, 9}, 9, 4},
		{"single element", []int{42}, 42, 0},
		{"empty array", []int{}, 0, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotIdx := FindMax(tt.arr)
			if gotVal != tt.wantVal || gotIdx != tt.wantIdx {
				t.Errorf("FindMax() = (%v, %v), want (%v, %v)", gotVal, gotIdx, tt.wantVal, tt.wantIdx)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   bool
	}{
		{"contains", []int{1, 2, 3, 4, 5}, 3, true},
		{"not contains", []int{1, 2, 3, 4, 5}, 6, false},
		{"empty array", []int{}, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmarks
func BenchmarkBinarySearch(b *testing.B) {
	arr := make([]int, 10000)
	for i := range arr {
		arr[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BinarySearch(arr, 5000)
	}
}

func BenchmarkLinearSearch(b *testing.B) {
	arr := make([]int, 10000)
	for i := range arr {
		arr[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearSearch(arr, 5000)
	}
}

func BenchmarkBinarySearchRecursive(b *testing.B) {
	arr := make([]int, 10000)
	for i := range arr {
		arr[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BinarySearchRecursive(arr, 5000, 0, len(arr)-1)
	}
}
