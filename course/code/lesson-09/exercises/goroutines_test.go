package exercises

import (
	"testing"
	"time"
)

func TestBasicGoroutine(t *testing.T) {
	ch := make(chan int)
	
	go BasicGoroutine(ch)
	
	expected := []int{1, 2, 3, 4, 5}
	for i, exp := range expected {
		select {
		case got := <-ch:
			if got != exp {
				t.Errorf("Expected %d at position %d, got %d", exp, i, got)
			}
		case <-time.After(time.Second):
			t.Fatalf("Timeout waiting for value")
		}
	}
	
	// Check channel is closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("Channel should be closed")
		}
	default:
	}
}

func TestBufferedChannelExample(t *testing.T) {
	sent, received := BufferedChannelExample()
	if sent != 3 {
		t.Errorf("Expected sent=3, got %d", sent)
	}
	if received != 3 {
		t.Errorf("Expected received=3, got %d", received)
	}
}

func TestSelectFirst(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	
	// Send to ch2 first
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch2 <- 42
	}()
	
	result := SelectFirst(ch1, ch2)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestSelectWithTimeout(t *testing.T) {
	ch := make(chan int)
	
	// Should timeout
	_, ok := SelectWithTimeout(ch, 50*time.Millisecond)
	if ok {
		t.Error("Expected timeout")
	}
	
	// Should receive value
	ch2 := make(chan int)
	go func() {
		ch2 <- 100
	}()
	
	val, ok := SelectWithTimeout(ch2, time.Second)
	if !ok {
		t.Error("Expected to receive value")
	}
	if val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}
}

func TestReceiverOnly(t *testing.T) {
	ch := make(chan string)
	
	go func() {
		ch <- "hello"
	}()
	
	result := ReceiverOnly(ch)
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestFanOut(t *testing.T) {
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	
	// Send jobs
	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	close(jobs)
	
	// Run fan-out with 3 workers
	FanOut(jobs, results, 3)
	
	close(results)
	
	sum := 0
	for r := range results {
		sum += r
	}
	
	// Sum of 1+2+...+10 = 55 (each worker processes numbers)
	// Each number should be processed once
	if sum != 55 {
		t.Errorf("Expected sum=55, got %d", sum)
	}
}

func TestPipeline(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	
	pipeline := Square(Generate(nums...))
	result := Sum(pipeline)
	
	// 1^2 + 2^2 + 3^2 + 4^2 + 5^2 = 1 + 4 + 9 + 16 + 25 = 55
	if result != 55 {
		t.Errorf("Expected 55, got %d", result)
	}
}

func TestGenerator(t *testing.T) {
	// Test that Generate works correctly
	ch := Generate(10, 20, 30)
	
	expected := []int{10, 20, 30}
	for _, exp := range expected {
		select {
		case got := <-ch:
			if got != exp {
				t.Errorf("Expected %d, got %d", exp, got)
			}
		case <-time.After(time.Second):
			t.Fatalf("Timeout")
		}
	}
}

func TestSquare(t *testing.T) {
	// Test Square function
	ch := make(chan int)
	go func() {
		ch <- 3
		ch <- 4
		close(ch)
	}()
	
	squared := Square(ch)
	
	val1 := <-squared
	if val1 != 9 {
		t.Errorf("Expected 9, got %d", val1)
	}
	
	val2 := <-squared
	if val2 != 16 {
		t.Errorf("Expected 16, got %d", val2)
	}
}
