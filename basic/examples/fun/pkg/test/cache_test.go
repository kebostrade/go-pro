package test

import (
	"errors"
	"testing"
	"time"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/cache"
)

// TestCacheBasicOperations tests basic cache operations
func TestCacheBasicOperations(t *testing.T) {
	c := cache.NewCache[string, int](5*time.Second, 100)
	defer c.Stop()

	// Test Set and Get
	c.Set("key1", 100)
	if val, found := c.Get("key1"); !found || val != 100 {
		t.Errorf("Expected to find key1 with value 100, got %d, %v", val, found)
	}

	// Test missing key
	if _, found := c.Get("missing"); found {
		t.Error("Expected not to find missing key")
	}

	// Test Size
	c.Set("key2", 200)
	c.Set("key3", 300)
	if c.Size() != 3 {
		t.Errorf("Expected size 3, got %d", c.Size())
	}

	// Test Delete
	c.Delete("key2")
	if c.Size() != 2 {
		t.Errorf("Expected size 2 after delete, got %d", c.Size())
	}

	// Test Clear
	c.Clear()
	if c.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", c.Size())
	}
}

// TestCacheTTL tests TTL expiration
func TestCacheTTL(t *testing.T) {
	c := cache.NewCache[string, string](500*time.Millisecond, 100)
	defer c.Stop()

	c.Set("temp", "value")

	// Should exist immediately
	if _, found := c.Get("temp"); !found {
		t.Error("Expected to find temp immediately")
	}

	// Wait for expiration
	time.Sleep(600 * time.Millisecond)

	// Should be expired
	if _, found := c.Get("temp"); found {
		t.Error("Expected temp to be expired")
	}
}

// TestCacheCustomTTL tests custom TTL
func TestCacheCustomTTL(t *testing.T) {
	c := cache.NewCache[string, string](5*time.Second, 100)
	defer c.Stop()

	c.SetWithTTL("short", "value", 200*time.Millisecond)
	c.SetWithTTL("long", "value", 2*time.Second)

	// Both should exist
	if !c.Has("short") || !c.Has("long") {
		t.Error("Expected both keys to exist")
	}

	// Wait for short to expire
	time.Sleep(300 * time.Millisecond)

	if c.Has("short") {
		t.Error("Expected short to be expired")
	}
	if !c.Has("long") {
		t.Error("Expected long to still exist")
	}
}

// TestCacheStatistics tests cache statistics
func TestCacheStatistics(t *testing.T) {
	c := cache.NewCache[int, string](5*time.Second, 100)
	defer c.Stop()

	// Perform operations
	c.Set(1, "one")
	c.Set(2, "two")
	c.Get(1)  // Hit
	c.Get(1)  // Hit
	c.Get(99) // Miss
	c.Delete(2)

	stats := c.GetStats()

	if stats.Sets != 2 {
		t.Errorf("Expected 2 sets, got %d", stats.Sets)
	}
	if stats.Hits != 2 {
		t.Errorf("Expected 2 hits, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	if stats.Deletes != 1 {
		t.Errorf("Expected 1 delete, got %d", stats.Deletes)
	}

	hitRate := c.HitRate()
	expectedRate := 2.0 / 3.0
	if hitRate < expectedRate-0.01 || hitRate > expectedRate+0.01 {
		t.Errorf("Expected hit rate ~%.2f, got %.2f", expectedRate, hitRate)
	}
}

// TestCacheGetOrSet tests GetOrSet functionality
func TestCacheGetOrSet(t *testing.T) {
	c := cache.NewCache[string, int](5*time.Second, 100)
	defer c.Stop()

	// First call should set
	val := c.GetOrSet("key", 100)
	if val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}

	// Second call should get existing
	val = c.GetOrSet("key", 200)
	if val != 100 {
		t.Errorf("Expected 100 (existing), got %d", val)
	}
}

// TestCacheGetOrCompute tests GetOrCompute functionality
func TestCacheGetOrCompute(t *testing.T) {
	c := cache.NewCache[string, int](5*time.Second, 100)
	defer c.Stop()

	computeCalls := 0
	compute := func() int {
		computeCalls++
		return 42
	}

	// First call should compute
	val := c.GetOrCompute("key", compute)
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
	if computeCalls != 1 {
		t.Errorf("Expected 1 compute call, got %d", computeCalls)
	}

	// Second call should use cached value
	val = c.GetOrCompute("key", compute)
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
	if computeCalls != 1 {
		t.Errorf("Expected still 1 compute call, got %d", computeCalls)
	}
}

// TestLoadingCache tests loading cache
func TestLoadingCache(t *testing.T) {
	loadCalls := 0
	loader := func(key int) (string, error) {
		loadCalls++
		if key < 0 {
			return "", errors.New("negative key")
		}
		return "value", nil
	}

	lc := cache.NewLoadingCache(5*time.Second, 100, loader)
	defer lc.Stop()

	// First get should load
	val, err := lc.Get(1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if val != "value" {
		t.Errorf("Expected 'value', got %s", val)
	}
	if loadCalls != 1 {
		t.Errorf("Expected 1 load call, got %d", loadCalls)
	}

	// Second get should use cache
	val, err = lc.Get(1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if loadCalls != 1 {
		t.Errorf("Expected still 1 load call, got %d", loadCalls)
	}

	// Test error case
	_, err = lc.Get(-1)
	if err == nil {
		t.Error("Expected error for negative key")
	}
}

// TestLRUCache tests LRU cache
func TestLRUCache(t *testing.T) {
	lru := cache.NewLRUCache[string, int](3)

	// Add items
	lru.Set("a", 1)
	lru.Set("b", 2)
	lru.Set("c", 3)

	if lru.Size() != 3 {
		t.Errorf("Expected size 3, got %d", lru.Size())
	}

	// Add 4th item (should evict 'a')
	lru.Set("d", 4)

	if lru.Size() != 3 {
		t.Errorf("Expected size 3 after eviction, got %d", lru.Size())
	}

	if _, found := lru.Get("a"); found {
		t.Error("Expected 'a' to be evicted")
	}

	if _, found := lru.Get("d"); !found {
		t.Error("Expected 'd' to exist")
	}
}

// TestLRUEvictionOrder tests LRU eviction order
func TestLRUEvictionOrder(t *testing.T) {
	lru := cache.NewLRUCache[int, string](3)

	lru.Set(1, "one")
	lru.Set(2, "two")
	lru.Set(3, "three")

	// Access 1 (moves to front)
	lru.Get(1)

	// Add 4 (should evict 2, not 1)
	lru.Set(4, "four")

	if _, found := lru.Get(2); found {
		t.Error("Expected 2 to be evicted")
	}

	if _, found := lru.Get(1); !found {
		t.Error("Expected 1 to still exist")
	}
}

// TestLRUPeek tests Peek doesn't update access order
func TestLRUPeek(t *testing.T) {
	lru := cache.NewLRUCache[int, string](2)

	lru.Set(1, "one")
	lru.Set(2, "two")

	// Peek at 1 (shouldn't update order)
	lru.Peek(1)

	// Add 3 (should still evict 1)
	lru.Set(3, "three")

	if _, found := lru.Get(1); found {
		t.Error("Expected 1 to be evicted (Peek shouldn't update order)")
	}
}

// TestLFUCache tests LFU cache
func TestLFUCache(t *testing.T) {
	lfu := cache.NewLFUCache[string, int](3)

	lfu.Set("a", 1)
	lfu.Set("b", 2)
	lfu.Set("c", 3)

	// Access 'a' multiple times
	lfu.Get("a")
	lfu.Get("a")
	lfu.Get("a")

	// Access 'b' once
	lfu.Get("b")

	// 'c' not accessed

	// Add 'd' (should evict 'c' - least frequently used)
	lfu.Set("d", 4)

	if _, found := lfu.Get("c"); found {
		t.Error("Expected 'c' to be evicted")
	}

	if _, found := lfu.Get("a"); !found {
		t.Error("Expected 'a' to still exist")
	}
}

// Benchmarks

func BenchmarkCacheSet(b *testing.B) {
	c := cache.NewCache[int, int](5*time.Second, 10000)
	defer c.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(i, i)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := cache.NewCache[int, int](5*time.Second, 10000)
	defer c.Stop()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		c.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(i % 1000)
	}
}

func BenchmarkLRUSet(b *testing.B) {
	lru := cache.NewLRUCache[int, int](1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Set(i, i)
	}
}

func BenchmarkLRUGet(b *testing.B) {
	lru := cache.NewLRUCache[int, int](1000)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		lru.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Get(i % 1000)
	}
}
