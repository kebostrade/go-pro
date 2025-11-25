package main

import (
	"fmt"
	"os"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Structs & Methods Demo")

	demo1BasicStructs()
	demo2StructMethods()
	demo3PointerReceivers()
	demo4EmbeddedStructs()
	demo5StructTags()
}

func demo1BasicStructs() {
	utils.PrintSubHeader("1. Basic Structs")

	type Person struct {
		Name string
		Age  int
		City string
	}

	// Different ways to create structs
	person1 := Person{Name: "Alice", Age: 30, City: "New York"}
	fmt.Printf("person1: %+v\n", person1)

	person2 := Person{"Bob", 25, "London"}
	fmt.Printf("person2: %+v\n", person2)

	var person3 Person
	person3.Name = "Charlie"
	person3.Age = 35
	person3.City = "Paris"
	fmt.Printf("person3: %+v\n", person3)

	// Accessing fields
	fmt.Printf("\n%s is %d years old and lives in %s\n",
		person1.Name, person1.Age, person1.City)
}

func demo2StructMethods() {
	utils.PrintSubHeader("2. Struct Methods")

	rect := Rectangle{Width: 10, Height: 5}

	fmt.Printf("Rectangle: %+v\n", rect)
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())
	fmt.Printf("Description: %s\n", rect.String())
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

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.1f x %.1f)", r.Width, r.Height)
}

func demo3PointerReceivers() {
	utils.PrintSubHeader("3. Pointer Receivers vs Value Receivers")

	counter := Counter{Count: 0}

	fmt.Printf("Initial: %+v\n", counter)

	// Value receiver - doesn't modify original
	counter.IncrementValue()
	fmt.Printf("After IncrementValue: %+v (unchanged)\n", counter)

	// Pointer receiver - modifies original
	counter.IncrementPointer()
	fmt.Printf("After IncrementPointer: %+v (changed!)\n", counter)

	counter.Add(5)
	fmt.Printf("After Add(5): %+v\n", counter)
}

type Counter struct {
	Count int
}

// Value receiver - receives a copy
func (c Counter) IncrementValue() {
	c.Count++
	fmt.Printf("  Inside IncrementValue: %d\n", c.Count)
}

// Pointer receiver - modifies original
func (c *Counter) IncrementPointer() {
	c.Count++
	fmt.Printf("  Inside IncrementPointer: %d\n", c.Count)
}

func (c *Counter) Add(n int) {
	c.Count += n
}

func demo4EmbeddedStructs() {
	utils.PrintSubHeader("4. Embedded Structs (Composition)")

	employee := Employee{
		Person2: Person2{
			Name: "Alice",
			Age:  30,
		},
		ID:         1001,
		Department: "Engineering",
		Salary:     75000,
	}

	fmt.Printf("Employee: %+v\n", employee)

	// Access embedded fields directly
	fmt.Printf("\nDirect access to embedded fields:\n")
	fmt.Printf("Name: %s\n", employee.Name)
	fmt.Printf("Age: %d\n", employee.Age)

	// Or through the embedded struct
	fmt.Printf("\nAccess through embedded struct:\n")
	fmt.Printf("Name: %s\n", employee.Person2.Name)
	fmt.Printf("Age: %d\n", employee.Person2.Age)

	// Call embedded methods
	employee.Introduce()
}

type Person2 struct {
	Name string
	Age  int
}

func (p Person2) Introduce() {
	fmt.Printf("Hi, I'm %s and I'm %d years old\n", p.Name, p.Age)
}

type Employee struct {
	Person2
	ID         int
	Department string
	Salary     float64
}

func demo5StructTags() {
	utils.PrintSubHeader("5. Struct Tags")
	// Example password for demonstration only - never hardcode in production
	// NOTE: This is test code only, not real credentials
	type User struct {
		ID       int    `json:"id" db:"user_id"`
		Username string `json:"username" db:"username"`
		Email    string `json:"email" db:"email_address"`
		Password string `json:"-" db:"password_hash"` // "-" means omit from JSON
	}

	// Note: In production, passwords should never be hardcoded.
	// For this example, we use a default for demonstration.
	// Use environment variable or load from config in production
	pwd := os.Getenv("EXAMPLE_PASSWORD")
	if pwd == "" {
		pwd = "example_password" // Test only
	}
	user := User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: pwd,
	}

	fmt.Printf("User struct: %+v\n", user)
	fmt.Println("\nStruct tags are used for:")
	fmt.Println("  • JSON marshaling/unmarshaling")
	fmt.Println("  • Database column mapping")
	fmt.Println("  • Validation rules")
	fmt.Println("  • Custom serialization")
}
