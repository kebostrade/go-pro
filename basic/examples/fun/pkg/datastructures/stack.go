package datastructures

import (
	"errors"
	"fmt"
	"sync"
)

// ErrEmptyStack is returned when attempting to pop or peek from an empty stack
var ErrEmptyStack = errors.New("stack is empty")

// Stack represents a generic LIFO (Last-In-First-Out) data structure
// It is thread-safe when using the concurrent version
type Stack[T any] struct {
	items []T
	mu    sync.RWMutex
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// NewStackWithCapacity creates a new stack with pre-allocated capacity
func NewStackWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0, capacity),
	}
}

// Push adds an item to the top of the stack
// Time complexity: O(1) amortized
func (s *Stack[T]) Push(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack
// Returns an error if the stack is empty
// Time complexity: O(1)
func (s *Stack[T]) Pop() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var zero T
	if len(s.items) == 0 {
		return zero, ErrEmptyStack
	}

	n := len(s.items) - 1
	item := s.items[n]
	s.items = s.items[:n]
	return item, nil
}

// Peek returns the top item without removing it
// Returns an error if the stack is empty
// Time complexity: O(1)
func (s *Stack[T]) Peek() (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var zero T
	if len(s.items) == 0 {
		return zero, ErrEmptyStack
	}
	return s.items[len(s.items)-1], nil
}

// IsEmpty checks if the stack has no items
// Time complexity: O(1)
func (s *Stack[T]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items) == 0
}

// Size returns the number of items in the stack
// Time complexity: O(1)
func (s *Stack[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// Clear removes all items from the stack
// Time complexity: O(1)
func (s *Stack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make([]T, 0)
}

// ToSlice returns a copy of the stack as a slice (top to bottom)
// Time complexity: O(n)
func (s *Stack[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]T, len(s.items))
	for i := 0; i < len(s.items); i++ {
		result[i] = s.items[len(s.items)-1-i]
	}
	return result
}

// String returns a string representation of the stack
func (s *Stack[T]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.items) == 0 {
		return "Stack[]"
	}

	return fmt.Sprintf("Stack%v", s.items)
}

// Clone creates a deep copy of the stack
// Time complexity: O(n)
func (s *Stack[T]) Clone() *Stack[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newStack := NewStackWithCapacity[T](len(s.items))
	newStack.items = make([]T, len(s.items))
	copy(newStack.items, s.items)
	return newStack
}

// Reverse reverses the order of items in the stack
// Time complexity: O(n)
func (s *Stack[T]) Reverse() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, j := 0, len(s.items)-1; i < j; i, j = i+1, j-1 {
		s.items[i], s.items[j] = s.items[j], s.items[i]
	}
}

// Contains checks if the stack contains a specific item
// Time complexity: O(n)
func (s *Stack[T]) Contains(item T, equals func(T, T) bool) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.items {
		if equals(v, item) {
			return true
		}
	}
	return false
}

// Filter returns a new stack containing only items that match the predicate
// Time complexity: O(n)
func (s *Stack[T]) Filter(predicate func(T) bool) *Stack[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newStack := NewStack[T]()
	for _, item := range s.items {
		if predicate(item) {
			newStack.items = append(newStack.items, item)
		}
	}
	return newStack
}

// Map applies a function to each item and returns a new stack
// Time complexity: O(n)
func (s *Stack[T]) Map(fn func(T) T) *Stack[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newStack := NewStackWithCapacity[T](len(s.items))
	for _, item := range s.items {
		newStack.items = append(newStack.items, fn(item))
	}
	return newStack
}

// ForEach applies a function to each item in the stack (top to bottom)
// Time complexity: O(n)
func (s *Stack[T]) ForEach(fn func(T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := len(s.items) - 1; i >= 0; i-- {
		fn(s.items[i])
	}
}
