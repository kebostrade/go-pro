package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*
GOROUTINES - Lightweight Threads

Goroutines are functions that run concurrently with other functions.
They are extremely lightweight (2KB stack) and managed by the Go runtime.

Key Concepts:
- Launched with 'go' keyword
- Managed by Go scheduler
- Communicate via channels
- Share memory (use sync primitives)
*/

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                              ║")
	fmt.Println("║                  🔄 Goroutine Basics                         ║")
	fmt.Println("║                                                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 1. Basic Goroutine
	basicGoroutine()

	// 2. Multiple Goroutines
	multipleGoroutines()

	// 3. Goroutine with WaitGroup
	goroutineWithWaitGroup()

	// 4. Goroutine with Anonymous Function
	anonymousGoroutine()

	// 5. Goroutine Closure
	goroutineClosure()

	// 6. Goroutine Scheduler Info
	schedulerInfo()
}

// 1. Basic Goroutine
func basicGoroutine() {
	fmt.Println("\n1️⃣  BASIC GOROUTINE")
	fmt.Println("───────────────────────────────────────────────────────────────")

	// Without goroutine (sequential)
	fmt.Println("Sequential execution:")
	sayHello("Alice")
	sayHello("Bob")

	// With goroutine (concurrent)
	fmt.Println("\nConcurrent execution:")
	go sayHello("Charlie")
	go sayHello("Diana")

	// Wait for goroutines to complete
	time.Sleep(100 * time.Millisecond)
}

func sayHello(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// 2. Multiple Goroutines
func multipleGoroutines() {
	fmt.Println("\n2️⃣  MULTIPLE GOROUTINES")
	fmt.Println("───────────────────────────────────────────────────────────────")

	// Launch 5 goroutines
	for i := 1; i <= 5; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d started\n", id)
			time.Sleep(time.Duration(id*10) * time.Millisecond)
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
}

// 3. Goroutine with WaitGroup
func goroutineWithWaitGroup() {
	fmt.Println("\n3️⃣  GOROUTINE WITH WAITGROUP")
	fmt.Println("───────────────────────────────────────────────────────────────")

	var wg sync.WaitGroup

	// Launch 3 workers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	// Wait for all workers to complete
	wg.Wait()
	fmt.Println("✅ All workers completed")
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Duration(id*20) * time.Millisecond)
	fmt.Printf("Worker %d done\n", id)
}

// 4. Anonymous Goroutine
func anonymousGoroutine() {
	fmt.Println("\n4️⃣  ANONYMOUS GOROUTINE")
	fmt.Println("───────────────────────────────────────────────────────────────")

	var wg sync.WaitGroup

	// Anonymous function as goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Anonymous goroutine executing")
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Anonymous goroutine complete")
	}()

	wg.Wait()
}

// 5. Goroutine Closure (Common Pitfall)
func goroutineClosure() {
	fmt.Println("\n5️⃣  GOROUTINE CLOSURE")
	fmt.Println("───────────────────────────────────────────────────────────────")

	var wg sync.WaitGroup

	// ❌ WRONG: Closure captures loop variable
	fmt.Println("❌ Wrong way (closure captures loop variable):")
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Wrong: i = %d\n", i) // May print 4, 4, 4
		}()
	}
	wg.Wait()

	// ✅ CORRECT: Pass variable as parameter
	fmt.Println("\n✅ Correct way (pass as parameter):")
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Correct: id = %d\n", id) // Prints 1, 2, 3
		}(i)
	}
	wg.Wait()
}

// 6. Scheduler Info
func schedulerInfo() {
	fmt.Println("\n6️⃣  GOROUTINE SCHEDULER INFO")
	fmt.Println("───────────────────────────────────────────────────────────────")

	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	// Launch many goroutines
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		}()
	}

	fmt.Printf("Number of Goroutines (after launch): %d\n", runtime.NumGoroutine())
	wg.Wait()
	fmt.Printf("Number of Goroutines (after completion): %d\n", runtime.NumGoroutine())
}

/*
KEY TAKEAWAYS:

1. Goroutines are cheap - you can have thousands
2. Always wait for goroutines to complete (WaitGroup, channels, etc.)
3. Be careful with closures - pass variables as parameters
4. Use 'defer wg.Done()' to ensure cleanup
5. Goroutines share memory - use sync primitives or channels

COMMON PATTERNS:
- Worker pools
- Fan-out/Fan-in
- Pipeline processing
- Concurrent requests
*/

