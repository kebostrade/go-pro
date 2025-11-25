package main

import (
	"fmt"
	"math"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Interfaces Demo")

	demo1BasicInterfaces()
	demo2InterfaceValues()
	demo3EmptyInterface()
	demo4TypeAssertions()
	demo5TypeSwitches()
}

func demo1BasicInterfaces() {
	utils.PrintSubHeader("1. Basic Interfaces")

	// Different types implementing the same interface
	var s Shape

	s = Circle{Radius: 5}
	fmt.Printf("Circle: Area = %.2f, Perimeter = %.2f\n", s.Area(), s.Perimeter())

	s = Rectangle2{Width: 10, Height: 5}
	fmt.Printf("Rectangle: Area = %.2f, Perimeter = %.2f\n", s.Area(), s.Perimeter())

	// Using interface as parameter
	printShapeInfo(Circle{Radius: 3})
	printShapeInfo(Rectangle2{Width: 6, Height: 4})
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Rectangle2 struct {
	Width  float64
	Height float64
}

func (r Rectangle2) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle2) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func printShapeInfo(s Shape) {
	fmt.Printf("Shape - Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func demo2InterfaceValues() {
	utils.PrintSubHeader("2. Interface Values")

	var w Writer

	fmt.Println("Interface value before assignment:")
	fmt.Printf("  Value: %v, Type: %T\n", w, w)

	w = ConsoleWriter{}
	fmt.Println("\nAfter assigning ConsoleWriter:")
	fmt.Printf("  Value: %v, Type: %T\n", w, w)
	w.Write("Hello from ConsoleWriter")

	w = FileWriter{Filename: "output.txt"}
	fmt.Println("\nAfter assigning FileWriter:")
	fmt.Printf("  Value: %v, Type: %T\n", w, w)
	w.Write("Hello from FileWriter")
}

type Writer interface {
	Write(data string)
}

type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(data string) {
	fmt.Printf("  [Console] %s\n", data)
}

type FileWriter struct {
	Filename string
}

func (fw FileWriter) Write(data string) {
	fmt.Printf("  [File: %s] %s\n", fw.Filename, data)
}

func demo3EmptyInterface() {
	utils.PrintSubHeader("3. Empty Interface (interface{})")

	// Empty interface can hold any type
	var anything interface{}

	anything = 42
	fmt.Printf("Integer: %v (type: %T)\n", anything, anything)

	anything = "Hello"
	fmt.Printf("String: %v (type: %T)\n", anything, anything)

	anything = true
	fmt.Printf("Boolean: %v (type: %T)\n", anything, anything)

	anything = []int{1, 2, 3}
	fmt.Printf("Slice: %v (type: %T)\n", anything, anything)

	// Function accepting any type
	printAnything(42)
	printAnything("Hello")
	printAnything(3.14)
	printAnything([]string{"a", "b", "c"})
}

func printAnything(value interface{}) {
	fmt.Printf("  Value: %v, Type: %T\n", value, value)
}

func demo4TypeAssertions() {
	utils.PrintSubHeader("4. Type Assertions")

	var i interface{} = "Hello, World!"

	// Type assertion
	s, ok := i.(string)
	if ok {
		fmt.Printf("String value: %s (length: %d)\n", s, len(s))
	}

	// Failed type assertion
	n, ok := i.(int)
	if !ok {
		fmt.Printf("Not an integer: %v\n", n)
	}

	// Type assertion without ok (panics if wrong type)
	// s2 := i.(string) // Safe - correct type
	// n2 := i.(int)    // Would panic!

	// Practical example
	processValue("Hello")
	processValue(42)
	processValue(3.14)
}

func processValue(value interface{}) {
	switch v := value.(type) {
	case string:
		fmt.Printf("  String: %s (length: %d)\n", v, len(v))
	case int:
		fmt.Printf("  Integer: %d (doubled: %d)\n", v, v*2)
	case float64:
		fmt.Printf("  Float: %.2f (squared: %.2f)\n", v, v*v)
	default:
		fmt.Printf("  Unknown type: %T\n", v)
	}
}

func demo5TypeSwitches() {
	utils.PrintSubHeader("5. Type Switches")

	values := []interface{}{
		42,
		"Hello",
		3.14,
		true,
		[]int{1, 2, 3},
		Circle{Radius: 5},
	}

	for i, val := range values {
		fmt.Printf("%d. ", i+1)
		describeType(val)
	}
}

func describeType(value interface{}) {
	switch v := value.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %q\n", v)
	case float64:
		fmt.Printf("Float: %.2f\n", v)
	case bool:
		fmt.Printf("Boolean: %v\n", v)
	case []int:
		fmt.Printf("Integer slice: %v\n", v)
	case Shape:
		fmt.Printf("Shape with area: %.2f\n", v.Area())
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}
