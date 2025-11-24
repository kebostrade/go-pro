// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides curriculum lessons 7-10
package service

import "go-pro-backend/internal/domain"

// Lesson 7: Interfaces and Polymorphism
func (s *curriculumService) getComprehensiveLessonData7() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          7,
		Title:       "Interfaces and Polymorphism",
		Description: "Master interface definitions, type assertions, composition, and polymorphic design patterns in Go.",
		Duration:    "6-7 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Intermediate",
		Objectives: []string{
			"Define and implement interfaces correctly",
			"Use type assertions and type switches",
			"Understand interface composition and embedding",
			"Work with empty interfaces and reflection",
			"Implement common Go interfaces (io.Reader, io.Writer, fmt.Stringer)",
			"Apply polymorphism patterns in real code",
			"Design with interfaces for flexibility and testability",
		},
		Theory: `# Interfaces and Polymorphism

## What Are Interfaces?

An interface is a set of method signatures that define a contract. Any type that implements all the methods of an interface automatically satisfies that interface, without explicit declaration. This is called **structural subtyping** or **duck typing**.

In traditional OOP languages like Java or C++, you explicitly declare that a class implements an interface. Go takes a different approach: if it walks like a duck and quacks like a duck, it's a duck. This makes Go's interfaces incredibly flexible and powerful.

## Interface Definition

An interface is defined using the ` + "`interface`" + ` keyword:

` + "```go" + `
type Reader interface {
    Read(p []byte) (n int, err error)
}
` + "```" + `

This interface defines a contract: "Anything that has a Read method with this signature is a Reader."

## Implementing Interfaces

There's no explicit "implements" keyword in Go. If your type has all the methods defined by an interface, it automatically implements that interface:

` + "```go" + `
type File struct {
    name string
    data []byte
}

func (f *File) Read(p []byte) (n int, err error) {
    // Implementation
    copy(p, f.data)
    return len(f.data), nil
}

// File now implements Reader - no explicit declaration needed
var r Reader = &File{}
` + "```" + `

## Method Receivers Matter

For interface implementation, the receiver type matters:

- If an interface method has a pointer receiver in your type, only pointers to that type satisfy the interface
- If an interface method has a value receiver, both values and pointers satisfy the interface

` + "```go" + `
type Writer interface {
    Write(p []byte) (n int, err error)
}

type Buffer struct {
    data []byte
}

// Pointer receiver - only *Buffer satisfies Writer
func (b *Buffer) Write(p []byte) (n int, err error) {
    b.data = append(b.data, p...)
    return len(p), nil
}

// This works:
var w Writer = &Buffer{}

// This does NOT work:
// var w Writer = Buffer{}  // Compile error
` + "```" + `

## The Empty Interface

The empty interface ` + "`interface{}`" + ` has no methods, so every type satisfies it. It's useful for accepting any value:

` + "```go" + `
func Print(v interface{}) {
    fmt.Println(v)
}

Print(42)           // Works: int implements interface{}
Print("hello")      // Works: string implements interface{}
Print([]int{1, 2})  // Works: slice implements interface{}
` + "```" + `

In Go 1.18+, you can use the ` + "`any`" + ` alias instead of ` + "`interface{}`" + `:

` + "```go" + `
func Print(v any) {
    fmt.Println(v)
}
` + "```" + `

## Type Assertions

Type assertions allow you to extract the concrete type from an interface value:

` + "```go" + `
var i interface{} = "hello"

s := i.(string)           // Extract string - no error handling
s, ok := i.(string)       // Safe extraction with error check
i, ok := i.(int)          // Returns false, doesn't panic
` + "```" + `

Always use the two-value form when you're not certain of the type to avoid panics.

## Type Switches

Type switches allow checking multiple types:

` + "```go" + `
func describe(i interface{}) {
    switch v := i.(type) {
    case string:
        fmt.Printf("String: %v\\n", v)
    case int:
        fmt.Printf("Integer: %v\\n", v)
    case float64:
        fmt.Printf("Float: %v\\n", v)
    default:
        fmt.Printf("Unknown type: %T\\n", i)
    }
}
` + "```" + `

## Interface Composition

Just like structs, interfaces can be composed:

` + "```go" + `
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// ReadCloser combines both interfaces
type ReadCloser interface {
    Reader
    Closer
}
` + "```" + `

Any type that implements both Read and Close automatically implements ReadCloser.

## Common Go Interfaces

### io.Reader

` + "```go" + `
type Reader interface {
    Read(p []byte) (n int, err error)
}
` + "```" + `

Reader is one of the most important interfaces in Go. It represents the ability to read data.

### io.Writer

` + "```go" + `
type Writer interface {
    Write(p []byte) (n int, err error)
}
` + "```" + `

Writer represents the ability to write data.

### fmt.Stringer

` + "```go" + `
type Stringer interface {
    String() string
}
` + "```" + `

If your type implements Stringer, it controls how it's printed:

` + "```go" + `
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

p := Person{"Alice", 30}
fmt.Println(p)  // Prints: Alice (30 years old)
` + "```" + `

## Polymorphism in Action

Polymorphism allows writing generic code that works with different types:

` + "```go" + `
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Works with any Shape
func PrintShape(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\\n", s.Area(), s.Perimeter())
}

func main() {
    shapes := []Shape{
        Circle{Radius: 5},
        Rectangle{Width: 3, Height: 4},
    }

    for _, s := range shapes {
        PrintShape(s)
    }
}
` + "```" + `

## Interface Design Best Practices

### 1. Keep Interfaces Small

Small, focused interfaces are more flexible. The best interfaces often have just 1-2 methods:

` + "```go" + `
// Good: Small, focused
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Less flexible: Too many methods
type FileOperation interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
    Close() error
    Seek(offset int64, whence int) (int64, error)
}
` + "```" + `

### 2. Use Interface Composition

Combine small interfaces instead of creating large ones:

` + "```go" + `
type ReadCloser interface {
    Reader
    Closer
}
` + "```" + `

### 3. Define Interfaces at the Point of Use

Don't define interfaces in the package that provides the implementation. Define them in the package that uses them:

` + "```go" + `
// Bad: Implementation package defines interface
package storage
type Reader interface { ... }

// Good: Usage package defines interface
package myapp
type Reader interface { ... }
` + "```" + `

This keeps packages loosely coupled.

### 4. Prefer Concrete Types When Possible

Interfaces add abstraction but also complexity. Use them when you need the flexibility:

` + "```go" + `
// Only use interface if you need multiple implementations
func ProcessFile(f *os.File) error { ... }  // Better if you only work with files

// Use interface if you want flexibility
func ProcessData(r io.Reader) error { ... }  // Works with any reader
` + "```" + `

## Gotchas and Common Mistakes

### Nil Interfaces vs Nil Values

A nil interface is different from an interface holding a nil pointer:

` + "```go" + `
var i interface{}
var p *int
i = p  // i is NOT nil! It's an interface holding a nil pointer

if i == nil {  // FALSE
    fmt.Println("i is nil")
}

if p == nil {  // TRUE
    fmt.Println("p is nil")
}
` + "```" + `

### Method Receivers and Interfaces

` + "```go" + `
type Reader interface {
    Read([]byte) (int, error)
}

type File struct{}

// With pointer receiver
func (f *File) Read(p []byte) (int, error) { ... }

var r Reader = File{}    // ERROR: File does not implement Reader
var r Reader = &File{}   // OK
` + "```" + `

### Type Assertion Panics

` + "```go" + `
var i interface{} = "hello"

// This panics!
n := i.(int)

// This is safe
n, ok := i.(int)
if ok {
    fmt.Println(n)
} else {
    fmt.Println("Not an int")
}
` + "```" + `

## Summary

- Interfaces define contracts through method sets
- Go uses structural subtyping - automatic implementation
- Type assertions and type switches extract concrete types
- The empty interface accepts any type
- Small, focused interfaces are more powerful
- Common interfaces like Reader and Writer are used throughout Go
- Interfaces enable polymorphism and make code testable and flexible
`,
		CodeExample: `package main

import (
	"fmt"
	"math"
)

// Interface defines a contract
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle implements Shape
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Rectangle implements Shape
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Generic function works with any Shape
func PrintShape(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// Type assertion
func GetRadius(s interface{}) {
	if c, ok := s.(Circle); ok {
		fmt.Printf("Radius: %.2f\n", c.Radius)
	}
}

// Type switch
func Describe(i interface{}) {
	switch v := i.(type) {
	case Circle:
		fmt.Printf("Circle with radius %.2f\n", v.Radius)
	case Rectangle:
		fmt.Printf("Rectangle %fx%f\n", v.Width, v.Height)
	default:
		fmt.Printf("Unknown: %T\n", i)
	}
}

func main() {
	shapes := []Shape{
		Circle{Radius: 3},
		Rectangle{Width: 4, Height: 5},
	}

	for _, s := range shapes {
		PrintShape(s)
	}

	// Type assertion
	circle := Circle{Radius: 2}
	GetRadius(circle)

	// Type switch
	Describe(Circle{Radius: 5})
	Describe(Rectangle{Width: 3, Height: 4})
}
`,
		Solution: `package main

import (
	"fmt"
	"io"
	"strings"
)

// Animal interface
type Animal interface {
	Speak() string
	Move() string
}

// Dog implements Animal
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + " says: Woof!"
}

func (d Dog) Move() string {
	return d.Name + " runs on four legs"
}

// Bird implements Animal
type Bird struct {
	Name string
}

func (b Bird) Speak() string {
	return b.Name + " says: Tweet!"
}

func (b Bird) Move() string {
	return b.Name + " flies in the air"
}

// Reader and Writer interfaces
type Logger interface {
	Log(msg string) error
}

type StringLogger struct {
	builder strings.Builder
}

func (sl *StringLogger) Log(msg string) error {
	_, err := sl.builder.WriteString(msg + "\n")
	return err
}

func (sl *StringLogger) String() string {
	return sl.builder.String()
}

// Describe animal using type switch
func AnimalInfo(a Animal) {
	fmt.Printf("%s: %s, %s\n",
		getTypeName(a), a.Speak(), a.Move())
}

func getTypeName(a Animal) string {
	switch a.(type) {
	case Dog:
		return "Dog"
	case Bird:
		return "Bird"
	default:
		return "Unknown Animal"
	}
}

// Write to io.Writer
func WriteAnimals(w io.Writer, animals []Animal) {
	for _, a := range animals {
		fmt.Fprintf(w, "%v\n", a)
	}
}

func main() {
	// Create animals
	animals := []Animal{
		Dog{Name: "Buddy"},
		Bird{Name: "Tweety"},
		Dog{Name: "Max"},
	}

	// Use polymorphism
	for _, animal := range animals {
		AnimalInfo(animal)
	}

	// Use with io.Writer
	logger := &StringLogger{}
	for _, animal := range animals {
		logger.Log(animal.Speak())
	}

	fmt.Println("\nLogger output:")
	fmt.Println(logger.String())
}
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "interface-basics",
				Title:       "Basic Interface Implementation",
				Description: "Implement a Vehicle interface with multiple types.",
				Requirements: []string{
					"Define a Vehicle interface with Start() and Stop() methods",
					"Create Car struct implementing Vehicle",
					"Create Motorcycle struct implementing Vehicle",
					"Write function that works with any Vehicle",
					"Test with multiple types",
				},
				InitialCode: `package main

import "fmt"

// TODO: Define Vehicle interface

// TODO: Define Car struct

// TODO: Define Motorcycle struct

// TODO: Implement Start and Stop for Car

// TODO: Implement Start and Stop for Motorcycle

func main() {
	// Test your implementation
}
`,
				Solution: `package main

import "fmt"

type Vehicle interface {
	Start() string
	Stop() string
}

type Car struct {
	Make string
	Model string
}

func (c Car) Start() string {
	return fmt.Sprintf("%s %s engine started", c.Make, c.Model)
}

func (c Car) Stop() string {
	return fmt.Sprintf("%s %s engine stopped", c.Make, c.Model)
}

type Motorcycle struct {
	Make string
	Model string
}

func (m Motorcycle) Start() string {
	return fmt.Sprintf("%s %s engine roared to life", m.Make, m.Model)
}

func (m Motorcycle) Stop() string {
	return fmt.Sprintf("%s %s engine shut down", m.Make, m.Model)
}

func OperateVehicle(v Vehicle) {
	fmt.Println(v.Start())
	fmt.Println(v.Stop())
}

func main() {
	car := Car{Make: "Toyota", Model: "Camry"}
	moto := Motorcycle{Make: "Harley", Model: "Davidson"}

	vehicles := []Vehicle{car, moto}
	for _, v := range vehicles {
		OperateVehicle(v)
		fmt.Println()
	}
}
`,
			},
			{
				ID:          "type-assertions",
				Title:       "Type Assertions and Type Switches",
				Description: "Practice extracting concrete types from interfaces.",
				Requirements: []string{
					"Accept interface{} values",
					"Use type assertions to extract specific types",
					"Use type switch for multiple types",
					"Handle unknown types gracefully",
				},
				InitialCode: `package main

import "fmt"

func ProcessData(data interface{}) {
	// TODO: Use type assertion to check for string
	// TODO: Use type assertion to check for int
	// TODO: Use type switch for other types
}

func main() {
	ProcessData("hello")
	ProcessData(42)
	ProcessData(3.14)
	ProcessData([]int{1, 2, 3})
}
`,
				Solution: `package main

import "fmt"

func ProcessData(data interface{}) {
	// Safe type assertion
	if str, ok := data.(string); ok {
		fmt.Printf("String: %s (length: %d)\n", str, len(str))
		return
	}

	if num, ok := data.(int); ok {
		fmt.Printf("Integer: %d (doubled: %d)\n", num, num*2)
		return
	}

	// Type switch for multiple types
	switch v := data.(type) {
	case float64:
		fmt.Printf("Float: %.2f\n", v)
	case []int:
		fmt.Printf("Slice: %v (sum: %d)\n", v, sum(v))
	case bool:
		fmt.Printf("Boolean: %v\n", v)
	default:
		fmt.Printf("Unknown type: %T with value %v\n", data, data)
	}
}

func sum(slice []int) int {
	total := 0
	for _, v := range slice {
		total += v
	}
	return total
}

func main() {
	ProcessData("hello")
	ProcessData(42)
	ProcessData(3.14)
	ProcessData([]int{1, 2, 3})
	ProcessData(true)
}
`,
			},
			{
				ID:          "interface-composition",
				Title:       "Interface Composition",
				Description: "Create and use composed interfaces.",
				Requirements: []string{
					"Define Reader interface with Read() method",
					"Define Writer interface with Write() method",
					"Compose ReadWriter from Reader and Writer",
					"Implement for a type",
					"Use composed interface",
				},
				InitialCode: `package main

import (
	"fmt"
)

// TODO: Define Reader interface

// TODO: Define Writer interface

// TODO: Define ReadWriter by composing Reader and Writer

type Buffer struct {
	data string
}

// TODO: Implement Read method

// TODO: Implement Write method

func main() {
	var rw ReadWriter = &Buffer{}
	fmt.Println(rw)
}
`,
				Solution: `package main

import (
	"fmt"
)

type Reader interface {
	Read() string
}

type Writer interface {
	Write(data string) error
}

type ReadWriter interface {
	Reader
	Writer
}

type Buffer struct {
	data string
}

func (b *Buffer) Read() string {
	return b.data
}

func (b *Buffer) Write(data string) error {
	b.data = data
	fmt.Println("Data written:", data)
	return nil
}

func ProcessBuffer(rw ReadWriter) {
	rw.Write("Hello, World!")
	fmt.Println("Data in buffer:", rw.Read())
}

func main() {
	buf := &Buffer{}
	ProcessBuffer(buf)
}
`,
			},
			{
				ID:          "stringer-interface",
				Title:       "Implementing fmt.Stringer",
				Description: "Implement String() method to control printing.",
				Requirements: []string{
					"Create Person struct with Name, Age, and City",
					"Implement fmt.Stringer interface",
					"Override default string representation",
					"Use in fmt functions",
				},
				InitialCode: `package main

import "fmt"

type Person struct {
	Name string
	Age  int
	City string
}

// TODO: Implement String() method for Person

func main() {
	p := Person{Name: "Alice", Age: 30, City: "NYC"}
	fmt.Println(p)
	fmt.Printf("Person: %s\n", p)
}
`,
				Solution: `package main

import "fmt"

type Person struct {
	Name string
	Age  int
	City string
}

func (p Person) String() string {
	return fmt.Sprintf("%s, %d years old, from %s", p.Name, p.Age, p.City)
}

func main() {
	people := []Person{
		{Name: "Alice", Age: 30, City: "NYC"},
		{Name: "Bob", Age: 25, City: "LA"},
		{Name: "Charlie", Age: 35, City: "Chicago"},
	}

	for _, p := range people {
		fmt.Println(p)
	}
}
`,
			},
			{
				ID:          "polymorphic-container",
				Title:       "Polymorphic Container",
				Description: "Store different types in a single slice using interfaces.",
				Requirements: []string{
					"Define Employee interface with GetSalary()",
					"Create Manager and Developer types",
					"Store in slice of interface{}",
					"Calculate total salary using type assertion",
					"Handle different employee types",
				},
				InitialCode: `package main

import "fmt"

type Employee interface {
	GetSalary() float64
	GetName() string
}

// TODO: Define Manager struct

// TODO: Define Developer struct

// TODO: Implement Employee for Manager

// TODO: Implement Employee for Developer

func main() {
	// TODO: Create employees and add to slice
	// TODO: Calculate total salary
}
`,
				Solution: `package main

import "fmt"

type Employee interface {
	GetSalary() float64
	GetName() string
}

type Manager struct {
	Name     string
	BaseSalary float64
	Bonus    float64
}

func (m Manager) GetSalary() float64 {
	return m.BaseSalary + m.Bonus
}

func (m Manager) GetName() string {
	return m.Name
}

type Developer struct {
	Name    string
	Salary  float64
}

func (d Developer) GetSalary() float64 {
	return d.Salary
}

func (d Developer) GetName() string {
	return d.Name
}

func main() {
	employees := []Employee{
		Manager{Name: "Alice", BaseSalary: 80000, Bonus: 20000},
		Developer{Name: "Bob", Salary: 100000},
		Manager{Name: "Charlie", BaseSalary: 75000, Bonus: 15000},
		Developer{Name: "David", Salary: 95000},
	}

	totalSalary := 0.0
	for _, emp := range employees {
		fmt.Printf("%s: $%.2f\n", emp.GetName(), emp.GetSalary())
		totalSalary += emp.GetSalary()
	}

	fmt.Printf("\nTotal Payroll: $%.2f\n", totalSalary)
}
`,
			},
			{
				ID:          "reader-writer",
				Title:       "Implementing io.Reader and io.Writer",
				Description: "Create custom types implementing standard Go interfaces.",
				Requirements: []string{
					"Implement io.Reader for a custom type",
					"Implement io.Writer for a custom type",
					"Work with standard library functions",
					"Handle multiple reads/writes correctly",
				},
				InitialCode: `package main

import (
	"fmt"
	"io"
)

type StringReader struct {
	data string
	pos  int
}

// TODO: Implement Read method for StringReader

type FileWriter struct {
	content string
}

// TODO: Implement Write method for FileWriter

func main() {
	reader := &StringReader{data: "Hello, Go!"}
	p := make([]byte, 5)

	for {
		n, err := reader.Read(p)
		if err != nil {
			break
		}
		fmt.Printf("Read %d bytes: %s\n", n, p[:n])
	}
}
`,
				Solution: `package main

import (
	"fmt"
	"io"
)

type StringReader struct {
	data string
	pos  int
}

func (sr *StringReader) Read(p []byte) (n int, err error) {
	if sr.pos >= len(sr.data) {
		return 0, io.EOF
	}

	n = copy(p, sr.data[sr.pos:])
	sr.pos += n
	return n, nil
}

type FileWriter struct {
	content string
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	fw.content += string(p)
	return len(p), nil
}

func main() {
	// Test StringReader
	reader := &StringReader{data: "Hello, Go!"}
	p := make([]byte, 5)

	fmt.Println("Reading from StringReader:")
	for {
		n, err := reader.Read(p)
		if err != nil {
			break
		}
		fmt.Printf("Read %d bytes: %s\n", n, p[:n])
	}

	// Test FileWriter
	fmt.Println("\nWriting with FileWriter:")
	writer := &FileWriter{}
	writer.Write([]byte("Hello "))
	writer.Write([]byte("World!"))
	fmt.Printf("Written content: %s\n", writer.content)
}
`,
			},
			{
				ID:          "interface-nil-gotcha",
				Title:       "Understanding Nil Interfaces",
				Description: "Understand the difference between nil interfaces and interfaces holding nil values.",
				Requirements: []string{
					"Create interface and check nil interface",
					"Assign nil value to interface",
					"Distinguish nil interface from nil value",
					"Handle properly in conditionals",
				},
				InitialCode: `package main

import "fmt"

type Reader interface {
	Read([]byte) (int, error)
}

func CheckNil(r Reader) {
	// TODO: Check if r is nil
	// TODO: Handle the difference
}

func main() {
	var r Reader
	CheckNil(r)

	var p *int
	fmt.Printf("r == nil: %v\n", r == nil)
	fmt.Printf("p == nil: %v\n", p == nil)
}
`,
				Solution: `package main

import "fmt"

type Reader interface {
	Read([]byte) (int, error)
}

type File struct {
	name string
}

func (f *File) Read(p []byte) (int, error) {
	return 0, nil
}

func CheckNil(r Reader) {
	if r == nil {
		fmt.Println("Interface is nil")
		return
	}

	fmt.Println("Interface has a value (even if that value is nil)")

	// Extract the concrete value
	if f, ok := r.(*File); ok && f == nil {
		fmt.Println("Interface holds a nil *File pointer")
	}
}

func main() {
	// Nil interface
	var r Reader
	fmt.Println("Nil interface:")
	CheckNil(r)

	// Interface holding nil pointer
	var f *File
	fmt.Println("\nInterface holding nil pointer:")
	r = f
	CheckNil(r)

	// Interface holding value
	fmt.Println("\nInterface holding value:")
	r = &File{name: "test.txt"}
	CheckNil(r)
}
`,
			},
			{
				ID:          "multiple-interfaces",
				Title:       "Implementing Multiple Interfaces",
				Description: "Create types that implement multiple interfaces.",
				Requirements: []string{
					"Create Reader interface",
					"Create Writer interface",
					"Create Closer interface",
					"Implement all three for a single type",
					"Use all interfaces in code",
				},
				InitialCode: `package main

import "fmt"

type Reader interface {
	Read() string
}

type Writer interface {
	Write(string) error
}

type Closer interface {
	Close() error
}

// TODO: Create a type implementing all three interfaces

func main() {
	// Test implementation
}
`,
				Solution: `package main

import (
	"fmt"
)

type Reader interface {
	Read() string
}

type Writer interface {
	Write(string) error
}

type Closer interface {
	Close() error
}

type Connection struct {
	data   string
	closed bool
}

func (c *Connection) Read() string {
	if c.closed {
		return ""
	}
	return c.data
}

func (c *Connection) Write(s string) error {
	if c.closed {
		return fmt.Errorf("connection closed")
	}
	c.data = s
	fmt.Println("Wrote:", s)
	return nil
}

func (c *Connection) Close() error {
	c.closed = true
	fmt.Println("Connection closed")
	return nil
}

func UseConnection(r Reader, w Writer, cl Closer) {
	w.Write("test data")
	fmt.Println("Read:", r.Read())
	cl.Close()
}

func main() {
	conn := &Connection{}
	UseConnection(conn, conn, conn)
}
`,
			},
		},
		NextLessonID: func() *int { i := 8; return &i }(),
		PrevLessonID: func() *int { i := 6; return &i }(),
	}
}

// Lesson 8: Error Handling Patterns
func (s *curriculumService) getComprehensiveLessonData8() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          8,
		Title:       "Error Handling Patterns",
		Description: "Master error handling in Go including custom errors, error wrapping, panic/recover, and best practices.",
		Duration:    "4-5 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Intermediate",
		Objectives: []string{
			"Understand the error interface",
			"Create custom error types",
			"Use error wrapping with fmt.Errorf and %w",
			"Use errors.Is() and errors.As()",
			"Understand panic and recover",
			"Follow Go error handling best practices",
			"Design error handling strategies",
		},
		Theory: `# Error Handling Patterns

## The Error Interface

In Go, errors are values. The built-in ` + "`error`" + ` interface is extremely simple:

` + "```go" + `
type error interface {
    Error() string
}
` + "```" + `

Any type with an Error() method is an error. This simplicity is powerful.

## Creating Errors

### Using errors.New()

` + "```go" + `
import "errors"

var ErrDivisionByZero = errors.New("division by zero")

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, ErrDivisionByZero
    }
    return a / b, nil
}
` + "```" + `

### Using fmt.Errorf()

For dynamic error messages:

` + "```go" + `
import "fmt"

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide %d by zero", a)
    }
    return a / b, nil
}
` + "```" + `

## Error Wrapping (Go 1.13+)

Error wrapping preserves the error chain, allowing caller to understand what went wrong:

` + "```go" + `
result, err := someFunction()
if err != nil {
    // Wrap error with additional context
    return fmt.Errorf("failed to process: %w", err)
}
` + "```" + `

The ` + "`%w`" + ` verb wraps the error. Without it:

` + "```go" + `
// Without wrapping - loses original error
return fmt.Errorf("failed to process: %v", err)

// With wrapping - preserves error chain
return fmt.Errorf("failed to process: %w", err)
` + "```" + `

## Examining Errors: errors.Is() and errors.As()

### errors.Is()

Check if an error in the chain matches a specific error:

` + "```go" + `
import "errors"

result, err := processFile("test.txt")
if err != nil {
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File not found")
    }
}
` + "```" + `

### errors.As()

Extract a specific error type from the chain:

` + "```go" + `
var syntaxErr *json.SyntaxError
if errors.As(err, &syntaxErr) {
    fmt.Printf("Syntax error at offset %d\n", syntaxErr.Offset)
}
` + "```" + `

## Custom Error Types

For more control, create custom error types:

` + "```go" + `
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func ValidateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return &ValidationError{
            Field:   "email",
            Message: "must contain @",
        }
    }
    return nil
}
` + "```" + `

Custom error types allow extracting additional information:

` + "```go" + `
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Field %s failed validation: %s\n", ve.Field, ve.Message)
}
` + "```" + `

## Sentinel Errors

Named error values that represent specific conditions:

` + "```go" + `
var (
    ErrInvalidInput = errors.New("invalid input")
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

func GetUser(id int) (*User, error) {
    if id < 0 {
        return nil, ErrInvalidInput
    }
    if user == nil {
        return nil, ErrNotFound
    }
    return user, nil
}

// Using sentinel errors
if err != nil {
    if errors.Is(err, ErrNotFound) {
        // Handle not found
    } else if errors.Is(err, ErrInvalidInput) {
        // Handle invalid input
    }
}
` + "```" + `

## Panic and Recover

### When to Use Panic

Panic is for truly exceptional circumstances, not normal error handling:

- Programming errors (nil dereference, array out of bounds, etc.)
- Fatal startup errors
- Impossible conditions indicating a bug in the code

### Using Panic and Recover

` + "```go" + `
func safeDiv(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("runtime error: %v", r)
        }
    }()

    if b == 0 {
        panic("division by zero")
    }

    return a / b, nil
}

// Usage
result, err := safeDiv(10, 0)
if err != nil {
    fmt.Println("Error:", err)
}
` + "```" + `

### Key Points about Panic/Recover

1. Only use recover() in defer functions
2. Recover returns nil if no panic occurred
3. Recovered panics stop the panic propagation
4. Most Go code doesn't use panic/recover
5. Use error returns for normal error handling

## Error Handling Patterns

### Pattern 1: Return Error for Possible Failures

` + "```go" + `
func OpenFile(name string) (*os.File, error) {
    file, err := os.Open(name)
    if err != nil {
        return nil, fmt.Errorf("open file: %w", err)
    }
    return file, nil
}
` + "```" + `

### Pattern 2: Check Errors Immediately

` + "```go" + `
file, err := os.Open("data.txt")
if err != nil {
    // Handle error immediately
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()
` + "```" + `

### Pattern 3: Don't Ignore Errors

` + "```go" + `
// Bad: Ignoring error
file.Close()

// Good: Handle error if needed
if err := file.Close(); err != nil {
    return fmt.Errorf("failed to close: %w", err)
}
` + "```" + `

### Pattern 4: Add Context When Wrapping

` + "```go" + `
data, err := readFile("config.txt")
if err != nil {
    // Add context about what you were doing
    return fmt.Errorf("loading configuration: %w", err)
}
` + "```" + `

### Pattern 5: Use Type Assertions for Special Handling

` + "```go" + `
if err := operation(); err != nil {
    if opErr, ok := err.(*os.PathError); ok {
        fmt.Printf("Operation %s failed on path %s\n",
            opErr.Op, opErr.Path)
    }
}
` + "```" + `

## Error Handling Best Practices

### 1. Return Errors, Don't Log and Continue

` + "```go" + `
// Bad: Logs error but continues
file, err := os.Open("data.txt")
if err != nil {
    log.Println("Could not open file")
    file = fallbackFile
}

// Good: Return error to let caller decide
file, err := os.Open("data.txt")
if err != nil {
    return fmt.Errorf("open file: %w", err)
}
` + "```" + `

### 2. Use Errors.Is() for Sentinel Errors

` + "```go" + `
// Good
if errors.Is(err, os.ErrNotExist) {
    // Handle not found
}

// Avoid
if err == os.ErrNotExist {
    // Doesn't work if error is wrapped
}
` + "```" + `

### 3. Wrap Errors with Context

` + "```go" + `
// Wrap error with %w to preserve chain
if err != nil {
    return fmt.Errorf("processing user %d: %w", userID, err)
}
` + "```" + `

### 4. Create Custom Types for Complex Scenarios

` + "```go" + `
type RequestError struct {
    Code    int
    Message string
    Err     error
}

func (e *RequestError) Error() string {
    return fmt.Sprintf("request error %d: %s", e.Code, e.Message)
}

func (e *RequestError) Unwrap() error {
    return e.Err
}
` + "```" + `

### 5. Don't Use Panic for Normal Errors

` + "```go" + `
// Bad: Using panic for expected errors
func ParseInt(s string) int {
    n, err := strconv.Atoi(s)
    if err != nil {
        panic(err)  // Wrong!
    }
    return n
}

// Good: Return error
func ParseInt(s string) (int, error) {
    return strconv.Atoi(s)
}
` + "```" + `

## Error Decoration vs Information Loss

When wrapping errors, you preserve the error chain:

` + "```go" + `
// Original error
if err := readFile(); err != nil {
    return fmt.Errorf("initialization: %w", err)
}

// Chain: "initialization: " -> "open file: " -> "permission denied"
` + "```" + `

Compare with %v:

` + "```go" + `
// Using %v loses the original type information
if err := readFile(); err != nil {
    return fmt.Errorf("initialization: %v", err)
}
// Chain is broken - original error info is just a string
` + "```" + `

## Summary

- The error interface is simple: just Error() string
- Use errors.New() for simple errors
- Wrap errors with %w to preserve the chain
- Use errors.Is() and errors.As() to examine errors
- Create custom error types for complex scenarios
- Panic is for exceptional conditions, not normal errors
- Return errors and let callers decide how to handle them
- Add context when wrapping errors
- Don't ignore errors without good reason
`,
		CodeExample: `package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Sentinel errors
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrDivideByZero = errors.New("cannot divide by zero")
)

// Custom error type
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// Error wrapping example
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("divide operation: %w", ErrDivideByZero)
	}
	return a / b, nil
}

// Using errors.Is()
func ProcessDivision(a, b int) error {
	result, err := Divide(a, b)
	if err != nil {
		if errors.Is(err, ErrDivideByZero) {
			return fmt.Errorf("calculation failed: %w", err)
		}
		return err
	}
	fmt.Printf("Result: %d\n", result)
	return nil
}

// Using errors.As()
func ValidateUser(name, email string) error {
	if name == "" {
		return &ValidationError{Field: "name", Message: "required"}
	}

	n, err := strconv.Atoi(email)
	if err != nil {
		if parseErr, ok := err.(*strconv.NumError); ok {
			return &ValidationError{
				Field:   "email",
				Message: fmt.Sprintf("invalid: %s", parseErr.Err),
			}
		}
	}

	return nil
}

// Panic and recover example
func SafeOperation(value int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()

	if value < 0 {
		panic("negative value not allowed")
	}
	return value * 2, nil
}

func main() {
	fmt.Println("=== Error Wrapping ===")
	if err := ProcessDivision(10, 0); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== Type Assertion ===")
	var ve *ValidationError
	err := ValidateUser("", "john")
	if errors.As(err, &ve) {
		fmt.Printf("Validation failed on %s: %s\n", ve.Field, ve.Message)
	}

	fmt.Println("\n=== Panic/Recover ===")
	result, err := SafeOperation(-5)
	if err != nil {
		fmt.Printf("Operation error: %v\n", err)
	}
	result, err = SafeOperation(5)
	if err == nil {
		fmt.Printf("Result: %d\n", result)
	}
}
`,
		Solution: `package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Sentinel errors
var (
	ErrEmptyInput      = errors.New("input cannot be empty")
	ErrInvalidAge      = errors.New("age must be positive")
	ErrInvalidEmail    = errors.New("email must contain @")
	ErrDatabaseConnect = errors.New("database connection failed")
)

// Custom error types
type InputError struct {
	Input   string
	Message string
}

func (e *InputError) Error() string {
	return fmt.Sprintf("input error: %q - %s", e.Input, e.Message)
}

type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s on table %s: %v",
		e.Operation, e.Table, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// Validation function
func ValidateAge(age string) error {
	n, err := strconv.Atoi(age)
	if err != nil {
		return &InputError{
			Input:   age,
			Message: "must be a valid number",
		}
	}

	if n < 0 || n > 150 {
		return fmt.Errorf("age validation: %w", ErrInvalidAge)
	}

	return nil
}

// Function with error wrapping
func ReadUserFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading user file: %w", err)
	}
	return data, nil
}

// Function using custom errors
func CreateUser(name, age, email string) error {
	if name == "" {
		return &InputError{Input: name, Message: "name required"}
	}

	if err := ValidateAge(age); err != nil {
		return fmt.Errorf("user validation: %w", err)
	}

	if !contains(email, "@") {
		return fmt.Errorf("email validation: %w", ErrInvalidEmail)
	}

	// Simulate database operation
	if err := saveToDatabase(name); err != nil {
		return &DatabaseError{
			Operation: "INSERT",
			Table:     "users",
			Err:       err,
		}
	}

	return nil
}

func saveToDatabase(name string) error {
	// Simulated error
	if name == "admin" {
		return ErrDatabaseConnect
	}
	return nil
}

// Error handling with type assertion
func HandleError(err error) {
	var inputErr *InputError
	var dbErr *DatabaseError

	switch {
	case errors.Is(err, ErrInvalidAge):
		fmt.Println("Please provide a valid age")
	case errors.Is(err, ErrInvalidEmail):
		fmt.Println("Email must contain @ symbol")
	case errors.As(err, &inputErr):
		fmt.Printf("Invalid input: %v\n", inputErr)
	case errors.As(err, &dbErr):
		fmt.Printf("Database failed: %s on %s\n",
			dbErr.Operation, dbErr.Table)
	default:
		fmt.Printf("Unhandled error: %v\n", err)
	}
}

func contains(s string, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("=== Testing Error Handling ===\n")

	// Test 1: Valid user
	fmt.Println("Test 1: Valid user")
	err := CreateUser("Alice", "30", "alice@example.com")
	if err != nil {
		HandleError(err)
	} else {
		fmt.Println("User created successfully")
	}

	// Test 2: Invalid age (not a number)
	fmt.Println("\nTest 2: Invalid age format")
	err = CreateUser("Bob", "abc", "bob@example.com")
	HandleError(err)

	// Test 3: Invalid age (out of range)
	fmt.Println("\nTest 3: Invalid age (out of range)")
	err = CreateUser("Charlie", "200", "charlie@example.com")
	HandleError(err)

	// Test 4: Invalid email
	fmt.Println("\nTest 4: Invalid email")
	err = CreateUser("David", "25", "david.example.com")
	HandleError(err)

	// Test 5: Database error
	fmt.Println("\nTest 5: Database error")
	err = CreateUser("admin", "40", "admin@example.com")
	HandleError(err)
}
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "sentinel-errors",
				Title:       "Sentinel Errors",
				Description: "Create and use sentinel errors for specific conditions.",
				Requirements: []string{
					"Define sentinel errors for different operations",
					"Return appropriate errors",
					"Use errors.Is() to check for specific errors",
				},
				InitialCode: `package main

import (
	"errors"
	"fmt"
)

// TODO: Define sentinel errors

func ParseInput(input string) (int, error) {
	if input == "" {
		// TODO: Return appropriate error
	}
	// TODO: Parse and validate
	return 0, nil
}

func main() {
	result, err := ParseInput("")
	// TODO: Check error with errors.Is()
}
`,
				Solution: `package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrEmptyInput    = errors.New("input cannot be empty")
	ErrInvalidNumber = errors.New("input is not a valid number")
)

func ParseInput(input string) (int, error) {
	if input == "" {
		return 0, ErrEmptyInput
	}

	n, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("parse input: %w", ErrInvalidNumber)
	}

	if n < 0 {
		return 0, fmt.Errorf("parse input: negative numbers not allowed")
	}

	return n, nil
}

func main() {
	tests := []string{"", "abc", "-5", "42"}

	for _, test := range tests {
		result, err := ParseInput(test)
		if err != nil {
			if errors.Is(err, ErrEmptyInput) {
				fmt.Printf("Error for %q: empty input\n", test)
			} else if errors.Is(err, ErrInvalidNumber) {
				fmt.Printf("Error for %q: invalid number\n", test)
			} else {
				fmt.Printf("Error for %q: %v\n", test, err)
			}
		} else {
			fmt.Printf("Parsed %q: %d\n", test, result)
		}
	}
}
`,
			},
			{
				ID:          "custom-errors",
				Title:       "Custom Error Types",
				Description: "Create custom error types with additional information.",
				Requirements: []string{
					"Define a custom error struct",
					"Implement Error() method",
					"Return custom error from function",
					"Extract using errors.As()",
				},
				InitialCode: `package main

import (
	"errors"
	"fmt"
)

// TODO: Define custom error type

type Config struct {
	Port int
}

func LoadConfig(port int) (*Config, error) {
	if port < 1 || port > 65535 {
		// TODO: Return custom error with details
	}
	return &Config{Port: port}, nil
}

func main() {
	_, err := LoadConfig(99999)
	if err != nil {
		// TODO: Extract and handle custom error
	}
}
`,
				Solution: `package main

import (
	"errors"
	"fmt"
)

type ConfigError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: field %q with value %v - %s",
		e.Field, e.Value, e.Message)
}

type Config struct {
	Port int
}

func LoadConfig(port int) (*Config, error) {
	if port < 1 || port > 65535 {
		return nil, &ConfigError{
			Field:   "port",
			Value:   port,
			Message: "must be between 1 and 65535",
		}
	}
	return &Config{Port: port}, nil
}

func main() {
	tests := []int{0, 8080, 99999}

	for _, port := range tests {
		cfg, err := LoadConfig(port)
		if err != nil {
			var cfgErr *ConfigError
			if errors.As(err, &cfgErr) {
				fmt.Printf("Failed: %s\n", cfgErr.Message)
			}
		} else {
			fmt.Printf("Config loaded: port %d\n", cfg.Port)
		}
	}
}
`,
			},
			{
				ID:          "error-wrapping",
				Title:       "Error Wrapping with Context",
				Description: "Wrap errors to add context while preserving the original error.",
				Requirements: []string{
					"Wrap errors using fmt.Errorf with %w",
					"Add meaningful context",
					"Verify error chain with errors.Is()",
				},
				InitialCode: `package main

import (
	"errors"
	"fmt"
)

var ErrFileNotFound = errors.New("file not found")

func ReadFile(filename string) (string, error) {
	if filename == "" {
		// TODO: Wrap error with context
		return "", ErrFileNotFound
	}
	return "file contents", nil
}

func ProcessData(filename string) (string, error) {
	// TODO: Call ReadFile and wrap any error
	data, err := ReadFile(filename)
	if err != nil {
		// TODO: Add context about what you were doing
	}
	return data, nil
}

func main() {
	result, err := ProcessData("")
	// TODO: Check error chain
}
`,
				Solution: `package main

import (
	"errors"
	"fmt"
)

var ErrFileNotFound = errors.New("file not found")

func ReadFile(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("read file: %w", ErrFileNotFound)
	}
	return "file contents", nil
}

func ProcessData(filename string) (string, error) {
	data, err := ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("processing data: %w", err)
	}
	return data, nil
}

func main() {
	_, err := ProcessData("")
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println()

		// Check if original error is in chain
		if errors.Is(err, ErrFileNotFound) {
			fmt.Println("Original error found in chain")
		}

		// Print formatted error
		fmt.Printf("Detailed error: %+v\n", err)
	}
}
`,
			},
			{
				ID:          "panic-recover",
				Title:       "Panic and Recover",
				Description: "Use panic and recover for exceptional conditions.",
				Requirements: []string{
					"Create function that panics",
					"Use defer and recover",
					"Convert panic to error return",
				},
				InitialCode: `package main

import "fmt"

func SafeSquareRoot(n float64) (result float64, err error) {
	// TODO: Set up defer with recover
	// TODO: Panic if negative
	// TODO: Return result
	return 0, nil
}

func main() {
	result, err := SafeSquareRoot(-4)
	// TODO: Handle error
}
`,
				Solution: `package main

import (
	"fmt"
	"math"
)

func SafeSquareRoot(n float64) (result float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	if n < 0 {
		panic("cannot take square root of negative number")
	}

	return math.Sqrt(n), nil
}

func main() {
	tests := []float64{4, -4, 16, -1, 0}

	for _, n := range tests {
		result, err := SafeSquareRoot(n)
		if err != nil {
			fmt.Printf("SafeSquareRoot(%.1f): error - %v\n", n, err)
		} else {
			fmt.Printf("SafeSquareRoot(%.1f): %.2f\n", n, result)
		}
	}
}
`,
			},
			{
				ID:          "error-assertion",
				Title:       "Type Assertion for Errors",
				Description: "Extract type-specific information from errors.",
				Requirements: []string{
					"Create custom error type",
					"Return custom error",
					"Use errors.As() to extract",
					"Access error-specific fields",
				},
				InitialCode: `package main

import (
	"errors"
	"fmt"
)

// TODO: Define custom error type with Code and Message

type API struct{}

func (a *API) FetchData(code int) (string, error) {
	if code < 200 || code >= 300 {
		// TODO: Return custom error
	}
	return "data", nil
}

func main() {
	api := &API{}
	_, err := api.FetchData(404)
	// TODO: Extract error and access fields
}
`,
				Solution: `package main

import (
	"errors"
	"fmt"
)

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

type API struct{}

func (a *API) FetchData(code int) (string, error) {
	if code < 200 || code >= 300 {
		return "", &HTTPError{
			Code:    code,
			Message: httpStatus(code),
		}
	}
	return "data", nil
}

func httpStatus(code int) string {
	switch code {
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	case 401:
		return "Unauthorized"
	default:
		return "Unknown Error"
	}
}

func main() {
	api := &API{}
	codes := []int{200, 404, 500, 401}

	for _, code := range codes {
		_, err := api.FetchData(code)
		if err != nil {
			var httpErr *HTTPError
			if errors.As(err, &httpErr) {
				fmt.Printf("Status %d (%s): %s\n",
					httpErr.Code,
					httpStatus(httpErr.Code),
					httpErr.Message)
			}
		} else {
			fmt.Printf("Code %d: Success\n", code)
		}
	}
}
`,
			},
			{
				ID:          "error-handling-strategy",
				Title:       "Error Handling Strategy",
				Description: "Design comprehensive error handling for a multi-function workflow.",
				Requirements: []string{
					"Create multiple functions with different error types",
					"Handle different errors appropriately",
					"Use error wrapping for context",
					"Implement error checking flow",
				},
				InitialCode: `package main

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUser     = errors.New("invalid user")
	ErrNotAuthorized   = errors.New("not authorized")
	ErrDataNotFound    = errors.New("data not found")
)

// TODO: Implement GetUser function
func GetUser(id int) (string, error) {
	return "", nil
}

// TODO: Implement CheckAuth function
func CheckAuth(user string) error {
	return nil
}

// TODO: Implement FetchData function
func FetchData(user string) (string, error) {
	return "", nil
}

func ProcessRequest(userID int) (string, error) {
	// TODO: Implement full workflow with error handling
	return "", nil
}

func main() {
	// TODO: Test with different scenarios
}
`,
				Solution: `package main

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUser   = errors.New("invalid user")
	ErrNotAuthorized = errors.New("not authorized")
	ErrDataNotFound  = errors.New("data not found")
)

func GetUser(id int) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("get user: %w", ErrInvalidUser)
	}
	if id == 999 {
		return "", fmt.Errorf("get user: %w", ErrDataNotFound)
	}
	return fmt.Sprintf("user_%d", id), nil
}

func CheckAuth(user string) error {
	if user == "guest" {
		return fmt.Errorf("auth check: %w", ErrNotAuthorized)
	}
	return nil
}

func FetchData(user string) (string, error) {
	if user == "" {
		return "", fmt.Errorf("fetch data: %w", ErrInvalidUser)
	}
	return "sensitive_data", nil
}

func ProcessRequest(userID int) (string, error) {
	user, err := GetUser(userID)
	if err != nil {
		return "", fmt.Errorf("process: %w", err)
	}

	if err := CheckAuth(user); err != nil {
		return "", fmt.Errorf("process: %w", err)
	}

	data, err := FetchData(user)
	if err != nil {
		return "", fmt.Errorf("process: %w", err)
	}

	return data, nil
}

func main() {
	tests := []int{-1, 1, 2, 999}

	for _, id := range tests {
		data, err := ProcessRequest(id)
		if err != nil {
			if errors.Is(err, ErrInvalidUser) {
				fmt.Printf("Request %d: Invalid user\n", id)
			} else if errors.Is(err, ErrNotAuthorized) {
				fmt.Printf("Request %d: Not authorized\n", id)
			} else if errors.Is(err, ErrDataNotFound) {
				fmt.Printf("Request %d: Data not found\n", id)
			} else {
				fmt.Printf("Request %d: %v\n", id, err)
			}
		} else {
			fmt.Printf("Request %d: Success - %s\n", id, data)
		}
	}
}
`,
			},
		},
		NextLessonID: func() *int { i := 9; return &i }(),
		PrevLessonID: func() *int { i := 7; return &i }(),
	}
}

// Lesson 9: Goroutines and Channels
func (s *curriculumService) getComprehensiveLessonData9() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          9,
		Title:       "Goroutines and Channels",
		Description: "Master concurrent programming with goroutines, channels, select statements, and common concurrency patterns.",
		Duration:    "7-8 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Intermediate",
		Objectives: []string{
			"Create and manage goroutines",
			"Understand goroutine lifecycle",
			"Use channels for communication",
			"Distinguish buffered and unbuffered channels",
			"Implement channel directions",
			"Master select statement patterns",
			"Avoid race conditions",
			"Implement common concurrency patterns",
		},
		Theory: `# Goroutines and Channels

## Goroutines

A goroutine is a lightweight thread managed by the Go runtime. Creating goroutines is cheap and easy:

` + "```go" + `
go function()      // Run function concurrently
go obj.Method()    // Run method concurrently
` + "```" + `

The key difference from operating system threads:
- Goroutines are multiplexed onto a small number of OS threads
- Starting thousands of goroutines is practical
- Context switching is efficient

## Simple Goroutine Example

` + "```go" + `
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        fmt.Println("Hello from goroutine")
    }()

    time.Sleep(time.Second)
    fmt.Println("Main function")
}
` + "```" + `

Important: The main goroutine must wait for others to complete. If main exits, all goroutines are terminated.

## Channels

Channels are typed conduits for communication between goroutines:

` + "```go" + `
make(chan int)           // Unbuffered channel
make(chan int, 5)        // Buffered channel with capacity 5
` + "```" + `

### Send and Receive

` + "```go" + `
ch := make(chan string)

// Send
ch <- "hello"

// Receive
msg := <-ch

// Non-blocking receive (rarely useful)
msg, ok := <-ch
` + "```" + `

### Unbuffered Channels

Unbuffered channels block until both sender and receiver are ready:

` + "```go" + `
func main() {
    ch := make(chan int)

    go func() {
        ch <- 42  // Blocks until main receives
    }()

    value := <-ch  // Blocks until goroutine sends
    fmt.Println(value)
}
` + "```" + `

Both operations must happen - the send blocks until receive, and vice versa.

### Buffered Channels

Buffered channels have a capacity. Send blocks only when buffer is full:

` + "```go" + `
ch := make(chan int, 2)

ch <- 1  // OK
ch <- 2  // OK
ch <- 3  // Blocks - buffer full!

value := <-ch  // Receive from buffer
` + "```" + `

### Closing Channels

Close a channel when no more values will be sent:

` + "```go" + `
ch := make(chan int)

go func() {
    ch <- 1
    ch <- 2
    close(ch)  // Signal no more values
}()

for value := range ch {
    fmt.Println(value)
}
` + "```" + `

Important:
- Only the sender should close (receiver closing is an error)
- Receiving from a closed channel returns zero value
- Sending on a closed channel panics

## Channel Directions

Restrict channels to send-only or receive-only:

` + "```go" + `
func Send(ch chan<- int) {
    ch <- 1  // OK
    val := <-ch  // ERROR
}

func Receive(ch <-chan int) {
    val := <-ch  // OK
    ch <- 1  // ERROR
}

func Bidirectional(ch chan int) {
    ch <- 1   // OK
    val := <-ch  // OK
}
` + "```" + `

Channels are automatically converted to directional types:

` + "```go" + `
ch := make(chan int)
Send(ch)      // Automatically chan int -> chan<- int
Receive(ch)   // Automatically chan int -> <-chan int
` + "```" + `

## The Select Statement

Select waits on multiple channel operations:

` + "```go" + `
select {
case msg := <-channel1:
    fmt.Println("From channel1:", msg)
case msg := <-channel2:
    fmt.Println("From channel2:", msg)
case channel3 <- value:
    fmt.Println("Sent on channel3")
default:
    fmt.Println("No operation ready")
}
` + "```" + `

### Select Behavior

- Waits until one case is ready
- If multiple cases are ready, chooses randomly
- default case executes if no channel is ready (non-blocking)
- Select in loop processes multiple operations

### Timeouts with Select

` + "```go" + `
select {
case msg := <-channel:
    fmt.Println(msg)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout!")
}
` + "```" + `

## Common Patterns

### Pattern 1: Worker Pool

` + "```go" + `
func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int)

    // Start workers
    for i := 0; i < 3; i++ {
        go worker(jobs, results)
    }

    // Send jobs
    go func() {
        for i := 1; i <= 9; i++ {
            jobs <- i
        }
        close(jobs)
    }()

    // Collect results
    for i := 0; i < 9; i++ {
        fmt.Println(<-results)
    }
}
` + "```" + `

### Pattern 2: Pipeline

` + "```go" + `
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    nums := generate(2, 3, 4)
    squares := square(nums)
    for result := range squares {
        fmt.Println(result)
    }
}
` + "```" + `

### Pattern 3: Fan-Out, Fan-In

Multiple goroutines process from one input, results combine:

` + "```go" + `
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            for n := range c {
                out <- n
            }
            wg.Done()
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
` + "```" + `

### Pattern 4: Cancellation

Use a channel to signal cancellation:

` + "```go" + `
func worker(jobs <-chan int, cancel <-chan struct{}) {
    for {
        select {
        case job := <-jobs:
            processJob(job)
        case <-cancel:
            return  // Stop when cancelled
        }
    }
}
` + "```" + `

## Race Conditions

Go includes a race detector to find concurrent access to shared data:

` + "```go" + `
go run -race main.go
go test -race ./...
` + "```" + `

### Common Race Condition

` + "```go" + `
var count int

go func() {
    count++  // Race!
}()

go func() {
    count++  // Race!
}()
` + "```" + `

Multiple goroutines accessing same variable without synchronization.

### Solution: Use Channels

` + "```go" + `
counter := make(chan int)

go func() {
    for {
        counter <- (<-counter) + 1
    }
}()

// Safe access through channel
counter <- 0
fmt.Println(<-counter)  // 1
` + "```" + `

Or use sync/atomic, sync.Mutex (covered later).

## Common Gotchas

### 1. Goroutine Leaks

` + "```go" + `
// Bad: Goroutine never exits
ch := make(chan int)
go func() {
    value := <-ch  // Waits forever
}()
// ch is never used

// Good: Provide way to stop
ch := make(chan int)
cancel := make(chan struct{})
go func() {
    select {
    case value := <-ch:
        fmt.Println(value)
    case <-cancel:
        return
    }
}()
close(cancel)
` + "```" + `

### 2. Sending on Closed Channel

` + "```go" + `
ch := make(chan int)
close(ch)
ch <- 1  // Panic!
` + "```" + `

Only sender should close channel.

### 3. Forgetting to Close Channel

` + "```go" + `
// Bad: Receiver waits forever
ch := make(chan int)
go func() {
    for range ch {  // Never completes if ch not closed
        // process
    }
}()

// Good: Close channel
ch := make(chan int)
go func() {
    ch <- 1
    close(ch)
}()
for range ch {
    // process
}
` + "```" + `

### 4. Deadlock

` + "```go" + `
// Deadlock: nobody receives
ch := make(chan int)
ch <- 1

// All goroutines blocked
` + "```" + `

Solution: receive or use goroutine.

## Summary

- Goroutines are lightweight threads
- Channels enable safe communication
- Unbuffered channels synchronize operations
- Buffered channels decouple send and receive
- Select waits on multiple operations
- Range iterates over channel values
- Close signals completion
- Race conditions require careful handling
- Common patterns: workers, pipelines, fan-out/fan-in
`,
		CodeExample: `package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker pattern
func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond)
		results <- job * 2
	}
}

// Pipeline pattern
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// Fan-in pattern
func merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("=== Worker Pattern ===")
	jobs := make(chan int, 100)
	results := make(chan int)

	for i := 1; i <= 3; i++ {
		go worker(i, jobs, results)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	for i := 0; i < 5; i++ {
		fmt.Printf("Result: %d\n", <-results)
	}

	fmt.Println("\n=== Pipeline Pattern ===")
	nums := generate(2, 3, 4)
	squares := square(nums)
	for result := range squares {
		fmt.Printf("Square: %d\n", result)
	}

	fmt.Println("\n=== Select Pattern ===")
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		ch1 <- "one"
	}()

	go func() {
		ch2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("From ch1: %s\n", msg)
		case msg := <-ch2:
			fmt.Printf("From ch2: %s\n", msg)
		}
	}
}
`,
		Solution: `package main

import (
	"fmt"
	"sync"
	"time"
)

// Datastore simulates a database
type Datastore struct {
	mu    sync.Mutex
	data  map[int]string
}

func (ds *Datastore) Get(id int) string {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	return ds.data[id]
}

// Request represents a data request
type Request struct {
	ID       int
	Response chan<- string
}

// RequestHandler processes requests concurrently
func RequestHandler(ds *Datastore, requests <-chan Request) {
	for req := range requests {
		// Simulate async work
		go func(r Request) {
			result := ds.Get(r.ID)
			r.Response <- result
		}(req)
	}
}

// CrawlUrls simulates concurrent web crawling
func CrawlUrls(urls []string) <-chan string {
	results := make(chan string)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			// Simulate network request
			time.Sleep(time.Millisecond * 100)
			results <- fmt.Sprintf("Data from %s", u)
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// Multiplexer combines multiple channels
func Multiplexer(ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for {
			select {
			case val := <-ch1:
				out <- val
			case val := <-ch2:
				out <- val
			}
		}
	}()

	return out
}

// Timeout example
func FetchWithTimeout(url string, timeout time.Duration) (string, error) {
	result := make(chan string)

	go func() {
		// Simulate network request
		time.Sleep(200 * time.Millisecond)
		result <- fmt.Sprintf("Data from %s", url)
	}()

	select {
	case res := <-result:
		return res, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout fetching %s", url)
	}
}

// RateLimiter limits concurrent operations
func RateLimiter(urls []string, maxConcurrent int) []string {
	semaphore := make(chan struct{}, maxConcurrent)
	results := make(chan string, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			// Simulate work
			time.Sleep(100 * time.Millisecond)
			results <- fmt.Sprintf("Processed %s", u)
		}(url)
	}

	wg.Wait()
	close(results)

	var output []string
	for result := range results {
		output = append(output, result)
	}
	return output
}

func main() {
	fmt.Println("=== Request Handler ===")
	ds := &Datastore{
		data: map[int]string{
			1: "Alice",
			2: "Bob",
			3: "Charlie",
		},
	}

	requests := make(chan Request, 3)
	go RequestHandler(ds, requests)

	for i := 1; i <= 3; i++ {
		response := make(chan string, 1)
		requests <- Request{ID: i, Response: response}
		fmt.Printf("ID %d: %s\n", i, <-response)
	}
	close(requests)

	fmt.Println("\n=== Web Crawling ===")
	urls := []string{"http://example.com", "http://golang.org", "http://github.com"}
	for result := range CrawlUrls(urls) {
		fmt.Println(result)
	}

	fmt.Println("\n=== Timeout Example ===")
	result, err := FetchWithTimeout("http://api.example.com", 300*time.Millisecond)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(result)
	}

	fmt.Println("\n=== Rate Limiting ===")
	testUrls := []string{"url1", "url2", "url3", "url4", "url5"}
	results := RateLimiter(testUrls, 2)
	for _, r := range results {
		fmt.Println(r)
	}
}
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "simple-goroutine",
				Title:       "Simple Goroutines",
				Description: "Create and coordinate goroutines.",
				Requirements: []string{
					"Start multiple goroutines",
					"Use channel to coordinate completion",
					"Wait for all goroutines to finish",
				},
				InitialCode: `package main

import (
	"fmt"
)

func main() {
	// TODO: Create a channel
	// TODO: Start 3 goroutines that send on the channel
	// TODO: Receive from channel 3 times
}
`,
				Solution: `package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)

	for i := 1; i <= 3; i++ {
		go func(n int) {
			ch <- fmt.Sprintf("Hello from goroutine %d", n)
		}(i)
	}

	for i := 0; i < 3; i++ {
		fmt.Println(<-ch)
	}
}
`,
			},
			{
				ID:          "unbuffered-channels",
				Title:       "Unbuffered Channels",
				Description: "Understand unbuffered channel synchronization.",
				Requirements: []string{
					"Create unbuffered channel",
					"Demonstrate blocking behavior",
					"Send and receive values",
				},
				InitialCode: `package main

import (
	"fmt"
)

func main() {
	// TODO: Create unbuffered channel
	// TODO: Send value from goroutine
	// TODO: Receive in main
}
`,
				Solution: `package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	go func() {
		fmt.Println("Goroutine: sending 42")
		ch <- 42
		fmt.Println("Goroutine: sent 42")
	}()

	fmt.Println("Main: waiting...")
	value := <-ch
	fmt.Printf("Main: received %d\n", value)
}
`,
			},
			{
				ID:          "buffered-channels",
				Title:       "Buffered Channels",
				Description: "Learn buffered channel capacity.",
				Requirements: []string{
					"Create buffered channel with capacity",
					"Send multiple values",
					"Understand buffer exhaustion",
				},
				InitialCode: `package main

import (
	"fmt"
)

func main() {
	// TODO: Create buffered channel with capacity 2
	// TODO: Send 2 values without goroutine
	// TODO: Receive both values
}
`,
				Solution: `package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 2)

	fmt.Println("Sending 1")
	ch <- 1
	fmt.Println("Sending 2")
	ch <- 2

	fmt.Println("Receiving...")
	fmt.Println("Got:", <-ch)
	fmt.Println("Got:", <-ch)
}
`,
			},
			{
				ID:          "channel-range",
				Title:       "Ranging Over Channels",
				Description: "Iterate over channel values until closed.",
				Requirements: []string{
					"Create generator function",
					"Close channel from sender",
					"Range over channel in receiver",
				},
				InitialCode: `package main

import (
	"fmt"
)

// TODO: Create generate function that returns <-chan int
// TODO: Send 1, 2, 3 and close

func main() {
	// TODO: Range over result of generate
}
`,
				Solution: `package main

import (
	"fmt"
)

func generate() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func main() {
	for value := range generate() {
		fmt.Println(value)
	}
}
`,
			},
			{
				ID:          "select-statement",
				Title:       "Select Statement",
				Description: "Wait on multiple channel operations.",
				Requirements: []string{
					"Create two channels",
					"Use select to wait for both",
					"Handle different channel cases",
				},
				InitialCode: `package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		ch1 <- "one"
	}()

	go func() {
		ch2 <- "two"
	}()

	// TODO: Use select to receive from either channel
}
`,
				Solution: `package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "one"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("From ch1:", msg)
		case msg := <-ch2:
			fmt.Println("From ch2:", msg)
		}
	}
}
`,
			},
			{
				ID:          "timeout-pattern",
				Title:       "Timeout Pattern",
				Description: "Implement timeout using select.",
				Requirements: []string{
					"Create channel operation",
					"Use time.After for timeout",
					"Handle timeout case",
				},
				InitialCode: `package main

import (
	"fmt"
	"time"
)

func SlowOperation() <-chan string {
	ch := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "done"
	}()
	return ch
}

func main() {
	// TODO: Use select with time.After(1 second)
	// TODO: Handle timeout
}
`,
				Solution: `package main

import (
	"fmt"
	"time"
)

func SlowOperation() <-chan string {
	ch := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "done"
	}()
	return ch
}

func main() {
	result := SlowOperation()

	select {
	case msg := <-result:
		fmt.Println(msg)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout!")
	}
}
`,
			},
			{
				ID:          "worker-pool",
				Title:       "Worker Pool Pattern",
				Description: "Implement concurrent worker pool.",
				Requirements: []string{
					"Create job and result channels",
					"Implement worker function",
					"Start multiple workers",
					"Distribute jobs and collect results",
				},
				InitialCode: `package main

import (
	"fmt"
)

// TODO: Define worker function
// Receives jobs from channel
// Sends results to results channel

func main() {
	jobs := make(chan int, 5)
	results := make(chan int)

	// TODO: Start 2 workers

	// TODO: Send jobs 1-5
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	// TODO: Close jobs channel

	// TODO: Collect and print results
}
`,
				Solution: `package main

import (
	"fmt"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		results <- job * 2
	}
}

func main() {
	jobs := make(chan int, 5)
	results := make(chan int)

	for i := 1; i <= 2; i++ {
		go worker(i, jobs, results)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	for i := 0; i < 5; i++ {
		fmt.Printf("Result: %d\n", <-results)
	}
}
`,
			},
			{
				ID:          "channel-directions",
				Title:       "Channel Directions",
				Description: "Use send-only and receive-only channels.",
				Requirements: []string{
					"Create send-only channel parameter",
					"Create receive-only channel parameter",
					"Use in different functions",
				},
				InitialCode: `package main

import (
	"fmt"
)

// TODO: Create Sender function with send-only channel

// TODO: Create Receiver function with receive-only channel

func main() {
	// TODO: Create bidirectional channel
	// TODO: Pass to Sender and Receiver
}
`,
				Solution: `package main

import (
	"fmt"
)

func Sender(ch chan<- string) {
	ch <- "message"
}

func Receiver(ch <-chan string) {
	msg := <-ch
	fmt.Println("Received:", msg)
}

func main() {
	ch := make(chan string)

	go Sender(ch)
	Receiver(ch)
}
`,
			},
		},
		NextLessonID: func() *int { i := 10; return &i }(),
		PrevLessonID: func() *int { i := 8; return &i }(),
	}
}

// Lesson 10: Packages and Modules
func (s *curriculumService) getComprehensiveLessonData10() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          10,
		Title:       "Packages and Modules",
		Description: "Master package organization, Go modules, dependency management, and best practices for structuring Go projects.",
		Duration:    "5-6 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Intermediate",
		Objectives: []string{
			"Understand package organization and naming",
			"Master exported vs unexported identifiers",
			"Create and use Go modules",
			"Manage dependencies with go.mod and go.sum",
			"Understand semantic versioning",
			"Use internal packages",
			"Initialize and use init functions",
			"Structure larger projects effectively",
		},
		Theory: `# Packages and Modules

## Packages

A package is a directory containing Go source files. All files in a directory must belong to the same package.

### Package Declaration

Each Go file starts with a package declaration:

` + "```go" + `
package main        // Executable package
package math        // Library package
package myapp       // Library package
` + "```" + `

The special package ` + "`main`" + ` is for executable programs. Other packages are libraries.

### Package Structure

` + "```go" + `
myproject/
├── main.go           // package main
├── math/
│   └── calc.go       // package math
└── utils/
    └── helper.go     // package utils
` + "```" + `

To use functions from other packages:

` + "```go" + `
package main

import (
    "myproject/math"
    "myproject/utils"
)

func main() {
    result := math.Add(2, 3)     // Use exported function
    helper := utils.Helper()
}
` + "```" + `

## Exported vs Unexported

Go visibility is controlled by capitalization:

- **Uppercase (Exported)**: Accessible from other packages
- **lowercase (Unexported)**: Only accessible within the same package

` + "```go" + `
package math

// Exported - can be used from other packages
func Add(a, b int) int {
    return a + b
}

// Unexported - only used within this package
func multiply(a, b int) int {
    return a * b
}

// Exported struct
type Calculator struct {
    Value int
}

// Unexported field (private)
type Config struct {
    apiKey string  // Not accessible from other packages
    Debug  bool    // Accessible from other packages
}
` + "```" + `

### Best Practices for Exported Names

1. **Use descriptive names**: ` + "`Reader`" + ` not ` + "`R`" + `
2. **Avoid redundancy**: In package ` + "`http`" + `, use ` + "`Server`" + ` not ` + "`HTTPServer`" + `
3. **Be consistent**: Parallel types should follow same pattern
4. **Document public API**: All exported items need comment

` + "```go" + `
package httpserver

// Server handles HTTP requests.
type Server struct {
    // ...
}

// Start begins listening for connections.
func (s *Server) Start() error {
    // ...
}
` + "```" + `

## Go Modules

A module is a collection of related packages with a ` + "`go.mod`" + ` file.

### Creating a Module

` + "```bash" + `
go mod init github.com/username/mymodule
` + "```" + `

This creates ` + "`go.mod`" + `:

` + "```" + `
module github.com/username/mymodule

go 1.23
` + "```" + `

### go.mod File

` + "```" + `
module github.com/mycompany/myapp

go 1.23

require (
    github.com/lib/pq v1.10.7
    github.com/golang/protobuf v1.5.2
)

require (
    github.com/stretchr/testify v1.7.0 // indirect
)

exclude github.com/bad/package v1.0.0

retract v1.0.0
` + "```" + `

- ` + "`require`" + `: Direct dependencies
- ` + "`indirect`" + `: Transitive dependencies
- ` + "`exclude`" + `: Prevent specific versions
- ` + "`retract`" + `: Remove buggy versions

### go.sum File

` + "`go.sum`" + ` contains checksums for verification:

` + "```" + `
github.com/lib/pq v1.10.7 h1:p7ZhMD+KsSRozJr34sdlHHk0SFG/a5IYzky9+5oBkc=
github.com/lib/pq v1.10.7/go.mod h1:AlVN5x4e4...
` + "```" + `

Always commit ` + "`go.sum`" + ` to version control for reproducible builds.

## Dependency Management

### Adding Dependencies

` + "```bash" + `
go get github.com/lib/pq
go get -u github.com/lib/pq          # Update to latest compatible
go get github.com/lib/pq@v1.10.7     # Specific version
` + "```" + `

### Removing Unused Dependencies

` + "```bash" + `
go mod tidy
` + "```" + `

### Viewing Dependencies

` + "```bash" + `
go list -m all              # List all modules
go list -m -u all          # Show available updates
` + "```" + `

## Semantic Versioning

Go modules follow semantic versioning: ` + "`MAJOR.MINOR.PATCH`" + `

- **MAJOR**: Breaking API changes (v1, v2, v3...)
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

` + "```" + `
v0.1.0  -> v0.1.1   Bug fix
v0.1.1  -> v0.2.0   New feature
v1.0.0  -> v2.0.0   Breaking change
` + "```" + `

### v2+ Modules

Breaking changes require new major version. Go imports them differently:

` + "```go" + `
import "github.com/user/mylib/v2"
` + "```" + `

## Package Organization

### Basic Structure

` + "```" + `
myproject/
├── go.mod
├── go.sum
├── main.go
├── config/
│   └── config.go
├── database/
│   └── db.go
└── api/
    └── server.go
` + "```" + `

### Internal Packages

Use ` + "`internal/`" + ` directory for packages not meant for external use:

` + "```" + `
myproject/
├── public/
│   └── api/
│       └── handler.go      # Exported API
└── internal/
    ├── auth/
    │   └── auth.go         # Not for external use
    └── database/
        └── db.go           # Not for external use
` + "```" + `

The ` + "`internal/`" + ` directory prevents accidental imports by other packages.

### Larger Projects

` + "```" + `
company/app/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── cli/
│       └── main.go
├── internal/
│   ├── api/
│   ├── database/
│   └── config/
├── pkg/
│   ├── metrics/
│   └── logging/
└── go.mod
` + "```" + `

- ` + "`cmd/`" + `: Executable entry points
- ` + "`internal/`" + `: Private application code
- ` + "`pkg/`" + `: Potentially reusable libraries

## Init Functions

Package init functions run automatically when package loads:

` + "```go" + `
package mypackage

import "fmt"

func init() {
    fmt.Println("Package initialized")
}

func init() {
    fmt.Println("Second init")
}

func main() {
    // init functions already ran
}
` + "```" + `

### Use Cases

1. **Setup**: Initialize databases, load configuration
2. **Validation**: Verify environment, check requirements
3. **Registration**: Register handlers, plugins, implementations

` + "```go" + `
package database

import "database/sql"

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", "...")
    if err != nil {
        panic(err)
    }
}

func GetDB() *sql.DB {
    return db
}
` + "```" + `

## Import Best Practices

### 1. Use Absolute Imports

` + "```go" + `
// Good
import "github.com/company/myapp/utils"

// Bad - relative imports not supported
import "./utils"
` + "```" + `

### 2. Group Imports

` + "```go" + `
import (
    // Standard library
    "fmt"
    "os"

    // Third-party
    "github.com/lib/pq"
    "github.com/stretchr/testify/assert"

    // Internal
    "myapp/utils"
    "myapp/database"
)
` + "```" + `

### 3. Use Aliases When Needed

` + "```go" + `
import (
    httpsrv "company.com/services/http"
    grpcsrv "company.com/services/grpc"
)

// Now can use httpsrv.Server and grpcsrv.Server
` + "```" + `

### 4. Avoid Wildcard Imports

` + "```go" + `
// Bad
import . "fmt"
fmt.Println()  // fmt. is implied

// Good
import "fmt"
fmt.Println()  // Clear where Println comes from
` + "```" + `

## Publishing Packages

When publishing to GitHub:

1. **Use descriptive README**: Explain purpose and usage
2. **Provide examples**: Show common use cases
3. **Document exports**: Add comments to all public items
4. **Use meaningful tags**: Release as versions (v1.0.0, v1.0.1, etc.)

` + "```bash" + `
git tag v1.0.0
git push origin v1.0.0
` + "```" + `

## Common Mistakes

### 1. Circular Imports

` + "```go" + `
// Package A imports B, B imports A = ERROR
// Solution: Restructure or use interfaces
` + "```" + `

### 2. Exposing Internal Details

` + "```go" + `
// Bad: Implementation leaks
type userStore struct { ... }
func (s *userStore) Find(...) { ... }

// Good: Interface-based
type UserStore interface {
    Find(...) (User, error)
}
` + "```" + `

### 3. Versioning Issues

` + "```go" + `
// Breaking change without version bump - breaks users' code
func OldFunction(a string) { ... }
func OldFunction(a string, b string) { ... }  // Breaking!

// Solution: Use v2.0.0
` + "```" + `

## Summary

- Packages organize code into reusable units
- Capitalization controls visibility (exported/unexported)
- Modules manage versioning and dependencies
- go.mod and go.sum handle dependency management
- Semantic versioning guides breaking changes
- internal/ directory prevents unwanted exports
- init() functions initialize packages
- Proper project structure enables scaling
- Good documentation makes packages usable
`,
		CodeExample: `package main

import (
	"fmt"

	"go-pro/internal/math"
	"go-pro/internal/utils"
)

func main() {
	fmt.Println("=== Using Packages ===")

	// Use exported functions from math package
	result := math.Add(5, 3)
	fmt.Printf("5 + 3 = %d\n", result)

	result = math.Multiply(4, 5)
	fmt.Printf("4 * 5 = %d\n", result)

	// Use exported functions from utils package
	message := utils.FormatMessage("Hello", "World")
	fmt.Println(message)

	// Use exported types
	calc := &math.Calculator{Value: 10}
	fmt.Printf("Calculator value: %d\n", calc.Value)
}

// Note: This example assumes the following structure:
// go-pro/
// ├── main.go
// ├── internal/
// │   ├── math/
// │   │   └── calc.go
// │   └── utils/
// │       └── helper.go
// └── go.mod
`,
		Solution: `package main

import (
	"fmt"
	"log"
)

// Database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Logger wrapper
type Logger struct {
	verbose bool
}

// Global instances initialized in init()
var (
	dbConfig *DatabaseConfig
	appLogger *Logger
)

// init() runs before main()
func init() {
	fmt.Println("Initializing application...")

	// Initialize database configuration
	dbConfig = &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "pass",
	}

	// Initialize logger
	appLogger = &Logger{verbose: true}

	fmt.Println("Application initialized successfully")
}

// GetDatabaseConfig returns the database configuration
func GetDatabaseConfig() *DatabaseConfig {
	return dbConfig
}

// GetLogger returns the application logger
func GetLogger() *Logger {
	return appLogger
}

// Log logs a message
func (l *Logger) Log(msg string) {
	if l.verbose {
		log.Println("[INFO]", msg)
	}
}

// Logf logs a formatted message
func (l *Logger) Logf(format string, args ...interface{}) {
	if l.verbose {
		log.Printf("[INFO] "+format, args...)
	}
}

// Module represents a loadable module
type Module interface {
	Name() string
	Initialize() error
	Shutdown() error
}

// AuthModule is a sample module
type AuthModule struct{}

func (m *AuthModule) Name() string {
	return "auth"
}

func (m *AuthModule) Initialize() error {
	fmt.Println("Initializing auth module")
	return nil
}

func (m *AuthModule) Shutdown() error {
	fmt.Println("Shutting down auth module")
	return nil
}

// ModuleRegistry manages modules
type ModuleRegistry struct {
	modules map[string]Module
}

// NewModuleRegistry creates new registry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]Module),
	}
}

// Register adds a module
func (r *ModuleRegistry) Register(m Module) error {
	GetLogger().Logf("Registering module: %s", m.Name())
	r.modules[m.Name()] = m
	return m.Initialize()
}

// Get retrieves a module
func (r *ModuleRegistry) Get(name string) (Module, bool) {
	m, ok := r.modules[name]
	return m, ok
}

func main() {
	fmt.Println("\n=== Module System ===")

	// Create module registry
	registry := NewModuleRegistry()

	// Register modules
	registry.Register(&AuthModule{})

	// Use database config
	config := GetDatabaseConfig()
	GetLogger().Logf("Database: %s:%d", config.Host, config.Port)

	// Check if auth module exists
	if module, ok := registry.Get("auth"); ok {
		GetLogger().Logf("Found module: %s", module.Name())
	}

	fmt.Println("\nApplication running...")
}
`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "basic-package",
				Title:       "Creating a Basic Package",
				Description: "Create a package with exported and unexported functions.",
				Requirements: []string{
					"Create math package with exported functions",
					"Add unexported helper function",
					"Use from main package",
					"Follow naming conventions",
				},
				InitialCode: `// File: main.go
package main

import (
	"fmt"
	"myapp/math"
)

func main() {
	// TODO: Use exported functions from math package
	fmt.Println("Using math package")
}

// File: math/calc.go
// TODO: Create math package
// TODO: Implement exported Add function
// TODO: Implement exported Multiply function
// TODO: Implement unexported helper
`,
				Solution: `// File: main.go
package main

import (
	"fmt"
	"myapp/math"
)

func main() {
	result1 := math.Add(10, 5)
	fmt.Printf("Add: 10 + 5 = %d\n", result1)

	result2 := math.Multiply(10, 5)
	fmt.Printf("Multiply: 10 * 5 = %d\n", result2)

	result3 := math.Power(2, 8)
	fmt.Printf("Power: 2^8 = %d\n", result3)
}

// File: math/calc.go
package math

func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func Power(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result = multiply(result, base)
	}
	return result
}

// unexported helper function
func multiply(a, b int) int {
	return a * b
}
`,
			},
			{
				ID:          "exported-unexported",
				Title:       "Controlling Visibility",
				Description: "Use capitalization to control what's exported.",
				Requirements: []string{
					"Create struct with exported and unexported fields",
					"Create exported methods",
					"Create unexported helper methods",
					"Demonstrate visibility rules",
				},
				InitialCode: `package database

// TODO: Define User struct with:
// - Exported Name field
// - Unexported password field

// TODO: Implement exported GetName method

// TODO: Implement unexported validate method

// TODO: Implement exported SetPassword with validation
`,
				Solution: `package database

import (
	"fmt"
)

// User represents an application user
type User struct {
	Name     string
	password string // unexported - private
	email    string // unexported - private
}

// GetName returns the user's name
func (u *User) GetName() string {
	return u.Name
}

// SetPassword sets and validates the password
func (u *User) SetPassword(pwd string) error {
	if !u.validate(pwd) {
		return fmt.Errorf("password too weak")
	}
	u.password = pwd
	return nil
}

// validate checks if password meets requirements
// unexported - only used internally
func (u *User) validate(pwd string) bool {
	return len(pwd) >= 8
}

// GetPassword is not exported - password is private
// This prevents external packages from reading it
`,
			},
			{
				ID:          "module-structure",
				Title:       "Module Organization",
				Description: "Create a well-organized module structure.",
				Requirements: []string{
					"Create main.go in cmd/app/",
					"Create internal packages",
					"Create public packages",
					"Organize files logically",
				},
				InitialCode: `// Directory structure to create:
// myapp/
// ├── go.mod
// ├── cmd/
// │   └── app/
// │       └── main.go
// ├── internal/
// │   ├── database/
// │   │   └── db.go
// │   └── auth/
// │       └── auth.go
// └── pkg/
//     └── utils/
//         └── helpers.go

// Files to implement:
// cmd/app/main.go
package main

func main() {
	// TODO: Use internal packages
}

// internal/database/db.go
package database

// TODO: Implement database package

// internal/auth/auth.go
package auth

// TODO: Implement auth package

// pkg/utils/helpers.go
package utils

// TODO: Implement utility package
`,
				Solution: `// cmd/app/main.go
package main

import (
	"fmt"
	"myapp/internal/auth"
	"myapp/internal/database"
	"myapp/pkg/utils"
)

func main() {
	fmt.Println("Starting application...")

	// Use database
	db := database.NewDB("localhost")
	fmt.Printf("Database: %s\n", db.Name())

	// Use auth
	token := auth.GenerateToken("user@example.com")
	fmt.Printf("Token: %s\n", token)

	// Use utils
	hash := utils.HashString("password")
	fmt.Printf("Hash: %s\n", hash)
}

// internal/database/db.go
package database

type Database struct {
	name string
}

func NewDB(name string) *Database {
	return &Database{name: name}
}

func (db *Database) Name() string {
	return db.name
}

// internal/auth/auth.go
package auth

import (
	"fmt"
)

func GenerateToken(email string) string {
	return fmt.Sprintf("token_%s", email)
}

// pkg/utils/helpers.go
package utils

import (
	"crypto/sha256"
	"fmt"
)

func HashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}
`,
			},
			{
				ID:          "init-functions",
				Title:       "Package Initialization",
				Description: "Use init() for package initialization.",
				Requirements: []string{
					"Create package with init() function",
					"Initialize global state",
					"Handle initialization errors",
					"Use multiple init functions",
				},
				InitialCode: `package config

import "log"

// Global configuration variables
var (
	AppName string
	Version string
	Debug   bool
)

// TODO: Implement init() to load configuration
// TODO: Handle initialization errors
`,
				Solution: `package config

import (
	"log"
	"os"
)

// Global configuration variables
var (
	AppName string
	Version string
	Debug   bool
)

// First init function
func init() {
	log.Println("Initializing configuration...")

	// Set defaults
	AppName = "MyApp"
	Version = "1.0.0"
}

// Second init function
func init() {
	// Load from environment
	if debug := os.Getenv("DEBUG"); debug == "true" {
		Debug = true
	}

	if name := os.Getenv("APP_NAME"); name != "" {
		AppName = name
	}

	log.Printf("Configuration loaded: %s v%s (Debug: %v)\n",
		AppName, Version, Debug)
}

// GetConfig returns current configuration
func GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":    AppName,
		"version": Version,
		"debug":   Debug,
	}
}
`,
			},
			{
				ID:          "internal-packages",
				Title:       "Using Internal Packages",
				Description: "Create and use internal packages for private code.",
				Requirements: []string{
					"Create internal/handlers package",
					"Create public API",
					"Prevent external imports of internal code",
				},
				InitialCode: `// File structure:
// myapp/
// ├── main.go
// ├── api.go (public API)
// └── internal/
//     └── handlers/
//         └── handlers.go

// main.go
package main

func main() {
	// TODO: Use public API
}

// api.go
// TODO: Define public API

// internal/handlers/handlers.go
package handlers

// TODO: Define internal handlers
`,
				Solution: `// main.go
package main

import (
	"fmt"
	"myapp/api"
)

func main() {
	server := api.NewServer()
	fmt.Println(server.Start())
}

// api.go
package main

import (
	"fmt"
	"myapp/internal/handlers"
)

// Server represents the application server
type Server struct {
	handler handlers.RequestHandler
}

// NewServer creates a new server
func NewServer() *Server {
	return &Server{
		handler: handlers.New(),
	}
}

// Start starts the server
func (s *Server) Start() string {
	return s.handler.Handle("request")
}

// internal/handlers/handlers.go
package handlers

import "fmt"

// RequestHandler handles requests
type RequestHandler struct{}

// New creates new handler
func New() RequestHandler {
	return RequestHandler{}
}

// Handle handles a request
func (h RequestHandler) Handle(req string) string {
	return fmt.Sprintf("Handled: %s", req)
}
`,
			},
		},
		NextLessonID: func() *int { i := 11; return &i }(),
		PrevLessonID: func() *int { i := 9; return &i }(),
	}
}
