//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task: Implement the classic Producer-Consumer pattern using goroutines and channels.
// Producers generate data and send it to a channel, while consumers receive and process it.

// Job represents a unit of work
type Job struct {
	ID   int
	Data string
}

// Result represents the output of processing a job
type Result struct {
	JobID  int
	Output string
}

// Producer generates jobs and sends them to the jobs channel
func Producer(id int, jobs chan<- Job, numJobs int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < numJobs; i++ {
		job := Job{
			ID:   id*1000 + i,
			Data: fmt.Sprintf("Data from Producer %d, Job %d", id, i),
		}

		fmt.Printf("Producer %d: Created job %d\n", id, job.ID)
		jobs <- job

		// Simulate variable production time
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	fmt.Printf("Producer %d: Finished producing\n", id)
}

// Consumer processes jobs from the jobs channel and sends results
func Consumer(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("  Consumer %d: Processing job %d\n", id, job.ID)

		// Simulate processing time
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

		result := Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Processed: %s by Consumer %d", job.Data, id),
		}

		results <- result
	}

	fmt.Printf("  Consumer %d: Finished consuming\n", id)
}

// BufferedProducerConsumer demonstrates using buffered channels
func BufferedProducerConsumer() {
	fmt.Println("\n1. Buffered Producer-Consumer Pattern")
	fmt.Println("-".repeat(60))

	const (
		numProducers    = 2
		numConsumers    = 3
		jobsPerProducer = 5
		bufferSize      = 10
	)

	jobs := make(chan Job, bufferSize)
	results := make(chan Result, bufferSize)

	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	// Start consumers
	for i := 1; i <= numConsumers; i++ {
		consumerWg.Add(1)
		go Consumer(i, jobs, results, &consumerWg)
	}

	// Start producers
	for i := 1; i <= numProducers; i++ {
		producerWg.Add(1)
		go Producer(i, jobs, jobsPerProducer, &producerWg)
	}

	// Close jobs channel when all producers are done
	go func() {
		producerWg.Wait()
		close(jobs)
	}()

	// Close results channel when all consumers are done
	go func() {
		consumerWg.Wait()
		close(results)
	}()

	// Collect results
	resultCount := 0
	for result := range results {
		fmt.Printf("    Result: Job %d completed\n", result.JobID)
		resultCount++
	}

	fmt.Printf("\nTotal jobs processed: %d\n", resultCount)
}

// WorkerPool demonstrates a worker pool pattern
func WorkerPool() {
	fmt.Println("\n2. Worker Pool Pattern")
	fmt.Println("-".repeat(60))

	const (
		numWorkers = 3
		numTasks   = 10
	)

	tasks := make(chan int, numTasks)
	results := make(chan int, numTasks)

	// Start worker pool
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("Worker %d: Processing task %d\n", workerID, task)
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				results <- task * 2 // Simple processing: multiply by 2
			}
		}(i)
	}

	// Send tasks
	go func() {
		for i := 1; i <= numTasks; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	// Close results when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Println("\nResults:")
	for result := range results {
		fmt.Printf("  Task result: %d\n", result)
	}
}

// PipelinePattern demonstrates a pipeline processing pattern
func PipelinePattern() {
	fmt.Println("\n3. Pipeline Pattern")
	fmt.Println("-".repeat(60))

	// Stage 1: Generate numbers
	generate := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			for _, n := range nums {
				out <- n
			}
			close(out)
		}()
		return out
	}

	// Stage 2: Square numbers
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			for n := range in {
				out <- n * n
			}
			close(out)
		}()
		return out
	}

	// Stage 3: Add 10 to numbers
	addTen := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			for n := range in {
				out <- n + 10
			}
			close(out)
		}()
		return out
	}

	// Build pipeline
	numbers := generate(1, 2, 3, 4, 5)
	squared := square(numbers)
	final := addTen(squared)

	// Consume results
	fmt.Println("Pipeline: Generate -> Square -> Add 10")
	for result := range final {
		fmt.Printf("  Result: %d\n", result)
	}
}

func main() {
	fmt.Println("Producer-Consumer Pattern Demo")
	fmt.Println("=".repeat(60))

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Demo 1: Buffered Producer-Consumer
	BufferedProducerConsumer()

	time.Sleep(500 * time.Millisecond)

	// Demo 2: Worker Pool
	WorkerPool()

	time.Sleep(500 * time.Millisecond)

	// Demo 3: Pipeline Pattern
	PipelinePattern()

	fmt.Println("\nProducer-Consumer demo completed!")
}

// Helper function to repeat strings
func (s string) repeat(count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
