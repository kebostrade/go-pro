package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/cache"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("LRU & LFU Cache Demo")

	demo1LRUBasics()
	demo2LRUEviction()
	demo3LRUStatistics()
	demo4LFUCache()
}

func demo1LRUBasics() {
	utils.PrintSubHeader("1. LRU Cache Basics")

	// Create LRU cache with capacity 5
	lru := cache.NewLRUCache[string, int](5)

	// Add items
	lru.Set("a", 1)
	lru.Set("b", 2)
	lru.Set("c", 3)
	lru.Set("d", 4)
	lru.Set("e", 5)

	fmt.Printf("Cache size: %d/%d\n", lru.Size(), lru.Capacity())
	fmt.Printf("Keys (most recent first): %v\n", lru.Keys())

	// Access some items (moves them to front)
	lru.Get("a")
	lru.Get("c")

	fmt.Printf("\nAfter accessing 'a' and 'c':\n")
	fmt.Printf("Keys (most recent first): %v\n", lru.Keys())

	// Get oldest and newest
	if oldKey, oldVal, ok := lru.GetOldest(); ok {
		fmt.Printf("Oldest: %s = %d\n", oldKey, oldVal)
	}

	if newKey, newVal, ok := lru.GetNewest(); ok {
		fmt.Printf("Newest: %s = %d\n", newKey, newVal)
	}
}

func demo2LRUEviction() {
	utils.PrintSubHeader("2. LRU Eviction Policy")

	// Create small LRU cache (capacity 3)
	lru := cache.NewLRUCache[int, string](3)

	fmt.Println("Adding items to cache (capacity 3):")
	lru.Set(1, "one")
	fmt.Printf("  Added 1, keys: %v\n", lru.Keys())

	lru.Set(2, "two")
	fmt.Printf("  Added 2, keys: %v\n", lru.Keys())

	lru.Set(3, "three")
	fmt.Printf("  Added 3, keys: %v\n", lru.Keys())

	fmt.Println("\nCache is full. Adding item 4 will evict oldest (1):")
	lru.Set(4, "four")
	fmt.Printf("  Added 4, keys: %v\n", lru.Keys())

	// Verify item 1 was evicted
	if _, found := lru.Get(1); !found {
		fmt.Println("  Item 1 was evicted ✓")
	}

	// Access item 2 (moves to front)
	fmt.Println("\nAccessing item 2 (moves to front):")
	lru.Get(2)
	fmt.Printf("  Keys: %v\n", lru.Keys())

	// Add item 5 (will evict 3, not 2)
	fmt.Println("\nAdding item 5 will evict oldest (3, not 2):")
	lru.Set(5, "five")
	fmt.Printf("  Keys: %v\n", lru.Keys())

	if _, found := lru.Get(3); !found {
		fmt.Println("  Item 3 was evicted ✓")
	}
	if _, found := lru.Get(2); found {
		fmt.Println("  Item 2 still in cache ✓")
	}
}

func demo3LRUStatistics() {
	utils.PrintSubHeader("3. LRU Cache Statistics")

	lru := cache.NewLRUCache[string, int](10)

	// Perform operations
	for i := 0; i < 5; i++ {
		lru.Set(fmt.Sprintf("key%d", i), i)
	}

	// Some hits
	for i := 0; i < 3; i++ {
		lru.Get(fmt.Sprintf("key%d", i))
		lru.Get(fmt.Sprintf("key%d", i)) // Hit twice
	}

	// Some misses
	for i := 10; i < 15; i++ {
		lru.Get(fmt.Sprintf("key%d", i))
	}

	// Get statistics
	stats := lru.GetStats()
	fmt.Printf("Sets: %d\n", stats.Sets)
	fmt.Printf("Hits: %d\n", stats.Hits)
	fmt.Printf("Misses: %d\n", stats.Misses)
	fmt.Printf("Evictions: %d\n", stats.Evictions)
	fmt.Printf("Hit Rate: %.2f%%\n", lru.HitRate()*100)

	// Demonstrate Peek (doesn't update access order)
	fmt.Println("\nPeek vs Get:")
	fmt.Printf("Keys before: %v\n", lru.Keys())

	lru.Peek("key4") // Doesn't change order
	fmt.Printf("After Peek('key4'): %v\n", lru.Keys())

	lru.Get("key4") // Changes order
	fmt.Printf("After Get('key4'): %v\n", lru.Keys())
}

func demo4LFUCache() {
	utils.PrintSubHeader("4. LFU Cache (Least Frequently Used)")

	// Create LFU cache with capacity 3
	lfu := cache.NewLFUCache[string, int](3)

	fmt.Println("Adding items to LFU cache (capacity 3):")
	lfu.Set("a", 1)
	lfu.Set("b", 2)
	lfu.Set("c", 3)
	fmt.Printf("  Size: %d\n", lfu.Size())

	// Access items different number of times
	fmt.Println("\nAccessing items:")
	lfu.Get("a") // freq = 2
	lfu.Get("a") // freq = 3
	lfu.Get("b") // freq = 2
	fmt.Println("  'a' accessed 3 times total")
	fmt.Println("  'b' accessed 2 times total")
	fmt.Println("  'c' accessed 1 time total")

	// Add new item (will evict 'c' - least frequently used)
	fmt.Println("\nAdding 'd' (will evict 'c' - least frequently used):")
	lfu.Set("d", 4)
	fmt.Printf("  Size: %d\n", lfu.Size())

	// Verify 'c' was evicted
	if _, found := lfu.Get("c"); !found {
		fmt.Println("  'c' was evicted ✓")
	}

	// Verify others still exist
	if _, found := lfu.Get("a"); found {
		fmt.Println("  'a' still in cache ✓")
	}
	if _, found := lfu.Get("b"); found {
		fmt.Println("  'b' still in cache ✓")
	}
	if _, found := lfu.Get("d"); found {
		fmt.Println("  'd' in cache ✓")
	}
}
