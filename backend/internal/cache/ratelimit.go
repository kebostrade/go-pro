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

// RedisRateLimiter implements rate limiting using Redis.
type RedisRateLimiter struct {
	client *redis.Client
	prefix string
}

// RateLimitInfo represents rate limit information.
type RateLimitInfo struct {
	Key       string        `json:"key"`
	Limit     int64         `json:"limit"`
	Remaining int64         `json:"remaining"`
	ResetTime time.Time     `json:"reset_time"`
	Window    time.Duration `json:"window"`
}

// NewRedisRateLimiter creates a new Redis rate limiter.
func NewRedisRateLimiter(client *redis.Client, prefix string) *RedisRateLimiter {
	if prefix == "" {
		prefix = "ratelimit:"
	}

	return &RedisRateLimiter{
		client: client,
		prefix: prefix,
	}
}

// rateLimitKey generates a rate limit key with prefix.
func (r *RedisRateLimiter) rateLimitKey(key string) string {
	return r.prefix + key
}

// Allow checks if a request is allowed under the rate limit.
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	return r.AllowN(ctx, key, 1, limit, window)
}

// AllowN checks if N requests are allowed under the rate limit.
func (r *RedisRateLimiter) AllowN(ctx context.Context, key string, n, limit int64, window time.Duration) (bool, error) {
	rateLimitKey := r.rateLimitKey(key)
	now := time.Now()
	windowStart := now.Truncate(window)

	// Lua script for atomic rate limiting using sliding window.
	rateLimitScript := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local window_size = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		local increment = tonumber(ARGV[4])
		local now = tonumber(ARGV[5])
		
		-- Remove expired entries (older than current window)
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start - window_size)
		
		-- Count current requests in the window
		local current = redis.call('ZCARD', key)
		
		-- Check if adding new requests would exceed limit
		if current + increment > limit then
			-- Set expiration for cleanup
			redis.call('EXPIRE', key, window_size)
			return {0, current, limit - current}
		end
		
		-- Add new requests with current timestamp as score
		for i = 1, increment do
			redis.call('ZADD', key, now + i, now + i)
		end
		
		-- Set expiration for cleanup
		redis.call('EXPIRE', key, window_size)
		
		-- Return success, current count, remaining
		local new_current = redis.call('ZCARD', key)
		return {1, new_current, limit - new_current}
	`

	result, err := r.client.Eval(ctx, rateLimitScript, []string{rateLimitKey},
		windowStart.Unix(),
		int64(window.Seconds()),
		limit,
		n,
		now.UnixNano(),
	).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check rate limit: %w", err)
	}

	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1

	return allowed, nil
}

// AllowWithInfo checks rate limit and returns detailed information.
func (r *RedisRateLimiter) AllowWithInfo(ctx context.Context, key string, limit int64, window time.Duration) (bool, *RateLimitInfo, error) {
	return r.AllowNWithInfo(ctx, key, 1, limit, window)
}

// AllowNWithInfo checks rate limit for N requests and returns detailed information.
func (r *RedisRateLimiter) AllowNWithInfo(ctx context.Context, key string, n, limit int64, window time.Duration) (bool, *RateLimitInfo, error) {
	rateLimitKey := r.rateLimitKey(key)
	now := time.Now()
	windowStart := now.Truncate(window)

	// Enhanced Lua script that returns more information.
	rateLimitScript := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local window_size = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		local increment = tonumber(ARGV[4])
		local now = tonumber(ARGV[5])
		
		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start - window_size)
		
		-- Count current requests
		local current = redis.call('ZCARD', key)
		
		-- Calculate reset time (end of current window)
		local reset_time = window_start + window_size
		
		-- Check if adding new requests would exceed limit
		if current + increment > limit then
			redis.call('EXPIRE', key, window_size)
			return {0, current, limit - current, reset_time}
		end
		
		-- Add new requests
		for i = 1, increment do
			redis.call('ZADD', key, now + i, now + i)
		end
		
		redis.call('EXPIRE', key, window_size)
		
		local new_current = redis.call('ZCARD', key)
		return {1, new_current, limit - new_current, reset_time}
	`

	result, err := r.client.Eval(ctx, rateLimitScript, []string{rateLimitKey},
		windowStart.Unix(),
		int64(window.Seconds()),
		limit,
		n,
		now.UnixNano(),
	).Result()
	if err != nil {
		return false, nil, fmt.Errorf("failed to check rate limit with info: %w", err)
	}

	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1
	_ = resultSlice[1].(int64) // current - not used but part of the result
	remaining := resultSlice[2].(int64)
	resetTime := time.Unix(resultSlice[3].(int64), 0)

	info := &RateLimitInfo{
		Key:       key,
		Limit:     limit,
		Remaining: remaining,
		ResetTime: resetTime,
		Window:    window,
	}

	return allowed, info, nil
}

// Remaining returns the number of remaining requests in the current window.
func (r *RedisRateLimiter) Remaining(ctx context.Context, key string, limit int64, window time.Duration) (int64, error) {
	rateLimitKey := r.rateLimitKey(key)
	now := time.Now()
	windowStart := now.Truncate(window)

	// Lua script to get remaining requests.
	remainingScript := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local window_size = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		
		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start - window_size)
		
		-- Count current requests
		local current = redis.call('ZCARD', key)
		
		-- Return remaining
		return limit - current
	`

	result, err := r.client.Eval(ctx, remainingScript, []string{rateLimitKey},
		windowStart.Unix(),
		int64(window.Seconds()),
		limit,
	).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get remaining requests: %w", err)
	}

	remaining := result.(int64)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// Reset clears the rate limit for a key.
func (r *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	rateLimitKey := r.rateLimitKey(key)

	err := r.client.Del(ctx, rateLimitKey).Err()
	if err != nil {
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}

	return nil
}

// GetInfo returns detailed rate limit information without consuming quota.
func (r *RedisRateLimiter) GetInfo(ctx context.Context, key string, limit int64, window time.Duration) (*RateLimitInfo, error) {
	rateLimitKey := r.rateLimitKey(key)
	now := time.Now()
	windowStart := now.Truncate(window)

	// Lua script to get info without consuming quota.
	infoScript := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local window_size = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		
		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start - window_size)
		
		-- Count current requests
		local current = redis.call('ZCARD', key)
		
		-- Calculate reset time
		local reset_time = window_start + window_size
		
		return {current, limit - current, reset_time}
	`

	result, err := r.client.Eval(ctx, infoScript, []string{rateLimitKey},
		windowStart.Unix(),
		int64(window.Seconds()),
		limit,
	).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit info: %w", err)
	}

	resultSlice := result.([]interface{})
	_ = resultSlice[0].(int64) // current - not used but part of the result
	remaining := resultSlice[1].(int64)
	resetTime := time.Unix(resultSlice[2].(int64), 0)

	if remaining < 0 {
		remaining = 0
	}

	info := &RateLimitInfo{
		Key:       key,
		Limit:     limit,
		Remaining: remaining,
		ResetTime: resetTime,
		Window:    window,
	}

	return info, nil
}

// TokenBucketAllow implements token bucket rate limiting.
func (r *RedisRateLimiter) TokenBucketAllow(ctx context.Context, key string, capacity, refillRate, tokens int64) (bool, error) {
	rateLimitKey := r.rateLimitKey(key)
	now := time.Now().Unix()

	// Lua script for token bucket algorithm.
	tokenBucketScript := `
		local key = KEYS[1]
		local capacity = tonumber(ARGV[1])
		local refill_rate = tonumber(ARGV[2])
		local tokens_requested = tonumber(ARGV[3])
		local now = tonumber(ARGV[4])
		
		-- Get current bucket state
		local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
		local current_tokens = tonumber(bucket[1]) or capacity
		local last_refill = tonumber(bucket[2]) or now
		
		-- Calculate tokens to add based on time elapsed
		local time_elapsed = now - last_refill
		local tokens_to_add = time_elapsed * refill_rate
		current_tokens = math.min(capacity, current_tokens + tokens_to_add)
		
		-- Check if we have enough tokens
		if current_tokens >= tokens_requested then
			current_tokens = current_tokens - tokens_requested
			
			-- Update bucket state
			redis.call('HMSET', key, 'tokens', current_tokens, 'last_refill', now)
			redis.call('EXPIRE', key, 3600) -- Expire after 1 hour of inactivity
			
			return {1, current_tokens}
		else
			-- Update last_refill time even if request is denied
			redis.call('HMSET', key, 'tokens', current_tokens, 'last_refill', now)
			redis.call('EXPIRE', key, 3600)
			
			return {0, current_tokens}
		end
	`

	result, err := r.client.Eval(ctx, tokenBucketScript, []string{rateLimitKey},
		capacity,
		refillRate,
		tokens,
		now,
	).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token bucket rate limit: %w", err)
	}

	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1

	return allowed, nil
}

// CleanupExpiredRateLimits removes expired rate limit entries.
func (r *RedisRateLimiter) CleanupExpiredRateLimits(ctx context.Context) error {
	// Use SCAN to iterate through all rate limit keys.
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, r.prefix+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan rate limit keys: %w", err)
		}

		// Check each key's TTL.
		for _, key := range keys {
			ttl, err := r.client.TTL(ctx, key).Result()
			if err != nil {
				continue
			}

			// If TTL is -2 (key doesn't exist), it's already cleaned up.
			// If TTL is -1 (no expiration), set a reasonable expiration.
			if ttl == -1 {
				r.client.Expire(ctx, key, time.Hour)
			} else if ttl == -2 {
				// Key doesn't exist, nothing to do.
				continue
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// ListRateLimits returns all active rate limits matching a pattern.
func (r *RedisRateLimiter) ListRateLimits(ctx context.Context, pattern string) ([]string, error) {
	rateLimitPattern := r.rateLimitKey(pattern)
	keys, err := r.client.Keys(ctx, rateLimitPattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list rate limits: %w", err)
	}

	// Remove prefix from keys.
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = key[len(r.prefix):]
	}

	return result, nil
}
