// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"time"
)

// ExecutorService defines the interface for code execution and testing.
type ExecutorService interface {
	ExecuteCode(ctx context.Context, req *ExecuteRequest) (*ExecuteResult, error)
}

// ExecuteRequest represents a code execution request.
type ExecuteRequest struct {
	Code      string        `json:"code" validate:"required"`
	Language  string        `json:"language" validate:"required,oneof=go python javascript"`
	TestCases []TestCase    `json:"test_cases" validate:"required,min=1"`
	Timeout   time.Duration `json:"timeout" validate:"required"`
}

// ExecuteResult represents the result of code execution.
type ExecuteResult struct {
	Passed        bool          `json:"passed"`
	Score         int           `json:"score"`
	Results       []TestResult  `json:"results"`
	ExecutionTime time.Duration `json:"execution_time"`
	Error         error         `json:"error,omitempty"`
}

// TestCase represents a single test case.
type TestCase struct {
	Name     string `json:"name" validate:"required"`
	Input    string `json:"input"`
	Expected string `json:"expected" validate:"required"`
}

// TestResult represents the result of a single test execution.
type TestResult struct {
	TestName string `json:"test_name"`
	Passed   bool   `json:"passed"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Error    string `json:"error,omitempty"`
}

// mockExecutorService is a placeholder implementation for Phase 1.
// Phase 2 will replace this with a real sandboxed execution engine.
type mockExecutorService struct{}

// NewMockExecutorService creates a mock executor service for testing.
// This is a placeholder implementation that simulates test results.
// Phase 2 will implement real code execution with Docker sandboxing.
func NewMockExecutorService() ExecutorService {
	return &mockExecutorService{}
}

// ExecuteCode executes code against test cases (mock implementation).
func (m *mockExecutorService) ExecuteCode(ctx context.Context, req *ExecuteRequest) (*ExecuteResult, error) {
	// Mock implementation: simulate test execution
	// Phase 2 will replace this with actual code execution

	results := make([]TestResult, len(req.TestCases))
	passedCount := 0

	for i, tc := range req.TestCases {
		// Mock: simulate tests passing/failing
		passed := i%2 == 0 // Simple pattern: alternating pass/fail for demo

		results[i] = TestResult{
			TestName: tc.Name,
			Passed:   passed,
			Expected: tc.Expected,
			Actual:   tc.Expected, // In mock, actual matches expected when passed
			Error:    "",
		}

		if !passed {
			results[i].Actual = "mock output (different)"
			results[i].Error = "Output does not match expected result"
		}

		if passed {
			passedCount++
		}
	}

	// Calculate score
	score := 0
	if len(req.TestCases) > 0 {
		score = (passedCount * 100) / len(req.TestCases)
	}

	return &ExecuteResult{
		Passed:        passedCount == len(req.TestCases),
		Score:         score,
		Results:       results,
		ExecutionTime: 45 * time.Millisecond, // Mock execution time
		Error:         nil,
	}, nil
}
