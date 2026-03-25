package main

import (
	"context"
	"fmt"
	"lesson-11/exercises"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Lesson 11: Advanced Concurrency Patterns ===")
	fmt.Println()

	// Exercise 1: Worker Pool
	fmt.Println("1. Worker Pool:")
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)
	
	exercises.WorkerPool(jobs, results, 2)
	close(results)
	
	for r := range results {
		fmt.Printf("   Result: %d\n", r)
	}
	fmt.Println()

	// Exercise 4: Pipeline
	fmt.Println("4. Pipeline:")
	sum := exercises.PipelineDemo([]int{1, 2, 3, 4, 5})
	fmt.Printf("   Sum of squares of evens: %d\n", sum) // 4^2 + 8^2? Wait - 1^2=1(odd),2^2=4(even),3^2=9(odd),4^2=16(even),5^2=25(odd) = 4+16=20
	fmt.Println()

	// Exercise 6: WaitGroup
	fmt.Println("6. WaitGroup:")
	tasks := []func(){
		func() { fmt.Println("   Task 1 done") },
		func() { fmt.Println("   Task 2 done") },
		func() { fmt.Println("   Task 3 done") },
	}
	exercises.WaitGroupDemo(tasks)
	fmt.Println()

	// Exercise 7: SafeCounter
	fmt.Println("7. SafeCounter:")
	counter := exercises.SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("   Counter value: %d\n", counter.Value())
	fmt.Println()

	// Exercise 8: SafeCache
	fmt.Println("8. SafeCache:")
	cache := exercises.NewSafeCache()
	cache.Set("name", "Alice")
	cache.Set("age", "30")
	name, _ := cache.Get("name")
	age, _ := cache.Get("age")
	fmt.Printf("   name=%s, age=%s\n", name, age)
	fmt.Println()

	// Exercise 9: Once
	fmt.Println("9. Once Executor:")
	executor := exercises.OnceExecutor{}
	executor.Do(func() { fmt.Println("   First call") })
	executor.Do(func() { fmt.Println("   Second call") })
	executor.Do(func() { fmt.Println("   Third call") })
	fmt.Println()

	// Exercise 11: Atomic Counter
	fmt.Println("11. Atomic Counter:")
	atomic := exercises.NewAtomicCounter()
	for i := 0; i < 1000; i++ {
		atomic.Increment()
	}
	fmt.Printf("   Atomic value: %d\n", atomic.Value())
	fmt.Println()

	fmt.Println("=== All exercises completed! ===")
}
