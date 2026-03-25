package main

import (
	"fmt"
	"lesson-09/exercises"
	"time"
)

func main() {
	fmt.Println("=== Lesson 09: Goroutines and Channels ===")
	fmt.Println()

	// Exercise 1: Basic Goroutine
	fmt.Println("1. Basic Goroutine:")
	ch := make(chan int)
	go exercises.BasicGoroutine(ch)
	for v := range ch {
		fmt.Printf("   Received: %d\n", v)
	}
	fmt.Println()

	// Exercise 2: Buffered Channel
	fmt.Println("2. Buffered Channel:")
	sent, received := exercises.BufferedChannelExample()
	fmt.Printf("   Sent: %d, Received: %d\n", sent, received)
	fmt.Println()

	// Exercise 3: Unbuffered Sync
	fmt.Println("3. Unbuffered Channel Sync:")
	done := make(chan bool)
	go exercises.UnbufferedSync(done)
	<-done
	fmt.Println("   Done!")
	fmt.Println()

	// Exercise 4: Select First
	fmt.Println("4. Select First:")
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() { ch1 <- 100 }()
	go func() { time.Sleep(10 * time.Millisecond); ch2 <- 200 }()
	result := exercises.SelectFirst(ch1, ch2)
	fmt.Printf("   First value: %d\n", result)
	fmt.Println()

	// Exercise 5: Select with Timeout
	fmt.Println("5. Select with Timeout:")
	ch3 := make(chan int)
	go func() { time.Sleep(50 * time.Millisecond); ch3 <- 300 }()
	val, ok := exercises.SelectWithTimeout(ch3, 100*time.Millisecond)
	fmt.Printf("   Value: %d, OK: %t\n", val, ok)
	fmt.Println()

	// Exercise 9: Pipeline
	fmt.Println("9. Pipeline Pattern:")
	pipeline := exercises.Square(exercises.Generate(1, 2, 3, 4, 5))
	sum := exercises.Sum(pipeline)
	fmt.Printf("   Sum of squares 1-5: %d\n", sum)
	fmt.Println()

	fmt.Println("=== All exercises completed! ===")
}
