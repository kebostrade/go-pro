// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides exercise evaluation functionality tests.
package service

import (
	"context"
	"testing"

	"go-pro-backend/internal/domain"
	"go-pro-backend/pkg/logger"
)

func TestExerciseEvaluator_ValidateGoCode(t *testing.T) {
	executor := NewLocalExecutor()
	logger := &noOpLogger{}
	evaluator := NewExerciseEvaluator(executor, logger)

	tests := []struct {
		name        string
		code        string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Go code",
			code: `package main

func main() {
	println("test")
}`,
			expectError: false,
		},
		{
			name: "Missing package declaration",
			code: `func main() {
	println("test")
}`,
			expectError: true,
			errorMsg:    "package main",
		},
		{
			name: "Missing main function",
			code: `package main

func helper() {
	println("test")
}`,
			expectError: true,
			errorMsg:    "func main()",
		},
		{
			name: "Code too large",
			code: string(make([]byte, 50001)),
			expectError: true,
			errorMsg:    "too large",
		},
		{
			name: "Dangerous os import",
			code: `package main

import "os"

func main() {
	os.Exit(0)
}`,
			expectError: true,
			errorMsg:    "not allowed",
		},
		{
			name: "Dangerous net import",
			code: `package main

import "net"

func main() {
}`,
			expectError: true,
			errorMsg:    "not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := evaluator.validateGoCode(tt.code)

			if tt.expectError {
				if err == nil {
					t.Errorf("validateGoCode() expected error containing %q, got nil", tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateGoCode() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestExerciseEvaluator_GetTestCasesForExercise(t *testing.T) {
	executor := NewLocalExecutor()
	logger := &noOpLogger{}
	evaluator := NewExerciseEvaluator(executor, logger)

	tests := []struct {
		name              string
		exercise          *domain.Exercise
		expectedTestCount int
		expectedTestName   string
	}{
		{
			name: "Hello World exercise",
			exercise: &domain.Exercise{
				Title:       "Hello World",
				Description: "Print Hello, World!",
			},
			expectedTestCount: 2,
			expectedTestName:  "Basic output",
		},
		{
			name: "FizzBuzz exercise",
			exercise: &domain.Exercise{
				Title:       "FizzBuzz Challenge",
				Description: "Implement FizzBuzz",
			},
			expectedTestCount: 4,
			expectedTestName:  "Test for 3",
		},
		{
			name: "Generic exercise",
			exercise: &domain.Exercise{
				Title:       "Some Exercise",
				Description: "Do something",
			},
			expectedTestCount: 2,
			expectedTestName:  "Basic functionality test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCases := evaluator.getTestCasesForExercise(tt.exercise)

			if len(testCases) != tt.expectedTestCount {
				t.Errorf("getTestCasesForExercise() returned %d tests, want %d", len(testCases), tt.expectedTestCount)
			}

			if len(testCases) > 0 && testCases[0].Name != tt.expectedTestName {
				t.Errorf("getTestCasesForExercise() first test = %v, want %v", testCases[0].Name, tt.expectedTestName)
			}
		})
	}
}

func TestGetPassedTests(t *testing.T) {
	results := []domain.TestResult{
		{TestName: "Test 1", Passed: true},
		{TestName: "Test 2", Passed: false},
		{TestName: "Test 3", Passed: true},
	}

	passed := GetPassedTests(results)

	if len(passed) != 2 {
		t.Errorf("GetPassedTests() returned %d, want 2", len(passed))
	}

	if passed[0] != "Test 1" || passed[1] != "Test 3" {
		t.Errorf("GetPassedTests() returned %v, want [Test 1, Test 3]", passed)
	}
}

func TestGetFailedTests(t *testing.T) {
	results := []domain.TestResult{
		{TestName: "Test 1", Passed: true},
		{TestName: "Test 2", Passed: false},
		{TestName: "Test 3", Passed: true},
	}

	failed := GetFailedTests(results)

	if len(failed) != 1 {
		t.Errorf("GetFailedTests() returned %d, want 1", len(failed))
	}

	if failed[0] != "Test 2" {
		t.Errorf("GetFailedTests() returned %v, want [Test 2]", failed)
	}
}

// noOpLogger is a simple no-op logger for testing
type noOpLogger struct{}

func (m *noOpLogger) Debug(ctx context.Context, msg string, keyvals ...interface{}) {}
func (m *noOpLogger) Info(ctx context.Context, msg string, keyvals ...interface{})  {}
func (m *noOpLogger) Warn(ctx context.Context, msg string, keyvals ...interface{})  {}
func (m *noOpLogger) Error(ctx context.Context, msg string, keyvals ...interface{}) {}
func (m *noOpLogger) With(keyvals ...interface{}) logger.Logger                                      { return m }
