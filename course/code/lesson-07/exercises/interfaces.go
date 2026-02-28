package exercises

import "fmt"

// Exercise 1: Define an interface
// Create an interface called "Shape" with methods Area() float64 and Perimeter() float64

// Exercise 2: Implement interface for Rectangle
// Define a Rectangle struct with Width and Height
// Implement Shape interface for Rectangle

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	// TODO: Implement
	return 0
}

func (r Rectangle) Perimeter() float64 {
	// TODO: Implement
	return 0
}

// Exercise 3: Implement interface for Circle
// Define a Circle struct with Radius
// Implement Shape interface for Circle

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	// TODO: Implement
	return 0
}

func (c Circle) Perimeter() float64 {
	// TODO: Implement
	return 0
}

// Exercise 4: Empty interface
// Create a function that accepts any type and returns the type as a string
func GetType(i interface{}) string {
	// TODO: Implement using type assertion or fmt.Sprintf("%T")
	return ""
}

// Exercise 5: Type assertion
// Create a function that takes an interface{} and returns the int value if possible
func GetInt(i interface{}) (int, bool) {
	// TODO: Implement using type assertion
	return 0, false
}

// Exercise 6: Type switch
// Create a function that returns a description of the value
func Describe(i interface{}) string {
	// TODO: Implement using type switch
	return ""
}

// Exercise 7: Interface composition
// Create an interface "Reader" with Read() method
// Create an interface "Writer" with Write() method
// Create a combined interface "ReadWriter" that embeds both

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

// Exercise 8: Common Go interfaces - Stringer
// Implement the fmt.Stringer interface for a Person struct
type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	// TODO: Implement
	return ""
}

// Exercise 9: Common Go interfaces - Error
// Define a custom error type "ValidationError" with Field and Message
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	// TODO: Implement
	return ""
}

// Exercise 10: Interface as parameter
// Create a function that prints area of any Shape
func PrintArea(s interface{}) string {
	// TODO: Use type assertion to get Area
	// Return the area as formatted string
	return ""
}
