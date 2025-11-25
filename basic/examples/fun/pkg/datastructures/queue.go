package datastructures

import (
	"errors"
	"fmt"
	"sync"
)

// ErrEmptyQueue is returned when attempting to dequeue from an empty queue
var ErrEmptyQueue = errors.New("queue is empty")

// Queue represents a generic FIFO (First-In-First-Out) data structure
// It is thread-safe when using the concurrent version
type Queue[T any] struct {
	items []T
	mu    sync.RWMutex
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// NewQueueWithCapacity creates a new queue with pre-allocated capacity
func NewQueueWithCapacity[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, capacity),
	}
}

// Enqueue adds an item to the back of the queue
// Time complexity: O(1) amortized
func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item from the queue
// Returns an error if the queue is empty
// Time complexity: O(n) due to slice reallocation
func (q *Queue[T]) Dequeue() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var zero T
	if len(q.items) == 0 {
		return zero, ErrEmptyQueue
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek returns the front item without removing it
// Returns an error if the queue is empty
// Time complexity: O(1)
func (q *Queue[T]) Peek() (T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	var zero T
	if len(q.items) == 0 {
		return zero, ErrEmptyQueue
	}
	return q.items[0], nil
}

// IsEmpty checks if the queue has no items
// Time complexity: O(1)
func (q *Queue[T]) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items) == 0
}

// Size returns the number of items in the queue
// Time complexity: O(1)
func (q *Queue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items)
}

// Clear removes all items from the queue
// Time complexity: O(1)
func (q *Queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = make([]T, 0)
}

// ToSlice returns a copy of the queue as a slice (front to back)
// Time complexity: O(n)
func (q *Queue[T]) ToSlice() []T {
	q.mu.RLock()
	defer q.mu.RUnlock()

	result := make([]T, len(q.items))
	copy(result, q.items)
	return result
}

// String returns a string representation of the queue
func (q *Queue[T]) String() string {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.items) == 0 {
		return "Queue[]"
	}

	return fmt.Sprintf("Queue%v", q.items)
}

// Clone creates a deep copy of the queue
// Time complexity: O(n)
func (q *Queue[T]) Clone() *Queue[T] {
	q.mu.RLock()
	defer q.mu.RUnlock()

	newQueue := NewQueueWithCapacity[T](len(q.items))
	newQueue.items = make([]T, len(q.items))
	copy(newQueue.items, q.items)
	return newQueue
}

// Contains checks if the queue contains a specific item
// Time complexity: O(n)
func (q *Queue[T]) Contains(item T, equals func(T, T) bool) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, v := range q.items {
		if equals(v, item) {
			return true
		}
	}
	return false
}

// Filter returns a new queue containing only items that match the predicate
// Time complexity: O(n)
func (q *Queue[T]) Filter(predicate func(T) bool) *Queue[T] {
	q.mu.RLock()
	defer q.mu.RUnlock()

	newQueue := NewQueue[T]()
	for _, item := range q.items {
		if predicate(item) {
			newQueue.items = append(newQueue.items, item)
		}
	}
	return newQueue
}

// ForEach applies a function to each item in the queue (front to back)
// Time complexity: O(n)
func (q *Queue[T]) ForEach(fn func(T)) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, item := range q.items {
		fn(item)
	}
}

// Priority represents the priority level for items in a priority queue
type Priority int

const (
	PriorityHigh   Priority = 1
	PriorityMedium Priority = 2
	PriorityLow    Priority = 3
)

// PriorityItem represents an item with a priority
type PriorityItem[T any] struct {
	Value    T
	Priority Priority
}

// PriorityQueue represents a queue with priority levels
type PriorityQueue[T any] struct {
	high   []T
	medium []T
	low    []T
	mu     sync.RWMutex
}

// NewPriorityQueue creates a new empty priority queue
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		high:   make([]T, 0),
		medium: make([]T, 0),
		low:    make([]T, 0),
	}
}

// Enqueue adds an item with a priority level
// Time complexity: O(1) amortized
func (pq *PriorityQueue[T]) Enqueue(item T, priority Priority) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	switch priority {
	case PriorityHigh:
		pq.high = append(pq.high, item)
	case PriorityMedium:
		pq.medium = append(pq.medium, item)
	case PriorityLow:
		pq.low = append(pq.low, item)
	default:
		pq.medium = append(pq.medium, item)
	}
}

// Dequeue removes and returns the highest priority item
// Returns an error if the queue is empty
// Time complexity: O(n) due to slice reallocation
func (pq *PriorityQueue[T]) Dequeue() (T, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	var zero T

	// Check high priority first
	if len(pq.high) > 0 {
		item := pq.high[0]
		pq.high = pq.high[1:]
		return item, nil
	}

	// Then medium priority
	if len(pq.medium) > 0 {
		item := pq.medium[0]
		pq.medium = pq.medium[1:]
		return item, nil
	}

	// Finally low priority
	if len(pq.low) > 0 {
		item := pq.low[0]
		pq.low = pq.low[1:]
		return item, nil
	}

	return zero, ErrEmptyQueue
}

// IsEmpty checks if the priority queue has no items
// Time complexity: O(1)
func (pq *PriorityQueue[T]) IsEmpty() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return len(pq.high) == 0 && len(pq.medium) == 0 && len(pq.low) == 0
}

// Size returns the total number of items in the priority queue
// Time complexity: O(1)
func (pq *PriorityQueue[T]) Size() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return len(pq.high) + len(pq.medium) + len(pq.low)
}

// Clear removes all items from the priority queue
// Time complexity: O(1)
func (pq *PriorityQueue[T]) Clear() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.high = make([]T, 0)
	pq.medium = make([]T, 0)
	pq.low = make([]T, 0)
}

// String returns a string representation of the priority queue
func (pq *PriorityQueue[T]) String() string {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return fmt.Sprintf("PriorityQueue{High: %v, Medium: %v, Low: %v}",
		pq.high, pq.medium, pq.low)
}
