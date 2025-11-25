//go:build ignore

package main

import "fmt"

// Stack represents a LIFO (Last-In-First-Out) data structure
type Stack struct {
	items []string
}

// Push adds an item to the top of the stack
func (s *Stack) Push(item string) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack
func (s *Stack) Pop() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}

	n := len(s.items) - 1
	item := s.items[n]
	s.items = s.items[:n]
	return item, true
}

// Peek returns the top item without removing it
func (s *Stack) Peek() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty checks if the stack has no items
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func main() {
	stack := &Stack{}

	// Adding elements
	stack.Push("world!")
	stack.Push("Hello ")

	// Print all elements in LIFO order
	for !stack.IsEmpty() {
		if item, ok := stack.Pop(); ok {
			fmt.Print(item)
		}
	}
	fmt.Println() // Output: Hello world!
}
