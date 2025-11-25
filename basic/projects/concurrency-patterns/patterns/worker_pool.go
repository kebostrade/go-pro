package patterns

import (
	"fmt"
	"sync"
	"time"
)

/*
WORKER POOL PATTERN

A worker pool is a collection of goroutines that process jobs from a shared queue.
This pattern limits concurrency and reuses goroutines for efficiency.

Benefits:
- Controlled concurrency
- Resource management
- Better performance than creating goroutines per task
*/

// Job represents a unit of work
type Job struct {
	ID      int
	Data    interface{}
	Process func(interface{}) (interface{}, error)
}

// Result represents the outcome of a job
type Result struct {
	JobID  int
	Output interface{}
	Error  error
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers, jobQueueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, jobQueueSize),
		results:    make(chan Result, jobQueueSize),
	}
}

// Start launches the worker pool
func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker processes jobs from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for job := range wp.jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job.ID)

		// Process the job
		output, err := job.Process(job.Data)

		// Send result
		wp.results <- Result{
			JobID:  job.ID,
			Output: output,
			Error:  err,
		}
	}

	fmt.Printf("Worker %d finished\n", id)
}

// Submit adds a job to the queue
func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

// Close closes the job queue and waits for workers to finish
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

// Example: Image Processing Worker Pool
type ImageJob struct {
	Filename string
	Width    int
	Height   int
}

func processImage(data interface{}) (interface{}, error) {
	img := data.(ImageJob)

	// Simulate image processing
	time.Sleep(100 * time.Millisecond)

	return fmt.Sprintf("Processed %s (%dx%d)", img.Filename, img.Width, img.Height), nil
}

// Example: Data Processing Worker Pool
func processData(data interface{}) (interface{}, error) {
	num := data.(int)

	// Simulate computation
	time.Sleep(50 * time.Millisecond)

	return num * num, nil
}

// DemoWorkerPool demonstrates the worker pool pattern
func DemoWorkerPool() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                              ║")
	fmt.Println("║                  👷 Worker Pool Pattern                      ║")
	fmt.Println("║                                                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create worker pool with 3 workers
	pool := NewWorkerPool(3, 10)

	// Start workers
	pool.Start()

	// Submit jobs
	numJobs := 10
	go func() {
		for i := 1; i <= numJobs; i++ {
			pool.Submit(Job{
				ID:   i,
				Data: i,
				Process: processData,
			})
		}
		pool.Close()
	}()

	// Collect results
	for result := range pool.Results() {
		if result.Error != nil {
			fmt.Printf("❌ Job %d failed: %v\n", result.JobID, result.Error)
		} else {
			fmt.Printf("✅ Job %d result: %v\n", result.JobID, result.Output)
		}
	}

	fmt.Println("\n✅ All jobs completed")
}

// Advanced: Worker Pool with Context
type ContextWorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
	done       chan struct{}
}

func NewContextWorkerPool(numWorkers int) *ContextWorkerPool {
	return &ContextWorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, numWorkers*2),
		results:    make(chan Result, numWorkers*2),
		done:       make(chan struct{}),
	}
}

func (cwp *ContextWorkerPool) Start() {
	for i := 1; i <= cwp.numWorkers; i++ {
		cwp.wg.Add(1)
		go cwp.worker(i)
	}
}

func (cwp *ContextWorkerPool) worker(id int) {
	defer cwp.wg.Done()

	for {
		select {
		case job, ok := <-cwp.jobs:
			if !ok {
				return
			}

			output, err := job.Process(job.Data)
			cwp.results <- Result{
				JobID:  job.ID,
				Output: output,
				Error:  err,
			}

		case <-cwp.done:
			return
		}
	}
}

func (cwp *ContextWorkerPool) Submit(job Job) {
	cwp.jobs <- job
}

func (cwp *ContextWorkerPool) Shutdown() {
	close(cwp.done)
	cwp.wg.Wait()
	close(cwp.results)
}

func (cwp *ContextWorkerPool) Results() <-chan Result {
	return cwp.results
}

/*
KEY TAKEAWAYS:

1. Worker pools limit concurrency and reuse goroutines
2. Use buffered channels for job and result queues
3. Close job channel when done submitting
4. Use WaitGroup to wait for workers to finish
5. Separate submission and result collection

WHEN TO USE:
- Processing large number of tasks
- Need to limit concurrent operations
- Want to reuse goroutines
- Need backpressure control

VARIATIONS:
- Dynamic worker pool (add/remove workers)
- Priority queue
- Rate-limited pool
- Context-aware pool
*/

