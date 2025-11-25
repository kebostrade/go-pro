package sandbox

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// DockerSandbox provides Docker-based code execution
type DockerSandbox struct {
	config SandboxConfig
}

// SandboxConfig holds sandbox configuration
type SandboxConfig struct {
	// ImagePrefix for Docker images (e.g., "coding-sandbox")
	ImagePrefix string

	// NetworkMode for containers (none, bridge, host)
	NetworkMode string

	// DefaultTimeout for execution
	DefaultTimeout time.Duration

	// TempDir for temporary files
	TempDir string

	// EnableLogging whether to log container output
	EnableLogging bool

	// MaxConcurrent maximum concurrent executions
	MaxConcurrent int
}

// NewDockerSandbox creates a new Docker sandbox
func NewDockerSandbox(config SandboxConfig) *DockerSandbox {
	if config.ImagePrefix == "" {
		config.ImagePrefix = "coding-sandbox"
	}
	if config.NetworkMode == "" {
		config.NetworkMode = "none" // No network by default
	}
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 30 * time.Second
	}
	if config.TempDir == "" {
		config.TempDir = os.TempDir()
	}
	if config.MaxConcurrent == 0 {
		config.MaxConcurrent = 10
	}

	return &DockerSandbox{
		config: config,
	}
}

// Execute runs code in a Docker container
func (s *DockerSandbox) Execute(ctx context.Context, request types.ExecutionRequest) (*types.ExecutionResult, error) {
	startTime := time.Now()

	// Create execution context with timeout
	timeout := s.config.DefaultTimeout
	if request.Timeout > 0 {
		timeout = time.Duration(request.Timeout) * time.Second
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create temporary directory for this execution
	execDir, err := os.MkdirTemp(s.config.TempDir, "sandbox-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(execDir)

	// Write code and files to temp directory
	if err := s.prepareFiles(execDir, request); err != nil {
		return nil, fmt.Errorf("failed to prepare files: %w", err)
	}

	// Get Docker image for language
	image := s.getImageForLanguage(request.Language)

	// Build Docker run command
	dockerCmd := s.buildDockerCommand(image, execDir, request)

	// Execute in Docker
	result, err := s.runDocker(execCtx, dockerCmd, request.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to run Docker: %w", err)
	}

	result.ExecutionTime = time.Since(startTime).Milliseconds()
	return result, nil
}

// prepareFiles writes code and additional files to the execution directory
func (s *DockerSandbox) prepareFiles(execDir string, request types.ExecutionRequest) error {
	// Determine main file name based on language
	mainFile := s.getMainFileName(request.Language)
	
	// Write main code file
	codePath := filepath.Join(execDir, mainFile)
	if err := os.WriteFile(codePath, []byte(request.Code), 0644); err != nil {
		return fmt.Errorf("failed to write code file: %w", err)
	}

	// Write additional files
	for filename, content := range request.Files {
		filePath := filepath.Join(execDir, filename)
		
		// Create subdirectories if needed
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}
	}

	return nil
}

// getMainFileName returns the main file name for a language
func (s *DockerSandbox) getMainFileName(language string) string {
	switch language {
	case "go":
		return "main.go"
	case "python":
		return "main.py"
	case "javascript":
		return "main.js"
	case "typescript":
		return "main.ts"
	case "rust":
		return "main.rs"
	case "java":
		return "Main.java"
	case "cpp":
		return "main.cpp"
	case "c":
		return "main.c"
	default:
		return "main.txt"
	}
}

// getImageForLanguage returns the Docker image for a language
func (s *DockerSandbox) getImageForLanguage(language string) string {
	switch language {
	case "go":
		return fmt.Sprintf("%s-go:latest", s.config.ImagePrefix)
	case "python":
		return fmt.Sprintf("%s-python:latest", s.config.ImagePrefix)
	case "javascript":
		return fmt.Sprintf("%s-node:latest", s.config.ImagePrefix)
	case "typescript":
		return fmt.Sprintf("%s-node:latest", s.config.ImagePrefix)
	case "rust":
		return fmt.Sprintf("%s-rust:latest", s.config.ImagePrefix)
	case "java":
		return fmt.Sprintf("%s-java:latest", s.config.ImagePrefix)
	case "cpp":
		return fmt.Sprintf("%s-cpp:latest", s.config.ImagePrefix)
	case "c":
		return fmt.Sprintf("%s-c:latest", s.config.ImagePrefix)
	default:
		return "alpine:latest"
	}
}

// buildDockerCommand builds the Docker run command
func (s *DockerSandbox) buildDockerCommand(image, execDir string, request types.ExecutionRequest) []string {
	cmd := []string{
		"docker", "run",
		"--rm",                                    // Remove container after execution
		"--network", s.config.NetworkMode,        // Network mode
		"-v", fmt.Sprintf("%s:/workspace", execDir), // Mount code directory
		"-w", "/workspace",                        // Working directory
	}

	// Add resource limits
	if request.ResourceLimits != nil {
		limits := request.ResourceLimits
		
		// Memory limit
		if limits.MaxMemoryMB > 0 {
			cmd = append(cmd, "--memory", fmt.Sprintf("%dm", limits.MaxMemoryMB))
		}
		
		// CPU limit
		if limits.MaxCPUTime > 0 {
			cmd = append(cmd, "--cpus", fmt.Sprintf("%.2f", float64(limits.MaxCPUTime)/100.0))
		}
		
		// Process limit
		if limits.MaxProcesses > 0 {
			cmd = append(cmd, "--pids-limit", fmt.Sprintf("%d", limits.MaxProcesses))
		}
	} else {
		// Default limits
		cmd = append(cmd,
			"--memory", "512m",
			"--cpus", "1.0",
			"--pids-limit", "50",
		)
	}

	// Security options
	cmd = append(cmd,
		"--security-opt", "no-new-privileges",
		"--cap-drop", "ALL",
		"--read-only",                           // Read-only root filesystem
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=100m", // Writable tmp
	)

	// Add image
	cmd = append(cmd, image)

	// Add execution command based on language
	cmd = append(cmd, s.getExecutionCommand(request.Language)...)

	return cmd
}

// getExecutionCommand returns the command to execute code
func (s *DockerSandbox) getExecutionCommand(language string) []string {
	switch language {
	case "go":
		return []string{"go", "run", "main.go"}
	case "python":
		return []string{"python3", "main.py"}
	case "javascript":
		return []string{"node", "main.js"}
	case "typescript":
		return []string{"ts-node", "main.ts"}
	case "rust":
		return []string{"sh", "-c", "rustc main.rs && ./main"}
	case "java":
		return []string{"sh", "-c", "javac Main.java && java Main"}
	case "cpp":
		return []string{"sh", "-c", "g++ -o main main.cpp && ./main"}
	case "c":
		return []string{"sh", "-c", "gcc -o main main.c && ./main"}
	default:
		return []string{"cat", "main.txt"}
	}
}

// runDocker executes the Docker command
func (s *DockerSandbox) runDocker(ctx context.Context, cmd []string, stdin string) (*types.ExecutionResult, error) {
	// Note: This is a simplified implementation
	// In production, use the Docker SDK for Go (github.com/docker/docker/client)
	
	result := &types.ExecutionResult{
		Success: false,
	}

	// For now, return a mock result
	// In production, this would actually execute Docker
	result.Success = true
	result.Output = "Docker execution not yet implemented. Use local execution for now."
	result.ExitCode = 0

	return result, nil
}

// ValidateImage checks if a Docker image exists
func (s *DockerSandbox) ValidateImage(language string) error {
	_ = s.getImageForLanguage(language) // Will be used when Docker SDK is implemented

	// In production, check if image exists using Docker SDK
	// For now, just validate the language is supported
	supportedLanguages := []string{"go", "python", "javascript", "typescript", "rust", "java", "cpp", "c"}
	
	for _, lang := range supportedLanguages {
		if language == lang {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported language: %s", language)
}

// Cleanup removes old temporary files
func (s *DockerSandbox) Cleanup() error {
	// Remove old temp directories
	pattern := filepath.Join(s.config.TempDir, "sandbox-*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}

		// Remove directories older than 1 hour
		if time.Since(info.ModTime()) > time.Hour {
			os.RemoveAll(match)
		}
	}

	return nil
}

