package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Pointers Demo")

	demo1BasicPointers()
	demo2PointerOperations()
	demo3PassByValue()
	demo4PassByReference()
	demo5PointerToStruct()
}

func demo1BasicPointers() {
	utils.PrintSubHeader("1. Basic Pointers")

	x := 42
	fmt.Printf("Value of x: %d\n", x)
	fmt.Printf("Address of x: %p\n", &x)

	// Pointer variable
	var ptr *int = &x
	fmt.Printf("\nPointer ptr points to: %p\n", ptr)
	fmt.Printf("Value at pointer: %d\n", *ptr)

	// Modify through pointer
	*ptr = 100
	fmt.Printf("\nAfter *ptr = 100:\n")
	fmt.Printf("Value of x: %d\n", x)
	fmt.Printf("Value at pointer: %d\n", *ptr)
}

func demo2PointerOperations() {
	utils.PrintSubHeader("2. Pointer Operations")

	a := 10
	b := 20

	ptrA := &a
	ptrB := &b

	fmt.Printf("a = %d, b = %d\n", a, b)
	fmt.Printf("ptrA points to: %p, ptrB points to: %p\n", ptrA, ptrB)

	// Swap values through pointers
	*ptrA, *ptrB = *ptrB, *ptrA

	fmt.Printf("\nAfter swap through pointers:\n")
	fmt.Printf("a = %d, b = %d\n", a, b)

	// Nil pointer
	var nilPtr *int
	fmt.Printf("\nNil pointer: %v\n", nilPtr)
	if nilPtr == nil {
		fmt.Println("Pointer is nil (not pointing to anything)")
	}
}

func demo3PassByValue() {
	utils.PrintSubHeader("3. Pass by Value (Copy)")

	x := 10
	fmt.Printf("Before function call: x = %d\n", x)

	tryToModifyValue(x)

	fmt.Printf("After function call: x = %d (unchanged)\n", x)
	fmt.Println("→ Function received a copy, original unchanged")
}

func tryToModifyValue(num int) {
	num = 100
	fmt.Printf("Inside function: num = %d\n", num)
}

func demo4PassByReference() {
	utils.PrintSubHeader("4. Pass by Reference (Pointer)")

	x := 10
	fmt.Printf("Before function call: x = %d\n", x)

	modifyThroughPointer(&x)

	fmt.Printf("After function call: x = %d (modified!)\n", x)
	fmt.Println("→ Function received pointer, original modified")

	// Swap example
	a, b := 5, 10
	fmt.Printf("\nBefore swap: a = %d, b = %d\n", a, b)
	swap(&a, &b)
	fmt.Printf("After swap: a = %d, b = %d\n", a, b)
}

func modifyThroughPointer(ptr *int) {
	*ptr = 100
	fmt.Printf("Inside function: *ptr = %d\n", *ptr)
}

func swap(x, y *int) {
	*x, *y = *y, *x
}

type Person struct {
	Name string
	Age  int
}

func demo5PointerToStruct() {
	utils.PrintSubHeader("5. Pointers to Structs")

	person := Person{Name: "Alice", Age: 30}
	fmt.Printf("Original: %+v\n", person)

	// Pointer to struct
	ptr := &person

	// Access fields through pointer (automatic dereferencing)
	ptr.Age = 31
	fmt.Printf("After ptr.Age = 31: %+v\n", person)

	// Explicit dereferencing
	(*ptr).Name = "Alice Smith"
	fmt.Printf("After (*ptr).Name = 'Alice Smith': %+v\n", person)

	// Function with pointer to struct
	updatePerson(&person, "Bob", 25)
	fmt.Printf("After updatePerson: %+v\n", person)
}

func updatePerson(p *Person, name string, age int) {
	p.Name = name
	p.Age = age
}
