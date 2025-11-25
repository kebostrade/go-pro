package cache

import (
	"container/list"
	"sync"
)

// LRUCache implements a Least Recently Used cache with generic types
type LRUCache[K comparable, V any] struct {
	capacity  int
	items     map[K]*list.Element
	evictList *list.List
	mu        sync.RWMutex
	stats     *CacheStats
}

// entry represents a key-value pair in the LRU cache
type entry[K comparable, V any] struct {
	key   K
	value V
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity:  capacity,
		items:     make(map[K]*list.Element),
		evictList: list.New(),
		stats:     &CacheStats{},
	}
}

// Get retrieves a value from the cache
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.evictList.MoveToFront(elem)

		c.stats.mu.Lock()
		c.stats.Hits++
		c.stats.mu.Unlock()

		return elem.Value.(*entry[K, V]).value, true
	}

	c.stats.mu.Lock()
	c.stats.Misses++
	c.stats.mu.Unlock()

	var zero V
	return zero, false
}

// Set adds or updates a value in the cache
func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if key already exists
	if elem, found := c.items[key]; found {
		c.evictList.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value

		c.stats.mu.Lock()
		c.stats.Sets++
		c.stats.mu.Unlock()
		return
	}

	// Add new entry
	ent := &entry[K, V]{key: key, value: value}
	elem := c.evictList.PushFront(ent)
	c.items[key] = elem

	c.stats.mu.Lock()
	c.stats.Sets++
	c.stats.mu.Unlock()

	// Evict oldest if over capacity
	if c.evictList.Len() > c.capacity {
		c.removeOldest()
	}
}

// Delete removes a key from the cache
func (c *LRUCache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.removeElement(elem)

		c.stats.mu.Lock()
		c.stats.Deletes++
		c.stats.mu.Unlock()
		return true
	}

	return false
}

// Peek retrieves a value without updating access order
func (c *LRUCache[K, V]) Peek(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, found := c.items[key]; found {
		return elem.Value.(*entry[K, V]).value, true
	}

	var zero V
	return zero, false
}

// Contains checks if a key exists in the cache
func (c *LRUCache[K, V]) Contains(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, found := c.items[key]
	return found
}

// Size returns the number of items in the cache
func (c *LRUCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.evictList.Len()
}

// Capacity returns the cache capacity
func (c *LRUCache[K, V]) Capacity() int {
	return c.capacity
}

// Clear removes all items from the cache
func (c *LRUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*list.Element)
	c.evictList.Init()
}

// Keys returns all keys in the cache (most recent first)
func (c *LRUCache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, c.evictList.Len())
	for elem := c.evictList.Front(); elem != nil; elem = elem.Next() {
		keys = append(keys, elem.Value.(*entry[K, V]).key)
	}

	return keys
}

// GetOldest returns the oldest (least recently used) key-value pair
func (c *LRUCache[K, V]) GetOldest() (K, V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem := c.evictList.Back()
	if elem == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	ent := elem.Value.(*entry[K, V])
	return ent.key, ent.value, true
}

// GetNewest returns the newest (most recently used) key-value pair
func (c *LRUCache[K, V]) GetNewest() (K, V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem := c.evictList.Front()
	if elem == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	ent := elem.Value.(*entry[K, V])
	return ent.key, ent.value, true
}

// GetStats returns cache statistics
func (c *LRUCache[K, V]) GetStats() CacheStats {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()

	return CacheStats{
		Hits:      c.stats.Hits,
		Misses:    c.stats.Misses,
		Sets:      c.stats.Sets,
		Deletes:   c.stats.Deletes,
		Evictions: c.stats.Evictions,
	}
}

// HitRate returns the cache hit rate
func (c *LRUCache[K, V]) HitRate() float64 {
	stats := c.GetStats()
	total := stats.Hits + stats.Misses
	if total == 0 {
		return 0.0
	}
	return float64(stats.Hits) / float64(total)
}

// ResetStats resets cache statistics
func (c *LRUCache[K, V]) ResetStats() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()

	c.stats.Hits = 0
	c.stats.Misses = 0
	c.stats.Sets = 0
	c.stats.Deletes = 0
	c.stats.Evictions = 0
}

// removeOldest removes the oldest item from the cache
func (c *LRUCache[K, V]) removeOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)

		c.stats.mu.Lock()
		c.stats.Evictions++
		c.stats.mu.Unlock()
	}
}

// removeElement removes a specific element from the cache
func (c *LRUCache[K, V]) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)
	ent := elem.Value.(*entry[K, V])
	delete(c.items, ent.key)
}

// ForEach iterates over all items in the cache (most recent first)
func (c *LRUCache[K, V]) ForEach(fn func(K, V)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for elem := c.evictList.Front(); elem != nil; elem = elem.Next() {
		ent := elem.Value.(*entry[K, V])
		fn(ent.key, ent.value)
	}
}

// Resize changes the cache capacity
// If new capacity is smaller, oldest items are evicted
func (c *LRUCache[K, V]) Resize(newCapacity int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.capacity = newCapacity

	// Evict items if over new capacity
	for c.evictList.Len() > c.capacity {
		c.removeOldest()
	}
}

// LFUCache implements a Least Frequently Used cache
type LFUCache[K comparable, V any] struct {
	capacity int
	items    map[K]*lfuEntry[K, V]
	freqList map[int]*list.List
	minFreq  int
	mu       sync.RWMutex
	stats    *CacheStats
}

// lfuEntry represents an entry in the LFU cache
type lfuEntry[K comparable, V any] struct {
	key   K
	value V
	freq  int
	elem  *list.Element
}

// NewLFUCache creates a new LFU cache
func NewLFUCache[K comparable, V any](capacity int) *LFUCache[K, V] {
	return &LFUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*lfuEntry[K, V]),
		freqList: make(map[int]*list.List),
		stats:    &CacheStats{},
	}
}

// Get retrieves a value from the LFU cache
func (c *LFUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ent, found := c.items[key]; found {
		c.incrementFreq(ent)

		c.stats.mu.Lock()
		c.stats.Hits++
		c.stats.mu.Unlock()

		return ent.value, true
	}

	c.stats.mu.Lock()
	c.stats.Misses++
	c.stats.mu.Unlock()

	var zero V
	return zero, false
}

// Set adds or updates a value in the LFU cache
func (c *LFUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity <= 0 {
		return
	}

	// Update existing entry
	if ent, found := c.items[key]; found {
		ent.value = value
		c.incrementFreq(ent)

		c.stats.mu.Lock()
		c.stats.Sets++
		c.stats.mu.Unlock()
		return
	}

	// Evict if at capacity
	if len(c.items) >= c.capacity {
		c.evictLFU()
	}

	// Add new entry
	ent := &lfuEntry[K, V]{
		key:   key,
		value: value,
		freq:  1,
	}

	if c.freqList[1] == nil {
		c.freqList[1] = list.New()
	}

	ent.elem = c.freqList[1].PushFront(ent)
	c.items[key] = ent
	c.minFreq = 1

	c.stats.mu.Lock()
	c.stats.Sets++
	c.stats.mu.Unlock()
}

// incrementFreq increases the frequency of an entry
func (c *LFUCache[K, V]) incrementFreq(ent *lfuEntry[K, V]) {
	freq := ent.freq
	c.freqList[freq].Remove(ent.elem)

	if c.freqList[freq].Len() == 0 {
		delete(c.freqList, freq)
		if c.minFreq == freq {
			c.minFreq++
		}
	}

	ent.freq++
	if c.freqList[ent.freq] == nil {
		c.freqList[ent.freq] = list.New()
	}

	ent.elem = c.freqList[ent.freq].PushFront(ent)
}

// evictLFU removes the least frequently used item
func (c *LFUCache[K, V]) evictLFU() {
	if lst := c.freqList[c.minFreq]; lst != nil {
		elem := lst.Back()
		if elem != nil {
			ent := elem.Value.(*lfuEntry[K, V])
			lst.Remove(elem)
			delete(c.items, ent.key)

			c.stats.mu.Lock()
			c.stats.Evictions++
			c.stats.mu.Unlock()
		}
	}
}

// Size returns the number of items in the cache
func (c *LFUCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
