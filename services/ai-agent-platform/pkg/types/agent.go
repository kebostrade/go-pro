package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Agent represents the core agent interface that all agents must implement
type Agent interface {
	// Run executes the agent with the given input and returns the output
	Run(ctx context.Context, input AgentInput) (*AgentOutput, error)

	// Stream executes the agent and streams events as they occur
	Stream(ctx context.Context, input AgentInput) (<-chan AgentEvent, error)

	// GetMemory returns the agent's memory system
	GetMemory() Memory

	// GetTools returns the tools available to the agent
	GetTools() []Tool

	// GetConfig returns the agent's configuration
	GetConfig() AgentConfig
}

// AgentInput represents the input to an agent
type AgentInput struct {
	// Query is the user's input/question
	Query string `json:"query"`

	// SessionID identifies the conversation session
	SessionID string `json:"session_id,omitempty"`

	// UserID identifies the user making the request
	UserID string `json:"user_id,omitempty"`

	// Context provides additional context for the agent
	Context map[string]interface{} `json:"context,omitempty"`

	// MaxSteps limits the number of reasoning steps
	MaxSteps int `json:"max_steps,omitempty"`

	// Temperature controls randomness in LLM responses (0.0 to 1.0)
	Temperature float32 `json:"temperature,omitempty"`

	// Metadata for tracking and logging
	Metadata map[string]string `json:"metadata,omitempty"`
}

// AgentOutput represents the output from an agent
type AgentOutput struct {
	// Output is the final response from the agent
	Output string `json:"output"`

	// Steps contains the reasoning steps taken by the agent
	Steps []AgentStep `json:"steps,omitempty"`

	// ToolCalls contains information about tools that were called
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// Metadata contains additional information about the execution
	Metadata AgentMetadata `json:"metadata"`

	// Error contains error information if the agent failed
	Error *AgentError `json:"error,omitempty"`
}

// AgentStep represents a single reasoning step in the agent's execution
type AgentStep struct {
	// StepNumber is the sequential number of this step
	StepNumber int `json:"step_number"`

	// Thought is the agent's reasoning at this step
	Thought string `json:"thought"`

	// Action is the action the agent decided to take
	Action string `json:"action,omitempty"`

	// ActionInput is the input to the action
	ActionInput map[string]interface{} `json:"action_input,omitempty"`

	// Observation is the result of the action
	Observation string `json:"observation,omitempty"`

	// Timestamp when this step occurred
	Timestamp time.Time `json:"timestamp"`

	// Duration of this step
	Duration time.Duration `json:"duration"`
}

// AgentMetadata contains metadata about the agent execution
type AgentMetadata struct {
	// ExecutionID uniquely identifies this execution
	ExecutionID string `json:"execution_id"`

	// AgentType identifies the type of agent
	AgentType string `json:"agent_type"`

	// StartTime when execution started
	StartTime time.Time `json:"start_time"`

	// EndTime when execution completed
	EndTime time.Time `json:"end_time"`

	// Duration of the entire execution
	Duration time.Duration `json:"duration"`

	// TokensUsed tracks LLM token usage
	TokensUsed TokenUsage `json:"tokens_used"`

	// Cost estimated cost of the execution
	Cost float64 `json:"cost,omitempty"`

	// ModelUsed the LLM model that was used
	ModelUsed string `json:"model_used,omitempty"`
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// AgentError represents an error that occurred during agent execution
type AgentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// AgentEvent represents an event during streaming execution
type AgentEvent struct {
	// Type of event (thought, action, observation, output, error)
	Type AgentEventType `json:"type"`

	// Data associated with the event
	Data interface{} `json:"data"`

	// Timestamp when the event occurred
	Timestamp time.Time `json:"timestamp"`

	// StepNumber if this event is part of a step
	StepNumber int `json:"step_number,omitempty"`
}

// AgentEventType defines the types of events that can occur
type AgentEventType string

const (
	EventTypeThought      AgentEventType = "thought"
	EventTypeAction       AgentEventType = "action"
	EventTypeObservation  AgentEventType = "observation"
	EventTypeOutput       AgentEventType = "output"
	EventTypeError        AgentEventType = "error"
	EventTypeToolCall     AgentEventType = "tool_call"
	EventTypeToolResult   AgentEventType = "tool_result"
	EventTypeStreamChunk  AgentEventType = "stream_chunk"
)

// AgentConfig holds configuration for an agent
type AgentConfig struct {
	// Name of the agent
	Name string `json:"name"`

	// Description of what the agent does
	Description string `json:"description"`

	// MaxSteps maximum number of reasoning steps
	MaxSteps int `json:"max_steps"`

	// MaxTokens maximum tokens for LLM responses
	MaxTokens int `json:"max_tokens"`

	// Temperature for LLM responses
	Temperature float32 `json:"temperature"`

	// Timeout for agent execution
	Timeout time.Duration `json:"timeout"`

	// EnableStreaming whether to support streaming
	EnableStreaming bool `json:"enable_streaming"`

	// EnableMemory whether to use memory
	EnableMemory bool `json:"enable_memory"`

	// VerboseLogging whether to log detailed information
	VerboseLogging bool `json:"verbose_logging"`

	// CustomPrompt custom system prompt for the agent
	CustomPrompt string `json:"custom_prompt,omitempty"`
}

// NewAgentInput creates a new AgentInput with defaults
func NewAgentInput(query string) AgentInput {
	return AgentInput{
		Query:       query,
		SessionID:   uuid.New().String(),
		MaxSteps:    5,
		Temperature: 0.7,
		Context:     make(map[string]interface{}),
		Metadata:    make(map[string]string),
	}
}

// NewAgentOutput creates a new AgentOutput
func NewAgentOutput() *AgentOutput {
	return &AgentOutput{
		Steps:     make([]AgentStep, 0),
		ToolCalls: make([]ToolCall, 0),
		Metadata: AgentMetadata{
			ExecutionID: uuid.New().String(),
			StartTime:   time.Now(),
			TokensUsed:  TokenUsage{},
		},
	}
}

// AddStep adds a step to the agent output
func (o *AgentOutput) AddStep(step AgentStep) {
	o.Steps = append(o.Steps, step)
}

// SetError sets an error on the agent output
func (o *AgentOutput) SetError(code, message, details string) {
	o.Error = &AgentError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Finalize completes the agent output metadata
func (o *AgentOutput) Finalize() {
	o.Metadata.EndTime = time.Now()
	o.Metadata.Duration = o.Metadata.EndTime.Sub(o.Metadata.StartTime)
}

// AddTokenUsage adds token usage to the metadata
func (o *AgentOutput) AddTokenUsage(prompt, completion int) {
	o.Metadata.TokensUsed.PromptTokens += prompt
	o.Metadata.TokensUsed.CompletionTokens += completion
	o.Metadata.TokensUsed.TotalTokens += (prompt + completion)
}

