package programming

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/common"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// CodeAnalysisTool performs static code analysis
type CodeAnalysisTool struct {
	languageRegistry *common.LanguageRegistry
}

// NewCodeAnalysisTool creates a new code analysis tool
func NewCodeAnalysisTool(registry *common.LanguageRegistry) *CodeAnalysisTool {
	return &CodeAnalysisTool{
		languageRegistry: registry,
	}
}

// Name returns the tool name
func (t *CodeAnalysisTool) Name() string {
	return "code_analysis"
}

// Description returns the tool description
func (t *CodeAnalysisTool) Description() string {
	return "Analyzes code for quality, security issues, performance problems, and best practice violations. Supports multiple programming languages including Go, Python, JavaScript, and more."
}

// GetSchema returns the tool's input schema
func (t *CodeAnalysisTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"code": {
				Type:        "string",
				Description: "The code to analyze",
			},
			"language": {
				Type:        "string",
				Description: "Programming language (go, python, javascript, etc.)",
			},
			"check_security": {
				Type:        "boolean",
				Description: "Whether to check for security vulnerabilities",
				Default:     true,
			},
			"check_performance": {
				Type:        "boolean",
				Description: "Whether to check for performance issues",
				Default:     true,
			},
			"check_best_practices": {
				Type:        "boolean",
				Description: "Whether to check for best practice violations",
				Default:     true,
			},
		},
		Required: []string{"code", "language"},
	}
}

// Validate validates the input
func (t *CodeAnalysisTool) Validate(input types.ToolInput) error {
	code, ok := input.GetString("code")
	if !ok || code == "" {
		return &types.ToolError{
			Code:    "MISSING_CODE",
			Message: "Code parameter is required",
		}
	}

	language, ok := input.GetString("language")
	if !ok || language == "" {
		return &types.ToolError{
			Code:    "MISSING_LANGUAGE",
			Message: "Language parameter is required",
		}
	}

	if !t.languageRegistry.IsSupported(language) {
		return &types.ToolError{
			Code:    "UNSUPPORTED_LANGUAGE",
			Message: fmt.Sprintf("Language '%s' is not supported", language),
		}
	}

	return nil
}

// Execute performs code analysis
func (t *CodeAnalysisTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
	startTime := time.Now()

	// Validate input
	if err := t.Validate(input); err != nil {
		return &types.ToolOutput{
			Error: err.(*types.ToolError),
			Metadata: types.ToolMetadata{
				ExecutionTime: time.Since(startTime).Milliseconds(),
				Success:       false,
			},
		}, nil
	}

	// Get parameters
	code, _ := input.GetString("code")
	language, _ := input.GetString("language")

	// Get language provider
	provider, err := t.languageRegistry.Get(language)
	if err != nil {
		return &types.ToolOutput{
			Error: &types.ToolError{
				Code:    "LANGUAGE_ERROR",
				Message: err.Error(),
			},
			Metadata: types.ToolMetadata{
				ExecutionTime: time.Since(startTime).Milliseconds(),
				Success:       false,
			},
		}, nil
	}

	// Perform analysis
	analysis, err := provider.Analyze(ctx, code)
	if err != nil {
		return &types.ToolOutput{
			Error: &types.ToolError{
				Code:    "ANALYSIS_ERROR",
				Message: fmt.Sprintf("Failed to analyze code: %v", err),
			},
			Metadata: types.ToolMetadata{
				ExecutionTime: time.Since(startTime).Milliseconds(),
				Success:       false,
			},
		}, nil
	}

	// Filter results based on options
	checkSecurity, _ := input.GetBool("check_security")
	checkPerformance, _ := input.GetBool("check_performance")
	checkBestPractices, _ := input.GetBool("check_best_practices")

	if !checkSecurity {
		analysis.SecurityIssues = nil
	}
	if !checkPerformance {
		analysis.PerformanceIssues = nil
	}
	if !checkBestPractices {
		analysis.BestPractices = nil
	}

	// Convert to JSON for result
	resultJSON, _ := json.MarshalIndent(analysis, "", "  ")

	return &types.ToolOutput{
		Result: string(resultJSON),
		Metadata: types.ToolMetadata{
			ExecutionTime: time.Since(startTime).Milliseconds(),
			Success:       true,
			AdditionalInfo: map[string]interface{}{
				"language":      language,
				"overall_score": analysis.OverallScore,
				"issues_found":  len(analysis.Issues),
			},
		},
	}, nil
}

