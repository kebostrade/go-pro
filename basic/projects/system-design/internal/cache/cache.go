// Package cache provides a simple in-memory cache implementation.
package cache

import (
	"sync"
	"time"
)

// Cache is a simple in-memory cache with expiration.
type Cache struct {
	items      map[string]*cacheItem
	mu         sync.RWMutex
	defaultTTL time.Duration
	maxSize    int
	onEvict    func(key string, value interface{})
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// DefaultExpiration is the default TTL for cache items.
const DefaultExpiration = 10 * time.Minute

// MaxSize is the maximum number of items in the cache.
const MaxSize = 10000

// New creates a new cache.
func New(defaultTTL time.Duration, maxSize int) *Cache {
	return &Cache{
		items:      make(map[string]*cacheItem),
		defaultTTL: defaultTTL,
		maxSize:    maxSize,
	}
}

// Set adds an item to the cache.
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL adds an item to the cache with a specific TTL.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict if at capacity
	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}

	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Get retrieves an item from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
}

// Size returns the number of items in the cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// evictOldest removes the oldest item from the cache.
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestTime.IsZero() || item.expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.expiration
		}
	}

	if oldestKey != "" {
		if c.onEvict != nil {
			if item, exists := c.items[oldestKey]; exists {
				c.onEvict(oldestKey, item.value)
			}
		}
		delete(c.items, oldestKey)
	}
}

// OnEvict sets a callback to be called when items are evicted.
func (c *Cache) OnEvict(fn func(key string, value interface{})) {
	c.onEvict = fn
}
