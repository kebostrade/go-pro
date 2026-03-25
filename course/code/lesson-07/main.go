package main

import (
	"fmt"
	"lesson-07/exercises"
)

func main() {
	fmt.Println("=== Lesson 07: Interfaces and Polymorphism ===")
	fmt.Println()

	// Rectangle
	r := exercises.Rectangle{Width: 5, Height: 3}
	fmt.Printf("Rectangle: Width=%.1f, Height=%.1f\n", r.Width, r.Height)
	fmt.Printf("  Area: %.2f\n", r.Area())
	fmt.Printf("  Perimeter: %.2f\n", r.Perimeter())
	fmt.Println()

	// Circle
	c := exercises.Circle{Radius: 2}
	fmt.Printf("Circle: Radius=%.1f\n", c.Radius)
	fmt.Printf("  Area: %.2f\n", c.Area())
	fmt.Printf("  Perimeter: %.2f\n", c.Perimeter())
	fmt.Println()

	// Empty interface
	fmt.Println("Type checking:")
	fmt.Printf("  GetType(42) = %s\n", exercises.GetType(42))
	fmt.Printf("  GetType(\"hello\") = %s\n", exercises.GetType("hello"))
	fmt.Printf("  GetType(3.14) = %s\n", exercises.GetType(3.14))
	fmt.Println()

	// Type assertion
	fmt.Println("Type assertion:")
	val, ok := exercises.GetInt(42)
	fmt.Printf("  GetInt(42) = %d, ok=%t\n", val, ok)
	val2, ok2 := exercises.GetInt("hello")
	fmt.Printf("  GetInt(\"hello\") = %d, ok=%t\n", val2, ok2)
	fmt.Println()

	// Type switch
	fmt.Println("Type switch:")
	fmt.Printf("  Describe(42) = %s\n", exercises.Describe(42))
	fmt.Printf("  Describe(\"hello\") = %s\n", exercises.Describe("hello"))
	fmt.Printf("  Describe(3.14) = %s\n", exercises.Describe(3.14))
	fmt.Println()

	// Stringer
	p := exercises.Person{Name: "Alice", Age: 30}
	fmt.Printf("  Person: %s\n", p.String())
	fmt.Println()

	// Custom error
	err := exercises.ValidationError{Field: "email", Message: "invalid format"}
	fmt.Printf("  Error: %s\n", err.Error())
	fmt.Println()

	// PrintArea with interface
	fmt.Println("PrintArea with interface:")
	fmt.Printf("  PrintArea(Rectangle{4,5}) = %s\n", exercises.PrintArea(exercises.Rectangle{Width: 4, Height: 5}))
	fmt.Printf("  PrintArea(Circle{1}) = %s\n", exercises.PrintArea(exercises.Circle{Radius: 1}))
	fmt.Printf("  PrintArea(\"invalid\") = %s\n", exercises.PrintArea("invalid"))
	fmt.Println()

	fmt.Println("=== All exercises completed! ===")
}
