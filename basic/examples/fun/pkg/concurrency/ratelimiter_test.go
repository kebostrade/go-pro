package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	t.Run("Allow basic rate limiting", func(t *testing.T) {
		rl := NewRateLimiter(2, time.Second)
		defer rl.Stop()

		// Should allow first 2
		if !rl.Allow() {
			t.Error("Allow() = false; want true for first request")
		}
		if !rl.Allow() {
			t.Error("Allow() = false; want true for second request")
		}

		// Should block third
		if rl.Allow() {
			t.Error("Allow() = true; want false for third request")
		}
	})

	t.Run("Wait blocks and releases", func(t *testing.T) {
		rl := NewRateLimiter(1, 100*time.Millisecond)
		defer rl.Stop()

		// Consume the token
		rl.Wait()

		// This should block until refill
		start := time.Now()
		rl.Wait()
		duration := time.Since(start)

		if duration < 80*time.Millisecond {
			t.Errorf("Wait() returned too quickly: %v", duration)
		}
	})

	t.Run("WaitWithContext cancellation", func(t *testing.T) {
		rl := NewRateLimiter(1, time.Second)
		defer rl.Stop()

		// Consume the token
		rl.Wait()

		// Create cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := rl.WaitWithContext(ctx)
		if err != context.Canceled {
			t.Errorf("WaitWithContext() error = %v; want context.Canceled", err)
		}
	})

	t.Run("TryWait with timeout", func(t *testing.T) {
		rl := NewRateLimiter(1, time.Second)
		defer rl.Stop()

		// Consume the token
		rl.Wait()

		// Should timeout
		if rl.TryWait(50 * time.Millisecond) {
			t.Error("TryWait() = true; want false (timeout)")
		}
	})
}

func TestSlidingWindowRateLimiter(t *testing.T) {
	t.Run("Allow within window", func(t *testing.T) {
		rl := NewSlidingWindowRateLimiter(3, time.Second)

		// Should allow first 3
		for i := 0; i < 3; i++ {
			if !rl.Allow() {
				t.Errorf("Allow() = false for request %d; want true", i+1)
			}
		}

		// Should block 4th
		if rl.Allow() {
			t.Error("Allow() = true for 4th request; want false")
		}
	})

	t.Run("Count requests in window", func(t *testing.T) {
		rl := NewSlidingWindowRateLimiter(5, time.Second)

		rl.Allow()
		rl.Allow()
		rl.Allow()

		count := rl.Count()
		if count != 3 {
			t.Errorf("Count() = %d; want 3", count)
		}
	})

	t.Run("Requests expire after window", func(t *testing.T) {
		rl := NewSlidingWindowRateLimiter(2, 100*time.Millisecond)

		// Fill the window
		rl.Allow()
		rl.Allow()

		// Wait for window to expire
		time.Sleep(150 * time.Millisecond)

		// Should allow again
		if !rl.Allow() {
			t.Error("Allow() = false after window expiry; want true")
		}
	})
}

func TestLeakyBucket(t *testing.T) {
	t.Run("Add and leak", func(t *testing.T) {
		lb := NewLeakyBucket(2, 100*time.Millisecond)
		defer lb.Stop()

		// Should allow first 2
		if !lb.Add() {
			t.Error("Add() = false; want true for first request")
		}
		if !lb.Add() {
			t.Error("Add() = false; want true for second request")
		}

		// Should block third (bucket full)
		if lb.Add() {
			t.Error("Add() = true when bucket full; want false")
		}

		// Wait for leak
		time.Sleep(150 * time.Millisecond)

		// Should allow again
		if !lb.Add() {
			t.Error("Add() = false after leak; want true")
		}
	})
}

func TestAdaptiveRateLimiter(t *testing.T) {
	t.Run("Increase rate on success", func(t *testing.T) {
		rl := NewAdaptiveRateLimiter(1, 10, time.Second)
		defer rl.Stop()

		initialRate := rl.CurrentRate()

		// Record many successes
		for i := 0; i < 15; i++ {
			rl.RecordSuccess()
		}

		newRate := rl.CurrentRate()
		if newRate <= initialRate {
			t.Errorf("CurrentRate() = %d; want > %d after successes", newRate, initialRate)
		}
	})

	t.Run("Decrease rate on failure", func(t *testing.T) {
		rl := NewAdaptiveRateLimiter(5, 10, time.Second)
		defer rl.Stop()

		initialRate := rl.CurrentRate()

		// Record failures
		for i := 0; i < 5; i++ {
			rl.RecordFailure()
		}

		newRate := rl.CurrentRate()
		if newRate >= initialRate {
			t.Errorf("CurrentRate() = %d; want < %d after failures", newRate, initialRate)
		}
	})

	t.Run("Respects min and max rates", func(t *testing.T) {
		rl := NewAdaptiveRateLimiter(1, 5, time.Second)
		defer rl.Stop()

		// Try to decrease below min
		for i := 0; i < 20; i++ {
			rl.RecordFailure()
		}

		if rl.CurrentRate() < 1 {
			t.Errorf("CurrentRate() = %d; should not go below min 1", rl.CurrentRate())
		}

		// Try to increase above max
		for i := 0; i < 50; i++ {
			rl.RecordSuccess()
		}

		if rl.CurrentRate() > 5 {
			t.Errorf("CurrentRate() = %d; should not go above max 5", rl.CurrentRate())
		}
	})
}

// Benchmarks
func BenchmarkRateLimiterAllow(b *testing.B) {
	rl := NewRateLimiter(1000, time.Second)
	defer rl.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow()
	}
}

func BenchmarkSlidingWindowAllow(b *testing.B) {
	rl := NewSlidingWindowRateLimiter(1000, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow()
	}
}

func BenchmarkLeakyBucketAdd(b *testing.B) {
	lb := NewLeakyBucket(1000, time.Millisecond)
	defer lb.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.Add()
	}
}
