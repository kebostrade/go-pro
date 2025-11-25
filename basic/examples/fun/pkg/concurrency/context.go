package concurrency

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("operation timed out")
	// ErrCancelled is returned when an operation is cancelled
	ErrCancelled = errors.New("operation cancelled")
)

// WithTimeout executes a function with a timeout
func WithTimeout(timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- fn(ctx)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ErrTimeout
	}
}

// WithDeadline executes a function with a deadline
func WithDeadline(deadline time.Time, fn func(context.Context) error) error {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- fn(ctx)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ErrTimeout
	}
}

// WithCancel executes a function with cancellation support
func WithCancel(fn func(context.Context) error) (error, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error, 1)

	go func() {
		errChan <- fn(ctx)
	}()

	go func() {
		select {
		case <-ctx.Done():
			errChan <- ErrCancelled
		}
	}()

	return <-errChan, cancel
}

// ParallelWithTimeout executes multiple functions in parallel with a timeout
func ParallelWithTimeout(timeout time.Duration, fns ...func(context.Context) error) []error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errors := make([]error, len(fns))
	var wg sync.WaitGroup

	for i, fn := range fns {
		wg.Add(1)
		go func(index int, f func(context.Context) error) {
			defer wg.Done()

			errChan := make(chan error, 1)
			go func() {
				errChan <- f(ctx)
			}()

			select {
			case err := <-errChan:
				errors[index] = err
			case <-ctx.Done():
				errors[index] = ctx.Err()
			}
		}(i, fn)
	}

	wg.Wait()
	return errors
}

// RetryWithContext retries a function with exponential backoff
func RetryWithContext(ctx context.Context, maxAttempts int, initialDelay time.Duration, fn func() error) error {
	var lastErr error
	delay := initialDelay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Try the operation
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// Don't sleep after the last attempt
		if attempt < maxAttempts {
			select {
			case <-time.After(delay):
				delay *= 2 // Exponential backoff
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

// ContextValue represents a typed context value
type ContextValue[T any] struct {
	key string
}

// NewContextValue creates a new typed context value
func NewContextValue[T any](key string) *ContextValue[T] {
	return &ContextValue[T]{key: key}
}

// WithValue adds a value to the context
func (cv *ContextValue[T]) WithValue(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, cv.key, value)
}

// Value retrieves a value from the context
func (cv *ContextValue[T]) Value(ctx context.Context) (T, bool) {
	val := ctx.Value(cv.key)
	if val == nil {
		var zero T
		return zero, false
	}

	typedVal, ok := val.(T)
	return typedVal, ok
}

// MustValue retrieves a value from the context or panics
func (cv *ContextValue[T]) MustValue(ctx context.Context) T {
	val, ok := cv.Value(ctx)
	if !ok {
		panic(fmt.Sprintf("context value not found: %s", cv.key))
	}
	return val
}

// CancellableTask represents a task that can be cancelled
type CancellableTask struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
	err    error
	mu     sync.RWMutex
}

// NewCancellableTask creates a new cancellable task
func NewCancellableTask(fn func(context.Context) error) *CancellableTask {
	ctx, cancel := context.WithCancel(context.Background())

	task := &CancellableTask{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	go func() {
		defer close(task.done)
		task.mu.Lock()
		task.err = fn(ctx)
		task.mu.Unlock()
	}()

	return task
}

// Cancel cancels the task
func (t *CancellableTask) Cancel() {
	t.cancel()
}

// Wait waits for the task to complete
func (t *CancellableTask) Wait() error {
	<-t.done
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.err
}

// WaitWithTimeout waits for the task with a timeout
func (t *CancellableTask) WaitWithTimeout(timeout time.Duration) error {
	select {
	case <-t.done:
		t.mu.RLock()
		defer t.mu.RUnlock()
		return t.err
	case <-time.After(timeout):
		t.Cancel()
		return ErrTimeout
	}
}

// IsDone checks if the task is done
func (t *CancellableTask) IsDone() bool {
	select {
	case <-t.done:
		return true
	default:
		return false
	}
}

// TaskGroup manages a group of cancellable tasks
type TaskGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.Mutex
	errors []error
}

// NewTaskGroup creates a new task group
func NewTaskGroup(ctx context.Context) *TaskGroup {
	groupCtx, cancel := context.WithCancel(ctx)
	return &TaskGroup{
		ctx:    groupCtx,
		cancel: cancel,
		errors: make([]error, 0),
	}
}

// Go runs a function in the task group
func (tg *TaskGroup) Go(fn func(context.Context) error) {
	tg.wg.Add(1)

	go func() {
		defer tg.wg.Done()

		if err := fn(tg.ctx); err != nil {
			tg.mu.Lock()
			tg.errors = append(tg.errors, err)
			tg.mu.Unlock()

			// Cancel all other tasks on first error
			tg.cancel()
		}
	}()
}

// Wait waits for all tasks to complete
func (tg *TaskGroup) Wait() error {
	tg.wg.Wait()

	tg.mu.Lock()
	defer tg.mu.Unlock()

	if len(tg.errors) > 0 {
		return tg.errors[0]
	}

	return nil
}

// Errors returns all errors from the task group
func (tg *TaskGroup) Errors() []error {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	errorsCopy := make([]error, len(tg.errors))
	copy(errorsCopy, tg.errors)
	return errorsCopy
}

// Cancel cancels all tasks in the group
func (tg *TaskGroup) Cancel() {
	tg.cancel()
}
