// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package cache provides functionality for the GO-PRO Learning Platform.
package cache

import (
	"context"
	"errors"
	"time"
)

// Common cache errors.
var (
	ErrCacheMiss     = errors.New("cache miss")
	ErrCacheTimeout  = errors.New("cache timeout")
	ErrCacheNotFound = errors.New("cache key not found")
)

// Cache defines the interface for caching operations.
type Cache interface {
	// Basic operations.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Expiration operations.
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Atomic operations.
	Increment(ctx context.Context, key string) (int64, error)
	IncrementBy(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	DecrementBy(ctx context.Context, key string, value int64) (int64, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	GetSet(ctx context.Context, key string, value interface{}) (string, error)

	// Batch operations.
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	MSet(ctx context.Context, pairs ...interface{}) error

	// Key operations.
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)

	// Utility operations.
	FlushDB(ctx context.Context) error
	HealthCheck(ctx context.Context) error
	Close() error
}

// SessionStore defines the interface for session storage.
type SessionStore interface {
	// Session operations.
	CreateSession(ctx context.Context, sessionID string, data map[string]interface{}, expiration time.Duration) error
	GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error)
	UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error
	DeleteSession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string, expiration time.Duration) error

	// Session management.
	ListUserSessions(ctx context.Context, userID string) ([]string, error)
	DeleteUserSessions(ctx context.Context, userID string) error
	CleanupExpiredSessions(ctx context.Context) error
}

// DistributedLock defines the interface for distributed locking.
type DistributedLock interface {
	// Lock operations.
	Lock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
	Extend(ctx context.Context, key string, expiration time.Duration) error

	// Lock information.
	IsLocked(ctx context.Context, key string) (bool, error)
	GetLockTTL(ctx context.Context, key string) (time.Duration, error)
}

// RateLimiter defines the interface for rate limiting.
type RateLimiter interface {
	// Rate limiting operations.
	Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error)
	AllowN(ctx context.Context, key string, n, limit int64, window time.Duration) (bool, error)

	// Rate limit information.
	Remaining(ctx context.Context, key string, limit int64, window time.Duration) (int64, error)
	Reset(ctx context.Context, key string) error
}

// PubSub defines the interface for publish/subscribe operations.
type PubSub interface {
	// Publishing.
	Publish(ctx context.Context, channel string, message interface{}) error

	// Subscribing.
	Subscribe(ctx context.Context, channels ...string) (<-chan Message, error)
	Unsubscribe(ctx context.Context, channels ...string) error

	// Pattern subscribing.
	PSubscribe(ctx context.Context, patterns ...string) (<-chan Message, error)
	PUnsubscribe(ctx context.Context, patterns ...string) error

	// Close subscription.
	Close() error
}

// Message represents a pub/sub message.
type Message struct {
	Channel string
	Pattern string
	Payload string
}

// CacheManager combines all cache-related interfaces.
type CacheManager interface {
	Cache
	SessionStore
	DistributedLock
	RateLimiter
	PubSub
}

// CacheStats represents cache statistics.
type CacheStats struct {
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
	Sets        int64     `json:"sets"`
	Deletes     int64     `json:"deletes"`
	Errors      int64     `json:"errors"`
	Connections int       `json:"connections"`
	LastReset   time.Time `json:"last_reset"`
}

// CacheMetrics defines the interface for cache metrics.
type CacheMetrics interface {
	// Metrics operations.
	GetStats() *CacheStats
	ResetStats()

	// Metric recording.
	RecordHit()
	RecordMiss()
	RecordSet()
	RecordDelete()
	RecordError()
}

// CacheOptions represents options for cache operations.
type CacheOptions struct {
	Expiration time.Duration
	Tags       []string
	Namespace  string
	Compress   bool
	Serialize  bool
}

// TaggedCache defines the interface for tagged caching.
type TaggedCache interface {
	Cache

	// Tagged operations.
	SetWithTags(ctx context.Context, key string, value interface{}, tags []string, expiration time.Duration) error
	InvalidateTag(ctx context.Context, tag string) error
	InvalidateTags(ctx context.Context, tags []string) error
	GetKeysByTag(ctx context.Context, tag string) ([]string, error)
}

// NamespacedCache defines the interface for namespaced caching.
type NamespacedCache interface {
	Cache

	// Namespace operations.
	SetNamespace(namespace string)
	GetNamespace() string
	FlushNamespace(ctx context.Context, namespace string) error
}

// CompressedCache defines the interface for compressed caching.
type CompressedCache interface {
	Cache

	// Compression operations.
	SetCompressed(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetCompressed(ctx context.Context, key string, dest interface{}) error
}

// SerializedCache defines the interface for serialized caching.
type SerializedCache interface {
	Cache

	// Serialization operations.
	SetSerialized(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetSerialized(ctx context.Context, key string, dest interface{}) error
}

// CacheFactory defines the interface for creating cache instances.
type CacheFactory interface {
	// Cache creation.
	CreateCache(config interface{}) (Cache, error)
	CreateSessionStore(config interface{}) (SessionStore, error)
	CreateDistributedLock(config interface{}) (DistributedLock, error)
	CreateRateLimiter(config interface{}) (RateLimiter, error)
	CreatePubSub(config interface{}) (PubSub, error)

	// Health check.
	HealthCheck(ctx context.Context) error
}

// CacheMiddleware defines the interface for cache middleware.
type CacheMiddleware interface {
	// Middleware operations.
	Before(ctx context.Context, operation, key string, value interface{}) error
	After(ctx context.Context, operation, key string, value interface{}, err error) error
}

// CacheHook defines the interface for cache hooks.
type CacheHook interface {
	// Hook operations.
	OnSet(ctx context.Context, key string, value interface{}) error
	OnGet(ctx context.Context, key string, value interface{}) error
	OnDelete(ctx context.Context, key string) error
	OnExpire(ctx context.Context, key string) error
}

// CacheObserver defines the interface for cache observation.
type CacheObserver interface {
	// Observation operations.
	NotifySet(ctx context.Context, key string, value interface{})
	NotifyGet(ctx context.Context, key string, hit bool)
	NotifyDelete(ctx context.Context, key string)
	NotifyExpire(ctx context.Context, key string)
	NotifyError(ctx context.Context, operation, key string, err error)
}
