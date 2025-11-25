package golang

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// Executor provides Go code execution
type Executor struct {
	language types.Language
	tempDir  string
}

// NewExecutor creates a new Go executor
func NewExecutor() *Executor {
	lang, _ := types.GetLanguage("go")
	return &Executor{
		language: lang,
		tempDir:  os.TempDir(),
	}
}

// GetLanguage returns the Go language
func (e *Executor) GetLanguage() types.Language {
	return e.language
}

// Execute runs Go code and returns the result
func (e *Executor) Execute(ctx context.Context, request types.ExecutionRequest) (*types.ExecutionResult, error) {
	startTime := time.Now()

	// Validate code first
	if err := e.ValidateCode(ctx, request.Code); err != nil {
		return &types.ExecutionResult{
			Success:       false,
			Error:         err.Error(),
			ExitCode:      1,
			ExecutionTime: time.Since(startTime).Milliseconds(),
		}, nil
	}

	// Create temporary directory for execution
	execDir, err := os.MkdirTemp(e.tempDir, "go-exec-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(execDir)

	// Write code to file
	codeFile := filepath.Join(execDir, "main.go")
	if err := os.WriteFile(codeFile, []byte(request.Code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	// Write additional files if provided
	for filename, content := range request.Files {
		filePath := filepath.Join(execDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file %s: %w", filename, err)
		}
	}

	// Initialize go.mod if not present
	if _, exists := request.Files["go.mod"]; !exists {
		modContent := "module temp\n\ngo 1.22\n"
		modFile := filepath.Join(execDir, "go.mod")
		if err := os.WriteFile(modFile, []byte(modContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write go.mod: %w", err)
		}
	}

	// Set timeout
	timeout := 30 * time.Second
	if request.Timeout > 0 {
		timeout = time.Duration(request.Timeout) * time.Second
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Run the code
	cmd := exec.CommandContext(execCtx, "go", "run", "main.go")
	cmd.Dir = execDir

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range request.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Set stdin if provided
	if request.Input != "" {
		cmd.Stdin = strings.NewReader(request.Input)
	}

	// Capture output
	output, err := cmd.CombinedOutput()
	executionTime := time.Since(startTime).Milliseconds()

	result := &types.ExecutionResult{
		ExecutionTime: executionTime,
	}

	if err != nil {
		// Check if it was a timeout
		if execCtx.Err() == context.DeadlineExceeded {
			result.Success = false
			result.Error = "Execution timeout exceeded"
			result.ExitCode = 124 // Standard timeout exit code
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			result.Success = false
			result.Error = string(output)
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.Success = false
			result.Error = err.Error()
			result.ExitCode = 1
		}
	} else {
		result.Success = true
		result.Output = string(output)
		result.ExitCode = 0
	}

	return result, nil
}

// ValidateCode checks if Go code is safe to execute
func (e *Executor) ValidateCode(ctx context.Context, code string) error {
	if code == "" {
		return &types.ToolError{
			Code:    "EMPTY_CODE",
			Message: "Code cannot be empty",
		}
	}

	// Check for dangerous operations
	dangerousPatterns := []string{
		"os.RemoveAll",
		"os.Remove",
		"exec.Command",
		"syscall",
		"unsafe",
		"net/http",
		"net.Dial",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(code, pattern) {
			return &types.ToolError{
				Code:    "UNSAFE_CODE",
				Message: fmt.Sprintf("Code contains potentially unsafe operation: %s", pattern),
				Details: "For security reasons, certain operations are restricted",
			}
		}
	}

	return nil
}

// GetResourceLimits returns default resource limits for Go execution
func (e *Executor) GetResourceLimits() types.ResourceLimits {
	return types.ResourceLimits{
		MaxMemoryMB:      512,
		MaxCPUTime:       30,
		MaxProcesses:     5,
		MaxFileSize:      10 * 1024 * 1024, // 10MB
		MaxOutputSize:    1 * 1024 * 1024,  // 1MB
		NetworkAccess:    false,
		FileSystemAccess: false,
	}
}

// SupportsInteractive returns whether interactive execution is supported
func (e *Executor) SupportsInteractive() bool {
	return false
}

