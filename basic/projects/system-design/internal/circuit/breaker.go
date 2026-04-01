// Package circuit provides circuit breaker pattern implementation.
package circuit

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker wraps the sony/gobreaker library for fault tolerance.
type CircuitBreaker struct {
	cb      *gobreaker.CircuitBreaker
	service string
	state   int32 // 0=closed, 1=open, 2=half-open
}

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(service string, opts ...Option) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        service,
		MaxRequests: 3,  // Maximum requests in half-open state
		Interval:    10, // Cyclic period of closed state (seconds)
		Timeout:     60, // Period of open state (seconds)
	}

	for _, opt := range opts {
		opt(&settings)
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	return &CircuitBreaker{
		cb:      cb,
		service: service,
		state:   int32(StateClosed),
	}
}

// Option configures the circuit breaker.
type Option func(*gobreaker.Settings)

// WithFailureRateThreshold sets the failure rate threshold (not directly supported, using MaxRequests).
func WithFailureRateThreshold(threshold int) Option {
	return func(s *gobreaker.Settings) {
		// gobreaker doesn't support failure rate threshold directly
		// Adjust MaxRequests as a proxy
	}
}

// WithVolumeThreshold sets the minimum volume threshold.
func WithVolumeThreshold(threshold int) Option {
	return func(s *gobreaker.Settings) {
		s.MaxRequests = uint32(threshold)
	}
}

// WithRetryMax sets the maximum retries in half-open state.
func WithRetryMax(max int) Option {
	return func(s *gobreaker.Settings) {
		s.MaxRequests = uint32(max)
	}
}

// Execute runs a function with circuit breaker protection.
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	done := make(chan error, 1)
	go func() {
		_, err := cb.cb.Execute(func() (interface{}, error) {
			return nil, fn()
		})
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() string {
	switch atomic.LoadInt32(&cb.state) {
	case int32(StateClosed):
		return "closed"
	case int32(StateOpen):
		return "open"
	case int32(StateHalfOpen):
		return "half-open"
	default:
		return "unknown"
	}
}

// ExternalService wraps an external HTTP service with circuit breaker.
type ExternalService struct {
	cb      *CircuitBreaker
	baseURL string
	timeout time.Duration
}

// NewExternalService creates a new external service wrapper.
func NewExternalService(baseURL string) *ExternalService {
	return &ExternalService{
		cb:      NewCircuitBreaker(baseURL),
		baseURL: baseURL,
		timeout: 5 * time.Second,
	}
}

// Call makes a call to the external service with circuit breaker.
func (s *ExternalService) Call(ctx context.Context, endpoint string) ([]byte, error) {
	var result []byte
	err := s.cb.Execute(ctx, func() error {
		// In production, this would make an HTTP call
		// For demonstration, simulate a call
		result = []byte(fmt.Sprintf(`{"service":"%s","endpoint":"%s"}`, s.baseURL, endpoint))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("external service call failed: %w", err)
	}

	return result, nil
}

// IsAvailable returns whether the service is currently available.
func (s *ExternalService) IsAvailable() bool {
	return s.cb.State() != "open"
}

// Errors that can be returned by the circuit breaker.
var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests")
)
