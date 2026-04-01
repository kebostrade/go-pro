// Package concurrency provides concurrency patterns including worker pool.
package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkItem represents a unit of work to be processed.
type WorkItem struct {
	ID      string
	Payload interface{}
}

// Result represents the result of processing a work item.
type Result struct {
	ID     string
	Output interface{}
	Error  error
}

// Processor is a function that processes a work item.
type Processor func(WorkItem) Result

// WorkerPool manages a pool of workers for concurrent processing.
type WorkerPool struct {
	workers   int
	jobCh     chan WorkItem
	resultCh  chan Result
	processor Processor
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	started   bool
	mu        sync.Mutex
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(workers int, processor Processor) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		jobCh:     make(chan WorkItem, workers*10),
		resultCh:  make(chan Result, workers*10),
		processor: processor,
	}
}

// Start starts the worker pool.
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.started {
		return
	}

	wp.ctx, wp.cancel = context.WithCancel(context.Background())

	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	wp.started = true
}

// worker is the main worker loop.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.jobCh:
			if !ok {
				return
			}
			result := wp.processor(job)
			select {
			case wp.resultCh <- result:
			case <-wp.ctx.Done():
				return
			}
		}
	}
}

// Submit submits a work item for processing.
func (wp *WorkerPool) Submit(item WorkItem) error {
	select {
	case wp.jobCh <- item:
		return nil
	default:
		return fmt.Errorf("worker pool job channel full")
	}
}

// SubmitAndWait submits a work item and waits for its result.
func (wp *WorkerPool) SubmitAndWait(ctx context.Context, item WorkItem) (Result, error) {
	if err := wp.Submit(item); err != nil {
		return Result{}, err
	}

	select {
	case result := <-wp.resultCh:
		return result, nil
	case <-ctx.Done():
		return Result{}, ctx.Err()
	}
}

// Stop stops the worker pool gracefully.
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if !wp.started {
		wp.mu.Unlock()
		return
	}
	wp.mu.Unlock()

	if wp.cancel != nil {
		wp.cancel()
	}
	close(wp.jobCh)
	wp.wg.Wait()
	close(wp.resultCh)
}

// Stats returns current worker pool statistics.
type Stats struct {
	Workers     int
	PendingJobs int
}

// Stats returns current statistics.
func (wp *WorkerPool) Stats() Stats {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	return Stats{
		Workers:     wp.workers,
		PendingJobs: len(wp.jobCh),
	}
}

// ImageProcessor is an example processor for image processing.
func ImageProcessor(item WorkItem) Result {
	// Simulate image processing
	time.Sleep(10 * time.Millisecond)
	return Result{
		ID:     item.ID,
		Output: fmt.Sprintf("processed-%v", item.Payload),
	}
}

// EmailProcessor is an example processor for email sending.
func EmailProcessor(item WorkItem) Result {
	// Simulate email sending
	time.Sleep(5 * time.Millisecond)
	return Result{
		ID:     item.ID,
		Output: fmt.Sprintf("sent-to-%v", item.Payload),
	}
}
