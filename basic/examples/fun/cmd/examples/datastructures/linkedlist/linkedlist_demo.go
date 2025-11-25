package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/datastructures"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Linked List Data Structure Demo")

	// Demo 1: Basic Operations
	demo1BasicOperations()

	// Demo 2: Insertion at Different Positions
	demo2Insertions()

	// Demo 3: Deletion Operations
	demo3Deletions()

	// Demo 4: Search and Access
	demo4SearchAndAccess()

	// Demo 5: List Manipulation
	demo5Manipulation()
}

func demo1BasicOperations() {
	utils.PrintSubHeader("1. Basic Linked List Operations")

	list := datastructures.NewLinkedList[int]()

	fmt.Println("Inserting at beginning: 3, 2, 1")
	list.InsertAtBeginning(3)
	list.InsertAtBeginning(2)
	list.InsertAtBeginning(1)
	fmt.Printf("List: %s\n", list.String())

	fmt.Println("\nInserting at end: 4, 5, 6")
	list.InsertAtEnd(4)
	list.InsertAtEnd(5)
	list.InsertAtEnd(6)
	fmt.Printf("List: %s\n", list.String())

	fmt.Printf("\nSize: %d\n", list.Size())
	fmt.Printf("Is empty: %v\n", list.IsEmpty())
}

func demo2Insertions() {
	utils.PrintSubHeader("2. Insertion at Different Positions")

	list := datastructures.NewLinkedList[string]()

	// Build initial list
	words := []string{"apple", "banana", "cherry"}
	for _, word := range words {
		list.InsertAtEnd(word)
	}
	fmt.Printf("Initial list: %s\n", list.String())

	// Insert at position 0 (beginning)
	list.InsertAtPosition("aardvark", 0)
	fmt.Printf("After inserting 'aardvark' at position 0: %s\n", list.String())

	// Insert at position 2 (middle)
	list.InsertAtPosition("avocado", 2)
	fmt.Printf("After inserting 'avocado' at position 2: %s\n", list.String())

	// Insert at end
	list.InsertAtPosition("date", list.Size())
	fmt.Printf("After inserting 'date' at end: %s\n", list.String())

	// Try invalid position
	if err := list.InsertAtPosition("invalid", 100); err != nil {
		fmt.Printf("\nError inserting at invalid position: %v\n", err)
	}
}

func demo3Deletions() {
	utils.PrintSubHeader("3. Deletion Operations")

	list := datastructures.NewLinkedList[int]()

	// Build list: 1 -> 2 -> 3 -> 4 -> 5
	for i := 1; i <= 5; i++ {
		list.InsertAtEnd(i)
	}
	fmt.Printf("Initial list: %s\n", list.String())

	// Delete at beginning
	if val, err := list.DeleteAtBeginning(); err == nil {
		fmt.Printf("Deleted from beginning: %d\n", val)
		fmt.Printf("List: %s\n", list.String())
	}

	// Delete at end
	if val, err := list.DeleteAtEnd(); err == nil {
		fmt.Printf("\nDeleted from end: %d\n", val)
		fmt.Printf("List: %s\n", list.String())
	}

	// Delete at position 1 (middle)
	if val, err := list.DeleteAtPosition(1); err == nil {
		fmt.Printf("\nDeleted at position 1: %d\n", val)
		fmt.Printf("List: %s\n", list.String())
	}

	// Delete all remaining
	fmt.Println("\nDeleting all remaining elements:")
	for !list.IsEmpty() {
		if val, err := list.DeleteAtBeginning(); err == nil {
			fmt.Printf("  Deleted: %d\n", val)
		}
	}
	fmt.Printf("List after clearing: %s (isEmpty: %v)\n", list.String(), list.IsEmpty())
}

func demo4SearchAndAccess() {
	utils.PrintSubHeader("4. Search and Access Operations")

	type Person struct {
		Name string
		Age  int
	}

	list := datastructures.NewLinkedList[Person]()

	// Add people
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
		{Name: "Diana", Age: 28},
	}

	for _, person := range people {
		list.InsertAtEnd(person)
	}

	fmt.Printf("List: %s\n", list.String())

	// Search for a person
	searchName := "Charlie"
	equals := func(a, b Person) bool {
		return a.Name == b.Name
	}

	if index, found := list.Search(Person{Name: searchName}, equals); found {
		fmt.Printf("\nFound '%s' at index %d\n", searchName, index)
	} else {
		fmt.Printf("\n'%s' not found\n", searchName)
	}

	// Get element at position
	if person, err := list.Get(2); err == nil {
		fmt.Printf("Person at index 2: %s (Age: %d)\n", person.Name, person.Age)
	}

	// Convert to slice
	slice := list.ToSlice()
	fmt.Printf("\nAs slice: %v\n", slice)

	// ForEach
	fmt.Println("\nIterating with ForEach:")
	list.ForEach(func(p Person) {
		fmt.Printf("  %s is %d years old\n", p.Name, p.Age)
	})
}

func demo5Manipulation() {
	utils.PrintSubHeader("5. List Manipulation")

	list := datastructures.NewLinkedList[int]()

	// Build list
	for i := 1; i <= 5; i++ {
		list.InsertAtEnd(i)
	}
	fmt.Printf("Original list: %s\n", list.String())

	// Reverse
	list.Reverse()
	fmt.Printf("After reverse: %s\n", list.String())

	// Reverse again to restore
	list.Reverse()
	fmt.Printf("After second reverse: %s\n", list.String())

	// Clear
	list.Clear()
	fmt.Printf("After clear: %s (size: %d)\n", list.String(), list.Size())

	// Rebuild and demonstrate edge cases
	fmt.Println("\nEdge cases:")

	// Single element
	list.InsertAtEnd(42)
	fmt.Printf("Single element: %s\n", list.String())
	list.Reverse()
	fmt.Printf("After reverse: %s\n", list.String())

	// Two elements
	list.Clear()
	list.InsertAtEnd(1)
	list.InsertAtEnd(2)
	fmt.Printf("\nTwo elements: %s\n", list.String())
	list.Reverse()
	fmt.Printf("After reverse: %s\n", list.String())
}
