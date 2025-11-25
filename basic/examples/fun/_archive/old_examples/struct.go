//go:build ignore

package main

import "fmt"

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
}

func main() {
	employee := Employee{
		ID:         1,
		Name:       "John Doe",
		Department: "IT",
		Salary:     5000.0,
	}

	fmt.Printf("Employee ID: %d\n", employee.ID)
	fmt.Printf("Employee Name: %s\n", employee.Name)
	fmt.Printf("Employee Department: %s\n", employee.Department)
	fmt.Printf("Employee Salary: %.2f\n", employee.Salary)
}
