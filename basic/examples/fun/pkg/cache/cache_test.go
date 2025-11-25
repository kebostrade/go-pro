package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 0)
		defer c.Stop()

		c.Set("key1", 42)
		val, found := c.Get("key1")

		if !found || val != 42 {
			t.Errorf("Get(\"key1\") = %v, %v; want 42, true", val, found)
		}
	})

	t.Run("Get non-existent", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 0)
		defer c.Stop()

		_, found := c.Get("missing")
		if found {
			t.Error("Get(\"missing\") found = true; want false")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		c := NewCache[string, int](50*time.Millisecond, 0)
		defer c.Stop()

		c.Set("key1", 42)
		time.Sleep(100 * time.Millisecond)

		_, found := c.Get("key1")
		if found {
			t.Error("Expired key should not be found")
		}
	})

	t.Run("MaxSize eviction", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 2)
		defer c.Stop()

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Set("key3", 3)

		if c.Size() > 2 {
			t.Errorf("Size() = %d; want <= 2", c.Size())
		}
	})

	t.Run("GetOrSet", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 0)
		defer c.Stop()

		val := c.GetOrSet("key1", 42)
		if val != 42 {
			t.Errorf("GetOrSet() = %d; want 42", val)
		}

		val = c.GetOrSet("key1", 100)
		if val != 42 {
			t.Errorf("GetOrSet() should return existing value 42, got %d", val)
		}
	})

	t.Run("GetOrCompute", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 0)
		defer c.Stop()

		called := false
		val := c.GetOrCompute("key1", func() int {
			called = true
			return 42
		})

		if !called || val != 42 {
			t.Errorf("GetOrCompute() = %d, called = %v; want 42, true", val, called)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 0)
		defer c.Stop()

		c.Set("key1", 42)
		deleted := c.Delete("key1")

		if !deleted {
			t.Error("Delete() = false; want true")
		}

		_, found := c.Get("key1")
		if found {
			t.Error("Key should be deleted")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 0)
		defer c.Stop()

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Clear()

		if c.Size() != 0 {
			t.Errorf("Clear() Size() = %d; want 0", c.Size())
		}
	})

	t.Run("Stats", func(t *testing.T) {
		c := NewCache[string, int](time.Minute, 0)
		defer c.Stop()

		c.Set("key1", 42)
		c.Get("key1")
		c.Get("missing")

		stats := c.GetStats()
		if stats.Hits != 1 || stats.Misses != 1 || stats.Sets != 1 {
			t.Errorf("Stats = %+v; want Hits=1, Misses=1, Sets=1", stats)
		}

		hitRate := c.HitRate()
		if hitRate != 0.5 {
			t.Errorf("HitRate() = %f; want 0.5", hitRate)
		}
	})
}

func TestLRUCache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		c := NewLRUCache[string, int](3)

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Set("key3", 3)

		val, found := c.Get("key1")
		if !found || val != 1 {
			t.Errorf("Get(\"key1\") = %v, %v; want 1, true", val, found)
		}
	})

	t.Run("Eviction", func(t *testing.T) {
		c := NewLRUCache[string, int](2)

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Set("key3", 3)

		_, found := c.Get("key1")
		if found {
			t.Error("LRU key1 should be evicted")
		}
	})

	t.Run("Access updates order", func(t *testing.T) {
		c := NewLRUCache[string, int](2)

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Get("key1")
		c.Set("key3", 3)

		_, found := c.Get("key2")
		if found {
			t.Error("key2 should be evicted after accessing key1")
		}
	})

	t.Run("Peek does not update order", func(t *testing.T) {
		c := NewLRUCache[string, int](2)

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Peek("key1")
		c.Set("key3", 3)

		_, found := c.Get("key1")
		if found {
			t.Error("Peek should not affect eviction order")
		}
	})

	t.Run("GetOldest and GetNewest", func(t *testing.T) {
		c := NewLRUCache[string, int](3)

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Set("key3", 3)

		key, val, found := c.GetNewest()
		if !found || key != "key3" || val != 3 {
			t.Errorf("GetNewest() = %v, %v, %v; want key3, 3, true", key, val, found)
		}

		key, val, found = c.GetOldest()
		if !found || key != "key1" || val != 1 {
			t.Errorf("GetOldest() = %v, %v, %v; want key1, 1, true", key, val, found)
		}
	})

	t.Run("Resize", func(t *testing.T) {
		c := NewLRUCache[string, int](3)

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Set("key3", 3)

		c.Resize(2)

		if c.Size() != 2 {
			t.Errorf("After Resize(2) Size() = %d; want 2", c.Size())
		}
	})
}

func TestLFUCache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		c := NewLFUCache[string, int](3)

		c.Set("key1", 1)
		c.Set("key2", 2)

		val, found := c.Get("key1")
		if !found || val != 1 {
			t.Errorf("Get(\"key1\") = %v, %v; want 1, true", val, found)
		}
	})

	t.Run("Eviction based on frequency", func(t *testing.T) {
		c := NewLFUCache[string, int](2)

		c.Set("key1", 1)
		c.Set("key2", 2)
		c.Get("key1")
		c.Get("key1")
		c.Set("key3", 3)

		_, found := c.Get("key2")
		if found {
			t.Error("LFU key2 should be evicted (lowest frequency)")
		}
	})
}

func TestLoadingCache(t *testing.T) {
	t.Run("Get with loader", func(t *testing.T) {
		loadCount := 0
		lc := NewLoadingCache[string, int](time.Minute, 0, func(key string) (int, error) {
			loadCount++
			return 42, nil
		})
		defer lc.Stop()

		val, err := lc.Get("key1")
		if err != nil || val != 42 || loadCount != 1 {
			t.Errorf("Get() = %v, %v, loadCount = %d; want 42, nil, 1", val, err, loadCount)
		}

		val, err = lc.Get("key1")
		if err != nil || val != 42 || loadCount != 1 {
			t.Errorf("Second Get() should use cache, loadCount = %d; want 1", loadCount)
		}
	})

	t.Run("Get with loader error", func(t *testing.T) {
		lc := NewLoadingCache[string, int](time.Minute, 0, func(key string) (int, error) {
			return 0, errors.New("load error")
		})
		defer lc.Stop()

		_, err := lc.Get("key1")
		if err == nil {
			t.Error("Get() error = nil; want error")
		}
	})

	t.Run("GetWithContext", func(t *testing.T) {
		lc := NewLoadingCache[string, int](time.Minute, 0, func(key string) (int, error) {
			return 42, nil
		})
		defer lc.Stop()

		ctx := context.Background()
		val, err := lc.GetWithContext(ctx, "key1")

		if err != nil || val != 42 {
			t.Errorf("GetWithContext() = %v, %v; want 42, nil", val, err)
		}
	})

	t.Run("Invalidate", func(t *testing.T) {
		lc := NewLoadingCache[string, int](time.Minute, 0, func(key string) (int, error) {
			return 42, nil
		})
		defer lc.Stop()

		lc.Get("key1")
		lc.Invalidate("key1")

		// Verify key is no longer in cache by checking if loader is called again
		loadCount := 0
		lc = NewLoadingCache[string, int](time.Minute, 0, func(key string) (int, error) {
			loadCount++
			return 42, nil
		})
		defer lc.Stop()

		lc.Get("key1")
		lc.Invalidate("key1")
		lc.Get("key1")

		if loadCount != 2 {
			t.Errorf("After Invalidate, loadCount = %d; want 2", loadCount)
		}
	})
}

// Benchmarks
func BenchmarkCacheSet(b *testing.B) {
	c := NewCache[int, int](time.Minute, 0)
	defer c.Stop()

	for i := 0; i < b.N; i++ {
		c.Set(i, i)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := NewCache[int, int](time.Minute, 0)
	defer c.Stop()

	for i := 0; i < 1000; i++ {
		c.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(i % 1000)
	}
}

func BenchmarkLRUCacheSet(b *testing.B) {
	c := NewLRUCache[int, int](1000)

	for i := 0; i < b.N; i++ {
		c.Set(i, i)
	}
}

func BenchmarkLRUCacheGet(b *testing.B) {
	c := NewLRUCache[int, int](1000)

	for i := 0; i < 1000; i++ {
		c.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(i % 1000)
	}
}
