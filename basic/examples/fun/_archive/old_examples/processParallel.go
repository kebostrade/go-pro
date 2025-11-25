//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func processData(v int, workerID int) int {
	sleep := time.Duration(rand.Intn(10)) * time.Second
	fmt.Printf("Worker %d processing value %d (sleeping for %v)\n", workerID, v, sleep)
	time.Sleep(sleep)
	return v + 2
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	in := make(chan int)
	out := make(chan int, 10) // Buffered channel to prevent blocking

	go func() {
		for i := 0; i < 10; i++ {
			in <- i
		}
		close(in)
	}()

	start := time.Now()
	go processParallel(in, out, 5)

	for v := range out {
		fmt.Printf("Received result: %d\n", v)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

func processParallel(in <-chan int, out chan<- int, n int) {
	var wg sync.WaitGroup

	// Create worker pool
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for v := range in {
				out <- processData(v, id)
			}
		}(i)
	}

	// Wait for all workers to finish and close output channel
	go func() {
		wg.Wait()
		close(out)
	}()
}
