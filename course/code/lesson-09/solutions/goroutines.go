package exercises

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Exercise 1: Basic goroutine
func BasicGoroutine(ch chan int) {
	defer close(ch)
	for i := 1; i <= 5; i++ {
		ch <- i
	}
}

// Exercise 2: Buffered channel
func BufferedChannelExample() (sent int, received int) {
	ch := make(chan int, 3)
	
	// Send 3 values without blocking
	ch <- 1
	ch <- 2
	ch <- 3
	sent = 3
	
	// Receive 3 values
	<-ch
	<-ch
	<-ch
	received = 3
	
	return sent, received
}

// Exercise 3: Unbuffered channel synchronization
func UnbufferedSync(done chan bool) {
	// Simulate work
	time.Sleep(10 * time.Millisecond)
	done <- true
}

// Exercise 4: Select statement
func SelectFirst(ch1, ch2 <-chan int) int {
	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v
	}
}

// Exercise 5: Select with timeout
func SelectWithTimeout(ch <-chan int, timeout time.Duration) (int, bool) {
	select {
	case v, ok := <-ch:
		return v, ok
	case <-time.After(timeout):
		return 0, false
	}
}

// Exercise 6: Channel direction
func ReceiverOnly(ch <-chan string) string {
	return <-ch
}

// Exercise 7: Fan-out pattern
func FanOut(jobs <-chan int, results chan<- int, numWorkers int) {
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
	
	go func() {
		wg.Wait()
		close(results)
	}()
}

// Exercise 8: Fan-in pattern
func FanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	
	output := func(ch <-chan int) {
		defer wg.Done()
		for n := range ch {
			out <- n
		}
	}
	
	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}
	
	go func() {
		wg.Wait()
		close(out)
	}()
	
	return out
}

// Exercise 9: Pipeline pattern
func Generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func Sum(in <-chan int) int {
	sum := 0
	for n := range in {
		sum += n
	}
	return sum
}

// Exercise 10: Context-like cancellation
func WithCancellation(done chan bool, f func()) {
	go func() {
		f()
	}()
	
	select {
	case <-done:
		// Cancellation signal received
	case <-context.Background().Done():
	}
}
