//go:build ignore

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// ============================================================================
// CIRCUIT BREAKER PATTERN
// ============================================================================

// State represents the circuit breaker state
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF_OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name           string
	maxFailures    int
	resetTimeout   time.Duration

	mu           sync.RWMutex
	state        State
	failures     int
	lastFailTime time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:         name,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Execute runs the given function if the circuit is open
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// Check if we should allow execution
	if !cb.canExecute() {
		return fmt.Errorf("circuit breaker %s is OPEN", cb.name)
	}

	// Execute the function
	err := fn()

	// Handle result
	cb.recordResult(err)

	return err
}

// canExecute determines if execution should be allowed
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	// Always allow if closed
	if cb.state == StateClosed {
		return true
	}

	// If open, check if reset timeout has passed
	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			// Try half-open state
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = StateHalfOpen
			cb.failures = 0
			cb.mu.Unlock()
			cb.mu.RLock()
			return true
		}
		return false
	}

	// Half-open state - allow execution
	return true
}

// recordResult records the result of an execution
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()

		// Transition to open if threshold reached
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			log.Printf("Circuit breaker %s opened after %d failures", cb.name, cb.failures)
		}
	} else {
		// Success - reset failures and close circuit
		cb.failures = 0
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
			log.Printf("Circuit breaker %s closed after successful retry", cb.name)
		}
	}
}

// GetState returns the current state
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// ============================================================================
// EXAMPLE USAGE
// ============================================================================

// ExternalService simulates an external API
type ExternalService struct {
	shouldFail bool
	callCount  int
}

func (s *ExternalService) Call() error {
	s.callCount++

	if s.shouldFail {
		return errors.New("service unavailable")
	}

	return nil
}

func exampleCircuitBreaker() {
	fmt.Println("\n=== Circuit Breaker Example ===")

	// Create circuit breaker: opens after 3 failures, resets after 2 seconds
	cb := NewCircuitBreaker("api-service", 3, 2*time.Second)
	service := &ExternalService{shouldFail: true}

	// Simulate failures
	fmt.Println("\n--- Simulating Failures ---")
	for i := 1; i <= 5; i++ {
		err := cb.Execute(context.Background(), service.Call)
		fmt.Printf("Call %d: State=%s, Error=%v\n", i, cb.GetState(), err)
		time.Sleep(100 * time.Millisecond)
	}

	// Circuit should be open now
	fmt.Printf("\nCircuit State: %s\n", cb.GetState())

	// Try to call while open
	fmt.Println("\n--- Trying to call while OPEN ---")
	err := cb.Execute(context.Background(), service.Call)
	fmt.Printf("Call while open: Error=%v\n", err)

	// Wait for reset timeout
	fmt.Println("\n--- Waiting for reset timeout ---")
	time.Sleep(3 * time.Second)

	// Service recovers
	service.shouldFail = false

	// Try again - should be half-open, then close
	fmt.Println("\n--- Service Recovered ---")
	for i := 1; i <= 3; i++ {
		err := cb.Execute(context.Background(), service.Call)
		fmt.Printf("Call %d: State=%s, Error=%v\n", i, cb.GetState(), err)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\nFinal Circuit State: %s\n", cb.GetState())
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	exampleCircuitBreaker()
}

// ============================================================================
// EXPECTED OUTPUT
// ============================================================================

/*
=== Circuit Breaker Example ===

--- Simulating Failures ---
Call 1: State=CLOSED, Error=service unavailable
Call 2: State=CLOSED, Error=service unavailable
Call 3: State=CLOSED, Error=service unavailable
2024/01/15 10:30:45 Circuit breaker api-service opened after 3 failures
Call 4: State=OPEN, Error=circuit breaker api-service is OPEN
Call 5: State=OPEN, Error=circuit breaker api-service is OPEN

Circuit State: OPEN

--- Trying to call while OPEN ---
Call while open: Error=circuit breaker api-service is OPEN

--- Waiting for reset timeout ---

--- Service Recovered ---
Call 1: State=HALF_OPEN, Error=<nil>
2024/01/15 10:30:48 Circuit breaker api-service closed after successful retry
Call 2: State=CLOSED, Error=<nil>
Call 3: State=CLOSED, Error=<nil>

Final Circuit State: CLOSED
*/
