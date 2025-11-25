package types

import (
	"context"
)

// Language represents a programming language
type Language struct {
	// Name of the language (go, python, javascript, etc.)
	Name string `json:"name"`

	// DisplayName human-readable name
	DisplayName string `json:"display_name"`

	// Version language version
	Version string `json:"version,omitempty"`

	// FileExtensions common file extensions
	FileExtensions []string `json:"file_extensions"`

	// Supported whether this language is supported
	Supported bool `json:"supported"`

	// Features supported features for this language
	Features LanguageFeatures `json:"features"`
}

// LanguageFeatures defines what features are supported for a language
type LanguageFeatures struct {
	// Execution whether code can be executed
	Execution bool `json:"execution"`

	// StaticAnalysis whether static analysis is available
	StaticAnalysis bool `json:"static_analysis"`

	// Linting whether linting is available
	Linting bool `json:"linting"`

	// Formatting whether code formatting is available
	Formatting bool `json:"formatting"`

	// Testing whether test generation is supported
	Testing bool `json:"testing"`

	// Debugging whether debugging assistance is available
	Debugging bool `json:"debugging"`

	// Documentation whether doc generation is supported
	Documentation bool `json:"documentation"`
}

// LanguageAnalyzer defines the interface for language-specific analysis
type LanguageAnalyzer interface {
	// GetLanguage returns the language this analyzer handles
	GetLanguage() Language

	// Analyze performs static analysis on code
	Analyze(ctx context.Context, code string) (*CodeAnalysis, error)

	// Lint runs linting on code
	Lint(ctx context.Context, code string) ([]CodeIssue, error)

	// Format formats code according to language standards
	Format(ctx context.Context, code string) (string, error)

	// Validate checks if code is syntactically valid
	Validate(ctx context.Context, code string) error

	// ExtractImports extracts import/require statements
	ExtractImports(ctx context.Context, code string) ([]string, error)

	// GetComplexity calculates cyclomatic complexity
	GetComplexity(ctx context.Context, code string) (int, error)
}

// LanguageExecutor defines the interface for code execution
type LanguageExecutor interface {
	// GetLanguage returns the language this executor handles
	GetLanguage() Language

	// Execute runs code and returns the result
	Execute(ctx context.Context, request ExecutionRequest) (*ExecutionResult, error)

	// ValidateCode checks if code is safe to execute
	ValidateCode(ctx context.Context, code string) error

	// GetResourceLimits returns default resource limits
	GetResourceLimits() ResourceLimits

	// SupportsInteractive whether interactive execution is supported
	SupportsInteractive() bool
}

// ExecutionRequest represents a code execution request
type ExecutionRequest struct {
	// Code to execute
	Code string `json:"code"`

	// Language programming language
	Language string `json:"language"`

	// Input stdin input for the program
	Input string `json:"input,omitempty"`

	// Arguments command-line arguments
	Arguments []string `json:"arguments,omitempty"`

	// Environment variables
	Environment map[string]string `json:"environment,omitempty"`

	// Timeout execution timeout in seconds
	Timeout int `json:"timeout,omitempty"`

	// ResourceLimits resource constraints
	ResourceLimits *ResourceLimits `json:"resource_limits,omitempty"`

	// WorkingDirectory working directory for execution
	WorkingDirectory string `json:"working_directory,omitempty"`

	// Files additional files needed for execution
	Files map[string]string `json:"files,omitempty"`
}

// ResourceLimits defines resource constraints for code execution
type ResourceLimits struct {
	// MaxMemoryMB maximum memory in megabytes
	MaxMemoryMB int `json:"max_memory_mb"`

	// MaxCPUTime maximum CPU time in seconds
	MaxCPUTime int `json:"max_cpu_time"`

	// MaxProcesses maximum number of processes
	MaxProcesses int `json:"max_processes"`

	// MaxFileSize maximum file size in bytes
	MaxFileSize int64 `json:"max_file_size"`

	// MaxOutputSize maximum output size in bytes
	MaxOutputSize int64 `json:"max_output_size"`

	// NetworkAccess whether network access is allowed
	NetworkAccess bool `json:"network_access"`

	// FileSystemAccess whether file system access is allowed
	FileSystemAccess bool `json:"file_system_access"`
}

// SupportedLanguages defines all supported programming languages
var SupportedLanguages = map[string]Language{
	"go": {
		Name:           "go",
		DisplayName:    "Go",
		Version:        "1.22",
		FileExtensions: []string{".go"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        true,
			Debugging:      true,
			Documentation:  true,
		},
	},
	"python": {
		Name:           "python",
		DisplayName:    "Python",
		Version:        "3.11",
		FileExtensions: []string{".py"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        true,
			Debugging:      true,
			Documentation:  true,
		},
	},
	"javascript": {
		Name:           "javascript",
		DisplayName:    "JavaScript",
		Version:        "ES2023",
		FileExtensions: []string{".js", ".mjs"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        true,
			Debugging:      true,
			Documentation:  true,
		},
	},
	"typescript": {
		Name:           "typescript",
		DisplayName:    "TypeScript",
		Version:        "5.0",
		FileExtensions: []string{".ts"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        true,
			Debugging:      true,
			Documentation:  true,
		},
	},
	"rust": {
		Name:           "rust",
		DisplayName:    "Rust",
		Version:        "1.75",
		FileExtensions: []string{".rs"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        true,
			Debugging:      true,
			Documentation:  true,
		},
	},
	"java": {
		Name:           "java",
		DisplayName:    "Java",
		Version:        "21",
		FileExtensions: []string{".java"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        true,
			Debugging:      true,
			Documentation:  true,
		},
	},
	"cpp": {
		Name:           "cpp",
		DisplayName:    "C++",
		Version:        "C++20",
		FileExtensions: []string{".cpp", ".cc", ".cxx"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        false,
			Debugging:      true,
			Documentation:  false,
		},
	},
	"c": {
		Name:           "c",
		DisplayName:    "C",
		Version:        "C17",
		FileExtensions: []string{".c", ".h"},
		Supported:      true,
		Features: LanguageFeatures{
			Execution:      true,
			StaticAnalysis: true,
			Linting:        true,
			Formatting:     true,
			Testing:        false,
			Debugging:      true,
			Documentation:  false,
		},
	},
}

// GetLanguage retrieves a language by name
func GetLanguage(name string) (Language, bool) {
	lang, ok := SupportedLanguages[name]
	return lang, ok
}

// IsLanguageSupported checks if a language is supported
func IsLanguageSupported(name string) bool {
	lang, ok := SupportedLanguages[name]
	return ok && lang.Supported
}

// GetSupportedLanguageNames returns a list of all supported language names
func GetSupportedLanguageNames() []string {
	names := make([]string, 0, len(SupportedLanguages))
	for name, lang := range SupportedLanguages {
		if lang.Supported {
			names = append(names, name)
		}
	}
	return names
}

// DetectLanguage attempts to detect the programming language from code
func DetectLanguage(code string) string {
	// Simple heuristic-based detection
	// In production, use a proper language detection library
	
	// Check for common patterns
	if containsPattern(code, []string{"package main", "func main()", "import ("}) {
		return "go"
	}
	if containsPattern(code, []string{"def ", "import ", "print("}) {
		return "python"
	}
	if containsPattern(code, []string{"function ", "const ", "let ", "var "}) {
		return "javascript"
	}
	if containsPattern(code, []string{"interface ", "type ", ": string", ": number"}) {
		return "typescript"
	}
	if containsPattern(code, []string{"fn main()", "let mut", "impl "}) {
		return "rust"
	}
	if containsPattern(code, []string{"public class", "public static void main", "System.out"}) {
		return "java"
	}
	if containsPattern(code, []string{"#include <", "std::", "int main("}) {
		return "cpp"
	}
	if containsPattern(code, []string{"#include <", "printf(", "int main("}) {
		return "c"
	}
	
	return "unknown"
}

// containsPattern checks if code contains any of the patterns
func containsPattern(code string, patterns []string) bool {
	for _, pattern := range patterns {
		if len(code) > 0 && len(pattern) > 0 {
			// Simple substring check
			// In production, use regex or AST parsing
			for i := 0; i <= len(code)-len(pattern); i++ {
				if code[i:i+len(pattern)] == pattern {
					return true
				}
			}
		}
	}
	return false
}

