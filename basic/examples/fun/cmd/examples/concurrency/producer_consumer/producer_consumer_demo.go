package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/concurrency"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Producer-Consumer Pattern Demo")

	demo1BasicProducerConsumer()
	demo2WorkerPool()
	demo3Pipeline()
	demo4RateLimitedProducer()
}

func demo1BasicProducerConsumer() {
	utils.PrintSubHeader("1. Basic Producer-Consumer")

	ctx := context.Background()

	// Create producer-consumer pool
	pool := concurrency.NewProducerConsumerPool(
		2,  // 2 producers
		3,  // 3 consumers
		10, // buffer size
		// Generate function
		func(producerID, index int) string {
			return fmt.Sprintf("Data from Producer %d, Item %d", producerID, index)
		},
		// Process function
		func(data string) (string, error) {
			time.Sleep(50 * time.Millisecond) // Simulate processing
			return fmt.Sprintf("Processed: %s", data), nil
		},
	)

	fmt.Println("Starting 2 producers and 3 consumers...")
	fmt.Println("Each producer will create 5 jobs")

	start := time.Now()
	results, err := pool.Run(ctx, 5)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nCompleted %d jobs in %v\n", len(results), duration)
	fmt.Println("\nFirst 5 results:")
	for i := 0; i < 5 && i < len(results); i++ {
		fmt.Printf("  %s\n", results[i])
	}
}

func demo2WorkerPool() {
	utils.PrintSubHeader("2. Worker Pool Pattern")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create worker pool
	pool := concurrency.NewWorkerPool(4, func(n int) int {
		time.Sleep(100 * time.Millisecond) // Simulate work
		return n * n
	})

	// Start workers
	pool.Start(ctx)

	// Submit jobs
	fmt.Println("Submitting 10 jobs to worker pool with 4 workers...")
	go func() {
		for i := 1; i <= 10; i++ {
			pool.Submit(i)
		}
		pool.Close()
	}()

	// Collect results
	fmt.Println("\nResults:")
	for result := range pool.Results() {
		fmt.Printf("  %d\n", result)
	}
}

func demo3Pipeline() {
	utils.PrintSubHeader("3. Pipeline Pattern")

	// Create input channel
	input := make(chan int)

	// Stage 1: Double the number
	stage1 := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * 2
			}
		}()
		return out
	}

	// Stage 2: Add 10
	stage2 := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n + 10
			}
		}()
		return out
	}

	// Stage 3: Square the number
	stage3 := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	// Create pipeline
	output := concurrency.CreatePipeline(input, stage1, stage2, stage3)

	// Send input
	go func() {
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)
	}()

	// Collect results
	fmt.Println("Pipeline: (n * 2 + 10)²")
	fmt.Println("Input -> Double -> Add 10 -> Square")

	for result := range output {
		fmt.Printf("Result: %d\n", result)
	}
}

func demo4RateLimitedProducer() {
	utils.PrintSubHeader("4. Rate-Limited Producer")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create rate-limited producer (1 item per 200ms)
	producer := concurrency.NewRateLimitedProducer(
		200*time.Millisecond,
		func(i int) string {
			return fmt.Sprintf("Item %d", i)
		},
	)

	fmt.Println("Producing items at rate of 1 per 200ms...")
	fmt.Println("(Will produce 10 items)")

	start := time.Now()
	output := producer.Start(ctx, 10)

	count := 0
	for item := range output {
		count++
		elapsed := time.Since(start)
		fmt.Printf("[%v] %s\n", elapsed.Round(time.Millisecond), item)
	}

	totalDuration := time.Since(start)
	fmt.Printf("\nProduced %d items in %v\n", count, totalDuration)
	fmt.Printf("Average rate: %.2f items/second\n", float64(count)/totalDuration.Seconds())
}
