package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/langchain/pkg/llm"
	"github.com/DimaJoyti/go-pro/services/langchain/pkg/schema"
	"github.com/google/uuid"
)

const (
	MaxIterations = 10
	ReActPrompt   = `You are a helpful AI assistant that can use tools to answer questions.

Available tools: %s

Instructions:
- For each step, provide your reasoning (Thought)
- If you need to use a tool, specify the action (Action) and input (Action Input)
- After using a tool, note the result (Observation)
- When you have the final answer, provide it (Answer)

`
)

type ReActAgent struct {
	name        string
	description string
	llm         *llm.LLMProvider
	tools       []schema.Tool
	maxSteps    int
	verbose     bool
}

type ReActAgentOption func(*ReActAgent)

func WithMaxSteps(maxSteps int) ReActAgentOption {
	return func(a *ReActAgent) {
		a.maxSteps = maxSteps
	}
}

func WithVerbose(verbose bool) ReActAgentOption {
	return func(a *ReActAgent) {
		a.verbose = verbose
	}
}

func NewReActAgent(name string, llmProvider *llm.LLMProvider, tools []schema.Tool, opts ...ReActAgentOption) *ReActAgent {
	agent := &ReActAgent{
		name:        name,
		description: "ReAct (Reasoning + Acting) agent",
		llm:         llmProvider,
		tools:       tools,
		maxSteps:    5,
		verbose:     false,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func (a *ReActAgent) Run(ctx context.Context, input schema.AgentInput) (*schema.AgentOutput, error) {
	output := &schema.AgentOutput{
		SessionID:   input.SessionID,
		ExecutionID: uuid.New().String(),
		StartTime:   time.Now(),
	}

	if input.MaxSteps > 0 {
		a.maxSteps = input.MaxSteps
	}

	messages := append([]schema.Message{
		schema.NewSystemMessage(fmt.Sprintf(ReActPrompt, a.getToolDescriptions())),
	}, input.Messages...)

	if input.Query != "" {
		messages = append(messages, schema.NewUserMessage(input.Query))
	}

	var lastMessage schema.Message
	var err error

	for i := 0; i < a.maxSteps; i++ {
		if len(a.tools) > 0 {
			lastMessage, err = a.llm.GenerateWithTools(ctx, messages, a.tools)
		} else {
			lastMessage, err = a.llm.Generate(ctx, messages)
		}

		if err != nil {
			output.Error = fmt.Errorf("LLM error at step %d: %w", i, err)
			break
		}

		output.AddMessage(lastMessage)

		if lastMessage.Content != "" && !strings.Contains(lastMessage.Content, "Action:") {
			output.Output = lastMessage.Content
			break
		}

		if len(lastMessage.ToolCalls) > 0 {
			for _, tc := range lastMessage.ToolCalls {
				output.AddToolCall(tc)

				tool := a.findTool(tc.Name)
				if tool != nil {
					result, err := tool.Invoke(ctx, tc.Arguments)
					var resultStr string
					if err != nil {
						resultStr = fmt.Sprintf("Error: %v", err)
					} else {
						resultStr = fmt.Sprintf("%v", result)
					}

					toolResult := schema.ToolResult{
						ToolCallID: tc.ID,
						ToolName:   tc.Name,
						Output:     result,
					}
					if err != nil {
						toolResult.Error = err.Error()
					}
					output.AddToolResult(toolResult)

					messages = append(messages, schema.NewToolMessage(resultStr, tc.ID))

					if a.verbose {
						fmt.Printf("Step %d: Tool %s result: %s\n", i+1, tc.Name, resultStr)
					}
				}
			}
		}

		messages = append(messages, lastMessage)
	}

	if output.Output == "" && len(output.Messages) > 0 {
		output.Output = output.Messages[len(output.Messages)-1].Content
	}

	output.Finalize()
	return output, nil
}

func (a *ReActAgent) getToolDescriptions() string {
	if len(a.tools) == 0 {
		return "No tools available"
	}

	var descs []string
	for _, tool := range a.tools {
		descs = append(descs, fmt.Sprintf("- %s: %s", tool.Name, tool.Description))
	}
	return strings.Join(descs, "\n")
}

func (a *ReActAgent) findTool(name string) *schema.Tool {
	for i := range a.tools {
		if a.tools[i].Name == name {
			return &a.tools[i]
		}
	}
	return nil
}
