//go:build ignore

package main

import "fmt"

func main() {
	foo := 42
	fmt.Println("Value of foo:", foo)
	// Pointers: Variables that store memory addresses of other variables
	pointerFoo := &foo
	fmt.Println("Address of foo:", pointerFoo)
	fmt.Println("Value at pointer:", *pointerFoo)

	// Modify value through pointer
	*pointerFoo = 100
	fmt.Println("Modified value of foo:", foo)
}
