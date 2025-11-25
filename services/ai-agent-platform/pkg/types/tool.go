package types

import (
	"context"
	"encoding/json"
)

// Tool represents a tool that an agent can use
type Tool interface {
	// Name returns the name of the tool
	Name() string

	// Description returns a description of what the tool does
	Description() string

	// Execute runs the tool with the given input
	Execute(ctx context.Context, input ToolInput) (*ToolOutput, error)

	// GetSchema returns the JSON schema for the tool's input parameters
	GetSchema() ToolSchema

	// Validate validates the input before execution
	Validate(input ToolInput) error
}

// ToolInput represents input to a tool
type ToolInput struct {
	// Parameters for the tool execution
	Parameters map[string]interface{} `json:"parameters"`

	// Context provides additional context
	Context map[string]interface{} `json:"context,omitempty"`

	// UserID of the user making the request
	UserID string `json:"user_id,omitempty"`

	// SessionID of the current session
	SessionID string `json:"session_id,omitempty"`
}

// ToolOutput represents output from a tool
type ToolOutput struct {
	// Result is the output from the tool
	Result interface{} `json:"result"`

	// Error if the tool execution failed
	Error *ToolError `json:"error,omitempty"`

	// Metadata about the execution
	Metadata ToolMetadata `json:"metadata"`
}

// ToolError represents an error from tool execution
type ToolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *ToolError) Error() string {
	if e.Details != "" {
		return e.Message + ": " + e.Details
	}
	return e.Message
}

// ToolMetadata contains metadata about tool execution
type ToolMetadata struct {
	// ExecutionTime how long the tool took to execute
	ExecutionTime int64 `json:"execution_time_ms"`

	// Success whether the tool executed successfully
	Success bool `json:"success"`

	// Cached whether the result was cached
	Cached bool `json:"cached,omitempty"`

	// AdditionalInfo any additional information
	AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
}

// ToolSchema defines the schema for a tool's parameters
type ToolSchema struct {
	// Type is typically "object"
	Type string `json:"type"`

	// Properties defines the parameters
	Properties map[string]PropertySchema `json:"properties"`

	// Required lists required parameter names
	Required []string `json:"required,omitempty"`

	// Description of the schema
	Description string `json:"description,omitempty"`
}

// PropertySchema defines a single parameter's schema
type PropertySchema struct {
	// Type of the property (string, number, boolean, object, array)
	Type string `json:"type"`

	// Description of the property
	Description string `json:"description"`

	// Enum possible values (for string types)
	Enum []string `json:"enum,omitempty"`

	// Items schema for array items
	Items *PropertySchema `json:"items,omitempty"`

	// Properties for nested objects
	Properties map[string]PropertySchema `json:"properties,omitempty"`

	// Default value
	Default interface{} `json:"default,omitempty"`

	// Minimum value (for numbers)
	Minimum *float64 `json:"minimum,omitempty"`

	// Maximum value (for numbers)
	Maximum *float64 `json:"maximum,omitempty"`
}

// ToolCall represents a call to a tool
type ToolCall struct {
	// ID unique identifier for this tool call
	ID string `json:"id"`

	// ToolName name of the tool that was called
	ToolName string `json:"tool_name"`

	// Input that was provided to the tool
	Input ToolInput `json:"input"`

	// Output from the tool
	Output *ToolOutput `json:"output,omitempty"`

	// StartTime when the call started
	StartTime int64 `json:"start_time"`

	// EndTime when the call completed
	EndTime int64 `json:"end_time,omitempty"`

	// Error if the call failed
	Error *ToolError `json:"error,omitempty"`
}

// ToolRegistry manages available tools
type ToolRegistry interface {
	// Register adds a tool to the registry
	Register(tool Tool) error

	// Get retrieves a tool by name
	Get(name string) (Tool, error)

	// List returns all registered tools
	List() []Tool

	// Unregister removes a tool from the registry
	Unregister(name string) error

	// GetByCategory returns tools in a specific category
	GetByCategory(category string) []Tool
}

// NewToolInput creates a new ToolInput
func NewToolInput(params map[string]interface{}) ToolInput {
	return ToolInput{
		Parameters: params,
		Context:    make(map[string]interface{}),
	}
}

// NewToolOutput creates a new ToolOutput
func NewToolOutput(result interface{}) *ToolOutput {
	return &ToolOutput{
		Result: result,
		Metadata: ToolMetadata{
			Success: true,
		},
	}
}

// NewToolError creates a new ToolError
func NewToolError(code, message, details string) *ToolError {
	return &ToolError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// GetString safely gets a string parameter
func (ti ToolInput) GetString(key string) (string, bool) {
	val, ok := ti.Parameters[key]
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

// GetInt safely gets an int parameter
func (ti ToolInput) GetInt(key string) (int, bool) {
	val, ok := ti.Parameters[key]
	if !ok {
		return 0, false
	}
	
	switch v := val.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case json.Number:
		i, err := v.Int64()
		return int(i), err == nil
	default:
		return 0, false
	}
}

// GetFloat safely gets a float parameter
func (ti ToolInput) GetFloat(key string) (float64, bool) {
	val, ok := ti.Parameters[key]
	if !ok {
		return 0, false
	}
	
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case json.Number:
		f, err := v.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

// GetBool safely gets a bool parameter
func (ti ToolInput) GetBool(key string) (bool, bool) {
	val, ok := ti.Parameters[key]
	if !ok {
		return false, false
	}
	b, ok := val.(bool)
	return b, ok
}

// GetMap safely gets a map parameter
func (ti ToolInput) GetMap(key string) (map[string]interface{}, bool) {
	val, ok := ti.Parameters[key]
	if !ok {
		return nil, false
	}
	m, ok := val.(map[string]interface{})
	return m, ok
}

// GetSlice safely gets a slice parameter
func (ti ToolInput) GetSlice(key string) ([]interface{}, bool) {
	val, ok := ti.Parameters[key]
	if !ok {
		return nil, false
	}
	s, ok := val.([]interface{})
	return s, ok
}

