// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package circuitbreaker provides circuit breaker pattern implementation
// for protecting external service calls.
package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	// StateClosed means the circuit breaker is allowing requests.
	StateClosed State = iota
	// StateOpen means the circuit breaker is rejecting requests.
	StateOpen
	// StateHalfOpen means the circuit breaker is testing if the service is recovered.
	StateHalfOpen
)

func (s State) String() string {
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

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// Config holds circuit breaker configuration.
type Config struct {
	// Name is the identifier for this circuit breaker
	Name string
	// MaxRequests is the maximum number of requests allowed in half-open state
	MaxRequests uint32
	// Interval is the time window for counting failures
	Interval time.Duration
	// Timeout is how long the circuit stays open before transitioning to half-open
	Timeout time.Duration
	// ReadyToTrip is called when a failure is recorded; returns true if the circuit should open
	ReadyToTrip func(counts Counts) bool
	// OnStateChange is called when the circuit breaker state changes
	OnStateChange func(name string, from, to State)
	// FailureThreshold is the number of failures before opening the circuit
	FailureThreshold uint32
	// SuccessThreshold is the number of successes before closing the circuit
	SuccessThreshold uint32
}

// Counts holds the counts of requests and their outcomes.
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// DefaultConfig returns a default circuit breaker configuration.
func DefaultConfig(name string) *Config {
	return &Config{
		Name:            name,
		MaxRequests:     5,
		Interval:        10 * time.Second,
		Timeout:         30 * time.Second,
		FailureThreshold: 5,
		SuccessThreshold: 2,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from, to State) {
			fmt.Printf("[CircuitBreaker] %s: state changed from %s to %s\n", name, from, to)
		},
	}
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	config      *Config
	state       State
	counts      Counts
	lastStateChange time.Time
	lastFailure time.Time
	mutex       sync.RWMutex
}

// New creates a new circuit breaker with the given configuration.
func New(config *Config) *CircuitBreaker {
	if config == nil {
		config = DefaultConfig("default")
	}

	return &CircuitBreaker{
		config:            config,
		state:             StateClosed,
		lastStateChange:   time.Now(),
	}
}

// Execute runs the given function with circuit breaker protection.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	return cb.ExecuteWithContext(context.Background(), fn)
}

// ExecuteWithContext runs the given function with circuit breaker protection and context.
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, fn func() error) error {
	// Check if we can make the request
	if !cb.allowRequest() {
		return ErrCircuitOpen
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Execute the function
	err := fn()

	// Record the result
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// allowRequest checks if a request should be allowed.
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if timeout has elapsed
		if now.Sub(cb.lastStateChange) >= cb.config.Timeout {
			cb.setState(StateHalfOpen, now)
			return true
		}
		return false
	case StateHalfOpen:
		// Allow limited requests in half-open state
		return cb.counts.Requests < cb.config.MaxRequests
	default:
		return false
	}
}

// recordSuccess records a successful request.
func (cb *CircuitBreaker) recordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.counts.Requests++
	cb.counts.TotalSuccesses++
	cb.counts.ConsecutiveSuccesses++
	cb.counts.ConsecutiveFailures = 0

	// In half-open state, check if we should close
	if cb.state == StateHalfOpen && cb.counts.ConsecutiveSuccesses >= cb.config.SuccessThreshold {
		cb.setState(StateClosed, time.Now())
	}
}

// recordFailure records a failed request.
func (cb *CircuitBreaker) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	cb.counts.Requests++
	cb.counts.TotalFailures++
	cb.counts.ConsecutiveFailures++
	cb.counts.ConsecutiveSuccesses = 0
	cb.lastFailure = now

	// Check if we should open the circuit
	if cb.config.ReadyToTrip(cb.counts) {
		cb.setState(StateOpen, now)
	}
}

// setState changes the circuit breaker state.
func (cb *CircuitBreaker) setState(newState State, now time.Time) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = now
	cb.counts = Counts{} // Reset counts on state change

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(cb.config.Name, oldState, newState)
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Counts returns the current counts.
func (cb *CircuitBreaker) Counts() Counts {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.counts
}

// Reset resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.setState(StateClosed, time.Now())
}

// Manager manages multiple circuit breakers.
type Manager struct {
	breakers map[string]*CircuitBreaker
	mutex    sync.RWMutex
}

// NewManager creates a new circuit breaker manager.
func NewManager() *Manager {
	return &Manager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// Get returns a circuit breaker by name, creating one if it doesn't exist.
func (m *Manager) Get(name string, config *Config) *CircuitBreaker {
	m.mutex.RLock()
	breaker, exists := m.breakers[name]
	m.mutex.RUnlock()

	if exists {
		return breaker
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Double check after acquiring write lock
	if breaker, exists := m.breakers[name]; exists {
		return breaker
	}

	if config == nil {
		config = DefaultConfig(name)
	} else {
		config.Name = name
	}

	breaker = New(config)
	m.breakers[name] = breaker
	return breaker
}

// GetAll returns all circuit breakers.
func (m *Manager) GetAll() map[string]*CircuitBreaker {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]*CircuitBreaker, len(m.breakers))
	for k, v := range m.breakers {
		result[k] = v
	}
	return result
}

// ResetAll resets all circuit breakers.
func (m *Manager) ResetAll() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, breaker := range m.breakers {
		breaker.Reset()
	}
}

// Default circuit breaker manager instance
var defaultManager = NewManager()

// GetCircuitBreaker returns a circuit breaker from the default manager.
func GetCircuitBreaker(name string, config *Config) *CircuitBreaker {
	return defaultManager.Get(name, config)
}

// ExecuteWithBreaker executes a function with a named circuit breaker.
func ExecuteWithBreaker(name string, fn func() error) error {
	breaker := GetCircuitBreaker(name, nil)
	return breaker.Execute(fn)
}

// ExecuteWithBreakerAndConfig executes a function with a named circuit breaker and custom config.
func ExecuteWithBreakerAndConfig(name string, config *Config, fn func() error) error {
	breaker := GetCircuitBreaker(name, config)
	return breaker.Execute(fn)
}
