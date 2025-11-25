package test

import (
	"testing"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/datastructures"
)

// Stack Tests

func TestStack_PushPop(t *testing.T) {
	stack := datastructures.NewStack[int]()

	// Test push
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.Size() != 3 {
		t.Errorf("Expected size 3, got %d", stack.Size())
	}

	// Test pop
	val, err := stack.Pop()
	if err != nil || val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}

	val, err = stack.Pop()
	if err != nil || val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	if stack.Size() != 1 {
		t.Errorf("Expected size 1, got %d", stack.Size())
	}
}

func TestStack_EmptyPop(t *testing.T) {
	stack := datastructures.NewStack[int]()

	_, err := stack.Pop()
	if err == nil {
		t.Error("Expected error when popping from empty stack")
	}
}

func TestStack_Peek(t *testing.T) {
	stack := datastructures.NewStack[string]()

	stack.Push("first")
	stack.Push("second")

	val, err := stack.Peek()
	if err != nil || val != "second" {
		t.Errorf("Expected 'second', got '%s'", val)
	}

	// Size should not change after peek
	if stack.Size() != 2 {
		t.Errorf("Expected size 2 after peek, got %d", stack.Size())
	}
}

// Queue Tests

func TestQueue_EnqueueDequeue(t *testing.T) {
	queue := datastructures.NewQueue[int]()

	// Test enqueue
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	if queue.Size() != 3 {
		t.Errorf("Expected size 3, got %d", queue.Size())
	}

	// Test dequeue (FIFO)
	val, err := queue.Dequeue()
	if err != nil || val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	val, err = queue.Dequeue()
	if err != nil || val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	if queue.Size() != 1 {
		t.Errorf("Expected size 1, got %d", queue.Size())
	}
}

func TestQueue_EmptyDequeue(t *testing.T) {
	queue := datastructures.NewQueue[int]()

	_, err := queue.Dequeue()
	if err == nil {
		t.Error("Expected error when dequeueing from empty queue")
	}
}

func TestPriorityQueue(t *testing.T) {
	pq := datastructures.NewPriorityQueue[string]()

	pq.Enqueue("low1", datastructures.PriorityLow)
	pq.Enqueue("high1", datastructures.PriorityHigh)
	pq.Enqueue("medium1", datastructures.PriorityMedium)
	pq.Enqueue("high2", datastructures.PriorityHigh)

	// Should dequeue high priority first
	val, err := pq.Dequeue()
	if err != nil || val != "high1" {
		t.Errorf("Expected 'high1', got '%s'", val)
	}

	val, err = pq.Dequeue()
	if err != nil || val != "high2" {
		t.Errorf("Expected 'high2', got '%s'", val)
	}

	// Then medium
	val, err = pq.Dequeue()
	if err != nil || val != "medium1" {
		t.Errorf("Expected 'medium1', got '%s'", val)
	}

	// Finally low
	val, err = pq.Dequeue()
	if err != nil || val != "low1" {
		t.Errorf("Expected 'low1', got '%s'", val)
	}
}

// LinkedList Tests

func TestLinkedList_InsertAtBeginning(t *testing.T) {
	list := datastructures.NewLinkedList[int]()

	list.InsertAtBeginning(3)
	list.InsertAtBeginning(2)
	list.InsertAtBeginning(1)

	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}

	val, _ := list.Get(0)
	if val != 1 {
		t.Errorf("Expected first element to be 1, got %d", val)
	}
}

func TestLinkedList_InsertAtEnd(t *testing.T) {
	list := datastructures.NewLinkedList[int]()

	list.InsertAtEnd(1)
	list.InsertAtEnd(2)
	list.InsertAtEnd(3)

	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}

	val, _ := list.Get(2)
	if val != 3 {
		t.Errorf("Expected last element to be 3, got %d", val)
	}
}

func TestLinkedList_DeleteAtBeginning(t *testing.T) {
	list := datastructures.NewLinkedList[int]()

	list.InsertAtEnd(1)
	list.InsertAtEnd(2)
	list.InsertAtEnd(3)

	val, err := list.DeleteAtBeginning()
	if err != nil || val != 1 {
		t.Errorf("Expected to delete 1, got %d", val)
	}

	if list.Size() != 2 {
		t.Errorf("Expected size 2, got %d", list.Size())
	}
}

func TestLinkedList_Reverse(t *testing.T) {
	list := datastructures.NewLinkedList[int]()

	for i := 1; i <= 5; i++ {
		list.InsertAtEnd(i)
	}

	list.Reverse()

	val, _ := list.Get(0)
	if val != 5 {
		t.Errorf("Expected first element to be 5 after reverse, got %d", val)
	}

	val, _ = list.Get(4)
	if val != 1 {
		t.Errorf("Expected last element to be 1 after reverse, got %d", val)
	}
}

func TestLinkedList_Search(t *testing.T) {
	list := datastructures.NewLinkedList[string]()

	list.InsertAtEnd("apple")
	list.InsertAtEnd("banana")
	list.InsertAtEnd("cherry")

	equals := func(a, b string) bool {
		return a == b
	}

	index, found := list.Search("banana", equals)
	if !found || index != 1 {
		t.Errorf("Expected to find 'banana' at index 1, got index %d, found %v", index, found)
	}

	index, found = list.Search("grape", equals)
	if found {
		t.Error("Expected not to find 'grape'")
	}
}

// Benchmark Tests

func BenchmarkStack_Push(b *testing.B) {
	stack := datastructures.NewStack[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkStack_Pop(b *testing.B) {
	stack := datastructures.NewStack[int]()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

func BenchmarkQueue_Enqueue(b *testing.B) {
	queue := datastructures.NewQueue[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
}

func BenchmarkLinkedList_InsertAtEnd(b *testing.B) {
	list := datastructures.NewLinkedList[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.InsertAtEnd(i)
	}
}
