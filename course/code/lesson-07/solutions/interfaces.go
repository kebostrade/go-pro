package exercises

import (
	"fmt"
	"strings"
)

// Exercise 2: Rectangle Area
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Exercise 2: Rectangle Perimeter
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Exercise 3: Circle Area
func (c Circle) Area() float64 {
	return 3.141592653589793 * c.Radius * c.Radius
}

// Exercise 3: Circle Perimeter (Circumference)
func (c Circle) Perimeter() float64 {
	return 2 * 3.141592653589793 * c.Radius
}

// Exercise 4: Empty interface - get type
func GetType(i interface{}) string {
	return fmt.Sprintf("%T", i)
}

// Exercise 5: Type assertion
func GetInt(i interface{}) (int, bool) {
	if v, ok := i.(int); ok {
		return v, true
	}
	return 0, false
}

// Exercise 6: Type switch
func Describe(i interface{}) string {
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int: %d", v)
	case string:
		return fmt.Sprintf("string: %s", v)
	case float64:
		return fmt.Sprintf("float64: %g", v)
	case bool:
		return fmt.Sprintf("bool: %t", v)
	case nil:
		return "nil: nil"
	default:
		return fmt.Sprintf("unknown: %T", v)
	}
}

// Exercise 8: Stringer implementation
func (p Person) String() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Age)
}

// Exercise 9: Error implementation
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Exercise 10: Interface as parameter
func PrintArea(s interface{}) string {
	if shape, ok := s.(interface{ Area() float64 }); ok {
		return fmt.Sprintf("%.2f", shape.Area())
	}
	return ""
}

// Helper for testing - make Rectangle implement Shape
type Shape interface {
	Area() float64
	Perimeter() float64
}

// This makes Rectangle and Circle implement Shape
var _ Shape = Rectangle{}
var _ Shape = Circle{}
