package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"

	"github.com/google/uuid"
)

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	config types.AgentConfig
	llm    types.LLMProvider
	tools  []types.Tool
	memory types.Memory
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(config types.AgentConfig, llm types.LLMProvider) *BaseAgent {
	return &BaseAgent{
		config: config,
		llm:    llm,
		tools:  make([]types.Tool, 0),
	}
}

// GetMemory returns the agent's memory
func (a *BaseAgent) GetMemory() types.Memory {
	return a.memory
}

// GetTools returns the agent's tools
func (a *BaseAgent) GetTools() []types.Tool {
	return a.tools
}

// GetConfig returns the agent's configuration
func (a *BaseAgent) GetConfig() types.AgentConfig {
	return a.config
}

// SetMemory sets the agent's memory
func (a *BaseAgent) SetMemory(memory types.Memory) {
	a.memory = memory
}

// AddTool adds a tool to the agent
func (a *BaseAgent) AddTool(tool types.Tool) {
	a.tools = append(a.tools, tool)
}

// AddTools adds multiple tools to the agent
func (a *BaseAgent) AddTools(tools []types.Tool) {
	a.tools = append(a.tools, tools...)
}

// GetTool retrieves a tool by name
func (a *BaseAgent) GetTool(name string) (types.Tool, error) {
	for _, tool := range a.tools {
		if tool.Name() == name {
			return tool, nil
		}
	}
	return nil, errors.NewToolNotFoundError(name)
}

// ExecuteTool executes a tool with the given input
func (a *BaseAgent) ExecuteTool(ctx context.Context, toolName string, input types.ToolInput) (*types.ToolOutput, error) {
	tool, err := a.GetTool(toolName)
	if err != nil {
		return nil, err
	}

	// Validate input
	if err := tool.Validate(input); err != nil {
		return nil, errors.NewToolInvalidInputError(toolName, err.Error())
	}

	// Execute tool
	startTime := time.Now()
	output, err := tool.Execute(ctx, input)
	if err != nil {
		return nil, errors.NewToolExecutionError(toolName, err)
	}

	// Update metadata
	output.Metadata.ExecutionTime = time.Since(startTime).Milliseconds()

	return output, nil
}

// GenerateLLMResponse generates a response from the LLM
func (a *BaseAgent) GenerateLLMResponse(ctx context.Context, messages []types.Message) (*types.LLMResponse, error) {
	request := types.NewLLMRequest(messages)
	request.Temperature = a.config.Temperature
	request.MaxTokens = a.config.MaxTokens

	// Add function definitions if tools are available
	if len(a.tools) > 0 && a.llm.SupportsFunctionCalling() {
		request.Functions = a.getToolFunctions()
	}

	response, err := a.llm.Generate(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// getToolFunctions converts tools to function definitions
func (a *BaseAgent) getToolFunctions() []types.Function {
	functions := make([]types.Function, len(a.tools))
	for i, tool := range a.tools {
		schema := tool.GetSchema()
		functions[i] = types.Function{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  convertSchemaToMap(schema),
		}
	}
	return functions
}

// convertSchemaToMap converts a ToolSchema to a map
func convertSchemaToMap(schema types.ToolSchema) map[string]interface{} {
	result := make(map[string]interface{})
	result["type"] = schema.Type
	result["properties"] = schema.Properties
	if len(schema.Required) > 0 {
		result["required"] = schema.Required
	}
	if schema.Description != "" {
		result["description"] = schema.Description
	}
	return result
}

// AddToMemory adds a message to memory if memory is enabled
func (a *BaseAgent) AddToMemory(ctx context.Context, message types.Message) error {
	if a.memory == nil || !a.config.EnableMemory {
		return nil
	}

	return a.memory.Add(ctx, message)
}

// GetFromMemory retrieves messages from memory
func (a *BaseAgent) GetFromMemory(ctx context.Context, limit int) ([]types.Message, error) {
	if a.memory == nil || !a.config.EnableMemory {
		return []types.Message{}, nil
	}

	return a.memory.Get(ctx, limit)
}

// BuildMessages builds the message list for LLM including memory
func (a *BaseAgent) BuildMessages(ctx context.Context, systemPrompt, userQuery string) ([]types.Message, error) {
	messages := []types.Message{}

	// Add system prompt
	if systemPrompt != "" {
		messages = append(messages, types.NewSystemMessage(systemPrompt))
	}

	// Add memory if enabled
	if a.config.EnableMemory && a.memory != nil {
		memoryMessages, err := a.GetFromMemory(ctx, 10)
		if err != nil {
			return nil, err
		}
		messages = append(messages, memoryMessages...)
	}

	// Add current user query
	messages = append(messages, types.NewUserMessage(userQuery))

	return messages, nil
}

// CreateAgentOutput creates a new agent output
func (a *BaseAgent) CreateAgentOutput() *types.AgentOutput {
	output := types.NewAgentOutput()
	output.Metadata.AgentType = a.config.Name
	return output
}

// ValidateInput validates agent input
func (a *BaseAgent) ValidateInput(input types.AgentInput) error {
	if input.Query == "" {
		return errors.NewAgentInvalidInputError("query cannot be empty")
	}

	if input.MaxSteps < 0 {
		return errors.NewAgentInvalidInputError("max_steps must be positive")
	}

	if input.Temperature < 0 || input.Temperature > 2 {
		return errors.NewAgentInvalidInputError("temperature must be between 0 and 2")
	}

	return nil
}

// CreateStep creates a new agent step
func CreateStep(stepNumber int, thought, action string, actionInput map[string]interface{}) types.AgentStep {
	return types.AgentStep{
		StepNumber:  stepNumber,
		Thought:     thought,
		Action:      action,
		ActionInput: actionInput,
		Timestamp:   time.Now(),
	}
}

// CreateToolCall creates a new tool call record
func CreateToolCall(toolName string, input types.ToolInput, output *types.ToolOutput) types.ToolCall {
	return types.ToolCall{
		ID:        uuid.New().String(),
		ToolName:  toolName,
		Input:     input,
		Output:    output,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix(),
	}
}

// LogVerbose logs if verbose logging is enabled
func (a *BaseAgent) LogVerbose(format string, args ...interface{}) {
	if a.config.VerboseLogging {
		fmt.Printf("[AGENT] "+format+"\n", args...)
	}
}

