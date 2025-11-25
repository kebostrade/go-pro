package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job represents a unit of work
type Job[T any] struct {
	ID   int
	Data T
}

// Result represents the output of processing a job
type Result[T any] struct {
	JobID  int
	Output T
	Error  error
}

// Producer generates jobs and sends them to the jobs channel
type Producer[T any] struct {
	ID       int
	jobs     chan<- Job[T]
	generate func(id, index int) T
}

// NewProducer creates a new producer
func NewProducer[T any](id int, jobs chan<- Job[T], generate func(id, index int) T) *Producer[T] {
	return &Producer[T]{
		ID:       id,
		jobs:     jobs,
		generate: generate,
	}
}

// Produce generates numJobs jobs
func (p *Producer[T]) Produce(ctx context.Context, numJobs int) error {
	for i := 0; i < numJobs; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			job := Job[T]{
				ID:   p.ID*1000 + i,
				Data: p.generate(p.ID, i),
			}

			select {
			case p.jobs <- job:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}

// Consumer processes jobs from the jobs channel
type Consumer[T, R any] struct {
	ID      int
	jobs    <-chan Job[T]
	results chan<- Result[R]
	process func(T) (R, error)
}

// NewConsumer creates a new consumer
func NewConsumer[T, R any](id int, jobs <-chan Job[T], results chan<- Result[R], process func(T) (R, error)) *Consumer[T, R] {
	return &Consumer[T, R]{
		ID:      id,
		jobs:    jobs,
		results: results,
		process: process,
	}
}

// Consume processes jobs until the jobs channel is closed
func (c *Consumer[T, R]) Consume(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job, ok := <-c.jobs:
			if !ok {
				return nil // Channel closed
			}

			output, err := c.process(job.Data)
			result := Result[R]{
				JobID:  job.ID,
				Output: output,
				Error:  err,
			}

			select {
			case c.results <- result:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// ProducerConsumerPool manages a pool of producers and consumers
type ProducerConsumerPool[T, R any] struct {
	numProducers int
	numConsumers int
	bufferSize   int
	jobs         chan Job[T]
	results      chan Result[R]
	generate     func(id, index int) T
	process      func(T) (R, error)
}

// NewProducerConsumerPool creates a new producer-consumer pool
func NewProducerConsumerPool[T, R any](
	numProducers, numConsumers, bufferSize int,
	generate func(id, index int) T,
	process func(T) (R, error),
) *ProducerConsumerPool[T, R] {
	return &ProducerConsumerPool[T, R]{
		numProducers: numProducers,
		numConsumers: numConsumers,
		bufferSize:   bufferSize,
		jobs:         make(chan Job[T], bufferSize),
		results:      make(chan Result[R], bufferSize),
		generate:     generate,
		process:      process,
	}
}

// Run starts the producer-consumer pool
func (p *ProducerConsumerPool[T, R]) Run(ctx context.Context, jobsPerProducer int) ([]Result[R], error) {
	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	// Start consumers
	for i := 1; i <= p.numConsumers; i++ {
		consumerWg.Add(1)
		consumer := NewConsumer(i, p.jobs, p.results, p.process)
		go func() {
			defer consumerWg.Done()
			consumer.Consume(ctx)
		}()
	}

	// Start producers
	for i := 1; i <= p.numProducers; i++ {
		producerWg.Add(1)
		producer := NewProducer(i, p.jobs, p.generate)
		go func() {
			defer producerWg.Done()
			producer.Produce(ctx, jobsPerProducer)
		}()
	}

	// Close jobs channel when all producers are done
	go func() {
		producerWg.Wait()
		close(p.jobs)
	}()

	// Close results channel when all consumers are done
	go func() {
		consumerWg.Wait()
		close(p.results)
	}()

	// Collect results
	results := make([]Result[R], 0)
	for result := range p.results {
		results = append(results, result)
	}

	return results, nil
}

// WorkerPool implements a worker pool pattern
type WorkerPool[T, R any] struct {
	workers int
	jobs    chan T
	results chan R
	process func(T) R
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool[T, R any](workers int, process func(T) R) *WorkerPool[T, R] {
	return &WorkerPool[T, R]{
		workers: workers,
		jobs:    make(chan T),
		results: make(chan R),
		process: process,
	}
}

// Start starts the worker pool
func (w *WorkerPool[T, R]) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < w.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-w.jobs:
					if !ok {
						return
					}
					result := w.process(job)
					select {
					case w.results <- result:
					case <-ctx.Done():
						return
					}
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(w.results)
	}()
}

// Submit submits a job to the worker pool
func (w *WorkerPool[T, R]) Submit(job T) {
	w.jobs <- job
}

// Close closes the jobs channel
func (w *WorkerPool[T, R]) Close() {
	close(w.jobs)
}

// Results returns the results channel
func (w *WorkerPool[T, R]) Results() <-chan R {
	return w.results
}

// PipelineStage represents a stage in a processing pipeline
type PipelineStage[T, R any] func(<-chan T) <-chan R

// CreatePipeline creates a processing pipeline
func CreatePipeline[T any](input <-chan T, stages ...func(<-chan T) <-chan T) <-chan T {
	current := input
	for _, stage := range stages {
		current = stage(current)
	}
	return current
}

// BufferedPipeline creates a pipeline with buffered channels
func BufferedPipeline[T any](bufferSize int, input <-chan T, stages ...func(<-chan T) <-chan T) <-chan T {
	current := input
	for _, stage := range stages {
		buffered := make(chan T, bufferSize)
		go func(in <-chan T, out chan<- T, s func(<-chan T) <-chan T) {
			defer close(out)
			result := s(in)
			for val := range result {
				out <- val
			}
		}(current, buffered, stage)
		current = buffered
	}
	return current
}

// RateLimitedProducer produces items at a limited rate
type RateLimitedProducer[T any] struct {
	rate     time.Duration
	generate func(int) T
	output   chan T
}

// NewRateLimitedProducer creates a new rate-limited producer
func NewRateLimitedProducer[T any](rate time.Duration, generate func(int) T) *RateLimitedProducer[T] {
	return &RateLimitedProducer[T]{
		rate:     rate,
		generate: generate,
		output:   make(chan T),
	}
}

// Start starts producing items
func (p *RateLimitedProducer[T]) Start(ctx context.Context, count int) <-chan T {
	go func() {
		defer close(p.output)
		ticker := time.NewTicker(p.rate)
		defer ticker.Stop()

		for i := 0; i < count; i++ {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				item := p.generate(i)
				select {
				case p.output <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return p.output
}

// String returns a string representation
func (j Job[T]) String() string {
	return fmt.Sprintf("Job{ID: %d, Data: %v}", j.ID, j.Data)
}

// String returns a string representation
func (r Result[T]) String() string {
	if r.Error != nil {
		return fmt.Sprintf("Result{JobID: %d, Error: %v}", r.JobID, r.Error)
	}
	return fmt.Sprintf("Result{JobID: %d, Output: %v}", r.JobID, r.Output)
}
