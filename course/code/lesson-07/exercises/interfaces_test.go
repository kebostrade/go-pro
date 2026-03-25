package exercises

import (
	"testing"
)

func TestRectangleArea(t *testing.T) {
	r := Rectangle{Width: 5, Height: 3}
	expected := 15.0
	result := r.Area()
	if result != expected {
		t.Errorf("Rectangle{Area} = %v; want %v", result, expected)
	}
}

func TestRectanglePerimeter(t *testing.T) {
	r := Rectangle{Width: 5, Height: 3}
	expected := 16.0
	result := r.Perimeter()
	if result != expected {
		t.Errorf("Rectangle{Perimeter} = %v; want %v", result, expected)
	}
}

func TestCircleArea(t *testing.T) {
	c := Circle{Radius: 2}
	expected := 3.141592653589793 * 4 // π * r²
	result := c.Area()
	if result != expected {
		t.Errorf("Circle{Area} = %v; want %v", result, expected)
	}
}

func TestCirclePerimeter(t *testing.T) {
	c := Circle{Radius: 2}
	expected := 3.141592653589793 * 4 // 2πr
	result := c.Perimeter()
	if result != expected {
		t.Errorf("Circle{Perimeter} = %v; want %v", result, expected)
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{42, "int"},
		{"hello", "string"},
		{3.14, "float64"},
		{[]int{1, 2, 3}, "[]int"},
		{nil, "<nil>"},
	}

	for _, tt := range tests {
		result := GetType(tt.input)
		if result != tt.expected {
			t.Errorf("GetType(%v) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		input       interface{}
		expectedVal int
		expectedOk  bool
	}{
		{42, 42, true},
		{"hello", 0, false},
		{3.14, 0, false},
		{nil, 0, false},
		{0, 0, true},
	}

	for _, tt := range tests {
		val, ok := GetInt(tt.input)
		if val != tt.expectedVal || ok != tt.expectedOk {
			t.Errorf("GetInt(%v) = (%d, %v); want (%d, %v)", tt.input, val, ok, tt.expectedVal, tt.expectedOk)
		}
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{42, "int: 42"},
		{"hello", "string: hello"},
		{3.14, "float64: 3.14"},
		{true, "bool: true"},
		{nil, "nil: nil"},
	}

	for _, tt := range tests {
		result := Describe(tt.input)
		if result != tt.expected {
			t.Errorf("Describe(%v) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestPersonString(t *testing.T) {
	p := Person{Name: "Alice", Age: 30}
	expected := "Alice (30)"
	result := p.String()
	if result != expected {
		t.Errorf("Person.String() = %s; want %s", result, expected)
	}
}

func TestValidationError(t *testing.T) {
	e := ValidationError{Field: "email", Message: "invalid format"}
	expected := "email: invalid format"
	result := e.Error()
	if result != expected {
		t.Errorf("ValidationError.Error() = %s; want %s", result, expected)
	}
}

func TestPrintArea(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{Rectangle{Width: 4, Height: 5}, "20.00"},
		{Circle{Radius: 1}, "3.14"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		result := PrintArea(tt.input)
		if result != tt.expected {
			t.Errorf("PrintArea(%v) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}
