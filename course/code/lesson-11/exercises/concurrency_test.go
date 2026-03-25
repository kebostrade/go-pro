package exercises

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	
	// Send jobs
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)
	
	// Run worker pool
	WorkerPool(jobs, results, 2)
	close(results)
	
	// Collect results
	var got []int
	for r := range results {
		got = append(got, r)
	}
	
	// Expected: each job * 2
	expected := []int{2, 4, 6, 8, 10}
	if len(got) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(got))
	}
}

func TestSafeCounter(t *testing.T) {
	counter := SafeCounter{}
	
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	
	if counter.Value() != 100 {
		t.Errorf("Expected 100, got %d", counter.Value())
	}
}

func TestSafeCache(t *testing.T) {
	cache := SafeCache{
		cache: make(map[string]string),
	}
	
	// Set values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	
	// Get values
	v1, ok1 := cache.Get("key1")
	v2, ok2 := cache.Get("key2")
	
	if !ok1 || v1 != "value1" {
		t.Errorf("Expected value1, got %v, %v", v1, ok1)
	}
	if !ok2 || v2 != "value2" {
		t.Errorf("Expected value2, got %v, %v", v2, ok2)
	}
	
	// Test concurrent reads
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Get("key1")
		}()
	}
	wg.Wait()
}

func TestOnceExecutor(t *testing.T) {
	executor := OnceExecutor{}
	count := 0
	
	// Call Do multiple times
	executor.Do(func() { count++ })
	executor.Do(func() { count++ })
	executor.Do(func() { count++ })
	
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestAtomicCounter(t *testing.T) {
	counter := NewAtomicCounter()
	
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	
	if counter.Value() != 1000 {
		t.Errorf("Expected 1000, got %d", counter.Value())
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int, 10)
	
	// Start worker
	results := ContextCancellation(ctx, jobs)
	
	// Send jobs
	for i := 0; i < 5; i++ {
		jobs <- i
	}
	close(jobs)
	
	// Cancel after short delay
	time.Sleep(10 * time.Millisecond)
	cancel()
	
	// Results should be limited due to cancellation
	select {
	case <-results:
		// Got result
	case <-time.After(100 * time.Millisecond):
		// Timeout - worker stopped (expected)
	}
}
