//go:build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// Add multiple goroutines to demonstrate concurrent execution
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate some work
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Goroutine %d completed\n", id)
		}(i)
	}

	fmt.Println("Main: waiting for goroutines to complete...")
	wg.Wait()
	fmt.Println("Main: all goroutines completed")
}
