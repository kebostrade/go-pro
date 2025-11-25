package main

import (
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/concurrency"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Parallel Processing Demo")

	demo1ParallelMap()
	demo2ParallelFilter()
	demo3ParallelReduce()
	demo4ParallelBatch()
}

func demo1ParallelMap() {
	utils.PrintSubHeader("1. Parallel Map")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("Input: %v\n", numbers)

	// Square each number in parallel
	fmt.Println("\nSquaring numbers with 4 workers...")
	start := time.Now()

	results := concurrency.ParallelMap(numbers, 4, func(n int) int {
		time.Sleep(50 * time.Millisecond) // Simulate work
		return n * n
	})

	duration := time.Since(start)
	fmt.Printf("Results: %v\n", results)
	fmt.Printf("Completed in %v\n", duration)

	// Compare with sequential
	fmt.Println("\nSequential processing for comparison...")
	start = time.Now()
	sequential := make([]int, len(numbers))
	for i, n := range numbers {
		time.Sleep(50 * time.Millisecond)
		sequential[i] = n * n
	}
	seqDuration := time.Since(start)

	fmt.Printf("Sequential completed in %v\n", seqDuration)
	fmt.Printf("Speedup: %.2fx\n", float64(seqDuration)/float64(duration))
}

func demo2ParallelFilter() {
	utils.PrintSubHeader("2. Parallel Filter")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	fmt.Printf("Input: %v\n", numbers)

	// Filter even numbers
	fmt.Println("\nFiltering even numbers with 4 workers...")
	results := concurrency.ParallelFilter(numbers, 4, func(n int) bool {
		return n%2 == 0
	})

	fmt.Printf("Even numbers: %v\n", results)

	// Filter numbers > 10
	fmt.Println("\nFiltering numbers > 10...")
	results = concurrency.ParallelFilter(numbers, 4, func(n int) bool {
		return n > 10
	})

	fmt.Printf("Numbers > 10: %v\n", results)
}

func demo3ParallelReduce() {
	utils.PrintSubHeader("3. Parallel Reduce")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("Input: %v\n", numbers)

	// Sum
	sum := concurrency.ParallelSum(numbers, 4)
	fmt.Printf("\nSum: %d\n", sum)

	// Max
	max := concurrency.ParallelMax(numbers, 4)
	fmt.Printf("Max: %d\n", max)

	// Min
	min := concurrency.ParallelMin(numbers, 4)
	fmt.Printf("Min: %d\n", min)

	// Custom reduce: product
	product := concurrency.ParallelReduce(
		numbers,
		4,
		1,
		func(acc int, val int) int { return acc * val },
		func(a int, b int) int { return a * b },
	)
	fmt.Printf("Product: %d\n", product)
}

func demo4ParallelBatch() {
	utils.PrintSubHeader("4. Parallel Batch Processing")

	// Create 20 items
	items := make([]int, 20)
	for i := range items {
		items[i] = i + 1
	}

	fmt.Printf("Processing %d items in batches of 5 with 3 workers...\n", len(items))

	// Process in batches
	results := concurrency.ParallelBatch(
		items,
		5, // batch size
		3, // workers
		func(batch []int) []string {
			// Process each batch
			batchResults := make([]string, len(batch))
			for i, item := range batch {
				time.Sleep(50 * time.Millisecond) // Simulate work
				batchResults[i] = fmt.Sprintf("Processed-%d", item)
			}
			return batchResults
		},
	)

	fmt.Printf("\nProcessed %d items\n", len(results))
	fmt.Println("First 10 results:")
	for i := 0; i < 10 && i < len(results); i++ {
		fmt.Printf("  %s\n", results[i])
	}
}
