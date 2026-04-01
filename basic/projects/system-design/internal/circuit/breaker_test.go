package circuit

import (
	"context"
	"testing"
	"time"
)

func TestNewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker("test-service")

	if cb == nil {
		t.Error("NewCircuitBreaker returned nil")
	}

	if cb.State() != "closed" {
		t.Errorf("Initial state = %s, want closed", cb.State())
	}
}

func TestCircuitBreaker_Execute(t *testing.T) {
	cb := NewCircuitBreaker("test-service")

	err := cb.Execute(context.Background(), func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
}

func TestCircuitBreaker_ExecuteWithError(t *testing.T) {
	cb := NewCircuitBreaker("test-service")

	err := cb.Execute(context.Background(), func() error {
		return context.DeadlineExceeded
	})

	if err == nil {
		t.Error("Expected error from Execute")
	}
}

func TestCircuitBreaker_ContextCancellation(t *testing.T) {
	cb := NewCircuitBreaker("test-service")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := cb.Execute(ctx, func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestExternalService(t *testing.T) {
	svc := NewExternalService("http://example.com")

	result, err := svc.Call(context.Background(), "/test")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected result from Call")
	}

	if !svc.IsAvailable() {
		t.Error("Service should be available")
	}
}

func TestWithOptions(t *testing.T) {
	cb := NewCircuitBreaker(
		"test-service",
		WithFailureRateThreshold(30),
		WithVolumeThreshold(20),
		WithRetryMax(5),
	)

	if cb == nil {
		t.Error("NewCircuitBreaker with options returned nil")
	}
}
