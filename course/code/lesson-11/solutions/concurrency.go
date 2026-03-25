package exercises

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Exercise 1: Worker Pool
func WorkerPool(jobs <-chan int, results chan<- int, numWorkers int) {
	var wg sync.WaitGroup
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- job * 2
			}
		}()
	}
	
	wg.Wait()
}

// Exercise 4: Pipeline
func PipelineDemo(input []int) int {
	// Stage 1: Generate
	numChan := make(chan int)
	go func() {
		defer close(numChan)
		for _, n := range input {
			numChan <- n
		}
	}()
	
	// Stage 2: Square
	squareChan := make(chan int)
	go func() {
		defer close(squareChan)
		for n := range numChan {
			squareChan <- n * n
		}
	}()
	
	// Stage 3: Filter even
	filterChan := make(chan int)
	go func() {
		defer close(filterChan)
		for n := range squareChan {
			if n%2 == 0 {
				filterChan <- n
			}
		}
	}()
	
	// Stage 4: Sum
	sum := 0
	for n := range filterChan {
		sum += n
	}
	
	return sum
}

// Exercise 5: Context cancellation
func ContextCancellation(ctx context.Context, jobs <-chan int) <-chan int {
	results := make(chan int)
	
	go func() {
		defer close(results)
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobs:
				if !ok {
					return
				}
				// Process job
				results <- job * 2
			}
		}
	}()
	
	return results
}

// Exercise 6: WaitGroup
func WaitGroupDemo(tasks []func()) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	
	for _, task := range tasks {
		go func(t func()) {
			defer wg.Done()
			t()
		}(task)
	}
	
	wg.Wait()
}

// Exercise 7: SafeCounter
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// Exercise 8: SafeCache
type SafeCache struct {
	mu    sync.RWMutex
	cache map[string]string
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		cache: make(map[string]string),
	}
}

func (c *SafeCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.cache[key]
	return val, ok
}

func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

// Exercise 9: Once
func (o *OnceExecutor) Do(fn func()) {
	o.once.Do(fn)
}

// Exercise 10: Latch
type Latch struct {
	mu    sync.Mutex
	count int
	total int
	cond  *sync.Cond
}

func NewLatch(total int) *Latch {
	l := &Latch{total: total, count: total}
	l.cond = sync.NewCond(&l.mu)
	return l
}

func (l *Latch) Done() {
	l.mu.Lock()
	l.count--
	if l.count == 0 {
		l.cond.Broadcast()
	}
	l.mu.Unlock()
}

func (l *Latch) Wait() {
	l.mu.Lock()
	for l.count > 0 {
		l.cond.Wait()
	}
	l.mu.Unlock()
}

// Exercise 11: Atomic
type AtomicCounter struct {
	value int64
}

func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{}
}

func (c *AtomicCounter) Increment() {
	c.value++
}

func (c *AtomicCounter) Value() int64 {
	return c.value
}
