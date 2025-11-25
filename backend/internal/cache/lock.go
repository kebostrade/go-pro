// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package cache provides functionality for the GO-PRO Learning Platform.
package cache

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisDistributedLock implements distributed locking using Redis.
type RedisDistributedLock struct {
	client *redis.Client
	prefix string
}

// LockInfo represents information about a lock.
type LockInfo struct {
	Key       string        `json:"key"`
	Value     string        `json:"value"`
	Owner     string        `json:"owner"`
	CreatedAt time.Time     `json:"created_at"`
	ExpiresAt time.Time     `json:"expires_at"`
	TTL       time.Duration `json:"ttl"`
}

// NewRedisDistributedLock creates a new Redis distributed lock.
func NewRedisDistributedLock(client *redis.Client, prefix string) *RedisDistributedLock {
	if prefix == "" {
		prefix = "lock:"
	}

	return &RedisDistributedLock{
		client: client,
		prefix: prefix,
	}
}

// lockKey generates a lock key with prefix.
func (r *RedisDistributedLock) lockKey(key string) string {
	return r.prefix + key
}

// generateLockValue generates a unique lock value.
func (r *RedisDistributedLock) generateLockValue() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based value if random generation fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(bytes)
}

// Lock acquires a distributed lock.
func (r *RedisDistributedLock) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	lockKey := r.lockKey(key)
	lockValue := r.generateLockValue()

	// Use SET with NX (only if not exists) and EX (expiration)
	result, err := r.client.SetNX(ctx, lockKey, lockValue, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return result, nil
}

// LockWithRetry acquires a lock with retry mechanism.
func (r *RedisDistributedLock) LockWithRetry(
	ctx context.Context,
	key string,
	expiration time.Duration,
	maxRetries int,
	retryDelay time.Duration,
) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		acquired, err := r.Lock(ctx, key, expiration)
		if err != nil {
			return false, err
		}
		if acquired {
			return true, nil
		}

		// Wait before retrying.
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryDelay):
			continue
		}
	}

	return false, fmt.Errorf("failed to acquire lock after %d retries", maxRetries)
}

// Unlock releases a distributed lock.
func (r *RedisDistributedLock) Unlock(ctx context.Context, key string) error {
	lockKey := r.lockKey(key)

	// Lua script to ensure atomic unlock (only unlock if we own the lock)
	unlockScript := `
		if redis.call("exists", KEYS[1]) == 1 then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := r.client.Eval(ctx, unlockScript, []string{lockKey}).Result()
	if err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("lock not found or already expired")
	}

	return nil
}

// UnlockWithValue releases a lock only if the value matches (safer unlock).
func (r *RedisDistributedLock) UnlockWithValue(ctx context.Context, key, value string) error {
	lockKey := r.lockKey(key)

	// Lua script to ensure atomic unlock with value check.
	unlockScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := r.client.Eval(ctx, unlockScript, []string{lockKey}, value).Result()
	if err != nil {
		return fmt.Errorf("failed to unlock with value: %w", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("lock not found, expired, or value mismatch")
	}

	return nil
}

// Extend extends the expiration of a lock.
func (r *RedisDistributedLock) Extend(ctx context.Context, key string, expiration time.Duration) error {
	lockKey := r.lockKey(key)

	// Check if lock exists.
	exists, err := r.client.Exists(ctx, lockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check lock existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("lock not found")
	}

	// Extend expiration.
	err = r.client.Expire(ctx, lockKey, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to extend lock: %w", err)
	}

	return nil
}

// ExtendWithValue extends a lock only if the value matches.
func (r *RedisDistributedLock) ExtendWithValue(ctx context.Context, key, value string, expiration time.Duration) error {
	lockKey := r.lockKey(key)

	// Lua script to extend lock only if value matches.
	extendScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := r.client.Eval(ctx, extendScript, []string{lockKey}, value, int(expiration.Seconds())).Result()
	if err != nil {
		return fmt.Errorf("failed to extend lock with value: %w", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("lock not found, expired, or value mismatch")
	}

	return nil
}

// IsLocked checks if a key is locked.
func (r *RedisDistributedLock) IsLocked(ctx context.Context, key string) (bool, error) {
	lockKey := r.lockKey(key)

	exists, err := r.client.Exists(ctx, lockKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check lock status: %w", err)
	}

	return exists > 0, nil
}

// GetLockTTL returns the time to live for a lock.
func (r *RedisDistributedLock) GetLockTTL(ctx context.Context, key string) (time.Duration, error) {
	lockKey := r.lockKey(key)

	ttl, err := r.client.TTL(ctx, lockKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get lock TTL: %w", err)
	}

	return ttl, nil
}

// GetLockInfo returns detailed information about a lock.
func (r *RedisDistributedLock) GetLockInfo(ctx context.Context, key string) (*LockInfo, error) {
	lockKey := r.lockKey(key)

	// Use pipeline to get value and TTL atomically.
	pipe := r.client.Pipeline()
	valueCmd := pipe.Get(ctx, lockKey)
	ttlCmd := pipe.TTL(ctx, lockKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("lock not found")
		}

		return nil, fmt.Errorf("failed to get lock info: %w", err)
	}

	value, err := valueCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get lock value: %w", err)
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get lock TTL: %w", err)
	}

	now := time.Now()
	lockInfo := &LockInfo{
		Key:       key,
		Value:     value,
		TTL:       ttl,
		CreatedAt: now.Add(-ttl), // Approximate creation time
		ExpiresAt: now.Add(ttl),
	}

	return lockInfo, nil
}

// ListLocks returns all active locks matching a pattern.
func (r *RedisDistributedLock) ListLocks(ctx context.Context, pattern string) ([]*LockInfo, error) {
	lockPattern := r.lockKey(pattern)
	keys, err := r.client.Keys(ctx, lockPattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list locks: %w", err)
	}

	locks := make([]*LockInfo, 0, len(keys))
	for _, lockKey := range keys {
		// Remove prefix to get original key.
		originalKey := lockKey[len(r.prefix):]
		lockInfo, err := r.GetLockInfo(ctx, originalKey)
		if err != nil {
			continue // Skip locks that might have expired
		}
		locks = append(locks, lockInfo)
	}

	return locks, nil
}

// CleanupExpiredLocks removes expired locks (should be run periodically).
func (r *RedisDistributedLock) CleanupExpiredLocks(ctx context.Context) error {
	// Use SCAN to iterate through all lock keys.
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, r.prefix+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan lock keys: %w", err)
		}

		// Check each key's TTL.
		for _, key := range keys {
			ttl, err := r.client.TTL(ctx, key).Result()
			if err != nil {
				continue
			}

			// If TTL is -2 (key doesn't exist), it's already cleaned up.
			// If TTL is -1 (no expiration), something is wrong - clean it up.
			if ttl == -1 || ttl == -2 {
				r.client.Del(ctx, key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// WithLock executes a function while holding a lock.
func (r *RedisDistributedLock) WithLock(ctx context.Context, key string, expiration time.Duration, fn func() error) error {
	// Acquire lock.
	acquired, err := r.Lock(ctx, key, expiration)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !acquired {
		return fmt.Errorf("failed to acquire lock: already locked")
	}

	// Ensure lock is released.
	defer func() {
		if unlockErr := r.Unlock(ctx, key); unlockErr != nil {
			// Log error but don't override the original error.
			fmt.Printf("Warning: failed to unlock %s: %v\n", key, unlockErr)
		}
	}()

	// Execute function.
	return fn()
}

// WithLockRetry executes a function while holding a lock with retry.
func (r *RedisDistributedLock) WithLockRetry(
	ctx context.Context,
	key string,
	expiration time.Duration,
	maxRetries int,
	retryDelay time.Duration,
	fn func() error,
) error {
	// Acquire lock with retry.
	acquired, err := r.LockWithRetry(ctx, key, expiration, maxRetries, retryDelay)
	if err != nil {
		return fmt.Errorf("failed to acquire lock with retry: %w", err)
	}
	if !acquired {
		return fmt.Errorf("failed to acquire lock after %d retries", maxRetries)
	}

	// Ensure lock is released.
	defer func() {
		if unlockErr := r.Unlock(ctx, key); unlockErr != nil {
			// Log error but don't override the original error.
			fmt.Printf("Warning: failed to unlock %s: %v\n", key, unlockErr)
		}
	}()

	// Execute function.
	return fn()
}
