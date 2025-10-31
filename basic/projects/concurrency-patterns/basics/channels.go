package main

import (
	"fmt"
	"time"
)

/*
CHANNELS - Communication Between Goroutines

Channels are typed conduits for sending and receiving values.
They enable safe communication between goroutines.

Key Concepts:
- Send: ch <- value
- Receive: value := <-ch
- Close: close(ch)
- Buffered vs Unbuffered
- Select statement
*/

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                              ║")
	fmt.Println("║                    📡 Channel Patterns                       ║")
	fmt.Println("║                                                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 1. Unbuffered Channel
	unbufferedChannel()

	// 2. Buffered Channel
	bufferedChannel()

	// 3. Channel Direction
	channelDirection()

	// 4. Range and Close
	rangeAndClose()

	// 5. Select Statement
	selectStatement()

	// 6. Timeout Pattern
	timeoutPattern()

	// 7. Non-blocking Operations
	nonBlockingOperations()
}

// 1. Unbuffered Channel
func unbufferedChannel() {
	fmt.Println("\n1️⃣  UNBUFFERED CHANNEL")
	fmt.Println("───────────────────────────────────────────────────────────────")

	ch := make(chan string)

	// Send in goroutine (would block if not in goroutine)
	go func() {
		ch <- "Hello from goroutine"
	}()

	// Receive (blocks until value is sent)
	msg := <-ch
	fmt.Printf("Received: %s\n", msg)
}

// 2. Buffered Channel
func bufferedChannel() {
	fmt.Println("\n2️⃣  BUFFERED CHANNEL")
	fmt.Println("───────────────────────────────────────────────────────────────")

	// Buffer size of 2
	ch := make(chan int, 2)

	// Can send 2 values without blocking
	ch <- 1
	ch <- 2
	fmt.Println("Sent 2 values to buffered channel")

	// Receive values
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("Received: %d\n", <-ch)
}

// 3. Channel Direction
func channelDirection() {
	fmt.Println("\n3️⃣  CHANNEL DIRECTION")
	fmt.Println("───────────────────────────────────────────────────────────────")

	ch := make(chan string, 1)

	// Send-only channel
	go sendOnly(ch)

	// Receive-only channel
	receiveOnly(ch)
}

func sendOnly(ch chan<- string) {
	ch <- "Message from send-only channel"
}

func receiveOnly(ch <-chan string) {
	msg := <-ch
	fmt.Printf("Received: %s\n", msg)
}

// 4. Range and Close
func rangeAndClose() {
	fmt.Println("\n4️⃣  RANGE AND CLOSE")
	fmt.Println("───────────────────────────────────────────────────────────────")

	ch := make(chan int, 5)

	// Send values
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch) // Close when done sending
	}()

	// Range over channel (stops when closed)
	for num := range ch {
		fmt.Printf("Received: %d\n", num)
	}
}

// 5. Select Statement
func selectStatement() {
	fmt.Println("\n5️⃣  SELECT STATEMENT")
	fmt.Println("───────────────────────────────────────────────────────────────")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// Select waits on multiple channels
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received %s\n", msg2)
		}
	}
}

// 6. Timeout Pattern
func timeoutPattern() {
	fmt.Println("\n6️⃣  TIMEOUT PATTERN")
	fmt.Println("───────────────────────────────────────────────────────────────")

	ch := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "result"
	}()

	select {
	case res := <-ch:
		fmt.Printf("Received: %s\n", res)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("⏱️  Timeout! Operation took too long")
	}
}

// 7. Non-blocking Operations
func nonBlockingOperations() {
	fmt.Println("\n7️⃣  NON-BLOCKING OPERATIONS")
	fmt.Println("───────────────────────────────────────────────────────────────")

	ch := make(chan int, 1)

	// Non-blocking send
	select {
	case ch <- 42:
		fmt.Println("✅ Sent value")
	default:
		fmt.Println("❌ Channel full, cannot send")
	}

	// Non-blocking receive
	select {
	case val := <-ch:
		fmt.Printf("✅ Received: %d\n", val)
	default:
		fmt.Println("❌ No value available")
	}
}

/*
KEY TAKEAWAYS:

1. Unbuffered channels block until both sender and receiver are ready
2. Buffered channels block only when buffer is full (send) or empty (receive)
3. Always close channels when done sending (sender's responsibility)
4. Range over channels to receive all values until closed
5. Select enables waiting on multiple channels
6. Use select with time.After for timeouts
7. Use select with default for non-blocking operations

CHANNEL PATTERNS:
- Producer-Consumer
- Fan-out/Fan-in
- Pipeline
- Worker Pool
- Pub/Sub
*/

