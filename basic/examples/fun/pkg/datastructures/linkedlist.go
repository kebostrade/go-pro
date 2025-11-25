package datastructures

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	// ErrEmptyList is returned when attempting operations on an empty list
	ErrEmptyList = errors.New("linked list is empty")
	// ErrIndexOutOfBounds is returned when accessing an invalid index
	ErrIndexOutOfBounds = errors.New("index out of bounds")
)

// Node represents a single node in the linked list
type Node[T any] struct {
	Data T
	Next *Node[T]
}

// LinkedList represents a generic singly linked list
type LinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
	size int
	mu   sync.RWMutex
}

// NewLinkedList creates a new empty linked list
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// InsertAtBeginning adds a new node at the start of the list
// Time complexity: O(1)
func (ll *LinkedList[T]) InsertAtBeginning(data T) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	newNode := &Node[T]{Data: data, Next: ll.head}
	ll.head = newNode

	if ll.tail == nil {
		ll.tail = newNode
	}

	ll.size++
}

// InsertAtEnd adds a new node at the end of the list
// Time complexity: O(1) with tail pointer
func (ll *LinkedList[T]) InsertAtEnd(data T) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	newNode := &Node[T]{Data: data, Next: nil}

	if ll.head == nil {
		ll.head = newNode
		ll.tail = newNode
	} else {
		ll.tail.Next = newNode
		ll.tail = newNode
	}

	ll.size++
}

// InsertAtPosition inserts a node at a specific position (0-indexed)
// Time complexity: O(n)
func (ll *LinkedList[T]) InsertAtPosition(data T, position int) error {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	if position < 0 || position > ll.size {
		return ErrIndexOutOfBounds
	}

	if position == 0 {
		newNode := &Node[T]{Data: data, Next: ll.head}
		ll.head = newNode
		if ll.tail == nil {
			ll.tail = newNode
		}
		ll.size++
		return nil
	}

	current := ll.head
	for i := 0; i < position-1; i++ {
		current = current.Next
	}

	newNode := &Node[T]{Data: data, Next: current.Next}
	current.Next = newNode

	if newNode.Next == nil {
		ll.tail = newNode
	}

	ll.size++
	return nil
}

// DeleteAtBeginning removes the first node
// Time complexity: O(1)
func (ll *LinkedList[T]) DeleteAtBeginning() (T, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	var zero T
	if ll.head == nil {
		return zero, ErrEmptyList
	}

	data := ll.head.Data
	ll.head = ll.head.Next
	ll.size--

	if ll.head == nil {
		ll.tail = nil
	}

	return data, nil
}

// DeleteAtEnd removes the last node
// Time complexity: O(n)
func (ll *LinkedList[T]) DeleteAtEnd() (T, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	var zero T
	if ll.head == nil {
		return zero, ErrEmptyList
	}

	if ll.head.Next == nil {
		data := ll.head.Data
		ll.head = nil
		ll.tail = nil
		ll.size--
		return data, nil
	}

	current := ll.head
	for current.Next.Next != nil {
		current = current.Next
	}

	data := current.Next.Data
	current.Next = nil
	ll.tail = current
	ll.size--
	return data, nil
}

// DeleteAtPosition removes a node at a specific position (0-indexed)
// Time complexity: O(n)
func (ll *LinkedList[T]) DeleteAtPosition(position int) (T, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	var zero T
	if position < 0 || position >= ll.size {
		return zero, ErrIndexOutOfBounds
	}

	if position == 0 {
		data := ll.head.Data
		ll.head = ll.head.Next
		ll.size--
		if ll.head == nil {
			ll.tail = nil
		}
		return data, nil
	}

	current := ll.head
	for i := 0; i < position-1; i++ {
		current = current.Next
	}

	data := current.Next.Data
	current.Next = current.Next.Next

	if current.Next == nil {
		ll.tail = current
	}

	ll.size--
	return data, nil
}

// Search finds the first occurrence of data in the list
// Time complexity: O(n)
func (ll *LinkedList[T]) Search(data T, equals func(T, T) bool) (int, bool) {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	current := ll.head
	index := 0

	for current != nil {
		if equals(current.Data, data) {
			return index, true
		}
		current = current.Next
		index++
	}

	return -1, false
}

// Get retrieves the data at a specific position
// Time complexity: O(n)
func (ll *LinkedList[T]) Get(position int) (T, error) {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	var zero T
	if position < 0 || position >= ll.size {
		return zero, ErrIndexOutOfBounds
	}

	current := ll.head
	for i := 0; i < position; i++ {
		current = current.Next
	}

	return current.Data, nil
}

// Size returns the number of nodes in the list
// Time complexity: O(1)
func (ll *LinkedList[T]) Size() int {
	ll.mu.RLock()
	defer ll.mu.RUnlock()
	return ll.size
}

// IsEmpty checks if the list has no nodes
// Time complexity: O(1)
func (ll *LinkedList[T]) IsEmpty() bool {
	ll.mu.RLock()
	defer ll.mu.RUnlock()
	return ll.size == 0
}

// Clear removes all nodes from the list
// Time complexity: O(1)
func (ll *LinkedList[T]) Clear() {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	ll.head = nil
	ll.tail = nil
	ll.size = 0
}

// ToSlice converts the linked list to a slice
// Time complexity: O(n)
func (ll *LinkedList[T]) ToSlice() []T {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	result := make([]T, 0, ll.size)
	current := ll.head

	for current != nil {
		result = append(result, current.Data)
		current = current.Next
	}

	return result
}

// Reverse reverses the linked list in place
// Time complexity: O(n)
func (ll *LinkedList[T]) Reverse() {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	if ll.head == nil || ll.head.Next == nil {
		return
	}

	var prev *Node[T]
	current := ll.head
	ll.tail = ll.head

	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	ll.head = prev
}

// ForEach applies a function to each node in the list
// Time complexity: O(n)
func (ll *LinkedList[T]) ForEach(fn func(T)) {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	current := ll.head
	for current != nil {
		fn(current.Data)
		current = current.Next
	}
}

// String returns a string representation of the linked list
func (ll *LinkedList[T]) String() string {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	if ll.head == nil {
		return "LinkedList[]"
	}

	var sb strings.Builder
	sb.WriteString("LinkedList[")

	current := ll.head
	for current != nil {
		sb.WriteString(fmt.Sprintf("%v", current.Data))
		if current.Next != nil {
			sb.WriteString(" -> ")
		}
		current = current.Next
	}

	sb.WriteString("]")
	return sb.String()
}
