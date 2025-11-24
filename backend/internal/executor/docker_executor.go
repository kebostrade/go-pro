// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package executor provides secure code execution using Docker containers.
package executor

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go-pro-backend/internal/service"
)

const (
	// Default execution constraints
	defaultTimeout     = 5 * time.Second
	defaultMemoryLimit = "128m"
	defaultCPULimit    = "0.5"
	defaultTmpfsSize   = "10m"
	maxCodeSize        = 65536 // 64KB max code size

	// Docker image
	goImage = "golang:1.23-alpine"

	// Security settings
	containerUser = "1000:1000"
)

var (
	// Dangerous imports that should be blocked (multiline with any whitespace/newlines)
	dangerousImports = regexp.MustCompile(`(?s)import\s+\([^)]*"(os|net|syscall|unsafe|runtime/debug)"`)

	// Alternative dangerous single imports
	dangerousSingleImport = regexp.MustCompile(`import\s+"(os|net|syscall|unsafe|runtime/debug)"`)
)

// DockerExecutor implements ExecutorService using Docker containers.
type DockerExecutor struct {
	image      string
	timeout    time.Duration
	memory     string
	cpuLimit   string
	tmpfsSize  string
}

// NewDockerExecutor creates a new Docker-based code executor.
func NewDockerExecutor() service.ExecutorService {
	return &DockerExecutor{
		image:     goImage,
		timeout:   defaultTimeout,
		memory:    defaultMemoryLimit,
		cpuLimit:  defaultCPULimit,
		tmpfsSize: defaultTmpfsSize,
	}
}

// ExecuteCode executes Go code in a sandboxed Docker container.
func (e *DockerExecutor) ExecuteCode(ctx context.Context, req *service.ExecuteRequest) (*service.ExecuteResult, error) {
	startTime := time.Now()

	// Validate code before execution
	if err := e.validateCode(req.Code); err != nil {
		return &service.ExecuteResult{
			Passed:        false,
			Score:         0,
			Results:       []service.TestResult{},
			ExecutionTime: time.Since(startTime),
			Error:         err,
		}, nil
	}

	// Create temporary directory for code
	tempDir, err := os.MkdirTemp("", "go-executor-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Always cleanup

	// Write code to main.go
	codeFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(codeFile, []byte(req.Code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	// Execute test cases
	testResults := e.runTestCases(ctx, tempDir, req.TestCases)

	// Calculate score
	passedCount := 0
	for _, result := range testResults {
		if result.Passed {
			passedCount++
		}
	}

	score := 0
	if len(req.TestCases) > 0 {
		score = (passedCount * 100) / len(req.TestCases)
	}

	return &service.ExecuteResult{
		Passed:        passedCount == len(req.TestCases),
		Score:         score,
		Results:       testResults,
		ExecutionTime: time.Since(startTime),
		Error:         nil,
	}, nil
}

// validateCode performs basic validation and security checks.
func (e *DockerExecutor) validateCode(code string) error {
	// Check code size
	if len(code) > maxCodeSize {
		return fmt.Errorf("code too large: max %d bytes allowed", maxCodeSize)
	}

	// Check for dangerous imports (multi-line)
	if dangerousImports.MatchString(code) {
		return fmt.Errorf("dangerous imports detected: os, net, syscall, unsafe, and runtime/debug are not allowed")
	}

	// Check for dangerous imports (single-line)
	if dangerousSingleImport.MatchString(code) {
		return fmt.Errorf("dangerous imports detected: os, net, syscall, unsafe, and runtime/debug are not allowed")
	}

	// Basic Go code structure validation
	if !strings.Contains(code, "package main") {
		return fmt.Errorf("code must contain 'package main'")
	}

	if !strings.Contains(code, "func main()") {
		return fmt.Errorf("code must contain 'func main()'")
	}

	return nil
}

// runTestCases executes code for each test case and compares outputs.
func (e *DockerExecutor) runTestCases(ctx context.Context, codeDir string, tests []service.TestCase) []service.TestResult {
	results := make([]service.TestResult, len(tests))

	for i, test := range tests {
		// Create context with timeout for each test
		testCtx, cancel := context.WithTimeout(ctx, e.timeout)

		// Execute code with test input
		output, err := e.runContainer(testCtx, codeDir, test.Input)
		cancel()

		// Create test result
		result := service.TestResult{
			TestName: test.Name,
			Expected: test.Expected,
			Actual:   output,
		}

		if err != nil {
			result.Passed = false
			result.Error = e.formatError(err)
		} else {
			// Compare output (trim whitespace for comparison)
			expectedTrimmed := strings.TrimSpace(test.Expected)
			actualTrimmed := strings.TrimSpace(output)

			result.Passed = expectedTrimmed == actualTrimmed

			if !result.Passed {
				result.Error = "Output does not match expected result"
			}
		}

		results[i] = result
	}

	return results
}

// runContainer executes code in a Docker container with security constraints.
func (e *DockerExecutor) runContainer(ctx context.Context, codeDir string, input string) (string, error) {
	// Create input file if provided
	var inputFile string
	if input != "" {
		inputFile = filepath.Join(codeDir, "input.txt")
		if err := os.WriteFile(inputFile, []byte(input), 0644); err != nil {
			return "", fmt.Errorf("failed to write input file: %w", err)
		}
	}

	// Build docker run command with security constraints
	args := []string{
		"run",
		"--rm",                                // Remove container after execution
		"--memory=" + e.memory,                // Memory limit
		"--cpus=" + e.cpuLimit,                // CPU limit
		"--network=none",                      // No network access
		"--read-only",                         // Read-only filesystem
		"--tmpfs=/tmp:rw,noexec,nosuid,size=" + e.tmpfsSize, // Temporary filesystem
		"--user=" + containerUser,             // Non-root user
		"-v", codeDir + ":/code:ro",           // Mount code directory as read-only
		"-w", "/code",                         // Working directory
		e.image,                               // Docker image
		"sh", "-c",                            // Shell command
	}

	// Build execution command
	var execCmd string
	if input != "" {
		// Redirect input from file
		execCmd = fmt.Sprintf("timeout %ds go run main.go < input.txt", int(e.timeout.Seconds()))
	} else {
		execCmd = fmt.Sprintf("timeout %ds go run main.go", int(e.timeout.Seconds()))
	}

	args = append(args, execCmd)

	// Execute Docker command
	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Handle execution errors
	if err != nil {
		// Check for timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("execution timeout")
		}

		// Check stderr for compilation or runtime errors
		stderrStr := stderr.String()
		if stderrStr != "" {
			return "", fmt.Errorf("execution error: %s", e.extractErrorMessage(stderrStr))
		}

		return "", fmt.Errorf("execution failed: %w", err)
	}

	return stdout.String(), nil
}

// extractErrorMessage extracts meaningful error messages from stderr.
func (e *DockerExecutor) extractErrorMessage(stderr string) string {
	// Split by lines and find relevant error messages
	scanner := bufio.NewScanner(strings.NewReader(stderr))
	var errorLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and progress messages
		if line == "" || strings.Contains(line, "go: downloading") {
			continue
		}

		// Compilation errors
		if strings.Contains(line, "syntax error") ||
		   strings.Contains(line, "undefined:") ||
		   strings.Contains(line, "cannot use") ||
		   strings.Contains(line, "not enough arguments") {
			errorLines = append(errorLines, line)
		}

		// Runtime errors
		if strings.Contains(line, "panic:") ||
		   strings.Contains(line, "runtime error:") {
			errorLines = append(errorLines, line)
		}
	}

	if len(errorLines) > 0 {
		// Return first few error lines (max 3)
		maxLines := 3
		if len(errorLines) < maxLines {
			maxLines = len(errorLines)
		}
		return strings.Join(errorLines[:maxLines], "\n")
	}

	// If no specific errors found, return truncated stderr
	if len(stderr) > 200 {
		return stderr[:200] + "..."
	}
	return stderr
}

// formatError formats error messages for user-friendly display.
func (e *DockerExecutor) formatError(err error) string {
	errMsg := err.Error()

	// Timeout errors
	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded") {
		return fmt.Sprintf("Code execution timed out (%v limit)", e.timeout)
	}

	// Compilation errors
	if strings.Contains(errMsg, "syntax error") {
		return "Compilation error: " + errMsg
	}

	// Runtime errors
	if strings.Contains(errMsg, "panic") || strings.Contains(errMsg, "runtime error") {
		return "Runtime error: " + errMsg
	}

	// Generic execution error
	if strings.Contains(errMsg, "execution error") {
		return errMsg
	}

	// Fallback
	return "Execution failed: " + errMsg
}
