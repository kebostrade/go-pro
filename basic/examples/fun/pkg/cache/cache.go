package cache

import (
	"context"
	"sync"
	"time"
)

// CacheItem represents a single cache entry with generic value type
type CacheItem[V any] struct {
	Value      V
	Expiration time.Time
	AccessTime time.Time
	HitCount   int64
}

// IsExpired checks if the cache item has expired
func (item *CacheItem[V]) IsExpired() bool {
	return !item.Expiration.IsZero() && time.Now().After(item.Expiration)
}

// Cache is a thread-safe in-memory cache with expiration and statistics
type Cache[K comparable, V any] struct {
	items       map[K]*CacheItem[V]
	mu          sync.RWMutex
	defaultTTL  time.Duration
	maxSize     int
	stopCleanup chan struct{}
	stats       *CacheStats
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	mu          sync.RWMutex
	Hits        int64
	Misses      int64
	Sets        int64
	Deletes     int64
	Evictions   int64
	Expirations int64
}

// NewCache creates a new cache with default TTL and optional max size
// maxSize of 0 means unlimited
func NewCache[K comparable, V any](defaultTTL time.Duration, maxSize int) *Cache[K, V] {
	cache := &Cache[K, V]{
		items:       make(map[K]*CacheItem[V]),
		defaultTTL:  defaultTTL,
		maxSize:     maxSize,
		stopCleanup: make(chan struct{}),
		stats:       &CacheStats{},
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Set adds or updates an item in the cache with default TTL
func (c *Cache[K, V]) Set(key K, value V) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL adds or updates an item with custom TTL
// ttl of 0 means no expiration
func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.items[key] = &CacheItem[V]{
		Value:      value,
		Expiration: expiration,
		AccessTime: time.Now(),
		HitCount:   0,
	}

	c.stats.mu.Lock()
	c.stats.Sets++
	c.stats.mu.Unlock()

	// Check if we need to evict items
	if c.maxSize > 0 && len(c.items) > c.maxSize {
		c.evictOldest()
	}
}

// Get retrieves an item from the cache
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		c.stats.mu.Lock()
		c.stats.Misses++
		c.stats.mu.Unlock()

		var zero V
		return zero, false
	}

	// Check if expired
	if item.IsExpired() {
		delete(c.items, key)
		c.stats.mu.Lock()
		c.stats.Misses++
		c.stats.Expirations++
		c.stats.mu.Unlock()

		var zero V
		return zero, false
	}

	// Update access time and hit count
	item.AccessTime = time.Now()
	item.HitCount++

	c.stats.mu.Lock()
	c.stats.Hits++
	c.stats.mu.Unlock()

	return item.Value, true
}

// GetOrSet retrieves an item or sets it if not found
func (c *Cache[K, V]) GetOrSet(key K, defaultValue V) V {
	if value, found := c.Get(key); found {
		return value
	}

	c.Set(key, defaultValue)
	return defaultValue
}

// GetOrCompute retrieves an item or computes it if not found
func (c *Cache[K, V]) GetOrCompute(key K, compute func() V) V {
	if value, found := c.Get(key); found {
		return value
	}

	value := compute()
	c.Set(key, value)
	return value
}

// Delete removes an item from the cache
func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; exists {
		delete(c.items, key)
		c.stats.mu.Lock()
		c.stats.Deletes++
		c.stats.mu.Unlock()
		return true
	}

	return false
}

// Clear removes all items from the cache
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*CacheItem[V])
}

// Size returns the number of items in the cache
func (c *Cache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// Keys returns all keys in the cache
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}

	return keys
}

// Has checks if a key exists in the cache
func (c *Cache[K, V]) Has(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	return !item.IsExpired()
}

// GetAll returns all non-expired items
func (c *Cache[K, V]) GetAll() map[K]V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[K]V)
	for key, item := range c.items {
		if !item.IsExpired() {
			result[key] = item.Value
		}
	}

	return result
}

// GetStats returns a copy of cache statistics
func (c *Cache[K, V]) GetStats() CacheStats {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()

	return CacheStats{
		Hits:        c.stats.Hits,
		Misses:      c.stats.Misses,
		Sets:        c.stats.Sets,
		Deletes:     c.stats.Deletes,
		Evictions:   c.stats.Evictions,
		Expirations: c.stats.Expirations,
	}
}

// HitRate returns the cache hit rate (0.0 to 1.0)
func (c *Cache[K, V]) HitRate() float64 {
	stats := c.GetStats()
	total := stats.Hits + stats.Misses
	if total == 0 {
		return 0.0
	}
	return float64(stats.Hits) / float64(total)
}

// ResetStats resets all statistics
func (c *Cache[K, V]) ResetStats() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()

	c.stats.Hits = 0
	c.stats.Misses = 0
	c.stats.Sets = 0
	c.stats.Deletes = 0
	c.stats.Evictions = 0
	c.stats.Expirations = 0
}

// Stop stops the cleanup goroutine
func (c *Cache[K, V]) Stop() {
	close(c.stopCleanup)
}

// cleanupExpired periodically removes expired items
func (c *Cache[K, V]) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCleanup:
			return
		case <-ticker.C:
			c.mu.Lock()
			expiredCount := 0
			for key, item := range c.items {
				if item.IsExpired() {
					delete(c.items, key)
					expiredCount++
				}
			}
			c.mu.Unlock()

			if expiredCount > 0 {
				c.stats.mu.Lock()
				c.stats.Expirations += int64(expiredCount)
				c.stats.mu.Unlock()
			}
		}
	}
}

// evictOldest removes the least recently accessed item
func (c *Cache[K, V]) evictOldest() {
	var oldestKey K
	var oldestTime time.Time
	first := true

	for key, item := range c.items {
		if first || item.AccessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.AccessTime
			first = false
		}
	}

	if !first {
		delete(c.items, oldestKey)
		c.stats.mu.Lock()
		c.stats.Evictions++
		c.stats.mu.Unlock()
	}
}

// LoadingCache wraps a cache with automatic loading on miss
type LoadingCache[K comparable, V any] struct {
	cache  *Cache[K, V]
	loader func(K) (V, error)
}

// NewLoadingCache creates a cache that automatically loads values on miss
func NewLoadingCache[K comparable, V any](
	defaultTTL time.Duration,
	maxSize int,
	loader func(K) (V, error),
) *LoadingCache[K, V] {
	return &LoadingCache[K, V]{
		cache:  NewCache[K, V](defaultTTL, maxSize),
		loader: loader,
	}
}

// Get retrieves an item, loading it if not found
func (lc *LoadingCache[K, V]) Get(key K) (V, error) {
	if value, found := lc.cache.Get(key); found {
		return value, nil
	}

	value, err := lc.loader(key)
	if err != nil {
		var zero V
		return zero, err
	}

	lc.cache.Set(key, value)
	return value, nil
}

// GetWithContext retrieves an item with context support
func (lc *LoadingCache[K, V]) GetWithContext(ctx context.Context, key K) (V, error) {
	if value, found := lc.cache.Get(key); found {
		return value, nil
	}

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		var zero V
		return zero, ctx.Err()
	default:
	}

	value, err := lc.loader(key)
	if err != nil {
		var zero V
		return zero, err
	}

	lc.cache.Set(key, value)
	return value, nil
}

// Invalidate removes an item from the cache
func (lc *LoadingCache[K, V]) Invalidate(key K) {
	lc.cache.Delete(key)
}

// Stop stops the underlying cache
func (lc *LoadingCache[K, V]) Stop() {
	lc.cache.Stop()
}
