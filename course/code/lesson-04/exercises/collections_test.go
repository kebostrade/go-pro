package exercises

import (
	"reflect"
	"sort"
	"testing"
)

// Test Array Practice
func TestGetFirstFivePrimes(t *testing.T) {
	expected := [5]int{2, 3, 5, 7, 11}
	result := GetFirstFivePrimes()

	if result != expected {
		t.Errorf("GetFirstFivePrimes() = %v, want %v", result, expected)
	}
}

func TestFindMaxInArray(t *testing.T) {
	arr := [10]int{3, 7, 2, 9, 1, 8, 4, 6, 5, 0}
	expectedMax := 9
	expectedIndex := 3

	max, index := FindMaxInArray(arr)

	if max != expectedMax || index != expectedIndex {
		t.Errorf("FindMaxInArray() = (%d, %d), want (%d, %d)",
			max, index, expectedMax, expectedIndex)
	}
}

// Test Slice Manipulation
func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 2, 3, 3, 3, 4}, []int{1, 2, 3, 4}},
		{[]int{1, 1, 1, 1}, []int{1}},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{[]int{}, []int{}},
	}

	for _, test := range tests {
		result := RemoveDuplicates(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v",
				test.input, result, test.expected)
		}
	}
}

func TestReverseSlice(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "b", "c", "d"}, []string{"d", "c", "b", "a"}},
		{[]string{"hello"}, []string{"hello"}},
		{[]string{}, []string{}},
		{[]string{"x", "y"}, []string{"y", "x"}},
	}

	for _, test := range tests {
		// Make a copy since ReverseSlice modifies in place
		input := make([]string, len(test.input))
		copy(input, test.input)

		ReverseSlice(input)
		if !reflect.DeepEqual(input, test.expected) {
			t.Errorf("ReverseSlice(%v) resulted in %v, want %v",
				test.input, input, test.expected)
		}
	}
}

func TestMergeSortedSlices(t *testing.T) {
	tests := []struct {
		slice1   []int
		slice2   []int
		expected []int
	}{
		{[]int{1, 3, 5}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{1, 2, 3}, []int{4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{}, []int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{1, 2, 3}, []int{}, []int{1, 2, 3}},
	}

	for _, test := range tests {
		result := MergeSortedSlices(test.slice1, test.slice2)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("MergeSortedSlices(%v, %v) = %v, want %v",
				test.slice1, test.slice2, result, test.expected)
		}
	}
}

// Test Map Operations
func TestCountCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected map[rune]int
	}{
		{"hello", map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}},
		{"aaa", map[rune]int{'a': 3}},
		{"", map[rune]int{}},
		{"abc", map[rune]int{'a': 1, 'b': 1, 'c': 1}},
	}

	for _, test := range tests {
		result := CountCharacters(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("CountCharacters(%q) = %v, want %v",
				test.input, result, test.expected)
		}
	}
}

func TestInvertMap(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3}
	expected := map[int]string{1: "a", 2: "b", 3: "c"}

	result := InvertMap(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InvertMap(%v) = %v, want %v", input, result, expected)
	}
}

func TestMergeMaps(t *testing.T) {
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"b": 3, "c": 4}
	expected := map[string]int{"a": 1, "b": 3, "c": 4}

	result := MergeMaps(map1, map2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MergeMaps(%v, %v) = %v, want %v", map1, map2, result, expected)
	}
}

// Test Advanced Operations
func TestFindIntersection(t *testing.T) {
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}
	expected := []int{3, 4}

	result := FindIntersection(slice1, slice2)
	sort.Ints(result) // Sort for consistent comparison

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FindIntersection(%v, %v) = %v, want %v",
			slice1, slice2, result, expected)
	}
}

func TestGroupByLength(t *testing.T) {
	words := []string{"cat", "dog", "elephant", "ant", "horse"}
	expected := map[int][]string{
		3: {"cat", "dog", "ant"},
		5: {"horse"},
		8: {"elephant"},
	}

	result := GroupByLength(words)

	// Sort slices for consistent comparison
	for key := range result {
		sort.Strings(result[key])
	}
	for key := range expected {
		sort.Strings(expected[key])
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GroupByLength(%v) = %v, want %v", words, result, expected)
	}
}

// Test Inventory Management
func TestInventoryOperations(t *testing.T) {
	inv := NewInventory()
	if inv == nil {
		t.Fatal("NewInventory() returned nil")
	}

	// Test adding products
	product1 := &Product{ID: "P001", Name: "Laptop", Price: 999.99, Stock: 10}
	product2 := &Product{ID: "P002", Name: "Mouse", Price: 29.99, Stock: 5}

	inv.AddProduct(product1)
	inv.AddProduct(product2)

	// Test getting products
	retrieved, exists := inv.GetProduct("P001")
	if !exists || retrieved.Name != "Laptop" {
		t.Errorf("GetProduct failed to retrieve correct product")
	}

	// Test updating stock
	success := inv.UpdateStock("P001", 15)
	if !success {
		t.Errorf("UpdateStock failed")
	}

	updated, _ := inv.GetProduct("P001")
	if updated.Stock != 15 {
		t.Errorf("Stock not updated correctly, got %d, want 15", updated.Stock)
	}

	// Test low stock products
	lowStock := inv.GetLowStockProducts(10)
	if len(lowStock) != 1 || lowStock[0].ID != "P002" {
		t.Errorf("GetLowStockProducts failed")
	}

	// Test total value
	expectedValue := 15*999.99 + 5*29.99
	totalValue := inv.GetTotalValue()
	delta := 0.01 // Allow small floating-point precision difference
	if (totalValue-expectedValue) > delta || (expectedValue-totalValue) > delta {
		t.Errorf("GetTotalValue() = %.2f, want %.2f", totalValue, expectedValue)
	}
}

// Test Memory Efficiency
func TestEfficientAppend(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	result := EfficientAppend(10, values)

	if len(result) != 5 {
		t.Errorf("EfficientAppend length = %d, want 5", len(result))
	}

	if cap(result) < 10 {
		t.Errorf("EfficientAppend capacity = %d, want at least 10", cap(result))
	}

	if !reflect.DeepEqual(result, values) {
		t.Errorf("EfficientAppend values = %v, want %v", result, values)
	}
}

func TestProcessInChunks(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunkSize := 3

	// Processor that sums each chunk
	processor := func(chunk []int) int {
		sum := 0
		for _, v := range chunk {
			sum += v
		}
		return sum
	}

	result := ProcessInChunks(data, chunkSize, processor)
	expected := []int{6, 15, 24, 10} // [1+2+3, 4+5+6, 7+8+9, 10]

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ProcessInChunks() = %v, want %v", result, expected)
	}
}
