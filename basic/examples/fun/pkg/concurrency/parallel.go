package concurrency

import (
	"context"
	"runtime"
	"sync"
)

// ParallelMap applies a function to each element in parallel
// Returns results in the same order as input
func ParallelMap[T, R any](input []T, workers int, fn func(T) R) []R {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	results := make([]R, len(input))
	jobs := make(chan struct {
		index int
		value T
	}, len(input))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results[job.index] = fn(job.value)
			}
		}()
	}

	// Send jobs
	for i, val := range input {
		jobs <- struct {
			index int
			value T
		}{i, val}
	}
	close(jobs)

	wg.Wait()
	return results
}

// ParallelMapWithContext applies a function to each element in parallel with context
func ParallelMapWithContext[T, R any](ctx context.Context, input []T, workers int, fn func(context.Context, T) (R, error)) ([]R, error) {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	results := make([]R, len(input))
	jobs := make(chan struct {
		index int
		value T
	}, len(input))

	errChan := make(chan error, workers)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				default:
					result, err := fn(ctx, job.value)
					if err != nil {
						errChan <- err
						return
					}
					results[job.index] = result
				}
			}
		}()
	}

	// Send jobs
	for i, val := range input {
		select {
		case <-ctx.Done():
			close(jobs)
			return nil, ctx.Err()
		case jobs <- struct {
			index int
			value T
		}{i, val}:
		}
	}
	close(jobs)

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ParallelFilter filters elements in parallel
func ParallelFilter[T any](input []T, workers int, predicate func(T) bool) []T {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	type result struct {
		index int
		keep  bool
	}

	results := make([]result, len(input))
	jobs := make(chan struct {
		index int
		value T
	}, len(input))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results[job.index] = result{
					index: job.index,
					keep:  predicate(job.value),
				}
			}
		}()
	}

	// Send jobs
	for i, val := range input {
		jobs <- struct {
			index int
			value T
		}{i, val}
	}
	close(jobs)

	wg.Wait()

	// Collect filtered results
	filtered := make([]T, 0)
	for _, r := range results {
		if r.keep {
			filtered = append(filtered, input[r.index])
		}
	}

	return filtered
}

// ParallelReduce reduces elements in parallel
func ParallelReduce[T, R any](input []T, workers int, initial R, reducer func(R, T) R, combiner func(R, R) R) R {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	if len(input) == 0 {
		return initial
	}

	// Divide work among workers
	chunkSize := (len(input) + workers - 1) / workers
	partialResults := make([]R, workers)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			start := workerID * chunkSize
			end := start + chunkSize
			if end > len(input) {
				end = len(input)
			}

			if start >= len(input) {
				partialResults[workerID] = initial
				return
			}

			result := initial
			for j := start; j < end; j++ {
				result = reducer(result, input[j])
			}
			partialResults[workerID] = result
		}(i)
	}

	wg.Wait()

	// Combine partial results
	finalResult := initial
	for _, partial := range partialResults {
		finalResult = combiner(finalResult, partial)
	}

	return finalResult
}

// ParallelForEach applies a function to each element in parallel
func ParallelForEach[T any](input []T, workers int, fn func(T)) {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	jobs := make(chan T, len(input))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				fn(item)
			}
		}()
	}

	// Send jobs
	for _, val := range input {
		jobs <- val
	}
	close(jobs)

	wg.Wait()
}

// ParallelBatch processes items in batches in parallel
func ParallelBatch[T, R any](input []T, batchSize int, workers int, fn func([]T) []R) []R {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	// Create batches
	batches := make([][]T, 0)
	for i := 0; i < len(input); i += batchSize {
		end := i + batchSize
		if end > len(input) {
			end = len(input)
		}
		batches = append(batches, input[i:end])
	}

	// Process batches in parallel
	type batchResult struct {
		index   int
		results []R
	}

	resultChan := make(chan batchResult, len(batches))
	jobs := make(chan struct {
		index int
		batch []T
	}, len(batches))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results := fn(job.batch)
				resultChan <- batchResult{
					index:   job.index,
					results: results,
				}
			}
		}()
	}

	// Send jobs
	for i, batch := range batches {
		jobs <- struct {
			index int
			batch []T
		}{i, batch}
	}
	close(jobs)

	// Wait and close results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect and order results
	batchResults := make([][]R, len(batches))
	for result := range resultChan {
		batchResults[result.index] = result.results
	}

	// Flatten results
	finalResults := make([]R, 0)
	for _, batch := range batchResults {
		finalResults = append(finalResults, batch...)
	}

	return finalResults
}

// ParallelSum sums numbers in parallel
func ParallelSum[T interface{ ~int | ~int64 | ~float64 }](input []T, workers int) T {
	return ParallelReduce(
		input,
		workers,
		T(0),
		func(acc T, val T) T { return acc + val },
		func(a T, b T) T { return a + b },
	)
}

// ParallelMax finds the maximum value in parallel
func ParallelMax[T interface{ ~int | ~int64 | ~float64 }](input []T, workers int) T {
	if len(input) == 0 {
		return T(0)
	}

	return ParallelReduce(
		input,
		workers,
		input[0],
		func(acc T, val T) T {
			if val > acc {
				return val
			}
			return acc
		},
		func(a T, b T) T {
			if b > a {
				return b
			}
			return a
		},
	)
}

// ParallelMin finds the minimum value in parallel
func ParallelMin[T interface{ ~int | ~int64 | ~float64 }](input []T, workers int) T {
	if len(input) == 0 {
		return T(0)
	}

	return ParallelReduce(
		input,
		workers,
		input[0],
		func(acc T, val T) T {
			if val < acc {
				return val
			}
			return acc
		},
		func(a T, b T) T {
			if b < a {
				return b
			}
			return a
		},
	)
}
