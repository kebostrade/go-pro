// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides exercise evaluation functionality.
package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go-pro-backend/internal/domain"
	"go-pro-backend/pkg/logger"
)

// ExerciseEvaluator handles code evaluation for exercises.
type ExerciseEvaluator struct {
	executor ExecutorService
	logger   logger.Logger
}

// NewExerciseEvaluator creates a new exercise evaluator.
func NewExerciseEvaluator(executor ExecutorService, logger logger.Logger) *ExerciseEvaluator {
	return &ExerciseEvaluator{
		executor: executor,
		logger:   logger,
	}
}

// EvaluateSubmission evaluates a code submission against test cases.
func (e *ExerciseEvaluator) EvaluateSubmission(
	ctx context.Context,
	exercise *domain.Exercise,
	req *domain.SubmitExerciseRequest,
) (*domain.ExerciseSubmissionResult, error) {
	startTime := time.Now()

	// Validate code structure based on language
	if err := e.validateCodeStructure(req.Code, req.Language); err != nil {
		return &domain.ExerciseSubmissionResult{
			Success:     false,
			ExerciseID:  exercise.ID,
			Score:       0,
			Passed:      false,
			Message:     fmt.Sprintf("Code validation failed: %s", err.Error()),
			Feedback:    "Please check your code structure and try again.",
			SubmittedAt: time.Now(),
		}, nil
	}

	// Get test cases for this exercise
	testCases := e.getTestCasesForExercise(exercise)

	// Execute code against test cases
	execReq := &ExecuteRequest{
		Code:      req.Code,
		Language:  req.Language,
		TestCases: convertDomainTestCases(testCases),
		Timeout:   30 * time.Second,
	}

	execResult, err := e.executor.ExecuteCode(ctx, execReq)
	if err != nil {
		e.logger.Error(ctx, "Failed to execute code", "error", err, "exercise_id", exercise.ID)
		return &domain.ExerciseSubmissionResult{
			Success:     false,
			ExerciseID:  exercise.ID,
			Score:       0,
			Passed:      false,
			Message:     fmt.Sprintf("Execution failed: %s", err.Error()),
			Feedback:    "There was an error executing your code. Please check for syntax errors or runtime issues.",
			SubmittedAt: time.Now(),
		}, nil
	}

	// Convert executor results to domain results
	domainResults := make([]domain.TestResult, len(execResult.Results))
	passedTests := make([]string, 0)
	failedTests := make([]string, 0)

	for i, r := range execResult.Results {
		domainResults[i] = domain.TestResult{
			TestName: r.TestName,
			Passed:   r.Passed,
			Expected: r.Expected,
			Actual:   r.Actual,
			Error:    r.Error,
		}
		if r.Passed {
			passedTests = append(passedTests, r.TestName)
		} else {
			failedTests = append(failedTests, r.TestName)
		}
	}

	// Build feedback message
	message := e.buildFeedbackMessage(execResult, domainResults)
	feedback := e.buildDetailedFeedback(execResult, domainResults, passedTests, failedTests)

	// Calculate execution time in milliseconds
	executionTimeMs := time.Since(startTime).Milliseconds()
	if execResult.ExecutionTime > 0 {
		executionTimeMs = execResult.ExecutionTime.Milliseconds()
	}

	return &domain.ExerciseSubmissionResult{
		Success:         true,
		ExerciseID:      exercise.ID,
		Score:           execResult.Score,
		Passed:          execResult.Passed,
		Message:         message,
		Feedback:        feedback,
		TestResults:     domainResults,
		ExecutionTimeMs: executionTimeMs,
		SubmittedAt:     time.Now(),
	}, nil
}

// convertDomainTestCases converts domain test cases to executor test cases
func convertDomainTestCases(testCases []domain.TestCaseForExecution) []TestCase {
	result := make([]TestCase, len(testCases))
	for i, tc := range testCases {
		result[i] = TestCase{
			Name:     tc.Name,
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}
	return result
}

// validateCodeStructure performs basic code structure validation.
func (e *ExerciseEvaluator) validateCodeStructure(code, language string) error {
	switch language {
	case "go":
		return e.validateGoCode(code)
	case "python":
		return e.validatePythonCode(code)
	case "javascript":
		return e.validateJavaScriptCode(code)
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}
}

// validateGoCode validates Go code structure.
func (e *ExerciseEvaluator) validateGoCode(code string) error {
	// Check code size
	if len(code) > 50000 {
		return fmt.Errorf("code too large: max 50KB allowed")
	}

	// Check for dangerous imports
	dangerousPatterns := []struct {
		pattern string
		reason  string
	}{
		{`(?s)import\s*\([^)]*"os"`, "direct os package imports are not allowed"},
		{`(?s)import\s*\([^)]*"net"`, "direct net package imports are not allowed"},
		{`(?s)import\s*\([^)]*"syscall"`, "syscall package imports are not allowed"},
		{`(?s)import\s*\([^)]*"unsafe"`, "unsafe package imports are not allowed"},
		{`(?s)import\s*\([^)]*"runtime/debug"`, "runtime/debug imports are not allowed"},
		{`import\s+"os"`, "direct os package imports are not allowed"},
		{`import\s+"net"`, "direct net package imports are not allowed"},
		{`import\s+"syscall"`, "syscall package imports are not allowed"},
		{`import\s+"unsafe"`, "unsafe package imports are not allowed"},
	}

	for _, dp := range dangerousPatterns {
		if matched, _ := regexp.MatchString(dp.pattern, code); matched {
			return fmt.Errorf("%s", dp.reason)
		}
	}

	// Basic structure validation
	if !strings.Contains(code, "package main") {
		return fmt.Errorf("code must contain 'package main'")
	}

	if !strings.Contains(code, "func main()") {
		return fmt.Errorf("code must contain 'func main()'")
	}

	return nil
}

// validatePythonCode validates Python code structure.
func (e *ExerciseEvaluator) validatePythonCode(code string) error {
	if len(code) > 50000 {
		return fmt.Errorf("code too large: max 50KB allowed")
	}

	// Check for dangerous imports
	dangerousPatterns := []struct {
		pattern string
		reason  string
	}{
		{`import\s+os\s*`, "os module imports require careful handling"},
		{`import\s+subprocess\s*`, "subprocess module is not allowed"},
		{`import\s+sys\s*`, "sys module usage is restricted"},
		{`from\s+os\s+import`, "os module imports require careful handling"},
		{`from\s+subprocess\s+import`, "subprocess module is not allowed"},
	}

	for _, dp := range dangerousPatterns {
		if matched, _ := regexp.MatchString(dp.pattern, code); matched {
			return fmt.Errorf("%s", dp.reason)
		}
	}

	return nil
}

// validateJavaScriptCode validates JavaScript code structure.
func (e *ExerciseEvaluator) validateJavaScriptCode(code string) error {
	if len(code) > 50000 {
		return fmt.Errorf("code too large: max 50KB allowed")
	}

	// Basic JavaScript validation
	dangerousPatterns := []struct {
		pattern string
		reason  string
	}{
		{`eval\s*\(`, "eval() is not allowed for security reasons"},
		{`require\s*\(\s*['"]child_process['"]`, "child_process module is not allowed"},
		{`require\s*\(\s*['"]fs['"]`, "fs module requires careful handling"},
	}

	for _, dp := range dangerousPatterns {
		if matched, _ := regexp.MatchString(dp.pattern, code); matched {
			return fmt.Errorf("%s", dp.reason)
		}
	}

	return nil
}

// getTestCasesForExercise returns test cases for an exercise.
func (e *ExerciseEvaluator) getTestCasesForExercise(exercise *domain.Exercise) []domain.TestCaseForExecution {
	lowerTitle := strings.ToLower(exercise.Title)
	lowerDesc := strings.ToLower(exercise.Description)

	switch {
	case strings.Contains(lowerTitle, "hello") || strings.Contains(lowerDesc, "hello"):
		return []domain.TestCaseForExecution{
			{Name: "Basic output", Input: "", Expected: "Hello, World!"},
			{Name: "No extra output", Input: "", Expected: "Hello, World!"},
		}

	case strings.Contains(lowerTitle, "fizzbuzz") || strings.Contains(lowerDesc, "fizzbuzz"):
		return []domain.TestCaseForExecution{
			{Name: "Test for 3", Input: "3", Expected: "Fizz"},
			{Name: "Test for 5", Input: "5", Expected: "Buzz"},
			{Name: "Test for 15", Input: "15", Expected: "FizzBuzz"},
			{Name: "Test for other number", Input: "7", Expected: "7"},
		}

	case strings.Contains(lowerTitle, "sum") || strings.Contains(lowerDesc, "sum"):
		return []domain.TestCaseForExecution{
			{Name: "Test positive numbers", Input: "5\n3", Expected: "8"},
			{Name: "Test with zero", Input: "5\n0", Expected: "5"},
		}

	case strings.Contains(lowerTitle, "factorial") || strings.Contains(lowerDesc, "factorial"):
		return []domain.TestCaseForExecution{
			{Name: "Test factorial of 5", Input: "5", Expected: "120"},
			{Name: "Test factorial of 0", Input: "0", Expected: "1"},
		}

	case strings.Contains(lowerTitle, "palindrome") || strings.Contains(lowerDesc, "palindrome"):
		return []domain.TestCaseForExecution{
			{Name: "Test simple palindrome", Input: "racecar", Expected: "true"},
			{Name: "Test non-palindrome", Input: "hello", Expected: "false"},
		}

	default:
		return []domain.TestCaseForExecution{
			{Name: "Basic functionality test", Input: "", Expected: ""},
			{Name: "Edge case test", Input: "", Expected: ""},
		}
	}
}

// buildFeedbackMessage builds a summary message based on test results.
func (e *ExerciseEvaluator) buildFeedbackMessage(execResult *ExecuteResult, results []domain.TestResult) string {
	passedCount := 0
	for _, r := range results {
		if r.Passed {
			passedCount++
		}
	}

	if execResult.Passed {
		return "All tests passed!"
	}
	return fmt.Sprintf("Tests passed: %d/%d", passedCount, len(results))
}

// buildDetailedFeedback builds detailed feedback based on test results.
func (e *ExerciseEvaluator) buildDetailedFeedback(execResult *ExecuteResult, results []domain.TestResult, passedTests, failedTests []string) string {
	var builder strings.Builder

	// Overall assessment
	if execResult.Passed {
		builder.WriteString("All tests passed! Your code correctly solves the exercise.\n\n")
	} else {
		builder.WriteString(fmt.Sprintf("Your code passed %d out of %d tests.\n\n", len(passedTests), len(results)))
	}

	// Specific feedback on failures
	if len(failedTests) > 0 {
		builder.WriteString("Issues found:\n")
		for _, r := range results {
			if !r.Passed {
				if r.Error != "" {
					builder.WriteString(fmt.Sprintf("  - %s: %s\n", r.TestName, r.Error))
				} else {
					builder.WriteString(fmt.Sprintf("  - %s: Expected '%s' but got '%s'\n", r.TestName, r.Expected, r.Actual))
				}
			}
		}
		builder.WriteString("\n")
	}

	// Learning suggestions
	if execResult.Score < 50 {
		builder.WriteString("Suggestions:\n")
		builder.WriteString("  - Review the exercise requirements carefully\n")
		builder.WriteString("  - Check your code syntax and logic\n")
		builder.WriteString("  - Try running your code with the test inputs manually\n")
	} else if execResult.Score < 100 {
		builder.WriteString("Suggestions:\n")
		builder.WriteString("  - You're close! Check the edge cases in the failing tests\n")
	}

	return builder.String()
}

// GetPassedTests returns a list of passed test names from results.
func GetPassedTests(results []domain.TestResult) []string {
	passed := make([]string, 0)
	for _, r := range results {
		if r.Passed {
			passed = append(passed, r.TestName)
		}
	}
	return passed
}

// GetFailedTests returns a list of failed test names from results.
func GetFailedTests(results []domain.TestResult) []string {
	failed := make([]string, 0)
	for _, r := range results {
		if !r.Passed {
			failed = append(failed, r.TestName)
		}
	}
	return failed
}
