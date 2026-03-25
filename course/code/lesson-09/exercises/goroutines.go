package exercises

import "time"

// Exercise 1: Basic goroutine
// Launch a goroutine that sends numbers 1-5 to a channel
func BasicGoroutine(ch chan int) {
	// TODO: Send 1, 2, 3, 4, 5 to the channel, then close it
	// Close channel when done
}

// Exercise 2: Buffered channel
// Create a buffered channel and demonstrate its behavior
func BufferedChannelExample() (sent int, received int) {
	// TODO: Create buffered channel with capacity 3
	// Send 3 values without blocking
	// Receive 3 values
	// Return sent and received counts
	return 0, 0
}

// Exercise 3: Unbuffered channel synchronization
// Use an unbuffered channel to synchronize goroutines
func UnbufferedSync(done chan bool) {
	// TODO: Do some work, then send true to done
	// This should block until someone receives
}

// Exercise 4: Select statement
// Implement a function that reads from two channels
// Return the first value received
func SelectFirst(ch1, ch2 <-chan int) int {
	// TODO: Use select to return first available value
	return 0
}

// Exercise 5: Select with timeout
// Implement a function with timeout using select
func SelectWithTimeout(ch <-chan int, timeout time.Duration) (int, bool) {
	// TODO: Use select to read from channel or timeout
	// Return value and whether it was received (not timed out)
	return 0, false
}

// Exercise 6: Channel direction
// Implement a function that only receives from a channel
func ReceiverOnly(ch <-chan string) string {
	// TODO: Receive and return value from channel
	return ""
}

// Exercise 7: Fan-out pattern
// Distribute work to multiple workers
func FanOut(jobs <-chan int, results chan<- int, numWorkers int) {
	// TODO: Launch numWorkers goroutines, each processing jobs
	// Send results to results channel
	// Close results when done
}

// Exercise 8: Fan-in pattern
// Combine multiple channels into one
func FanIn(channels ...<-chan int) <-chan int {
	// TODO: Combine all channels into one and return it
	return nil
}

// Exercise 9: Pipeline pattern
// Create a simple pipeline: generate -> square -> sum
func Generate(nums ...int) <-chan int {
	out := make(chan int)
	// TODO: Send numbers to out channel, close when done
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func Square(in <-chan int) <-chan int {
	out := make(chan int)
	// TODO: Square each number from in, send to out
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func Sum(in <-chan int) int {
	// TODO: Sum all values from in channel
	sum := 0
	for n := range in {
		sum += n
	}
	return sum
}

// Exercise 10: Context-like cancellation
// Implement a simple cancellation mechanism
func WithCancellation(done chan bool, f func()) {
	// TODO: Run f() in goroutine
	// Listen on done channel and cancel if signal received
}
