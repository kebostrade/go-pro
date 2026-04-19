// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package cache provides Redis-backed rate limiting.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter provides distributed rate limiting using Redis.
// This is suitable for multi-instance deployments.
type RedisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// NewRedisRateLimiter creates a new Redis-backed rate limiter.
func NewRedisRateLimiter(client *redis.Client, limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

// Allow checks if a request is allowed under rate limit.
// Uses sliding window algorithm with Redis INCR.
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-r.window)
	
	pipe := r.client.Pipeline()
	
	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixMilli()))
	
	// Count current requests in window
	pipe.ZCard(ctx, key)
	
	// Add new request
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixMilli()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})
	
	// Set expiry
	pipe.Expire(ctx, key, r.window)
	
	results, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, err
	}
	
	count := results[1].(*redis.IntCmd).Val()
	
	return count <= r.limit, nil
}

// Close closes the Redis connection.
func (r *RedisRateLimiter) Close() error {
	return r.client.Close()
}
