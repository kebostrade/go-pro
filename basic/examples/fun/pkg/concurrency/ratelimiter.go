package concurrency

import (
	"context"
	"sync"
	"time"
)

// RateLimiter controls the rate of operations using token bucket algorithm
type RateLimiter struct {
	tokens   chan struct{}
	maxRate  int
	interval time.Duration
	stopOnce sync.Once
	stopChan chan struct{}
}

// NewRateLimiter creates a new rate limiter
// maxRate: maximum number of operations allowed per interval
// interval: time window for the rate limit
func NewRateLimiter(maxRate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:   make(chan struct{}, maxRate),
		maxRate:  maxRate,
		interval: interval,
		stopChan: make(chan struct{}),
	}

	// Fill the token bucket initially
	for i := 0; i < maxRate; i++ {
		rl.tokens <- struct{}{}
	}

	// Start the token refill goroutine
	go rl.refillTokens()

	return rl
}

// refillTokens periodically adds tokens to the bucket
func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.interval / time.Duration(rl.maxRate))
	defer ticker.Stop()

	for {
		select {
		case <-rl.stopChan:
			return
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
				// Token added successfully
			default:
				// Bucket is full, skip
			}
		}
	}
}

// Allow checks if an operation is allowed (non-blocking)
// Returns true if allowed, false if rate limited
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Wait blocks until an operation is allowed
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// WaitWithContext blocks until an operation is allowed or context is cancelled
func (rl *RateLimiter) WaitWithContext(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryWait attempts to wait for a token with a timeout
// Returns true if token acquired, false if timeout
func (rl *RateLimiter) TryWait(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-rl.tokens:
		return true
	case <-timer.C:
		return false
	}
}

// Stop stops the rate limiter and releases resources
func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.stopChan)
	})
}

// SlidingWindowRateLimiter implements a sliding window rate limiter
type SlidingWindowRateLimiter struct {
	mu       sync.Mutex
	requests []time.Time
	maxRate  int
	window   time.Duration
}

// NewSlidingWindowRateLimiter creates a new sliding window rate limiter
func NewSlidingWindowRateLimiter(maxRate int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		requests: make([]time.Time, 0, maxRate),
		maxRate:  maxRate,
		window:   window,
	}
}

// Allow checks if a request is allowed under the sliding window
func (rl *SlidingWindowRateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Remove old requests outside the window
	validRequests := make([]time.Time, 0, len(rl.requests))
	for _, t := range rl.requests {
		if t.After(cutoff) {
			validRequests = append(validRequests, t)
		}
	}
	rl.requests = validRequests

	// Check if we can allow this request
	if len(rl.requests) < rl.maxRate {
		rl.requests = append(rl.requests, now)
		return true
	}

	return false
}

// Count returns the current number of requests in the window
func (rl *SlidingWindowRateLimiter) Count() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	count := 0
	for _, t := range rl.requests {
		if t.After(cutoff) {
			count++
		}
	}

	return count
}

// LeakyBucket implements a leaky bucket rate limiter
type LeakyBucket struct {
	capacity int
	rate     time.Duration
	bucket   chan struct{}
	stopChan chan struct{}
	stopOnce sync.Once
}

// NewLeakyBucket creates a new leaky bucket rate limiter
func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		capacity: capacity,
		rate:     rate,
		bucket:   make(chan struct{}, capacity),
		stopChan: make(chan struct{}),
	}

	go lb.leak()

	return lb
}

// leak continuously removes items from the bucket at the specified rate
func (lb *LeakyBucket) leak() {
	ticker := time.NewTicker(lb.rate)
	defer ticker.Stop()

	for {
		select {
		case <-lb.stopChan:
			return
		case <-ticker.C:
			select {
			case <-lb.bucket:
				// Leaked one item
			default:
				// Bucket is empty
			}
		}
	}
}

// Add attempts to add an item to the bucket
// Returns true if successful, false if bucket is full
func (lb *LeakyBucket) Add() bool {
	select {
	case lb.bucket <- struct{}{}:
		return true
	default:
		return false
	}
}

// Wait blocks until there's space in the bucket
func (lb *LeakyBucket) Wait() {
	lb.bucket <- struct{}{}
}

// Stop stops the leaky bucket
func (lb *LeakyBucket) Stop() {
	lb.stopOnce.Do(func() {
		close(lb.stopChan)
	})
}

// AdaptiveRateLimiter adjusts its rate based on success/failure
type AdaptiveRateLimiter struct {
	mu           sync.RWMutex
	currentRate  int
	minRate      int
	maxRate      int
	interval     time.Duration
	successCount int
	failureCount int
	limiter      *RateLimiter
}

// NewAdaptiveRateLimiter creates a new adaptive rate limiter
func NewAdaptiveRateLimiter(minRate, maxRate int, interval time.Duration) *AdaptiveRateLimiter {
	initialRate := (minRate + maxRate) / 2
	return &AdaptiveRateLimiter{
		currentRate: initialRate,
		minRate:     minRate,
		maxRate:     maxRate,
		interval:    interval,
		limiter:     NewRateLimiter(initialRate, interval),
	}
}

// Allow checks if an operation is allowed
func (rl *AdaptiveRateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

// RecordSuccess records a successful operation
func (rl *AdaptiveRateLimiter) RecordSuccess() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.successCount++

	// Increase rate if we have many successes
	if rl.successCount > 10 && rl.currentRate < rl.maxRate {
		rl.adjustRate(rl.currentRate + 1)
		rl.successCount = 0
		rl.failureCount = 0
	}
}

// RecordFailure records a failed operation
func (rl *AdaptiveRateLimiter) RecordFailure() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.failureCount++

	// Decrease rate if we have many failures
	if rl.failureCount > 3 && rl.currentRate > rl.minRate {
		rl.adjustRate(rl.currentRate - 1)
		rl.successCount = 0
		rl.failureCount = 0
	}
}

// adjustRate adjusts the current rate
func (rl *AdaptiveRateLimiter) adjustRate(newRate int) {
	if newRate < rl.minRate {
		newRate = rl.minRate
	}
	if newRate > rl.maxRate {
		newRate = rl.maxRate
	}

	if newRate != rl.currentRate {
		rl.currentRate = newRate
		rl.limiter.Stop()
		rl.limiter = NewRateLimiter(newRate, rl.interval)
	}
}

// CurrentRate returns the current rate
func (rl *AdaptiveRateLimiter) CurrentRate() int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.currentRate
}

// Stop stops the adaptive rate limiter
func (rl *AdaptiveRateLimiter) Stop() {
	rl.limiter.Stop()
}
