package exercises

import (
	"strings"
	"testing"
)

func TestSimpleGreeting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic greeting", "Alice", "Hello, Alice! Welcome to Go programming."},
		{"another name", "Bob", "Hello, Bob! Welcome to Go programming."},
		{"empty name", "", "Hello, ! Welcome to Go programming."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SimpleGreeting(tt.input)
			if got != tt.expected {
				t.Errorf("SimpleGreeting() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCalculator(t *testing.T) {
	tests := []struct {
		name       string
		a, b       int
		operation  string
		wantResult int
		wantError  bool
	}{
		{"addition", 5, 3, "add", 8, false},
		{"subtraction", 10, 4, "subtract", 6, false},
		{"multiplication", 6, 7, "multiply", 42, false},
		{"division", 15, 3, "divide", 5, false},
		{"division by zero", 10, 0, "divide", 0, true},
		{"unsupported operation", 5, 3, "modulo", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotError := Calculator(tt.a, tt.b, tt.operation)

			if tt.wantError {
				if gotError == nil {
					t.Errorf("Calculator() expected error, got nil")
				}
			} else {
				if gotError != nil {
					t.Errorf("Calculator() unexpected error: %v", gotError)
				}
				if gotResult != tt.wantResult {
					t.Errorf("Calculator() = %v, want %v", gotResult, tt.wantResult)
				}
			}
		})
	}
}

func TestMultipleReturns(t *testing.T) {
	tests := []struct {
		name                                  string
		x, y                                  float64
		wantSum, wantDiff, wantProd, wantQuot float64
	}{
		{"positive numbers", 10.0, 5.0, 15.0, 5.0, 50.0, 2.0},
		{"negative numbers", -6.0, 3.0, -3.0, -9.0, -18.0, -2.0},
		{"with zero", 8.0, 0.0, 8.0, 8.0, 0.0, 0.0}, // Note: division by zero handling varies
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, gotDiff, gotProd, gotQuot := MultipleReturns(tt.x, tt.y)

			if gotSum != tt.wantSum {
				t.Errorf("MultipleReturns() sum = %v, want %v", gotSum, tt.wantSum)
			}
			if gotDiff != tt.wantDiff {
				t.Errorf("MultipleReturns() diff = %v, want %v", gotDiff, tt.wantDiff)
			}
			if gotProd != tt.wantProd {
				t.Errorf("MultipleReturns() prod = %v, want %v", gotProd, tt.wantProd)
			}
			if tt.y != 0 && gotQuot != tt.wantQuot {
				t.Errorf("MultipleReturns() quot = %v, want %v", gotQuot, tt.wantQuot)
			}
		})
	}
}

func TestNamedReturns(t *testing.T) {
	tests := []struct {
		name                    string
		length, width           float64
		wantArea, wantPerimeter float64
	}{
		{"rectangle", 5.0, 3.0, 15.0, 16.0},
		{"square", 4.0, 4.0, 16.0, 16.0},
		{"unit rectangle", 1.0, 1.0, 1.0, 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArea, gotPerimeter := NamedReturns(tt.length, tt.width)
			if gotArea != tt.wantArea {
				t.Errorf("NamedReturns() area = %v, want %v", gotArea, tt.wantArea)
			}
			if gotPerimeter != tt.wantPerimeter {
				t.Errorf("NamedReturns() perimeter = %v, want %v", gotPerimeter, tt.wantPerimeter)
			}
		})
	}
}

func TestVariadicSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"single number", []int{5}, 5},
		{"multiple numbers", []int{1, 2, 3, 4, 5}, 15},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"mixed numbers", []int{10, -5, 3}, 8},
		{"empty", []int{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VariadicSum(tt.numbers...)
			if got != tt.expected {
				t.Errorf("VariadicSum() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestVariadicAverage(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		wantAvg   float64
		wantCount int
	}{
		{"single value", []float64{5.0}, 5.0, 1},
		{"multiple values", []float64{2.0, 4.0, 6.0}, 4.0, 3},
		{"empty", []float64{}, 0.0, 0},
		{"negative values", []float64{-2.0, -4.0}, -3.0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAvg, gotCount := VariadicAverage(tt.values...)
			if gotAvg != tt.wantAvg {
				t.Errorf("VariadicAverage() avg = %v, want %v", gotAvg, tt.wantAvg)
			}
			if gotCount != tt.wantCount {
				t.Errorf("VariadicAverage() count = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestStringJoiner(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		strings   []string
		expected  string
	}{
		{"basic join", "-", []string{"a", "b", "c"}, "a-b-c"},
		{"space separator", " ", []string{"Hello", "World"}, "Hello World"},
		{"empty separator", "", []string{"Go", "Pro"}, "GoPro"},
		{"single string", "-", []string{"alone"}, "alone"},
		{"empty strings", ",", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringJoiner(tt.separator, tt.strings...)
			if got != tt.expected {
				t.Errorf("StringJoiner() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFunctionAsParameter(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int
		operation func(int, int) int
		expected  int
	}{
		{"addition", 5, 3, Add, 8},
		{"multiplication", 4, 6, Multiply, 24},
		{"custom operation", 10, 2, func(x, y int) int { return x - y }, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FunctionAsParameter(tt.a, tt.b, tt.operation)
			if got != tt.expected {
				t.Errorf("FunctionAsParameter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReturnFunction(t *testing.T) {
	tests := []struct {
		name     string
		addValue int
		input    int
		expected int
	}{
		{"add 5", 5, 10, 15},
		{"add negative", -3, 10, 7},
		{"add zero", 0, 42, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adder := ReturnFunction(tt.addValue)
			if adder == nil {
				t.Errorf("ReturnFunction() returned nil")
				return
			}
			got := adder(tt.input)
			if got != tt.expected {
				t.Errorf("ReturnFunction()(%d) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestClosure(t *testing.T) {
	counter := Closure()
	if counter == nil {
		t.Errorf("Closure() returned nil")
		return
	}

	// Test that counter increments
	first := counter()
	second := counter()
	third := counter()

	if first != 1 || second != 2 || third != 3 {
		t.Errorf("Closure() counter sequence = %d, %d, %d, want 1, 2, 3", first, second, third)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		dividend   float64
		divisor    float64
		wantResult float64
		wantError  bool
		errorMsg   string
	}{
		{"valid division", 10.0, 2.0, 5.0, false, ""},
		{"division by zero", 10.0, 0.0, 0.0, true, "division by zero"},
		{"negative numbers", -6.0, 2.0, -3.0, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotError := ErrorHandling(tt.dividend, tt.divisor)

			if tt.wantError {
				if gotError == nil {
					t.Errorf("ErrorHandling() expected error, got nil")
				} else if !strings.Contains(gotError.Error(), tt.errorMsg) {
					t.Errorf("ErrorHandling() error = %v, want to contain %v", gotError, tt.errorMsg)
				}
			} else {
				if gotError != nil {
					t.Errorf("ErrorHandling() unexpected error: %v", gotError)
				}
				if gotResult != tt.wantResult {
					t.Errorf("ErrorHandling() = %v, want %v", gotResult, tt.wantResult)
				}
			}
		})
	}
}

func TestRecursiveFactorial(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"factorial 0", 0, 1},
		{"factorial 1", 1, 1},
		{"factorial 5", 5, 120},
		{"factorial 6", 6, 720},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecursiveFactorial(tt.n)
			if got != tt.expected {
				t.Errorf("RecursiveFactorial() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRecursiveFibonacci(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"fibonacci 0", 0, 0},
		{"fibonacci 1", 1, 1},
		{"fibonacci 2", 2, 1},
		{"fibonacci 5", 5, 5},
		{"fibonacci 8", 8, 21},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecursiveFibonacci(tt.n)
			if got != tt.expected {
				t.Errorf("RecursiveFibonacci() = %v, want %v", got, tt.expected)
			}
		})
	}
}
