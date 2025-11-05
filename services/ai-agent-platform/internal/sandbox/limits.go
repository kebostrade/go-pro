package sandbox

import (
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// ResourceLimitsManager manages resource limits for code execution
type ResourceLimitsManager struct {
	defaultLimits map[string]types.ResourceLimits
	customLimits  map[string]types.ResourceLimits
}

// NewResourceLimitsManager creates a new resource limits manager
func NewResourceLimitsManager() *ResourceLimitsManager {
	manager := &ResourceLimitsManager{
		defaultLimits: make(map[string]types.ResourceLimits),
		customLimits:  make(map[string]types.ResourceLimits),
	}

	// Initialize default limits for each language
	manager.initializeDefaultLimits()

	return manager
}

// initializeDefaultLimits sets up default resource limits
func (m *ResourceLimitsManager) initializeDefaultLimits() {
	// Go defaults
	m.defaultLimits["go"] = types.ResourceLimits{
		MaxMemoryMB:      512,
		MaxCPUTime:       30,
		MaxProcesses:     10,
		MaxFileSize:      10 * 1024 * 1024, // 10MB
		MaxOutputSize:    1 * 1024 * 1024,  // 1MB
		NetworkAccess:    false,
		FileSystemAccess: false,
	}

	// Python defaults
	m.defaultLimits["python"] = types.ResourceLimits{
		MaxMemoryMB:      512,
		MaxCPUTime:       30,
		MaxProcesses:     5,
		MaxFileSize:      10 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}

	// JavaScript defaults
	m.defaultLimits["javascript"] = types.ResourceLimits{
		MaxMemoryMB:      256,
		MaxCPUTime:       20,
		MaxProcesses:     5,
		MaxFileSize:      5 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}

	// TypeScript defaults
	m.defaultLimits["typescript"] = types.ResourceLimits{
		MaxMemoryMB:      256,
		MaxCPUTime:       20,
		MaxProcesses:     5,
		MaxFileSize:      5 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}

	// Rust defaults
	m.defaultLimits["rust"] = types.ResourceLimits{
		MaxMemoryMB:      1024, // Rust compilation needs more memory
		MaxCPUTime:       60,   // Compilation takes longer
		MaxProcesses:     10,
		MaxFileSize:      20 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}

	// Java defaults
	m.defaultLimits["java"] = types.ResourceLimits{
		MaxMemoryMB:      1024, // JVM needs more memory
		MaxCPUTime:       45,
		MaxProcesses:     10,
		MaxFileSize:      20 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}

	// C++ defaults
	m.defaultLimits["cpp"] = types.ResourceLimits{
		MaxMemoryMB:      512,
		MaxCPUTime:       30,
		MaxProcesses:     5,
		MaxFileSize:      10 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}

	// C defaults
	m.defaultLimits["c"] = types.ResourceLimits{
		MaxMemoryMB:      256,
		MaxCPUTime:       20,
		MaxProcesses:     5,
		MaxFileSize:      5 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}
}

// GetLimits returns resource limits for a language
func (m *ResourceLimitsManager) GetLimits(language string) types.ResourceLimits {
	// Check custom limits first
	if limits, ok := m.customLimits[language]; ok {
		return limits
	}

	// Fall back to default limits
	if limits, ok := m.defaultLimits[language]; ok {
		return limits
	}

	// Return conservative defaults if language not found
	return types.ResourceLimits{
		MaxMemoryMB:      256,
		MaxCPUTime:       20,
		MaxProcesses:     5,
		MaxFileSize:      5 * 1024 * 1024,
		MaxOutputSize:    1 * 1024 * 1024,
		NetworkAccess:    false,
		FileSystemAccess: false,
	}
}

// SetCustomLimits sets custom resource limits for a language
func (m *ResourceLimitsManager) SetCustomLimits(language string, limits types.ResourceLimits) error {
	// Validate limits
	if err := m.ValidateLimits(limits); err != nil {
		return err
	}

	m.customLimits[language] = limits
	return nil
}

// ValidateLimits validates resource limits
func (m *ResourceLimitsManager) ValidateLimits(limits types.ResourceLimits) error {
	// Check memory limits
	if limits.MaxMemoryMB < 64 {
		return fmt.Errorf("memory limit too low: minimum 64MB")
	}
	if limits.MaxMemoryMB > 4096 {
		return fmt.Errorf("memory limit too high: maximum 4096MB")
	}

	// Check CPU time limits
	if limits.MaxCPUTime < 1 {
		return fmt.Errorf("CPU time limit too low: minimum 1 second")
	}
	if limits.MaxCPUTime > 300 {
		return fmt.Errorf("CPU time limit too high: maximum 300 seconds")
	}

	// Check process limits
	if limits.MaxProcesses < 1 {
		return fmt.Errorf("process limit too low: minimum 1")
	}
	if limits.MaxProcesses > 100 {
		return fmt.Errorf("process limit too high: maximum 100")
	}

	// Check file size limits
	if limits.MaxFileSize > 100*1024*1024 {
		return fmt.Errorf("file size limit too high: maximum 100MB")
	}

	// Check output size limits
	if limits.MaxOutputSize > 10*1024*1024 {
		return fmt.Errorf("output size limit too high: maximum 10MB")
	}

	return nil
}

// MergeLimits merges custom limits with defaults
func (m *ResourceLimitsManager) MergeLimits(language string, custom *types.ResourceLimits) types.ResourceLimits {
	defaults := m.GetLimits(language)

	if custom == nil {
		return defaults
	}

	// Merge custom with defaults
	merged := defaults

	if custom.MaxMemoryMB > 0 {
		merged.MaxMemoryMB = custom.MaxMemoryMB
	}
	if custom.MaxCPUTime > 0 {
		merged.MaxCPUTime = custom.MaxCPUTime
	}
	if custom.MaxProcesses > 0 {
		merged.MaxProcesses = custom.MaxProcesses
	}
	if custom.MaxFileSize > 0 {
		merged.MaxFileSize = custom.MaxFileSize
	}
	if custom.MaxOutputSize > 0 {
		merged.MaxOutputSize = custom.MaxOutputSize
	}

	// Network and file system access are explicitly set
	merged.NetworkAccess = custom.NetworkAccess
	merged.FileSystemAccess = custom.FileSystemAccess

	return merged
}

// GetTimeoutDuration returns timeout as duration
func (m *ResourceLimitsManager) GetTimeoutDuration(limits types.ResourceLimits) time.Duration {
	// Add buffer to CPU time for overhead
	timeout := time.Duration(limits.MaxCPUTime+10) * time.Second
	
	// Cap at 5 minutes
	maxTimeout := 5 * time.Minute
	if timeout > maxTimeout {
		timeout = maxTimeout
	}

	return timeout
}

// EstimateCost estimates the cost of execution based on limits
func (m *ResourceLimitsManager) EstimateCost(limits types.ResourceLimits) float64 {
	// Simple cost estimation based on resources
	// In production, this would use actual cloud provider pricing
	
	memoryCost := float64(limits.MaxMemoryMB) * 0.0001  // $0.0001 per MB
	cpuCost := float64(limits.MaxCPUTime) * 0.001       // $0.001 per second
	
	return memoryCost + cpuCost
}

// GetRecommendedLimits returns recommended limits based on code complexity
func (m *ResourceLimitsManager) GetRecommendedLimits(language string, codeSize int, complexity int) types.ResourceLimits {
	base := m.GetLimits(language)

	// Adjust based on code size
	if codeSize > 10000 {
		base.MaxMemoryMB = int(float64(base.MaxMemoryMB) * 1.5)
		base.MaxCPUTime = int(float64(base.MaxCPUTime) * 1.5)
	}

	// Adjust based on complexity
	if complexity > 20 {
		base.MaxCPUTime = int(float64(base.MaxCPUTime) * 2.0)
	}

	// Ensure limits are within bounds
	if base.MaxMemoryMB > 2048 {
		base.MaxMemoryMB = 2048
	}
	if base.MaxCPUTime > 120 {
		base.MaxCPUTime = 120
	}

	return base
}

// GetLimitsSummary returns a human-readable summary of limits
func (m *ResourceLimitsManager) GetLimitsSummary(limits types.ResourceLimits) string {
	return fmt.Sprintf(
		"Memory: %dMB, CPU: %ds, Processes: %d, Network: %v, FileSystem: %v",
		limits.MaxMemoryMB,
		limits.MaxCPUTime,
		limits.MaxProcesses,
		limits.NetworkAccess,
		limits.FileSystemAccess,
	)
}

// CheckLimitsExceeded checks if execution exceeded limits
func (m *ResourceLimitsManager) CheckLimitsExceeded(result *types.ExecutionResult, limits types.ResourceLimits) []string {
	violations := make([]string, 0)

	// Check memory (if tracked)
	if result.MemoryUsed > 0 && result.MemoryUsed > int64(limits.MaxMemoryMB)*1024*1024 {
		violations = append(violations, fmt.Sprintf(
			"Memory limit exceeded: %d MB > %d MB",
			result.MemoryUsed/(1024*1024),
			limits.MaxMemoryMB,
		))
	}

	// Check CPU time (if tracked)
	if result.CPUTime > 0 && result.CPUTime > int64(limits.MaxCPUTime)*1000 {
		violations = append(violations, fmt.Sprintf(
			"CPU time limit exceeded: %d ms > %d ms",
			result.CPUTime,
			limits.MaxCPUTime*1000,
		))
	}

	// Check execution time
	if result.ExecutionTime > int64(limits.MaxCPUTime)*1000 {
		violations = append(violations, fmt.Sprintf(
			"Execution timeout: %d ms > %d ms",
			result.ExecutionTime,
			limits.MaxCPUTime*1000,
		))
	}

	// Check output size
	outputSize := int64(len(result.Output) + len(result.Error))
	if outputSize > limits.MaxOutputSize {
		violations = append(violations, fmt.Sprintf(
			"Output size exceeded: %d bytes > %d bytes",
			outputSize,
			limits.MaxOutputSize,
		))
	}

	return violations
}

