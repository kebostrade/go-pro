package programming

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/common"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// CodeExecutionTool executes code in a sandboxed environment
type CodeExecutionTool struct {
	languageRegistry *common.LanguageRegistry
}

// NewCodeExecutionTool creates a new code execution tool
func NewCodeExecutionTool(registry *common.LanguageRegistry) *CodeExecutionTool {
	return &CodeExecutionTool{
		languageRegistry: registry,
	}
}

// Name returns the tool name
func (t *CodeExecutionTool) Name() string {
	return "code_execution"
}

// Description returns the tool description
func (t *CodeExecutionTool) Description() string {
	return "Executes code in a secure sandboxed environment with resource limits. Supports multiple programming languages. Returns output, errors, and execution metrics."
}

// GetSchema returns the tool's input schema
func (t *CodeExecutionTool) GetSchema() types.ToolSchema {
	return types.ToolSchema{
		Type: "object",
		Properties: map[string]types.PropertySchema{
			"code": {
				Type:        "string",
				Description: "The code to execute",
			},
			"language": {
				Type:        "string",
				Description: "Programming language (go, python, javascript, etc.)",
			},
			"input": {
				Type:        "string",
				Description: "Standard input for the program",
			},
			"timeout": {
				Type:        "number",
				Description: "Execution timeout in seconds (default: 30)",
				Default:     30,
			},
			"arguments": {
				Type:        "array",
				Description: "Command-line arguments",
				Items: &types.PropertySchema{
					Type: "string",
				},
			},
		},
		Required: []string{"code", "language"},
	}
}

// Validate validates the input
func (t *CodeExecutionTool) Validate(input types.ToolInput) error {
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

// Execute runs the code
func (t *CodeExecutionTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
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
	stdin, _ := input.GetString("input")
	timeout, ok := input.GetInt("timeout")
	if !ok {
		timeout = 30
	}

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

	// Validate code safety
	if err := provider.ValidateCode(ctx, code); err != nil {
		return &types.ToolOutput{
			Error: &types.ToolError{
				Code:    "UNSAFE_CODE",
				Message: err.Error(),
			},
			Metadata: types.ToolMetadata{
				ExecutionTime: time.Since(startTime).Milliseconds(),
				Success:       false,
			},
		}, nil
	}

	// Create execution request
	execRequest := types.ExecutionRequest{
		Code:          code,
		Language:      language,
		Input:         stdin,
		Timeout:       timeout,
		Environment:   make(map[string]string),
		Files:         make(map[string]string),
	}

	// Get arguments if provided
	if args, ok := input.GetSlice("arguments"); ok {
		execRequest.Arguments = make([]string, len(args))
		for i, arg := range args {
			if str, ok := arg.(string); ok {
				execRequest.Arguments[i] = str
			}
		}
	}

	// Execute code
	result, err := provider.Execute(ctx, execRequest)
	if err != nil {
		return &types.ToolOutput{
			Error: &types.ToolError{
				Code:    "EXECUTION_ERROR",
				Message: fmt.Sprintf("Failed to execute code: %v", err),
			},
			Metadata: types.ToolMetadata{
				ExecutionTime: time.Since(startTime).Milliseconds(),
				Success:       false,
			},
		}, nil
	}

	// Convert result to JSON
	resultJSON, _ := json.MarshalIndent(result, "", "  ")

	return &types.ToolOutput{
		Result: string(resultJSON),
		Metadata: types.ToolMetadata{
			ExecutionTime: time.Since(startTime).Milliseconds(),
			Success:       result.Success,
			AdditionalInfo: map[string]interface{}{
				"language":       language,
				"exit_code":      result.ExitCode,
				"execution_time": result.ExecutionTime,
			},
		},
	}, nil
}

