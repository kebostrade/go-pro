package schema

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleTool      MessageRole = "tool"
)

type Message struct {
	Role       MessageRole `json:"role"`
	Content    string      `json:"content"`
	Name       string      `json:"name,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

func NewUserMessage(content string) Message {
	return Message{
		Role:      RoleUser,
		Content:   content,
		Timestamp: time.Now(),
	}
}

func NewSystemMessage(content string) Message {
	return Message{
		Role:      RoleSystem,
		Content:   content,
		Timestamp: time.Now(),
	}
}

func NewAssistantMessage(content string) Message {
	return Message{
		Role:      RoleAssistant,
		Content:   content,
		Timestamp: time.Now(),
	}
}

func NewToolMessage(content string, toolCallID string) Message {
	return Message{
		Role:       RoleTool,
		Content:    content,
		ToolCallID: toolCallID,
		Timestamp:  time.Now(),
	}
}

type ToolCall struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	ToolCallID string      `json:"tool_call_id"`
	ToolName   string      `json:"tool_name"`
	Output     interface{} `json:"output"`
	Error      string      `json:"error,omitempty"`
}

type AgentInput struct {
	Query       string
	SessionID   string
	UserID      string
	Context     map[string]interface{}
	Messages    []Message
	MaxSteps    int
	Temperature float32
}

type AgentOutput struct {
	Output      string
	Messages    []Message
	ToolCalls   []ToolCall
	ToolResults []ToolResult
	SessionID   string
	ExecutionID string
	StartTime   time.Time
	EndTime     time.Time
	TokenUsage  TokenUsage
	Error       error
}

type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

func NewAgentInput(query string) AgentInput {
	return AgentInput{
		Query:       query,
		SessionID:   uuid.New().String(),
		Context:     make(map[string]interface{}),
		Messages:    []Message{},
		MaxSteps:    5,
		Temperature: 0.7,
	}
}

func (a *AgentOutput) AddMessage(msg Message) {
	a.Messages = append(a.Messages, msg)
}

func (a *AgentOutput) AddToolCall(tc ToolCall) {
	a.ToolCalls = append(a.ToolCalls, tc)
}

func (a *AgentOutput) AddToolResult(tr ToolResult) {
	a.ToolResults = append(a.ToolResults, tr)
}

func (a *AgentOutput) Finalize() {
	a.EndTime = time.Now()
}

type Tool struct {
	Name        string
	Description string
	ArgsSchema  interface{}
	Func        func(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

func (t *Tool) Invoke(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if t.Func != nil {
		return t.Func(ctx, args)
	}
	return nil, nil
}

type Chain interface {
	Run(ctx context.Context, input interface{}) (interface{}, error)
}

type LLM interface {
	Generate(ctx context.Context, messages []Message) (Message, error)
	GenerateWithTools(ctx context.Context, messages []Message, tools []Tool) (Message, error)
}
