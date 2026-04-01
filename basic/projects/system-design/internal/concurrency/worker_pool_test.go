package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestNewWorkerPool(t *testing.T) {
	pool := NewWorkerPool(4, ImageProcessor)

	if pool == nil {
		t.Error("NewWorkerPool returned nil")
	}

	if pool.workers != 4 {
		t.Errorf("Workers = %d, want 4", pool.workers)
	}
}

func TestWorkerPool_StartStop(t *testing.T) {
	pool := NewWorkerPool(4, ImageProcessor)

	// Start should not block
	pool.Start()

	// Stop should not block
	pool.Stop()
}

func TestWorkerPool_Submit(t *testing.T) {
	pool := NewWorkerPool(2, ImageProcessor)
	pool.Start()
	defer pool.Stop()

	item := WorkItem{ID: "test-1", Payload: "test.jpg"}
	err := pool.Submit(item)

	if err != nil {
		t.Errorf("Submit failed: %v", err)
	}

	// Wait for result
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	select {
	case result := <-pool.resultCh:
		if result.ID != "test-1" {
			t.Errorf("Result ID = %s, want test-1", result.ID)
		}
	case <-ctx.Done():
		t.Error("Timed out waiting for result")
	}
}

func TestWorkerPool_SubmitAndWait(t *testing.T) {
	pool := NewWorkerPool(2, ImageProcessor)
	pool.Start()
	defer pool.Stop()

	item := WorkItem{ID: "test-2", Payload: "test2.jpg"}
	ctx := context.Background()

	result, err := pool.SubmitAndWait(ctx, item)
	if err != nil {
		t.Errorf("SubmitAndWait failed: %v", err)
	}

	if result.ID != "test-2" {
		t.Errorf("Result ID = %s, want test-2", result.ID)
	}
}

func TestWorkerPool_Stats(t *testing.T) {
	pool := NewWorkerPool(2, ImageProcessor)

	stats := pool.Stats()
	if stats.Workers != 2 {
		t.Errorf("Stats Workers = %d, want 2", stats.Workers)
	}

	pool.Start()
	defer pool.Stop()

	// Submit some items
	for i := 0; i < 5; i++ {
		pool.Submit(WorkItem{ID: "test", Payload: nil})
	}

	stats = pool.Stats()
	if stats.PendingJobs != 5 {
		t.Errorf("Stats PendingJobs = %d, want 5", stats.PendingJobs)
	}
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	pool := NewWorkerPool(2, func(item WorkItem) Result {
		time.Sleep(100 * time.Millisecond)
		return Result{ID: item.ID}
	})
	pool.Start()
	defer pool.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := pool.SubmitAndWait(ctx, WorkItem{ID: "test", Payload: nil})
	if err == nil {
		t.Error("Expected context deadline exceeded error")
	}
}

func TestEmailProcessor(t *testing.T) {
	item := WorkItem{ID: "email-1", Payload: "user@example.com"}
	result := EmailProcessor(item)

	if result.Error != nil {
		t.Errorf("EmailProcessor failed: %v", result.Error)
	}

	if result.Output != "sent-to-user@example.com" {
		t.Errorf("Output = %v, want sent-to-user@example.com", result.Output)
	}
}
