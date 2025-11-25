// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package middleware provides functionality for the GO-PRO Learning Platform.
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-pro-backend/internal/cache"
	"go-pro-backend/pkg/logger"

	apierrors "go-pro-backend/internal/errors"
)

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	RequestsPerWindow int           `json:"requests_per_window"`
	WindowDuration    time.Duration `json:"window_duration"`
	KeyPrefix         string        `json:"key_prefix"`
}

// RateLimiter provides rate limiting middleware.
type RateLimiter struct {
	cache  cache.CacheManager
	config *RateLimitConfig
	logger logger.Logger
}

// NewRateLimiter creates a new rate limiter middleware.
func NewRateLimiter(cache cache.CacheManager, config *RateLimitConfig, logger logger.Logger) *RateLimiter {
	if config.KeyPrefix == "" {
		config.KeyPrefix = "ratelimit"
	}
	if config.RequestsPerWindow == 0 {
		config.RequestsPerWindow = 100
	}
	if config.WindowDuration == 0 {
		config.WindowDuration = time.Minute
	}

	return &RateLimiter{
		cache:  cache,
		config: config,
		logger: logger,
	}
}

// Limit applies rate limiting based on IP address.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get client identifier (IP address)
		clientIP := getClientIP(r)
		key := fmt.Sprintf("%s:%s", rl.config.KeyPrefix, clientIP)

		// Check rate limit.
		allowed, err := rl.cache.Allow(
			ctx,
			key,
			int64(rl.config.RequestsPerWindow),
			rl.config.WindowDuration,
		)
		if err != nil {
			rl.logger.Error(ctx, "Rate limit check failed", "error", err, "client_ip", clientIP)
			// On error, allow the request but log it.
			next.ServeHTTP(w, r)

			return
		}

		if !allowed {
			// Get remaining count for headers.
			remaining, _ := rl.cache.Remaining(
				ctx,
				key,
				int64(rl.config.RequestsPerWindow),
				rl.config.WindowDuration,
			)

			// Set rate limit headers.
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.config.WindowDuration).Unix(), 10))

			rl.logger.Warn(ctx, "Rate limit exceeded",
				"client_ip", clientIP,
				"limit", rl.config.RequestsPerWindow,
				"window", rl.config.WindowDuration)

			apiErr := &apierrors.APIError{
				Type:       "RATE_LIMIT_EXCEEDED",
				Code:       "TOO_MANY_REQUESTS",
				Message:    "Rate limit exceeded. Please try again later.",
				StatusCode: http.StatusTooManyRequests,
			}
			writeRateLimitError(w, r, apiErr)

			return
		}

		// Get remaining count for headers.
		remaining, _ := rl.cache.Remaining(
			ctx,
			key,
			int64(rl.config.RequestsPerWindow),
			rl.config.WindowDuration,
		)

		// Set rate limit headers.
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.config.WindowDuration).Unix(), 10))

		next.ServeHTTP(w, r)
	})
}

// LimitByUser applies rate limiting based on authenticated user ID.
func (rl *RateLimiter) LimitByUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Try to get user ID from context.
		userID := getUserIDFromContext(ctx)
		if userID == "" {
			// Fall back to IP-based rate limiting.
			rl.Limit(next).ServeHTTP(w, r)
			return
		}

		key := fmt.Sprintf("%s:user:%s", rl.config.KeyPrefix, userID)

		// Check rate limit.
		allowed, err := rl.cache.Allow(
			ctx,
			key,
			int64(rl.config.RequestsPerWindow),
			rl.config.WindowDuration,
		)
		if err != nil {
			rl.logger.Error(ctx, "Rate limit check failed", "error", err, "user_id", userID)
			// On error, allow the request but log it.
			next.ServeHTTP(w, r)

			return
		}

		if !allowed {
			// Get remaining count for headers.
			remaining, _ := rl.cache.Remaining(
				ctx,
				key,
				int64(rl.config.RequestsPerWindow),
				rl.config.WindowDuration,
			)

			// Set rate limit headers.
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.config.WindowDuration).Unix(), 10))

			rl.logger.Warn(ctx, "Rate limit exceeded",
				"user_id", userID,
				"limit", rl.config.RequestsPerWindow,
				"window", rl.config.WindowDuration)

			apiErr := &apierrors.APIError{
				Type:       "RATE_LIMIT_EXCEEDED",
				Code:       "TOO_MANY_REQUESTS",
				Message:    "Rate limit exceeded. Please try again later.",
				StatusCode: http.StatusTooManyRequests,
			}
			writeRateLimitError(w, r, apiErr)

			return
		}

		// Get remaining count for headers.
		remaining, _ := rl.cache.Remaining(
			ctx,
			key,
			int64(rl.config.RequestsPerWindow),
			rl.config.WindowDuration,
		)

		// Set rate limit headers.
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.config.WindowDuration).Unix(), 10))

		next.ServeHTTP(w, r)
	})
}

// LimitByEndpoint applies different rate limits based on endpoint.
func (rl *RateLimiter) LimitByEndpoint(limits map[string]*RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			path := r.URL.Path

			// Find matching endpoint config.
			var config *RateLimitConfig
			for pattern, cfg := range limits {
				if matchesPattern(path, pattern) {
					config = cfg
					break
				}
			}

			// Use default config if no match.
			if config == nil {
				config = rl.config
			}

			// Get client identifier.
			clientIP := getClientIP(r)
			key := fmt.Sprintf("%s:%s:%s", config.KeyPrefix, path, clientIP)

			// Check rate limit.
			allowed, err := rl.cache.Allow(
				ctx,
				key,
				int64(config.RequestsPerWindow),
				config.WindowDuration,
			)
			if err != nil {
				rl.logger.Error(ctx, "Rate limit check failed", "error", err, "client_ip", clientIP, "path", path)
				next.ServeHTTP(w, r)

				return
			}

			if !allowed {
				remaining, _ := rl.cache.Remaining(
					ctx,
					key,
					int64(config.RequestsPerWindow),
					config.WindowDuration,
				)

				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerWindow))
				w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
				w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.WindowDuration).Unix(), 10))

				rl.logger.Warn(ctx, "Rate limit exceeded",
					"client_ip", clientIP,
					"path", path,
					"limit", config.RequestsPerWindow)

				apiErr := &apierrors.APIError{
					Type:       "RATE_LIMIT_EXCEEDED",
					Code:       "TOO_MANY_REQUESTS",
					Message:    "Rate limit exceeded. Please try again later.",
					StatusCode: http.StatusTooManyRequests,
				}
				writeRateLimitError(w, r, apiErr)

				return
			}

			remaining, _ := rl.cache.Remaining(
				ctx,
				key,
				int64(config.RequestsPerWindow),
				config.WindowDuration,
			)

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerWindow))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.WindowDuration).Unix(), 10))

			next.ServeHTTP(w, r)
		})
	}
}

// getUserIDFromContext extracts user ID from context.
func getUserIDFromContext(ctx context.Context) string {
	// This would integrate with your auth middleware.
	// For now, return empty string.
	return ""
}

// matchesPattern checks if a path matches a pattern.
func matchesPattern(path, pattern string) bool {
	// Simple prefix matching for now.
	// Could be enhanced with regex or path matching library.
	return len(path) >= len(pattern) && path[:len(pattern)] == pattern
}

// writeRateLimitError writes a rate limit error response.
func writeRateLimitError(w http.ResponseWriter, r *http.Request, apiErr *apierrors.APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", "60") // Suggest retry after 60 seconds
	w.WriteHeader(apiErr.StatusCode)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    apiErr.Code,
			"message": apiErr.Message,
			"type":    apiErr.Type,
		},
		"success":   false,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}
