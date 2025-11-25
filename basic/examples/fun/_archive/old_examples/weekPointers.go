//go:build ignore

package main

import (
	"fmt"
	"runtime"
	"time"
	"weak"
)

// Weak pointers are used to avoid memory leaks when storing pointers to objects that are no longer needed.
// Weak pointers allow you to reference objects without preventing garbage collection,
// enabling memory-efficient applications, caches, and smarter data structures.

type ResourceCache struct {
	items map[string]weak.Pointer[string]
}

func NewResourceCache() *ResourceCache {
	return &ResourceCache{
		items: make(map[string]weak.Pointer[string]),
	}
}

func (rc *ResourceCache) Add(key string, value *string) {
	weakPtr := weak.Make(value)
	rc.items[key] = weakPtr
}

func (rc *ResourceCache) Get(key string) (*string, bool) {
	weakPtr, exists := rc.items[key]
	if !exists {
		return nil, false
	}

	if ptr := weakPtr.Value(); ptr != nil {
		return ptr, true
	}

	delete(rc.items, key)
	return nil, false
}

func printMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf(
		"Memory: %.2f MB allocated\n",
		float64(m.Alloc)/1024/1024,
	)

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func main() {
	cache := NewResourceCache()

	createBigString := func() *string {
		s := make([]byte, 10<<20)
		str := string(s)
		return &str
	}

	bigData := createBigString()
	cache.Add("big", bigData)
	if _, ok := cache.Get("big"); ok {
		fmt.Println("Found value:")
	}

	bigData = nil

	fmt.Print("Before GC:")
	printMemoryUsage()

	runtime.GC()
	fmt.Print("After GC:")
	printMemoryUsage()

	if val, ok := cache.Get("big"); ok {
		fmt.Printf(
			"Still holding %.2f MB cache\n",
			float64(len(*val))/1024/1024,
		)
	}

	time.Sleep(time.Second)
}
