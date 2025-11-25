package test

import (
	"testing"
)

// Test basic variable operations
func TestVariables(t *testing.T) {
	// Integer
	var x int = 42
	if x != 42 {
		t.Errorf("Expected x to be 42, got %d", x)
	}

	// Type inference
	y := 100
	if y != 100 {
		t.Errorf("Expected y to be 100, got %d", y)
	}

	// Multiple assignment
	a, b := 1, 2
	if a != 1 || b != 2 {
		t.Errorf("Expected a=1, b=2, got a=%d, b=%d", a, b)
	}
}

// Test function operations
func TestFunctions(t *testing.T) {
	// Basic function
	result := add(5, 3)
	if result != 8 {
		t.Errorf("Expected 5+3=8, got %d", result)
	}

	// Multiple returns
	q, r := divMod(17, 5)
	if q != 3 || r != 2 {
		t.Errorf("Expected 17/5 = 3 remainder 2, got %d remainder %d", q, r)
	}

	// Variadic function
	sum := sumNumbers(1, 2, 3, 4, 5)
	if sum != 15 {
		t.Errorf("Expected sum=15, got %d", sum)
	}
}

func add(x, y int) int {
	return x + y
}

func divMod(a, b int) (int, int) {
	return a / b, a % b
}

func sumNumbers(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Test pointer operations
func TestPointers(t *testing.T) {
	x := 10
	ptr := &x

	// Check pointer points to x
	if *ptr != 10 {
		t.Errorf("Expected *ptr=10, got %d", *ptr)
	}

	// Modify through pointer
	*ptr = 20
	if x != 20 {
		t.Errorf("Expected x=20 after *ptr=20, got %d", x)
	}

	// Swap through pointers
	a, b := 5, 10
	swapPointers(&a, &b)
	if a != 10 || b != 5 {
		t.Errorf("Expected a=10, b=5 after swap, got a=%d, b=%d", a, b)
	}
}

func swapPointers(x, y *int) {
	*x, *y = *y, *x
}

// Test struct operations
func TestStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{Name: "Alice", Age: 30}

	if p.Name != "Alice" {
		t.Errorf("Expected name=Alice, got %s", p.Name)
	}

	if p.Age != 30 {
		t.Errorf("Expected age=30, got %d", p.Age)
	}

	// Modify struct
	p.Age = 31
	if p.Age != 31 {
		t.Errorf("Expected age=31 after modification, got %d", p.Age)
	}
}

// Test struct methods
func TestStructMethods(t *testing.T) {
	rect := Rectangle{Width: 10, Height: 5}

	area := rect.Area()
	if area != 50 {
		t.Errorf("Expected area=50, got %.2f", area)
	}

	perimeter := rect.Perimeter()
	if perimeter != 30 {
		t.Errorf("Expected perimeter=30, got %.2f", perimeter)
	}
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Test pointer receivers
func TestPointerReceivers(t *testing.T) {
	counter := Counter{Count: 0}

	counter.Increment()
	if counter.Count != 1 {
		t.Errorf("Expected count=1, got %d", counter.Count)
	}

	counter.Add(5)
	if counter.Count != 6 {
		t.Errorf("Expected count=6, got %d", counter.Count)
	}

	counter.Reset()
	if counter.Count != 0 {
		t.Errorf("Expected count=0 after reset, got %d", counter.Count)
	}
}

type Counter struct {
	Count int
}

func (c *Counter) Increment() {
	c.Count++
}

func (c *Counter) Add(n int) {
	c.Count += n
}

func (c *Counter) Reset() {
	c.Count = 0
}

// Test interfaces
func TestInterfaces(t *testing.T) {
	var s Shape

	s = Circle{Radius: 5}
	area := s.Area()
	expected := 3.14159 * 5 * 5
	if area < expected-0.01 || area > expected+0.01 {
		t.Errorf("Expected circle area ~%.2f, got %.2f", expected, area)
	}

	s = Rect{Width: 10, Height: 5}
	area = s.Area()
	if area != 50 {
		t.Errorf("Expected rectangle area=50, got %.2f", area)
	}
}

type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

type Rect struct {
	Width  float64
	Height float64
}

func (r Rect) Area() float64 {
	return r.Width * r.Height
}

// Test type assertions
func TestTypeAssertions(t *testing.T) {
	var i interface{} = "Hello"

	// Successful assertion
	s, ok := i.(string)
	if !ok {
		t.Error("Expected successful string assertion")
	}
	if s != "Hello" {
		t.Errorf("Expected 'Hello', got %s", s)
	}

	// Failed assertion
	_, ok = i.(int)
	if ok {
		t.Error("Expected failed int assertion")
	}
}

// Test loops
func TestLoops(t *testing.T) {
	// Basic for loop
	sum := 0
	for i := 1; i <= 5; i++ {
		sum += i
	}
	if sum != 15 {
		t.Errorf("Expected sum=15, got %d", sum)
	}

	// Range loop
	numbers := []int{1, 2, 3, 4, 5}
	sum = 0
	for _, n := range numbers {
		sum += n
	}
	if sum != 15 {
		t.Errorf("Expected sum=15, got %d", sum)
	}

	// While-style loop
	count := 0
	i := 1
	for i <= 5 {
		count++
		i++
	}
	if count != 5 {
		t.Errorf("Expected count=5, got %d", count)
	}
}

// Test constants and iota
func TestIota(t *testing.T) {
	const (
		First = iota
		Second
		Third
	)

	if First != 0 {
		t.Errorf("Expected First=0, got %d", First)
	}
	if Second != 1 {
		t.Errorf("Expected Second=1, got %d", Second)
	}
	if Third != 2 {
		t.Errorf("Expected Third=2, got %d", Third)
	}

	// Bit flags
	const (
		FlagRead    = 1 << iota // 1
		FlagWrite               // 2
		FlagExecute             // 4
	)

	if FlagRead != 1 {
		t.Errorf("Expected FlagRead=1, got %d", FlagRead)
	}
	if FlagWrite != 2 {
		t.Errorf("Expected FlagWrite=2, got %d", FlagWrite)
	}
	if FlagExecute != 4 {
		t.Errorf("Expected FlagExecute=4, got %d", FlagExecute)
	}
}
