package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// ReActAgent implements the ReAct (Reasoning + Acting) pattern
// This agent alternates between reasoning about what to do and taking actions
type ReActAgent struct {
	*BaseAgent
}

// ReActConfig holds configuration for ReAct agent
type ReActConfig struct {
	Name            string
	Description     string
	LLM             types.LLMProvider
	Tools           []types.Tool
	Memory          types.Memory
	MaxSteps        int
	Temperature     float32
	MaxTokens       int
	Timeout         time.Duration
	VerboseLogging  bool
	CustomPrompt    string
}

// NewReActAgent creates a new ReAct agent
func NewReActAgent(config ReActConfig) *ReActAgent {
	if config.MaxSteps == 0 {
		config.MaxSteps = 5
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 2000
	}
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Minute
	}

	agentConfig := types.AgentConfig{
		Name:            config.Name,
		Description:     config.Description,
		MaxSteps:        config.MaxSteps,
		MaxTokens:       config.MaxTokens,
		Temperature:     config.Temperature,
		Timeout:         config.Timeout,
		EnableStreaming: false,
		EnableMemory:    config.Memory != nil,
		VerboseLogging:  config.VerboseLogging,
		CustomPrompt:    config.CustomPrompt,
	}

	base := NewBaseAgent(agentConfig, config.LLM)
	base.SetMemory(config.Memory)
	base.AddTools(config.Tools)

	return &ReActAgent{
		BaseAgent: base,
	}
}

// Run executes the ReAct agent
func (a *ReActAgent) Run(ctx context.Context, input types.AgentInput) (*types.AgentOutput, error) {
	// Validate input
	if err := a.ValidateInput(input); err != nil {
		return nil, err
	}

	// Create output
	output := a.CreateAgentOutput()

	// Set timeout
	ctx, cancel := context.WithTimeout(ctx, a.config.Timeout)
	defer cancel()

	// Get max steps
	maxSteps := input.MaxSteps
	if maxSteps == 0 {
		maxSteps = a.config.MaxSteps
	}

	// Build system prompt
	systemPrompt := a.buildSystemPrompt()

	// Execute reasoning loop
	currentQuery := input.Query
	for step := 1; step <= maxSteps; step++ {
		a.LogVerbose("Step %d: Starting reasoning", step)

		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, errors.NewAgentTimeoutError("agent execution timeout")
		default:
		}

		// Build messages
		messages, err := a.BuildMessages(ctx, systemPrompt, currentQuery)
		if err != nil {
			return nil, err
		}

		// Generate response
		response, err := a.GenerateLLMResponse(ctx, messages)
		if err != nil {
			return nil, err
		}

		// Update token usage
		output.AddTokenUsage(response.Usage.PromptTokens, response.Usage.CompletionTokens)

		// Parse response
		thought, action, actionInput, finalAnswer, err := a.parseResponse(response.Content)
		if err != nil {
			a.LogVerbose("Failed to parse response: %v", err)
			output.SetError("PARSE_ERROR", "Failed to parse agent response", err.Error())
			break
		}

		// Create step
		agentStep := CreateStep(step, thought, action, actionInput)

		// Check if we have a final answer
		if finalAnswer != "" {
			a.LogVerbose("Final answer: %s", finalAnswer)
			agentStep.Observation = finalAnswer
			agentStep.Duration = time.Since(agentStep.Timestamp)
			output.AddStep(agentStep)
			output.Output = finalAnswer
			break
		}

		// Execute action
		if action != "" {
			a.LogVerbose("Executing action: %s", action)

			toolInput := types.NewToolInput(actionInput)
			toolInput.SessionID = input.SessionID
			toolInput.UserID = input.UserID

			toolOutput, err := a.ExecuteTool(ctx, action, toolInput)
			if err != nil {
				a.LogVerbose("Tool execution failed: %v", err)
				agentStep.Observation = fmt.Sprintf("Error: %v", err)
			} else {
				// Convert result to string
				resultStr, _ := json.Marshal(toolOutput.Result)
				agentStep.Observation = string(resultStr)

				// Record tool call
				toolCall := CreateToolCall(action, toolInput, toolOutput)
				output.ToolCalls = append(output.ToolCalls, toolCall)
			}

			// Update current query for next iteration
			currentQuery = fmt.Sprintf("Previous thought: %s\nAction taken: %s\nObservation: %s\n\nWhat should I do next?",
				thought, action, agentStep.Observation)
		}

		agentStep.Duration = time.Since(agentStep.Timestamp)
		output.AddStep(agentStep)

		// Check if this was the last step
		if step == maxSteps {
			output.SetError("MAX_STEPS", "Maximum steps reached", fmt.Sprintf("Reached maximum of %d steps", maxSteps))
			output.Output = "I apologize, but I've reached the maximum number of reasoning steps. Please try rephrasing your question or breaking it into smaller parts."
		}
	}

	// Finalize output
	output.Finalize()

	// Add to memory if enabled
	if a.config.EnableMemory {
		_ = a.AddToMemory(ctx, types.NewUserMessage(input.Query))
		_ = a.AddToMemory(ctx, types.NewAssistantMessage(output.Output))
	}

	return output, nil
}

// Stream executes the agent and streams events
func (a *ReActAgent) Stream(ctx context.Context, input types.AgentInput) (<-chan types.AgentEvent, error) {
	events := make(chan types.AgentEvent, 10)

	go func() {
		defer close(events)

		output, err := a.Run(ctx, input)
		if err != nil {
			events <- types.AgentEvent{
				Type:      types.EventTypeError,
				Data:      err.Error(),
				Timestamp: time.Now(),
			}
			return
		}

		// Send all steps as events
		for _, step := range output.Steps {
			events <- types.AgentEvent{
				Type:       types.EventTypeThought,
				Data:       step.Thought,
				Timestamp:  step.Timestamp,
				StepNumber: step.StepNumber,
			}

			if step.Action != "" {
				events <- types.AgentEvent{
					Type:       types.EventTypeAction,
					Data:       step.Action,
					Timestamp:  step.Timestamp,
					StepNumber: step.StepNumber,
				}
			}

			if step.Observation != "" {
				events <- types.AgentEvent{
					Type:       types.EventTypeObservation,
					Data:       step.Observation,
					Timestamp:  step.Timestamp,
					StepNumber: step.StepNumber,
				}
			}
		}

		// Send final output
		events <- types.AgentEvent{
			Type:      types.EventTypeOutput,
			Data:      output.Output,
			Timestamp: time.Now(),
		}
	}()

	return events, nil
}

// buildSystemPrompt builds the system prompt for the ReAct agent
func (a *ReActAgent) buildSystemPrompt() string {
	if a.config.CustomPrompt != "" {
		return a.config.CustomPrompt
	}

	toolDescriptions := ""
	for _, tool := range a.tools {
		toolDescriptions += fmt.Sprintf("- %s: %s\n", tool.Name(), tool.Description())
	}

	prompt := fmt.Sprintf(`You are a helpful AI assistant that uses the ReAct (Reasoning + Acting) framework to solve problems.

You have access to the following tools:
%s

To use a tool, you must follow this exact format:

Thought: [Your reasoning about what to do next]
Action: [The tool name to use]
Action Input: [JSON object with the tool parameters]

After using a tool, you will receive an observation. You can then continue reasoning and use more tools, or provide a final answer.

When you have enough information to answer the question, use this format:

Thought: [Your final reasoning]
Final Answer: [Your complete answer to the user's question]

Important guidelines:
1. Always think step-by-step
2. Use tools when you need information
3. Provide clear, concise final answers
4. If you're unsure, ask for clarification
5. For financial queries, be precise and accurate

Begin!`, toolDescriptions)

	return prompt
}

// parseResponse parses the agent's response to extract thought, action, and final answer
func (a *ReActAgent) parseResponse(content string) (thought, action string, actionInput map[string]interface{}, finalAnswer string, err error) {
	// Extract thought
	thoughtRegex := regexp.MustCompile(`(?i)Thought:\s*(.+?)(?:\n|$)`)
	if matches := thoughtRegex.FindStringSubmatch(content); len(matches) > 1 {
		thought = strings.TrimSpace(matches[1])
	}

	// Check for final answer
	finalAnswerRegex := regexp.MustCompile(`(?i)Final Answer:\s*(.+)`)
	if matches := finalAnswerRegex.FindStringSubmatch(content); len(matches) > 1 {
		finalAnswer = strings.TrimSpace(matches[1])
		return thought, "", nil, finalAnswer, nil
	}

	// Extract action
	actionRegex := regexp.MustCompile(`(?i)Action:\s*(.+?)(?:\n|$)`)
	if matches := actionRegex.FindStringSubmatch(content); len(matches) > 1 {
		action = strings.TrimSpace(matches[1])
	}

	// Extract action input
	actionInputRegex := regexp.MustCompile(`(?i)Action Input:\s*(\{.+?\})`)
	if matches := actionInputRegex.FindStringSubmatch(content); len(matches) > 1 {
		inputStr := strings.TrimSpace(matches[1])
		if err := json.Unmarshal([]byte(inputStr), &actionInput); err != nil {
			return thought, action, nil, "", fmt.Errorf("failed to parse action input: %w", err)
		}
	}

	return thought, action, actionInput, "", nil
}

