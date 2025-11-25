package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RunConcurrent executes multiple tasks concurrently and waits for all to complete
// Time complexity: O(max(task durations))
func RunConcurrent(tasks []func()) {
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			t()
		}(task)
	}

	wg.Wait()
}

// RunConcurrentWithContext executes tasks with context support for cancellation
func RunConcurrentWithContext(ctx context.Context, tasks []func(context.Context)) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t func(context.Context)) {
			defer wg.Done()

			// Check if context is already cancelled
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
			}

			t(ctx)
		}(task)
	}

	wg.Wait()
	close(errChan)

	// Return first error if any
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// RunWithTimeout executes a function with a timeout
func RunWithTimeout(timeout time.Duration, fn func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- fn()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// FanOut distributes work across multiple workers
// Returns a channel that receives results from all workers
func FanOut[T any, R any](input []T, workers int, process func(T) R) []R {
	jobs := make(chan T, len(input))
	results := make(chan R, len(input))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- process(job)
			}
		}()
	}

	// Send jobs
	for _, item := range input {
		jobs <- item
	}
	close(jobs)

	// Wait for workers and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	output := make([]R, 0, len(input))
	for result := range results {
		output = append(output, result)
	}

	return output
}

// FanIn merges multiple channels into a single channel
func FanIn[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan T) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}

	// Close output channel when all inputs are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Pipeline creates a processing pipeline with multiple stages
func Pipeline[T any](input <-chan T, stages ...func(<-chan T) <-chan T) <-chan T {
	current := input
	for _, stage := range stages {
		current = stage(current)
	}
	return current
}

// Broadcast sends a value to multiple channels
func Broadcast[T any](value T, channels ...chan<- T) {
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c chan<- T) {
			defer wg.Done()
			c <- value
		}(ch)
	}

	wg.Wait()
}

// Semaphore implements a counting semaphore using channels
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore creates a new semaphore with the given capacity
func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, capacity),
	}
}

// Acquire acquires a semaphore slot (blocks if full)
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

// Release releases a semaphore slot
func (s *Semaphore) Release() {
	<-s.sem
}

// TryAcquire attempts to acquire a slot without blocking
// Returns true if successful, false otherwise
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// Barrier synchronizes multiple goroutines at a specific point
type Barrier struct {
	count int
	ch    chan struct{}
	mu    sync.Mutex
}

// NewBarrier creates a new barrier for n goroutines
func NewBarrier(n int) *Barrier {
	return &Barrier{
		count: n,
		ch:    make(chan struct{}),
	}
}

// Wait blocks until all goroutines have called Wait
func (b *Barrier) Wait() {
	b.mu.Lock()
	b.count--

	if b.count == 0 {
		close(b.ch)
		b.mu.Unlock()
		return
	}

	ch := b.ch
	b.mu.Unlock()
	<-ch
}

// Once ensures a function is executed only once across multiple goroutines
type Once struct {
	done uint32
	mu   sync.Mutex
}

// Do executes the function f only once
func (o *Once) Do(f func()) {
	if o.done == 0 {
		o.mu.Lock()
		defer o.mu.Unlock()
		if o.done == 0 {
			defer func() { o.done = 1 }()
			f()
		}
	}
}

// SafeCounter is a thread-safe counter
type SafeCounter struct {
	mu    sync.RWMutex
	value int64
}

// NewSafeCounter creates a new safe counter
func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

// Increment increments the counter by 1
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Add adds delta to the counter
func (c *SafeCounter) Add(delta int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

// Value returns the current value
func (c *SafeCounter) Value() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Reset resets the counter to 0
func (c *SafeCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

// SafeMap is a thread-safe map
type SafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewSafeMap creates a new safe map
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

// Set sets a key-value pair
func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get retrieves a value by key
func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

// Delete removes a key-value pair
func (m *SafeMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len returns the number of items
func (m *SafeMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Keys returns all keys
func (m *SafeMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// ForEach applies a function to each key-value pair
func (m *SafeMap[K, V]) ForEach(fn func(K, V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.data {
		fn(k, v)
	}
}

// String returns a string representation
func (m *SafeMap[K, V]) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("SafeMap%v", m.data)
}
