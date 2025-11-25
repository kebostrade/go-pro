package main

import (
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/cache"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Generic Cache Demo")

	demo1BasicOperations()
	demo2TTLExpiration()
	demo3Statistics()
	demo4LoadingCache()
	demo5GetOrCompute()
}

func demo1BasicOperations() {
	utils.PrintSubHeader("1. Basic Cache Operations")

	// Create cache with 5 second TTL, max 100 items
	c := cache.NewCache[string, string](5*time.Second, 100)
	defer c.Stop()

	// Set values
	c.Set("name", "Alice")
	c.Set("city", "New York")
	c.Set("country", "USA")

	fmt.Printf("Cache size: %d\n", c.Size())

	// Get values
	if value, found := c.Get("name"); found {
		fmt.Printf("Found 'name': %s\n", value)
	}

	if value, found := c.Get("city"); found {
		fmt.Printf("Found 'city': %s\n", value)
	}

	// Check existence
	fmt.Printf("Has 'name': %v\n", c.Has("name"))
	fmt.Printf("Has 'missing': %v\n", c.Has("missing"))

	// Get all keys
	keys := c.Keys()
	fmt.Printf("All keys: %v\n", keys)

	// Delete
	c.Delete("city")
	fmt.Printf("After deleting 'city', size: %d\n", c.Size())
}

func demo2TTLExpiration() {
	utils.PrintSubHeader("2. TTL and Expiration")

	c := cache.NewCache[string, string](2*time.Second, 100)
	defer c.Stop()

	// Set with default TTL (2 seconds)
	c.Set("temp1", "expires in 2 seconds")
	fmt.Println("Set 'temp1' with 2 second TTL")

	// Set with custom TTL (1 second)
	c.SetWithTTL("temp2", "expires in 1 second", 1*time.Second)
	fmt.Println("Set 'temp2' with 1 second TTL")

	// Set with no expiration
	c.SetWithTTL("permanent", "never expires", 0)
	fmt.Println("Set 'permanent' with no expiration")

	// Check immediately
	fmt.Println("\nImmediately after setting:")
	fmt.Printf("  temp1: %v\n", c.Has("temp1"))
	fmt.Printf("  temp2: %v\n", c.Has("temp2"))
	fmt.Printf("  permanent: %v\n", c.Has("permanent"))

	// Wait 1.5 seconds
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("\nAfter 1.5 seconds:")
	fmt.Printf("  temp1: %v\n", c.Has("temp1"))
	fmt.Printf("  temp2: %v (expired)\n", c.Has("temp2"))
	fmt.Printf("  permanent: %v\n", c.Has("permanent"))

	// Wait another 1 second
	time.Sleep(1 * time.Second)
	fmt.Println("\nAfter 2.5 seconds total:")
	fmt.Printf("  temp1: %v (expired)\n", c.Has("temp1"))
	fmt.Printf("  permanent: %v\n", c.Has("permanent"))
}

func demo3Statistics() {
	utils.PrintSubHeader("3. Cache Statistics")

	c := cache.NewCache[int, string](5*time.Second, 100)
	defer c.Stop()

	// Perform operations
	for i := 0; i < 10; i++ {
		c.Set(i, fmt.Sprintf("value%d", i))
	}

	// Some hits
	for i := 0; i < 5; i++ {
		c.Get(i)
		c.Get(i) // Hit twice
	}

	// Some misses
	for i := 100; i < 105; i++ {
		c.Get(i)
	}

	// Some deletes
	c.Delete(0)
	c.Delete(1)

	// Get statistics
	stats := c.GetStats()
	fmt.Printf("Sets: %d\n", stats.Sets)
	fmt.Printf("Hits: %d\n", stats.Hits)
	fmt.Printf("Misses: %d\n", stats.Misses)
	fmt.Printf("Deletes: %d\n", stats.Deletes)
	fmt.Printf("Hit Rate: %.2f%%\n", c.HitRate()*100)
}

func demo4LoadingCache() {
	utils.PrintSubHeader("4. Loading Cache (Auto-Load on Miss)")

	// Create loading cache with loader function
	lc := cache.NewLoadingCache[int, string](
		5*time.Second,
		100,
		func(key int) (string, error) {
			// Simulate loading from database
			time.Sleep(100 * time.Millisecond)
			return fmt.Sprintf("loaded-value-%d", key), nil
		},
	)
	defer lc.Stop()

	fmt.Println("Getting key 1 (will load)...")
	start := time.Now()
	value, err := lc.Get(1)
	duration := time.Since(start)
	if err == nil {
		fmt.Printf("  Value: %s (took %v)\n", value, duration)
	}

	fmt.Println("\nGetting key 1 again (cached)...")
	start = time.Now()
	value, err = lc.Get(1)
	duration = time.Since(start)
	if err == nil {
		fmt.Printf("  Value: %s (took %v)\n", value, duration)
	}

	fmt.Println("\nGetting key 2 (will load)...")
	start = time.Now()
	value, err = lc.Get(2)
	duration = time.Since(start)
	if err == nil {
		fmt.Printf("  Value: %s (took %v)\n", value, duration)
	}
}

func demo5GetOrCompute() {
	utils.PrintSubHeader("5. GetOrCompute Pattern")

	c := cache.NewCache[string, int](5*time.Second, 100)
	defer c.Stop()

	// Expensive computation
	expensiveCompute := func() int {
		fmt.Println("  Computing expensive value...")
		time.Sleep(200 * time.Millisecond)
		return 42
	}

	fmt.Println("First call (will compute):")
	start := time.Now()
	value := c.GetOrCompute("expensive", expensiveCompute)
	duration := time.Since(start)
	fmt.Printf("  Result: %d (took %v)\n", value, duration)

	fmt.Println("\nSecond call (cached):")
	start = time.Now()
	value = c.GetOrCompute("expensive", expensiveCompute)
	duration = time.Since(start)
	fmt.Printf("  Result: %d (took %v)\n", value, duration)

	// GetOrSet pattern
	fmt.Println("\nGetOrSet pattern:")
	value2 := c.GetOrSet("default", 100)
	fmt.Printf("  First call: %d\n", value2)

	value2 = c.GetOrSet("default", 200)
	fmt.Printf("  Second call: %d (still returns first value)\n", value2)
}
