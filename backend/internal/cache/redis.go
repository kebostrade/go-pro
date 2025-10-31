// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package cache provides functionality for the GO-PRO Learning Platform.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache implements caching using Redis.
type RedisCache struct {
	client *redis.Client
	prefix string
}

// CacheConfig holds Redis cache configuration.
type CacheConfig struct {
	Host         string
	Port         int
	Password     string
	Database     int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	IdleTimeout  time.Duration
	Prefix       string
}

// DefaultCacheConfig returns default cache configuration.
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		Database:     0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
		IdleTimeout:  5 * time.Minute,
		Prefix:       "gopro:",
	}
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(config *CacheConfig) (*RedisCache, error) {
	if config == nil {
		config = DefaultCacheConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Database,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolTimeout:  config.PoolTimeout,
		IdleTimeout:  config.IdleTimeout,
	})

	// Test connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		prefix: config.Prefix,
	}, nil
}

// NewRedisCacheFromEnv creates a Redis cache from environment variables.
func NewRedisCacheFromEnv() (*RedisCache, error) {
	config := &CacheConfig{
		Host:         getEnvString("REDIS_HOST", "localhost"),
		Port:         getEnvInt("REDIS_PORT", 6379),
		Password:     getEnvString("REDIS_PASSWORD", ""),
		Database:     getEnvInt("REDIS_DB", 0),
		PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
		MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 5),
		MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
		DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		PoolTimeout:  getEnvDuration("REDIS_POOL_TIMEOUT", 4*time.Second),
		IdleTimeout:  getEnvDuration("REDIS_IDLE_TIMEOUT", 5*time.Minute),
		Prefix:       getEnvString("REDIS_PREFIX", "gopro:"),
	}

	return NewRedisCache(config)
}

// key adds the prefix to the cache key.
func (r *RedisCache) key(key string) string {
	return r.prefix + key
}

// Set stores a value in the cache with expiration.
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(ctx, r.key(key), data, expiration).Err()
}

// Get retrieves a value from the cache.
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, r.key(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}

		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a value from the cache.
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.key(key)).Err()
}

// Exists checks if a key exists in the cache.
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, r.key(key)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return count > 0, nil
}

// Expire sets expiration for a key.
func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, r.key(key), expiration).Err()
}

// TTL returns the time to live for a key.
func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, r.key(key)).Result()
}

// Increment atomically increments a counter.
func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, r.key(key)).Result()
}

// IncrementBy atomically increments a counter by a specific amount.
func (r *RedisCache) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, r.key(key), value).Result()
}

// Decrement atomically decrements a counter.
func (r *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, r.key(key)).Result()
}

// DecrementBy atomically decrements a counter by a specific amount.
func (r *RedisCache) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.DecrBy(ctx, r.key(key), value).Result()
}

// SetNX sets a value only if the key doesn't exist (atomic).
func (r *RedisCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.SetNX(ctx, r.key(key), data, expiration).Result()
}

// GetSet atomically sets a new value and returns the old value.
func (r *RedisCache) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.GetSet(ctx, r.key(key), data).Result()
}

// MGet gets multiple values at once.
func (r *RedisCache) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = r.key(key)
	}

	return r.client.MGet(ctx, prefixedKeys...).Result()
}

// MSet sets multiple values at once.
func (r *RedisCache) MSet(ctx context.Context, pairs ...interface{}) error {
	// Convert keys to prefixed keys.
	prefixedPairs := make([]interface{}, len(pairs))
	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			key := pairs[i].(string)
			value := pairs[i+1]

			data, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
			}

			prefixedPairs[i] = r.key(key)
			prefixedPairs[i+1] = data
		}
	}

	return r.client.MSet(ctx, prefixedPairs...).Err()
}

// FlushDB clears all keys in the current database.
func (r *RedisCache) FlushDB(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

// Keys returns all keys matching a pattern.
func (r *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.client.Keys(ctx, r.key(pattern)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	// Remove prefix from keys.
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = key[len(r.prefix):]
	}

	return result, nil
}

// Scan iterates over keys matching a pattern.
func (r *RedisCache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	keys, newCursor, err := r.client.Scan(ctx, cursor, r.key(match), count).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan keys: %w", err)
	}

	// Remove prefix from keys.
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = key[len(r.prefix):]
	}

	return result, newCursor, nil
}

// Pipeline creates a new pipeline for batch operations.
func (r *RedisCache) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// TxPipeline creates a new transaction pipeline.
func (r *RedisCache) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

// Close closes the Redis connection.
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// HealthCheck checks if Redis is healthy.
func (r *RedisCache) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Stats returns Redis connection statistics.
func (r *RedisCache) Stats() *redis.PoolStats {
	return r.client.PoolStats()
}

// Client returns the underlying Redis client (use with caution).
func (r *RedisCache) Client() *redis.Client {
	return r.client
}

// Helper functions for environment variables.

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}

	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}

	return defaultValue
}
