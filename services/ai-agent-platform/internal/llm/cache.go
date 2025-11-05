package llm

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// ResponseCache caches LLM responses
type ResponseCache struct {
	cache map[string]*CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// CacheEntry represents a cached response
type CacheEntry struct {
	Response  *types.LLMResponse
	Timestamp time.Time
}

// NewResponseCache creates a new response cache
func NewResponseCache() *ResponseCache {
	cache := &ResponseCache{
		cache: make(map[string]*CacheEntry),
		ttl:   1 * time.Hour,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached response
func (rc *ResponseCache) Get(ctx context.Context, request types.LLMRequest) (*types.LLMResponse, bool) {
	key := rc.generateKey(request)

	rc.mu.RLock()
	defer rc.mu.RUnlock()

	entry, exists := rc.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Since(entry.Timestamp) > rc.ttl {
		return nil, false
	}

	return entry.Response, true
}

// Set stores a response in the cache
func (rc *ResponseCache) Set(ctx context.Context, request types.LLMRequest, response *types.LLMResponse) {
	key := rc.generateKey(request)

	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache[key] = &CacheEntry{
		Response:  response,
		Timestamp: time.Now(),
	}
}

// Clear clears all cached responses
func (rc *ResponseCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache = make(map[string]*CacheEntry)
}

// SetTTL sets the time-to-live for cache entries
func (rc *ResponseCache) SetTTL(ttl time.Duration) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.ttl = ttl
}

// generateKey generates a cache key from a request
func (rc *ResponseCache) generateKey(request types.LLMRequest) string {
	// Create a deterministic representation of the request
	data := struct {
		Messages    []types.Message
		Model       string
		Temperature float32
		MaxTokens   int
	}{
		Messages:    request.Messages,
		Model:       request.Model,
		Temperature: request.Temperature,
		MaxTokens:   request.MaxTokens,
	}

	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash)
}

// cleanup periodically removes expired entries
func (rc *ResponseCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rc.mu.Lock()
		now := time.Now()
		for key, entry := range rc.cache {
			if now.Sub(entry.Timestamp) > rc.ttl {
				delete(rc.cache, key)
			}
		}
		rc.mu.Unlock()
	}
}

// Stats returns cache statistics
func (rc *ResponseCache) Stats() CacheStats {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return CacheStats{
		Size: len(rc.cache),
		TTL:  rc.ttl,
	}
}

// CacheStats represents cache statistics
type CacheStats struct {
	Size int
	TTL  time.Duration
}

