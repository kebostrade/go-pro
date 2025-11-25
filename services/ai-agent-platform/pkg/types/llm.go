package types

import (
	"context"
	"time"
)

// LLMProvider defines the interface for language model providers
type LLMProvider interface {
	// Generate generates a completion for the given prompt
	Generate(ctx context.Context, request LLMRequest) (*LLMResponse, error)

	// GenerateStream generates a completion and streams the response
	GenerateStream(ctx context.Context, request LLMRequest) (<-chan LLMStreamChunk, error)

	// GetModelInfo returns information about the model
	GetModelInfo() ModelInfo

	// GetProviderName returns the name of the provider
	GetProviderName() string

	// SupportsStreaming returns whether the provider supports streaming
	SupportsStreaming() bool

	// SupportsFunctionCalling returns whether the provider supports function calling
	SupportsFunctionCalling() bool
}

// LLMRequest represents a request to an LLM
type LLMRequest struct {
	// Messages in the conversation
	Messages []Message `json:"messages"`

	// Model to use for generation
	Model string `json:"model,omitempty"`

	// Temperature controls randomness (0.0 to 2.0)
	Temperature float32 `json:"temperature,omitempty"`

	// MaxTokens maximum tokens to generate
	MaxTokens int `json:"max_tokens,omitempty"`

	// TopP nucleus sampling parameter
	TopP float32 `json:"top_p,omitempty"`

	// FrequencyPenalty penalizes frequent tokens
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`

	// PresencePenalty penalizes tokens that have appeared
	PresencePenalty float32 `json:"presence_penalty,omitempty"`

	// Stop sequences that will stop generation
	Stop []string `json:"stop,omitempty"`

	// Functions available for function calling
	Functions []Function `json:"functions,omitempty"`

	// FunctionCall controls function calling behavior
	FunctionCall interface{} `json:"function_call,omitempty"`

	// Stream whether to stream the response
	Stream bool `json:"stream,omitempty"`

	// User identifier for tracking
	User string `json:"user,omitempty"`
}

// LLMResponse represents a response from an LLM
type LLMResponse struct {
	// Content is the generated text
	Content string `json:"content"`

	// FunctionCall if the model wants to call a function
	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// FinishReason why the generation stopped
	FinishReason string `json:"finish_reason"`

	// Usage token usage information
	Usage TokenUsage `json:"usage"`

	// Model that was used
	Model string `json:"model"`

	// Metadata additional response metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Latency time taken to generate
	Latency time.Duration `json:"latency"`
}

// LLMStreamChunk represents a chunk in a streaming response
type LLMStreamChunk struct {
	// Content is the text in this chunk
	Content string `json:"content"`

	// FunctionCall partial function call information
	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// FinishReason if this is the last chunk
	FinishReason string `json:"finish_reason,omitempty"`

	// Done indicates if streaming is complete
	Done bool `json:"done"`

	// Error if an error occurred
	Error error `json:"error,omitempty"`
}

// Message represents a message in a conversation
type Message struct {
	// Role of the message sender (system, user, assistant, function)
	Role MessageRole `json:"role"`

	// Content of the message
	Content string `json:"content"`

	// Name of the function (for function messages)
	Name string `json:"name,omitempty"`

	// FunctionCall if this message contains a function call
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

// MessageRole defines the role of a message sender
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleFunction  MessageRole = "function"
	RoleTool      MessageRole = "tool"
)

// Function represents a function that can be called by the LLM
type Function struct {
	// Name of the function
	Name string `json:"name"`

	// Description of what the function does
	Description string `json:"description"`

	// Parameters JSON schema for the function parameters
	Parameters map[string]interface{} `json:"parameters"`
}

// FunctionCall represents a function call from the LLM
type FunctionCall struct {
	// Name of the function to call
	Name string `json:"name"`

	// Arguments JSON string of arguments
	Arguments string `json:"arguments"`
}

// ModelInfo contains information about an LLM model
type ModelInfo struct {
	// Name of the model
	Name string `json:"name"`

	// Provider of the model
	Provider string `json:"provider"`

	// MaxTokens maximum context window
	MaxTokens int `json:"max_tokens"`

	// CostPer1KTokens pricing information
	CostPer1KTokens CostInfo `json:"cost_per_1k_tokens"`

	// SupportsStreaming whether streaming is supported
	SupportsStreaming bool `json:"supports_streaming"`

	// SupportsFunctions whether function calling is supported
	SupportsFunctions bool `json:"supports_functions"`

	// SupportsVision whether vision/image input is supported
	SupportsVision bool `json:"supports_vision"`
}

// CostInfo represents pricing information
type CostInfo struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

// NewSystemMessage creates a system message
func NewSystemMessage(content string) Message {
	return Message{
		Role:    RoleSystem,
		Content: content,
	}
}

// NewUserMessage creates a user message
func NewUserMessage(content string) Message {
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}

// NewAssistantMessage creates an assistant message
func NewAssistantMessage(content string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

// NewFunctionMessage creates a function result message
func NewFunctionMessage(name, content string) Message {
	return Message{
		Role:    RoleFunction,
		Name:    name,
		Content: content,
	}
}

// NewLLMRequest creates a new LLM request with defaults
func NewLLMRequest(messages []Message) LLMRequest {
	return LLMRequest{
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   2000,
		TopP:        1.0,
	}
}

