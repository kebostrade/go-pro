package datastructures

import (
	"testing"
)

// Stack Tests
func TestStack(t *testing.T) {
	t.Run("Push and Pop", func(t *testing.T) {
		s := NewStack[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)

		val, err := s.Pop()
		if err != nil || val != 3 {
			t.Errorf("Pop() = %v, %v; want 3, nil", val, err)
		}

		if s.Size() != 2 {
			t.Errorf("Size() = %d; want 2", s.Size())
		}
	})

	t.Run("Peek", func(t *testing.T) {
		s := NewStack[int]()
		s.Push(42)

		val, err := s.Peek()
		if err != nil || val != 42 {
			t.Errorf("Peek() = %v, %v; want 42, nil", val, err)
		}

		if s.Size() != 1 {
			t.Errorf("Peek should not remove element")
		}
	})

	t.Run("Empty stack", func(t *testing.T) {
		s := NewStack[int]()

		if !s.IsEmpty() {
			t.Error("IsEmpty() = false; want true")
		}

		_, err := s.Pop()
		if err != ErrEmptyStack {
			t.Errorf("Pop() error = %v; want ErrEmptyStack", err)
		}
	})

	t.Run("ToSlice", func(t *testing.T) {
		s := NewStack[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)

		slice := s.ToSlice()
		want := []int{3, 2, 1}

		if !intSliceEqual(slice, want) {
			t.Errorf("ToSlice() = %v; want %v", slice, want)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		s := NewStack[int]()
		s.Push(1)
		s.Push(2)
		s.Clear()

		if !s.IsEmpty() {
			t.Error("Clear() did not empty stack")
		}
	})

	t.Run("Clone", func(t *testing.T) {
		s := NewStack[int]()
		s.Push(1)
		s.Push(2)

		clone := s.Clone()
		clone.Push(3)

		if s.Size() == clone.Size() {
			t.Error("Clone should be independent")
		}
	})
}

// Queue Tests
func TestQueue(t *testing.T) {
	t.Run("Enqueue and Dequeue", func(t *testing.T) {
		q := NewQueue[int]()
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)

		val, err := q.Dequeue()
		if err != nil || val != 1 {
			t.Errorf("Dequeue() = %v, %v; want 1, nil", val, err)
		}

		if q.Size() != 2 {
			t.Errorf("Size() = %d; want 2", q.Size())
		}
	})

	t.Run("Peek", func(t *testing.T) {
		q := NewQueue[int]()
		q.Enqueue(42)

		val, err := q.Peek()
		if err != nil || val != 42 {
			t.Errorf("Peek() = %v, %v; want 42, nil", val, err)
		}

		if q.Size() != 1 {
			t.Error("Peek should not remove element")
		}
	})

	t.Run("Empty queue", func(t *testing.T) {
		q := NewQueue[int]()

		if !q.IsEmpty() {
			t.Error("IsEmpty() = false; want true")
		}

		_, err := q.Dequeue()
		if err != ErrEmptyQueue {
			t.Errorf("Dequeue() error = %v; want ErrEmptyQueue", err)
		}
	})

	t.Run("ToSlice", func(t *testing.T) {
		q := NewQueue[int]()
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)

		slice := q.ToSlice()
		want := []int{1, 2, 3}

		if !intSliceEqual(slice, want) {
			t.Errorf("ToSlice() = %v; want %v", slice, want)
		}
	})
}

// PriorityQueue Tests
func TestPriorityQueue(t *testing.T) {
	t.Run("Priority ordering", func(t *testing.T) {
		pq := NewPriorityQueue[string]()

		pq.Enqueue("low", PriorityLow)
		pq.Enqueue("high", PriorityHigh)
		pq.Enqueue("medium", PriorityMedium)

		val, err := pq.Dequeue()
		if err != nil || val != "high" {
			t.Errorf("Dequeue() = %v, %v; want \"high\", nil", val, err)
		}

		val, err = pq.Dequeue()
		if err != nil || val != "medium" {
			t.Errorf("Dequeue() = %v, %v; want \"medium\", nil", val, err)
		}

		val, err = pq.Dequeue()
		if err != nil || val != "low" {
			t.Errorf("Dequeue() = %v, %v; want \"low\", nil", val, err)
		}
	})

	t.Run("Empty priority queue", func(t *testing.T) {
		pq := NewPriorityQueue[int]()

		if !pq.IsEmpty() {
			t.Error("IsEmpty() = false; want true")
		}

		_, err := pq.Dequeue()
		if err != ErrEmptyQueue {
			t.Errorf("Dequeue() error = %v; want ErrEmptyQueue", err)
		}
	})
}

// LinkedList Tests
func TestLinkedList(t *testing.T) {
	t.Run("Insert at beginning", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtBeginning(1)
		ll.InsertAtBeginning(2)

		slice := ll.ToSlice()
		want := []int{2, 1}

		if !intSliceEqual(slice, want) {
			t.Errorf("ToSlice() = %v; want %v", slice, want)
		}
	})

	t.Run("Insert at end", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtEnd(1)
		ll.InsertAtEnd(2)

		slice := ll.ToSlice()
		want := []int{1, 2}

		if !intSliceEqual(slice, want) {
			t.Errorf("ToSlice() = %v; want %v", slice, want)
		}
	})

	t.Run("Insert at position", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtEnd(1)
		ll.InsertAtEnd(3)
		ll.InsertAtPosition(2, 1)

		slice := ll.ToSlice()
		want := []int{1, 2, 3}

		if !intSliceEqual(slice, want) {
			t.Errorf("ToSlice() = %v; want %v", slice, want)
		}
	})

	t.Run("Delete at beginning", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtEnd(1)
		ll.InsertAtEnd(2)

		val, err := ll.DeleteAtBeginning()
		if err != nil || val != 1 {
			t.Errorf("DeleteAtBeginning() = %v, %v; want 1, nil", val, err)
		}

		if ll.Size() != 1 {
			t.Errorf("Size() = %d; want 1", ll.Size())
		}
	})

	t.Run("Delete at end", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtEnd(1)
		ll.InsertAtEnd(2)

		val, err := ll.DeleteAtEnd()
		if err != nil || val != 2 {
			t.Errorf("DeleteAtEnd() = %v, %v; want 2, nil", val, err)
		}
	})

	t.Run("Get by position", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtEnd(1)
		ll.InsertAtEnd(2)
		ll.InsertAtEnd(3)

		val, err := ll.Get(1)
		if err != nil || val != 2 {
			t.Errorf("Get(1) = %v, %v; want 2, nil", val, err)
		}
	})

	t.Run("Search", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtEnd(1)
		ll.InsertAtEnd(2)
		ll.InsertAtEnd(3)

		idx, found := ll.Search(2, func(a, b int) bool { return a == b })
		if !found || idx != 1 {
			t.Errorf("Search(2) = %d, %v; want 1, true", idx, found)
		}
	})

	t.Run("Reverse", func(t *testing.T) {
		ll := NewLinkedList[int]()
		ll.InsertAtEnd(1)
		ll.InsertAtEnd(2)
		ll.InsertAtEnd(3)

		ll.Reverse()

		slice := ll.ToSlice()
		want := []int{3, 2, 1}

		if !intSliceEqual(slice, want) {
			t.Errorf("After Reverse ToSlice() = %v; want %v", slice, want)
		}
	})

	t.Run("Empty list operations", func(t *testing.T) {
		ll := NewLinkedList[int]()

		if !ll.IsEmpty() {
			t.Error("IsEmpty() = false; want true")
		}

		_, err := ll.DeleteAtBeginning()
		if err != ErrEmptyList {
			t.Errorf("DeleteAtBeginning() error = %v; want ErrEmptyList", err)
		}
	})
}

// Helper function
func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Benchmarks
func BenchmarkStackPush(b *testing.B) {
	s := NewStack[int]()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
}

func BenchmarkStackPop(b *testing.B) {
	s := NewStack[int]()
	for i := 0; i < 10000; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N && !s.IsEmpty(); i++ {
		s.Pop()
	}
}

func BenchmarkQueueEnqueue(b *testing.B) {
	q := NewQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

func BenchmarkLinkedListInsertAtEnd(b *testing.B) {
	ll := NewLinkedList[int]()
	for i := 0; i < b.N; i++ {
		ll.InsertAtEnd(i)
	}
}
