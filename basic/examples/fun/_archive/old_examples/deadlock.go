//go:build ignore

package main

import (
	"fmt"
	"sync"
)

// Gorutines are used to perform multiple tasks concurrently
// Gorutines are just functions that leave the main thread and run in te background and come back to join the main thread once the functions are finished too retunr value
// Gorutines do not stop the program flow and are non blocking
// Gorutines schedule tasks to be executed in parallel
// Deadlock is when gorutines are waiting for each other to finish and they are not able to finish because of a deadlock

func main() {
	var wg sync.WaitGroup
	nums := make(chan int)
	wg.Add(3)

	// Launch three goroutines to send values
	for i := 1; i <= 3; i++ {
		go func(val int) {
			fmt.Printf("Sending: %d\n", val)
			nums <- val
			fmt.Printf("Sent: %d\n", val)
			wg.Done()
		}(i)
	}

	// Close channel after all goroutines complete
	go func() {
		fmt.Println("Waiting for all sends to complete...")
		wg.Wait()
		fmt.Println("Closing channel...")
		close(nums)
	}()

	// Read until channel is closed
	fmt.Println("Starting to receive values:")
	for num := range nums {
		fmt.Printf("Received: %d\n", num)
	}
	fmt.Println("Channel closed, program ending")
}
