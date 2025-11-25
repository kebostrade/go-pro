package test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/concurrency"
)

// TestRunConcurrent tests concurrent execution
func TestRunConcurrent(t *testing.T) {
	counter := int32(0)

	tasks := []func(){
		func() { atomic.AddInt32(&counter, 1) },
		func() { atomic.AddInt32(&counter, 1) },
		func() { atomic.AddInt32(&counter, 1) },
	}

	concurrency.RunConcurrent(tasks)

	if counter != 3 {
		t.Errorf("Expected counter to be 3, got %d", counter)
	}
}

// TestSafeCounter tests thread-safe counter
func TestSafeCounter(t *testing.T) {
	counter := concurrency.NewSafeCounter()

	// Increment concurrently
	tasks := make([]func(), 100)
	for i := 0; i < 100; i++ {
		tasks[i] = func() {
			counter.Increment()
		}
	}

	concurrency.RunConcurrent(tasks)

	if counter.Value() != 100 {
		t.Errorf("Expected counter to be 100, got %d", counter.Value())
	}

	// Test Add
	counter.Reset()
	counter.Add(50)
	if counter.Value() != 50 {
		t.Errorf("Expected counter to be 50, got %d", counter.Value())
	}
}

// TestSafeMap tests thread-safe map
func TestSafeMap(t *testing.T) {
	m := concurrency.NewSafeMap[string, int]()

	// Set values concurrently
	tasks := make([]func(), 10)
	for i := 0; i < 10; i++ {
		key := string(rune('a' + i))
		value := i
		tasks[i] = func() {
			m.Set(key, value)
		}
	}

	concurrency.RunConcurrent(tasks)

	if m.Len() != 10 {
		t.Errorf("Expected map length to be 10, got %d", m.Len())
	}

	// Test Get
	val, ok := m.Get("a")
	if !ok || val != 0 {
		t.Errorf("Expected to get value 0 for key 'a', got %d, %v", val, ok)
	}

	// Test Delete
	m.Delete("a")
	if m.Len() != 9 {
		t.Errorf("Expected map length to be 9 after delete, got %d", m.Len())
	}
}

// TestRateLimiter tests token bucket rate limiter
func TestRateLimiter(t *testing.T) {
	limiter := concurrency.NewRateLimiter(5, time.Second)
	defer limiter.Stop()

	// Should allow first 5 requests
	allowed := 0
	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			allowed++
		}
	}

	if allowed != 5 {
		t.Errorf("Expected 5 requests to be allowed, got %d", allowed)
	}

	// Test Wait
	done := make(chan bool)
	go func() {
		limiter.Wait()
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("Wait() timed out")
	}
}

// TestSlidingWindowRateLimiter tests sliding window rate limiter
func TestSlidingWindowRateLimiter(t *testing.T) {
	limiter := concurrency.NewSlidingWindowRateLimiter(3, time.Second)

	// Should allow first 3 requests
	for i := 0; i < 3; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// 4th request should be denied
	if limiter.Allow() {
		t.Error("4th request should be denied")
	}

	// After window expires, should allow again
	time.Sleep(1100 * time.Millisecond)
	if !limiter.Allow() {
		t.Error("Request after window should be allowed")
	}
}

// TestProducerConsumerPool tests producer-consumer pattern
func TestProducerConsumerPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := concurrency.NewProducerConsumerPool(
		2,  // producers
		3,  // consumers
		10, // buffer
		func(id, index int) int { return id*100 + index },
		func(data int) (int, error) { return data * 2, nil },
	)

	results, err := pool.Run(ctx, 5)
	if err != nil {
		t.Fatalf("Pool.Run() failed: %v", err)
	}

	expectedJobs := 2 * 5 // 2 producers * 5 jobs each
	if len(results) != expectedJobs {
		t.Errorf("Expected %d results, got %d", expectedJobs, len(results))
	}
}

// TestWithTimeout tests timeout functionality
func TestWithTimeout(t *testing.T) {
	// Should succeed
	err := concurrency.WithTimeout(1*time.Second, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should timeout
	err = concurrency.WithTimeout(100*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		return nil
	})

	if err != concurrency.ErrTimeout {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

// TestCancellableTask tests cancellable task
func TestCancellableTask(t *testing.T) {
	executed := false

	task := concurrency.NewCancellableTask(func(ctx context.Context) error {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
		executed = true
		return nil
	})

	// Cancel after 200ms
	time.Sleep(200 * time.Millisecond)
	task.Cancel()

	err := task.Wait()
	if err == nil {
		t.Error("Expected cancellation error")
	}

	if executed {
		t.Error("Task should not have completed")
	}
}

// TestTaskGroup tests task group
func TestTaskGroup(t *testing.T) {
	ctx := context.Background()
	group := concurrency.NewTaskGroup(ctx)

	counter := int32(0)

	// Add tasks
	for i := 0; i < 5; i++ {
		group.Go(func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if counter != 5 {
		t.Errorf("Expected counter to be 5, got %d", counter)
	}
}

// TestTaskGroupFailFast tests task group fail-fast behavior
func TestTaskGroupFailFast(t *testing.T) {
	ctx := context.Background()
	group := concurrency.NewTaskGroup(ctx)

	// Add failing task
	group.Go(func(ctx context.Context) error {
		return errors.New("task failed")
	})

	// Add task that should be cancelled
	cancelled := false
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			cancelled = true
			return ctx.Err()
		case <-time.After(5 * time.Second):
			return nil
		}
	})

	err := group.Wait()
	if err == nil {
		t.Error("Expected error from task group")
	}

	if !cancelled {
		t.Error("Expected second task to be cancelled")
	}
}

// TestParallelMap tests parallel map
func TestParallelMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	results := concurrency.ParallelMap(input, 2, func(n int) int {
		return n * 2
	})

	expected := []int{2, 4, 6, 8, 10}
	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

// TestParallelFilter tests parallel filter
func TestParallelFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	results := concurrency.ParallelFilter(input, 2, func(n int) bool {
		return n%2 == 0
	})

	expected := []int{2, 4, 6, 8, 10}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}
}

// TestParallelSum tests parallel sum
func TestParallelSum(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sum := concurrency.ParallelSum(input, 2)
	expected := 55

	if sum != expected {
		t.Errorf("Expected sum to be %d, got %d", expected, sum)
	}
}

// Benchmarks

func BenchmarkSafeCounter(b *testing.B) {
	counter := concurrency.NewSafeCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkSafeMap(b *testing.B) {
	m := concurrency.NewSafeMap[int, int]()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(i, i)
			m.Get(i)
			i++
		}
	})
}

func BenchmarkParallelMap(b *testing.B) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		concurrency.ParallelMap(input, 4, func(n int) int {
			return n * 2
		})
	}
}
