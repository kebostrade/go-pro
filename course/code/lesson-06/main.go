package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Constants
const aliceJohnson = "Alice Johnson"

func main() {
	fmt.Println("=== Lesson 6: Structs and Methods ===")

	// Basic struct operations
	fmt.Println("1. Basic Struct Operations:")
	demonstrateBasicStructs()
	fmt.Println()

	// Methods with different receivers
	fmt.Println("2. Methods and Receivers:")
	demonstrateMethods()
	fmt.Println()

	// Struct embedding
	fmt.Println("3. Struct Embedding:")
	demonstrateEmbedding()
	fmt.Println()

	// Struct tags and JSON
	fmt.Println("4. Struct Tags and JSON:")
	demonstrateStructTags()
	fmt.Println()

	// Real-world example
	fmt.Println("5. Real-World Example - Employee Management:")
	demonstrateEmployeeSystem()
}

func demonstrateBasicStructs() {
	// Define and initialize structs
	type Person struct {
		Name   string
		Age    int
		Email  string
		Active bool
	}

	// Zero value
	var p1 Person
	fmt.Printf("Zero value: %+v\n", p1)

	// Struct literal with field names
	p2 := Person{
		Name:   "Alice",
		Age:    30,
		Email:  "alice@example.com",
		Active: true,
	}
	fmt.Printf("Named fields: %+v\n", p2)

	// Pointer to struct
	p3 := &Person{Name: "Bob", Age: 25}
	fmt.Printf("Pointer to struct: %+v\n", p3)

	// Accessing fields
	fmt.Printf("Name: %s, Age: %d\n", p2.Name, p2.Age)
	fmt.Printf("Pointer access: %s\n", p3.Name)
}

func demonstrateMethods() {
	type Rectangle struct {
		Width, Height float64
	}

	// Value receiver methods
	area := func(r Rectangle) float64 {
		return r.Width * r.Height
	}

	perimeter := func(r Rectangle) float64 {
		return 2 * (r.Width + r.Height)
	}

	// Pointer receiver method
	scale := func(r *Rectangle, factor float64) {
		r.Width *= factor
		r.Height *= factor
	}

	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Rectangle: %+v\n", rect)
	fmt.Printf("Area: %.2f\n", area(rect))
	fmt.Printf("Perimeter: %.2f\n", perimeter(rect))

	scale(&rect, 2)
	fmt.Printf("After scaling: %+v\n", rect)
	fmt.Printf("New area: %.2f\n", area(rect))
}

func demonstrateEmbedding() {
	type Address struct {
		Street  string
		City    string
		Country string
	}

	type Person struct {
		Name  string
		Age   int
		Email string
	}

	type Employee struct {
		Person     // Embedded struct
		Address    // Embedded struct
		EmployeeID string
		Department string
		Salary     float64
	}

	// Create employee with embedded fields
	emp := Employee{
		Person: Person{
			Name:  aliceJohnson,
			Age:   30,
			Email: "alice@company.com",
		},
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
		EmployeeID: "E001",
		Department: "Engineering",
		Salary:     75000,
	}

	fmt.Printf("Employee: %+v\n", emp)
	fmt.Printf("Direct access - Name: %s\n", emp.Name)
	fmt.Printf("Direct access - City: %s\n", emp.City)
	fmt.Printf("Explicit access - Person Name: %s\n", emp.Person.Name)
}

func demonstrateStructTags() {
	type User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"-"` // Excluded from JSON
		IsActive bool   `json:"is_active"`
	}

	// Note: In production, passwords should never be hardcoded.
	// For this example, we use a default for demonstration.
	defaultPassword := os.Getenv("EXAMPLE_PASSWORD")
	if defaultPassword == "" {
		defaultPassword = "example_password"
	}

	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: defaultPassword,
		IsActive: true,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}

	fmt.Printf("JSON: %s\n", string(jsonData))

	// Unmarshal from JSON
	jsonStr := `{"id":2,"name":"Jane Smith","email":"jane@example.com","is_active":false}`
	var user2 User
	err = json.Unmarshal([]byte(jsonStr), &user2)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return
	}

	fmt.Printf("Unmarshaled: %+v\n", user2)
}

func demonstrateEmployeeSystem() {
	// Employee management system using structs and methods
	type Department struct {
		Name      string
		Manager   string
		Budget    float64
		Employees []string
	}

	type Employee struct {
		ID         string
		Name       string
		Department string
		Salary     float64
		Active     bool
	}

	type Company struct {
		Name        string
		Departments map[string]*Department
		Employees   map[string]*Employee
	}

	// Create company
	company := &Company{
		Name:        "TechCorp",
		Departments: make(map[string]*Department),
		Employees:   make(map[string]*Employee),
	}

	// Add departments
	company.Departments["Engineering"] = &Department{
		Name:      "Engineering",
		Manager:   aliceJohnson,
		Budget:    500000,
		Employees: []string{},
	}

	company.Departments["Marketing"] = &Department{
		Name:      "Marketing",
		Manager:   "Bob Smith",
		Budget:    200000,
		Employees: []string{},
	}

	// Add employees
	employees := []*Employee{
		{ID: "E001", Name: aliceJohnson, Department: "Engineering", Salary: 90000, Active: true},
		{ID: "E002", Name: "Charlie Brown", Department: "Engineering", Salary: 75000, Active: true},
		{ID: "E003", Name: "Diana Prince", Department: "Marketing", Salary: 65000, Active: true},
	}

	for _, emp := range employees {
		company.Employees[emp.ID] = emp
		dept := company.Departments[emp.Department]
		dept.Employees = append(dept.Employees, emp.ID)
	}

	// Display company structure
	fmt.Printf("Company: %s\n", company.Name)
	fmt.Println("Departments:")
	for _, dept := range company.Departments {
		fmt.Printf("  %s (Manager: %s, Budget: $%.0f)\n",
			dept.Name, dept.Manager, dept.Budget)
		fmt.Printf("    Employees: %d\n", len(dept.Employees))

		totalSalary := 0.0
		for _, empID := range dept.Employees {
			emp := company.Employees[empID]
			fmt.Printf("      - %s (ID: %s, Salary: $%.0f)\n",
				emp.Name, emp.ID, emp.Salary)
			totalSalary += emp.Salary
		}
		fmt.Printf("    Total Salaries: $%.0f\n", totalSalary)
		fmt.Printf("    Budget Utilization: %.1f%%\n",
			(totalSalary/dept.Budget)*100)
	}
}

// Additional examples for method sets and interfaces
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

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}
