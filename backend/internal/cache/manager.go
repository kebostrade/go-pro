// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package cache provides functionality for the GO-PRO Learning Platform.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Manager combines all cache-related functionality.
type Manager struct {
	cache       *RedisCache
	sessions    *RedisSessionStore
	locks       *RedisDistributedLock
	rateLimiter *RedisRateLimiter
	pubsub      *RedisPubSub
	client      *redis.Client
	config      *CacheConfig
}

// ManagerConfig holds configuration for the cache manager.
type ManagerConfig struct {
	Redis     *CacheConfig
	Sessions  SessionConfig
	Locks     LockConfig
	RateLimit RateLimitConfig
	PubSub    PubSubConfig
}

// SessionConfig holds session-specific configuration.
type SessionConfig struct {
	Prefix            string
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

// LockConfig holds lock-specific configuration.
type LockConfig struct {
	Prefix          string
	DefaultTimeout  time.Duration
	CleanupInterval time.Duration
}

// RateLimitConfig holds rate limit-specific configuration.
type RateLimitConfig struct {
	Prefix          string
	CleanupInterval time.Duration
}

// PubSubConfig holds pub/sub-specific configuration.
type PubSubConfig struct {
	Prefix string
}

// DefaultManagerConfig returns default manager configuration.
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		Redis: DefaultCacheConfig(),
		Sessions: SessionConfig{
			Prefix:            "session:",
			DefaultExpiration: 24 * time.Hour,
			CleanupInterval:   time.Hour,
		},
		Locks: LockConfig{
			Prefix:          "lock:",
			DefaultTimeout:  30 * time.Second,
			CleanupInterval: 5 * time.Minute,
		},
		RateLimit: RateLimitConfig{
			Prefix:          "ratelimit:",
			CleanupInterval: 10 * time.Minute,
		},
		PubSub: PubSubConfig{
			Prefix: "pubsub:",
		},
	}
}

// NewManager creates a new cache manager.
func NewManager(config *ManagerConfig) (*Manager, error) {
	if config == nil {
		config = DefaultManagerConfig()
	}

	// Create Redis cache.
	cache, err := NewRedisCache(config.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis cache: %w", err)
	}

	// Get the underlying Redis client.
	client := cache.Client()

	// Create components.
	sessions := NewRedisSessionStore(client, config.Sessions.Prefix)
	locks := NewRedisDistributedLock(client, config.Locks.Prefix)
	rateLimiter := NewRedisRateLimiter(client, config.RateLimit.Prefix)
	pubsub := NewRedisPubSub(client, config.PubSub.Prefix)

	manager := &Manager{
		cache:       cache,
		sessions:    sessions,
		locks:       locks,
		rateLimiter: rateLimiter,
		pubsub:      pubsub,
		client:      client,
		config:      config.Redis,
	}

	// Start background cleanup routines.
	go manager.startCleanupRoutines(config)

	return manager, nil
}

// NewManagerFromEnv creates a cache manager from environment variables.
func NewManagerFromEnv() (*Manager, error) {
	config := DefaultManagerConfig()

	// Override with environment variables.
	config.Redis.Host = getEnvString("REDIS_HOST", config.Redis.Host)
	config.Redis.Port = getEnvInt("REDIS_PORT", config.Redis.Port)
	config.Redis.Password = getEnvString("REDIS_PASSWORD", config.Redis.Password)
	config.Redis.Database = getEnvInt("REDIS_DB", config.Redis.Database)
	config.Redis.Prefix = getEnvString("REDIS_PREFIX", config.Redis.Prefix)

	config.Sessions.Prefix = getEnvString("REDIS_SESSION_PREFIX", config.Sessions.Prefix)
	config.Sessions.DefaultExpiration = getEnvDuration("SESSION_DEFAULT_EXPIRATION", config.Sessions.DefaultExpiration)

	config.Locks.Prefix = getEnvString("REDIS_LOCK_PREFIX", config.Locks.Prefix)
	config.Locks.DefaultTimeout = getEnvDuration("LOCK_DEFAULT_TIMEOUT", config.Locks.DefaultTimeout)

	config.RateLimit.Prefix = getEnvString("REDIS_RATELIMIT_PREFIX", config.RateLimit.Prefix)

	return NewManager(config)
}

// Cache returns the cache interface.
func (m *Manager) Cache() Cache {
	return m.cache
}

// Sessions returns the session store interface.
func (m *Manager) Sessions() SessionStore {
	return m.sessions
}

// Locks returns the distributed lock interface.
func (m *Manager) Locks() DistributedLock {
	return m.locks
}

// RateLimiter returns the rate limiter interface.
func (m *Manager) RateLimiter() RateLimiter {
	return m.rateLimiter
}

// PubSub returns the pub/sub interface.
func (m *Manager) PubSub() PubSub {
	return m.pubsub
}

// Client returns the underlying Redis client.
func (m *Manager) Client() *redis.Client {
	return m.client
}

// HealthCheck checks the health of all cache components.
func (m *Manager) HealthCheck(ctx context.Context) error {
	// Check Redis connection.
	if err := m.cache.HealthCheck(ctx); err != nil {
		return fmt.Errorf("cache health check failed: %w", err)
	}

	// Test basic operations.
	testKey := "health_check_" + fmt.Sprintf("%d", time.Now().UnixNano())

	// Test cache.
	if err := m.cache.Set(ctx, testKey, "test", time.Minute); err != nil {
		return fmt.Errorf("cache set operation failed: %w", err)
	}

	var value string
	if err := m.cache.Get(ctx, testKey, &value); err != nil {
		return fmt.Errorf("cache get operation failed: %w", err)
	}

	if err := m.cache.Delete(ctx, testKey); err != nil {
		return fmt.Errorf("cache delete operation failed: %w", err)
	}

	// Test session store.
	sessionData := map[string]interface{}{
		"user_id": "test_user",
		"test":    true,
	}
	if err := m.sessions.CreateSession(ctx, testKey, sessionData, time.Minute); err != nil {
		return fmt.Errorf("session create operation failed: %w", err)
	}

	if err := m.sessions.DeleteSession(ctx, testKey); err != nil {
		return fmt.Errorf("session delete operation failed: %w", err)
	}

	// Test distributed lock.
	acquired, err := m.locks.Lock(ctx, testKey, time.Minute)
	if err != nil {
		return fmt.Errorf("lock acquire operation failed: %w", err)
	}
	if !acquired {
		return fmt.Errorf("lock acquire operation returned false")
	}

	if err := m.locks.Unlock(ctx, testKey); err != nil {
		return fmt.Errorf("lock release operation failed: %w", err)
	}

	// Test rate limiter.
	allowed, err := m.rateLimiter.Allow(ctx, testKey, 10, time.Minute)
	if err != nil {
		return fmt.Errorf("rate limiter operation failed: %w", err)
	}
	if !allowed {
		return fmt.Errorf("rate limiter operation returned false")
	}

	if err := m.rateLimiter.Reset(ctx, testKey); err != nil {
		return fmt.Errorf("rate limiter reset operation failed: %w", err)
	}

	return nil
}

// Stats returns comprehensive statistics about the cache manager.
func (m *Manager) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Redis connection stats.
	poolStats := m.cache.Stats()
	stats["redis_pool"] = map[string]interface{}{
		"hits":        poolStats.Hits,
		"misses":      poolStats.Misses,
		"timeouts":    poolStats.Timeouts,
		"total_conns": poolStats.TotalConns,
		"idle_conns":  poolStats.IdleConns,
		"stale_conns": poolStats.StaleConns,
	}

	// Redis info.
	info, err := m.client.Info(ctx, "memory", "stats", "keyspace").Result()
	if err == nil {
		stats["redis_info"] = info
	}

	// Session stats.
	sessionCount, err := m.sessions.GetSessionCount(ctx)
	if err == nil {
		stats["sessions"] = map[string]interface{}{
			"total_sessions": sessionCount,
		}
	}

	// Lock stats.
	locks, err := m.locks.ListLocks(ctx, "*")
	if err == nil {
		stats["locks"] = map[string]interface{}{
			"active_locks": len(locks),
		}
	}

	// Rate limit stats.
	rateLimits, err := m.rateLimiter.ListRateLimits(ctx, "*")
	if err == nil {
		stats["rate_limits"] = map[string]interface{}{
			"active_rate_limits": len(rateLimits),
		}
	}

	return stats, nil
}

// Close closes all connections and stops background routines.
func (m *Manager) Close() error {
	// Close pub/sub first.
	if err := m.pubsub.Close(); err != nil {
		return fmt.Errorf("failed to close pub/sub: %w", err)
	}

	// Close cache (which closes the Redis client)
	if err := m.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache: %w", err)
	}

	return nil
}

// startCleanupRoutines starts background cleanup routines.
func (m *Manager) startCleanupRoutines(config *ManagerConfig) {
	// Session cleanup.
	go func() {
		ticker := time.NewTicker(config.Sessions.CleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := m.sessions.CleanupExpiredSessions(ctx); err != nil {
				// Log error (in production, use proper logging)
				fmt.Printf("Session cleanup error: %v\n", err)
			}
			cancel()
		}
	}()

	// Lock cleanup.
	go func() {
		ticker := time.NewTicker(config.Locks.CleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := m.locks.CleanupExpiredLocks(ctx); err != nil {
				// Log error (in production, use proper logging)
				fmt.Printf("Lock cleanup error: %v\n", err)
			}
			cancel()
		}
	}()

	// Rate limit cleanup.
	go func() {
		ticker := time.NewTicker(config.RateLimit.CleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := m.rateLimiter.CleanupExpiredRateLimits(ctx); err != nil {
				// Log error (in production, use proper logging)
				fmt.Printf("Rate limit cleanup error: %v\n", err)
			}
			cancel()
		}
	}()
}

// Flush clears all data from all cache components.
func (m *Manager) Flush(ctx context.Context) error {
	return m.cache.FlushDB(ctx)
}

// Set stores a value in the cache.
func (m *Manager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.cache.Set(ctx, key, value, expiration)
}

// Get retrieves a value from the cache.
func (m *Manager) Get(ctx context.Context, key string, dest interface{}) error {
	return m.cache.Get(ctx, key, dest)
}

// Delete removes a value from the cache.
func (m *Manager) Delete(ctx context.Context, key string) error {
	return m.cache.Delete(ctx, key)
}

// Exists checks if a key exists in the cache.
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	return m.cache.Exists(ctx, key)
}

// Expire sets the expiration time for a key.
func (m *Manager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return m.cache.Expire(ctx, key, expiration)
}

// TTL returns the time to live for a key.
func (m *Manager) TTL(ctx context.Context, key string) (time.Duration, error) {
	return m.cache.TTL(ctx, key)
}

// Increment increments a counter.
func (m *Manager) Increment(ctx context.Context, key string) (int64, error) {
	return m.cache.Increment(ctx, key)
}

// IncrementBy increments a counter by a specific value.
func (m *Manager) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return m.cache.IncrementBy(ctx, key, value)
}

// Decrement decrements a counter.
func (m *Manager) Decrement(ctx context.Context, key string) (int64, error) {
	return m.cache.Decrement(ctx, key)
}

// DecrementBy decrements a counter by a specific value.
func (m *Manager) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return m.cache.DecrementBy(ctx, key, value)
}

// SetNX sets a value only if the key does not exist.
func (m *Manager) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return m.cache.SetNX(ctx, key, value, expiration)
}

// GetSet sets a new value and returns the old value.
func (m *Manager) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return m.cache.GetSet(ctx, key, value)
}

// MGet retrieves multiple values.
func (m *Manager) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return m.cache.MGet(ctx, keys...)
}

// MSet sets multiple values.
func (m *Manager) MSet(ctx context.Context, pairs ...interface{}) error {
	return m.cache.MSet(ctx, pairs...)
}

// Keys returns all keys matching a pattern.
func (m *Manager) Keys(ctx context.Context, pattern string) ([]string, error) {
	return m.cache.Keys(ctx, pattern)
}

// Scan scans keys matching a pattern.
func (m *Manager) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return m.cache.Scan(ctx, cursor, match, count)
}

// FlushDB clears all data from the cache.
func (m *Manager) FlushDB(ctx context.Context) error {
	return m.cache.FlushDB(ctx)
}

// Lock acquires a distributed lock.
func (m *Manager) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return m.locks.Lock(ctx, key, expiration)
}

// Unlock releases a distributed lock.
func (m *Manager) Unlock(ctx context.Context, key string) error {
	return m.locks.Unlock(ctx, key)
}

// Extend extends the expiration time of a lock.
func (m *Manager) Extend(ctx context.Context, key string, expiration time.Duration) error {
	return m.locks.Extend(ctx, key, expiration)
}

// IsLocked checks if a key is locked.
func (m *Manager) IsLocked(ctx context.Context, key string) (bool, error) {
	return m.locks.IsLocked(ctx, key)
}

// GetLockTTL returns the time to live for a lock.
func (m *Manager) GetLockTTL(ctx context.Context, key string) (time.Duration, error) {
	return m.locks.GetLockTTL(ctx, key)
}

// Allow checks if a request is allowed under the rate limit.
func (m *Manager) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	return m.rateLimiter.Allow(ctx, key, limit, window)
}

// AllowN checks if N requests are allowed under the rate limit.
func (m *Manager) AllowN(ctx context.Context, key string, n, limit int64, window time.Duration) (bool, error) {
	return m.rateLimiter.AllowN(ctx, key, n, limit, window)
}

// Remaining returns the number of requests remaining in the current window.
func (m *Manager) Remaining(ctx context.Context, key string, limit int64, window time.Duration) (int64, error) {
	return m.rateLimiter.Remaining(ctx, key, limit, window)
}

// Reset resets the rate limit for a key.
func (m *Manager) Reset(ctx context.Context, key string) error {
	return m.rateLimiter.Reset(ctx, key)
}

// CreateSession creates a new session.
func (m *Manager) CreateSession(ctx context.Context, sessionID string, data map[string]interface{}, expiration time.Duration) error {
	return m.sessions.CreateSession(ctx, sessionID, data, expiration)
}

// GetSession retrieves a session.
func (m *Manager) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	return m.sessions.GetSession(ctx, sessionID)
}

// UpdateSession updates a session.
func (m *Manager) UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	return m.sessions.UpdateSession(ctx, sessionID, data)
}

// DeleteSession deletes a session.
func (m *Manager) DeleteSession(ctx context.Context, sessionID string) error {
	return m.sessions.DeleteSession(ctx, sessionID)
}

// RefreshSession refreshes a session's expiration.
func (m *Manager) RefreshSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	return m.sessions.RefreshSession(ctx, sessionID, expiration)
}

// ListUserSessions lists all sessions for a user.
func (m *Manager) ListUserSessions(ctx context.Context, userID string) ([]string, error) {
	return m.sessions.ListUserSessions(ctx, userID)
}

// DeleteUserSessions deletes all sessions for a user.
func (m *Manager) DeleteUserSessions(ctx context.Context, userID string) error {
	return m.sessions.DeleteUserSessions(ctx, userID)
}

// CleanupExpiredSessions removes expired sessions.
func (m *Manager) CleanupExpiredSessions(ctx context.Context) error {
	return m.sessions.CleanupExpiredSessions(ctx)
}

// Publish publishes a message to a channel.
func (m *Manager) Publish(ctx context.Context, channel string, message interface{}) error {
	return m.pubsub.Publish(ctx, channel, message)
}

// Subscribe subscribes to channels.
func (m *Manager) Subscribe(ctx context.Context, channels ...string) (<-chan Message, error) {
	return m.pubsub.Subscribe(ctx, channels...)
}

// Unsubscribe unsubscribes from channels.
func (m *Manager) Unsubscribe(ctx context.Context, channels ...string) error {
	return m.pubsub.Unsubscribe(ctx, channels...)
}

// PSubscribe subscribes to channel patterns.
func (m *Manager) PSubscribe(ctx context.Context, patterns ...string) (<-chan Message, error) {
	return m.pubsub.PSubscribe(ctx, patterns...)
}

// PUnsubscribe unsubscribes from channel patterns.
func (m *Manager) PUnsubscribe(ctx context.Context, patterns ...string) error {
	return m.pubsub.PUnsubscribe(ctx, patterns...)
}

// FlushPattern clears all keys matching a pattern.
func (m *Manager) FlushPattern(ctx context.Context, pattern string) error {
	keys, err := m.cache.Keys(ctx, pattern)
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete keys in batches.
	batchSize := 100
	for i := 0; i < len(keys); i += batchSize {
		end := i + batchSize
		if end > len(keys) {
			end = len(keys)
		}

		batch := keys[i:end]
		pipe := m.client.Pipeline()
		for _, key := range batch {
			pipe.Del(ctx, m.cache.key(key))
		}

		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("failed to delete batch of keys: %w", err)
		}
	}

	return nil
}
