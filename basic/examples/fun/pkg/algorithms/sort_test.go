package algorithms

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"random order", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"single element", []int{42}, []int{42}},
		{"empty array", []int{}, []int{}},
		{"duplicates", []int{3, 3, 1, 1, 2, 2}, []int{1, 1, 2, 2, 3, 3}},
		{"negative numbers", []int{-3, -1, -5, -2}, []int{-5, -3, -2, -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeSort(tt.arr)
			if !slicesEqual(got, tt.want) {
				t.Errorf("MergeSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeSortConcurrent(t *testing.T) {
	tests := []struct {
		name  string
		arr   []int
		depth int
		want  []int
	}{
		{"large array", []int{5, 2, 8, 1, 9, 3, 7, 4, 6}, 3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"small array", []int{3, 1, 2}, 1, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeSortConcurrent(tt.arr, tt.depth)
			if !slicesEqual(got, tt.want) {
				t.Errorf("MergeSortConcurrent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuickSort(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"random order", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"single element", []int{42}, []int{42}},
		{"empty array", []int{}, []int{}},
		{"duplicates", []int{3, 3, 1, 1, 2, 2}, []int{1, 1, 2, 2, 3, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.arr))
			copy(arr, tt.arr)
			QuickSort(arr)
			if !slicesEqual(arr, tt.want) {
				t.Errorf("QuickSort() = %v, want %v", arr, tt.want)
			}
		})
	}
}

func TestBubbleSort(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"random order", []int{3, 1, 4, 2, 5}, []int{1, 2, 3, 4, 5}},
		{"single element", []int{42}, []int{42}},
		{"empty array", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.arr))
			copy(arr, tt.arr)
			BubbleSort(arr)
			if !slicesEqual(arr, tt.want) {
				t.Errorf("BubbleSort() = %v, want %v", arr, tt.want)
			}
		})
	}
}

func TestInsertionSort(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"random order", []int{3, 1, 4, 2, 5}, []int{1, 2, 3, 4, 5}},
		{"nearly sorted", []int{1, 2, 4, 3, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.arr))
			copy(arr, tt.arr)
			InsertionSort(arr)
			if !slicesEqual(arr, tt.want) {
				t.Errorf("InsertionSort() = %v, want %v", arr, tt.want)
			}
		})
	}
}

func TestSelectionSort(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"random order", []int{3, 1, 4, 2, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.arr))
			copy(arr, tt.arr)
			SelectionSort(arr)
			if !slicesEqual(arr, tt.want) {
				t.Errorf("SelectionSort() = %v, want %v", arr, tt.want)
			}
		})
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want bool
	}{
		{"sorted ascending", []int{1, 2, 3, 4, 5}, true},
		{"not sorted", []int{5, 4, 3, 2, 1}, false},
		{"single element", []int{42}, true},
		{"empty array", []int{}, true},
		{"duplicates sorted", []int{1, 1, 2, 2, 3}, true},
		{"duplicates not sorted", []int{1, 3, 2, 2, 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSorted(tt.arr)
			if got != tt.want {
				t.Errorf("IsSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSortedDescending(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want bool
	}{
		{"sorted descending", []int{5, 4, 3, 2, 1}, true},
		{"not sorted", []int{1, 2, 3, 4, 5}, false},
		{"single element", []int{42}, true},
		{"empty array", []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSortedDescending(tt.arr)
			if got != tt.want {
				t.Errorf("IsSortedDescending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"normal array", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"even length", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"odd length", []int{1, 2, 3}, []int{3, 2, 1}},
		{"single element", []int{42}, []int{42}},
		{"empty array", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.arr))
			copy(arr, tt.arr)
			Reverse(arr)
			if !slicesEqual(arr, tt.want) {
				t.Errorf("Reverse() = %v, want %v", arr, tt.want)
			}
		})
	}
}

func TestTopK(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		k    int
		want []int
	}{
		{"top 3", []int{5, 2, 8, 1, 9, 3}, 3, []int{9, 8, 5}},
		{"top 1", []int{5, 2, 8, 1, 9, 3}, 1, []int{9}},
		{"k larger than array", []int{5, 2, 3}, 5, []int{5, 3, 2}},
		{"k is 0", []int{5, 2, 3}, 0, []int{}},
		{"k negative", []int{5, 2, 3}, -1, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopK(tt.arr, tt.k)
			if !slicesEqual(got, tt.want) {
				t.Errorf("TopK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBottomK(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		k    int
		want []int
	}{
		{"bottom 3", []int{5, 2, 8, 1, 9, 3}, 3, []int{1, 2, 3}},
		{"bottom 1", []int{5, 2, 8, 1, 9, 3}, 1, []int{1}},
		{"k larger than array", []int{5, 2, 3}, 5, []int{2, 3, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BottomK(tt.arr, tt.k)
			if !slicesEqual(got, tt.want) {
				t.Errorf("BottomK() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func slicesEqual(a, b []int) bool {
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
func BenchmarkMergeSort(b *testing.B) {
	arr := make([]int, 1000)
	for i := range arr {
		arr[i] = 1000 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testArr := make([]int, len(arr))
		copy(testArr, arr)
		b.StartTimer()
		MergeSort(testArr)
	}
}

func BenchmarkQuickSort(b *testing.B) {
	arr := make([]int, 1000)
	for i := range arr {
		arr[i] = 1000 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testArr := make([]int, len(arr))
		copy(testArr, arr)
		b.StartTimer()
		QuickSort(testArr)
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	arr := make([]int, 100)
	for i := range arr {
		arr[i] = 100 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testArr := make([]int, len(arr))
		copy(testArr, arr)
		b.StartTimer()
		BubbleSort(testArr)
	}
}

func BenchmarkInsertionSort(b *testing.B) {
	arr := make([]int, 100)
	for i := range arr {
		arr[i] = 100 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testArr := make([]int, len(arr))
		copy(testArr, arr)
		b.StartTimer()
		InsertionSort(testArr)
	}
}

func BenchmarkMergeSortConcurrent(b *testing.B) {
	arr := make([]int, 10000)
	for i := range arr {
		arr[i] = 10000 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testArr := make([]int, len(arr))
		copy(testArr, arr)
		b.StartTimer()
		MergeSortConcurrent(testArr, 4)
	}
}
