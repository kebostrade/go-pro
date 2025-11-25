package calculator

import "testing"

// Exercise: Write table-driven tests
// Learn how to test multiple cases efficiently

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		a         int
		b         int
		operation string
		expected  int
	}{
		{"Add positive numbers", 5, 3, "+", 8},
		{"Add negative numbers", -5, -3, "+", -8},
		{"Subtract", 10, 4, "-", 6},
		{"Multiply", 6, 7, "*", 42},
		{"Divide", 20, 4, "/", 5},
		{"Divide by zero", 10, 0, "/", 0},
		{"Unknown operation", 5, 3, "%", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Calculate(tt.a, tt.b, tt.operation)
			if result != tt.expected {
				t.Errorf("Calculate(%d, %d, %s) = %d; want %d",
					tt.a, tt.b, tt.operation, result, tt.expected)
			}
		})
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Even positive", 4, true},
		{"Odd positive", 5, false},
		{"Even negative", -4, true},
		{"Odd negative", -5, false},
		{"Zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEven(tt.input)
			if result != tt.expected {
				t.Errorf("IsEven(%d) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"First larger", 10, 5, 10},
		{"Second larger", 5, 10, 10},
		{"Equal", 7, 7, 7},
		{"Negative numbers", -5, -10, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Max(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
