//go:build ignore

package main

import "fmt"

// Task: Implement a singly linked list data structure with common operations.
// Support insertion, deletion, search, and traversal operations.

// Node represents a single node in the linked list
type Node struct {
	Data int
	Next *Node
}

// LinkedList represents a singly linked list
type LinkedList struct {
	Head *Node
	Size int
}

// NewLinkedList creates a new empty linked list
func NewLinkedList() *LinkedList {
	return &LinkedList{
		Head: nil,
		Size: 0,
	}
}

// InsertAtBeginning adds a new node at the start of the list
func (ll *LinkedList) InsertAtBeginning(data int) {
	newNode := &Node{Data: data, Next: ll.Head}
	ll.Head = newNode
	ll.Size++
}

// InsertAtEnd adds a new node at the end of the list
func (ll *LinkedList) InsertAtEnd(data int) {
	newNode := &Node{Data: data, Next: nil}

	if ll.Head == nil {
		ll.Head = newNode
		ll.Size++
		return
	}

	current := ll.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
	ll.Size++
}

// InsertAtPosition inserts a node at a specific position (0-indexed)
func (ll *LinkedList) InsertAtPosition(data int, position int) bool {
	if position < 0 || position > ll.Size {
		return false
	}

	if position == 0 {
		ll.InsertAtBeginning(data)
		return true
	}

	newNode := &Node{Data: data}
	current := ll.Head
	for i := 0; i < position-1; i++ {
		current = current.Next
	}

	newNode.Next = current.Next
	current.Next = newNode
	ll.Size++
	return true
}

// DeleteAtBeginning removes the first node
func (ll *LinkedList) DeleteAtBeginning() (int, bool) {
	if ll.Head == nil {
		return 0, false
	}

	data := ll.Head.Data
	ll.Head = ll.Head.Next
	ll.Size--
	return data, true
}

// DeleteAtEnd removes the last node
func (ll *LinkedList) DeleteAtEnd() (int, bool) {
	if ll.Head == nil {
		return 0, false
	}

	if ll.Head.Next == nil {
		data := ll.Head.Data
		ll.Head = nil
		ll.Size--
		return data, true
	}

	current := ll.Head
	for current.Next.Next != nil {
		current = current.Next
	}

	data := current.Next.Data
	current.Next = nil
	ll.Size--
	return data, true
}

// DeleteByValue removes the first node with the given value
func (ll *LinkedList) DeleteByValue(value int) bool {
	if ll.Head == nil {
		return false
	}

	// If head node contains the value
	if ll.Head.Data == value {
		ll.Head = ll.Head.Next
		ll.Size--
		return true
	}

	current := ll.Head
	for current.Next != nil {
		if current.Next.Data == value {
			current.Next = current.Next.Next
			ll.Size--
			return true
		}
		current = current.Next
	}

	return false
}

// Search finds if a value exists in the list
func (ll *LinkedList) Search(value int) bool {
	current := ll.Head
	for current != nil {
		if current.Data == value {
			return true
		}
		current = current.Next
	}
	return false
}

// GetAt returns the value at a specific position
func (ll *LinkedList) GetAt(position int) (int, bool) {
	if position < 0 || position >= ll.Size {
		return 0, false
	}

	current := ll.Head
	for i := 0; i < position; i++ {
		current = current.Next
	}
	return current.Data, true
}

// Reverse reverses the linked list
func (ll *LinkedList) Reverse() {
	var prev *Node
	current := ll.Head
	var next *Node

	for current != nil {
		next = current.Next
		current.Next = prev
		prev = current
		current = next
	}

	ll.Head = prev
}

// Display prints all elements in the list
func (ll *LinkedList) Display() {
	if ll.Head == nil {
		fmt.Println("List is empty")
		return
	}

	current := ll.Head
	fmt.Print("List: ")
	for current != nil {
		fmt.Printf("%d", current.Data)
		if current.Next != nil {
			fmt.Print(" -> ")
		}
		current = current.Next
	}
	fmt.Println()
}

// ToSlice converts the linked list to a slice
func (ll *LinkedList) ToSlice() []int {
	result := make([]int, 0, ll.Size)
	current := ll.Head
	for current != nil {
		result = append(result, current.Data)
		current = current.Next
	}
	return result
}

// Clear removes all nodes from the list
func (ll *LinkedList) Clear() {
	ll.Head = nil
	ll.Size = 0
}

// repeatString repeats a string n times
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func main() {
	fmt.Println("Linked List Data Structure Demo")
	fmt.Println(repeatString("=", 60))

	// Demo 1: Basic Operations
	fmt.Println("\n1. Basic Insert and Display Operations")
	fmt.Println(repeatString("-", 60))

	list := NewLinkedList()

	list.InsertAtEnd(10)
	list.InsertAtEnd(20)
	list.InsertAtEnd(30)
	list.Display()
	fmt.Printf("Size: %d\n", list.Size)

	// Demo 2: Insert at Beginning
	fmt.Println("\n2. Insert at Beginning")
	fmt.Println(repeatString("-", 60))

	list.InsertAtBeginning(5)
	list.Display()

	// Demo 3: Insert at Position
	fmt.Println("\n3. Insert at Position")
	fmt.Println(repeatString("-", 60))

	list.InsertAtPosition(15, 2)
	list.Display()
	fmt.Printf("Size: %d\n", list.Size)

	// Demo 4: Search
	fmt.Println("\n4. Search Operations")
	fmt.Println(repeatString("-", 60))

	searchValues := []int{15, 25, 30}
	for _, val := range searchValues {
		if list.Search(val) {
			fmt.Printf("Value %d found in list\n", val)
		} else {
			fmt.Printf("Value %d not found in list\n", val)
		}
	}

	// Demo 5: Get at Position
	fmt.Println("\n5. Get Value at Position")
	fmt.Println(repeatString("-", 60))

	if value, ok := list.GetAt(2); ok {
		fmt.Printf("Value at position 2: %d\n", value)
	}

	// Demo 6: Delete Operations
	fmt.Println("\n6. Delete Operations")
	fmt.Println(repeatString("-", 60))

	list.Display()

	if value, ok := list.DeleteAtBeginning(); ok {
		fmt.Printf("Deleted from beginning: %d\n", value)
	}
	list.Display()

	if value, ok := list.DeleteAtEnd(); ok {
		fmt.Printf("Deleted from end: %d\n", value)
	}
	list.Display()

	if list.DeleteByValue(15) {
		fmt.Println("Deleted value: 15")
	}
	list.Display()

	// Demo 7: Reverse
	fmt.Println("\n7. Reverse List")
	fmt.Println(repeatString("-", 60))

	list.Clear()
	for i := 1; i <= 5; i++ {
		list.InsertAtEnd(i * 10)
	}

	fmt.Print("Original: ")
	list.Display()

	list.Reverse()
	fmt.Print("Reversed: ")
	list.Display()

	// Demo 8: Convert to Slice
	fmt.Println("\n8. Convert to Slice")
	fmt.Println(repeatString("-", 60))

	slice := list.ToSlice()
	fmt.Printf("As slice: %v\n", slice)

	fmt.Println("\nLinked List demo completed!")
}
