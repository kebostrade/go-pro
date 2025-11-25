package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// CodingExpertAgent is specialized for answering programming questions
type CodingExpertAgent struct {
	*BaseAgent
	systemPrompt string
}

// CodingExpertConfig holds configuration for the coding expert agent
type CodingExpertConfig struct {
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
	SupportedLangs  []string
}

// NewCodingExpertAgent creates a new coding expert agent
func NewCodingExpertAgent(config CodingExpertConfig) *CodingExpertAgent {
	if config.Name == "" {
		config.Name = "CodingExpert"
	}
	if config.Description == "" {
		config.Description = "Expert programming assistant for code questions, debugging, and best practices"
	}
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
		config.Timeout = 60 * time.Second
	}

	baseConfig := types.AgentConfig{
		Name:            config.Name,
		Description:     config.Description,
		MaxSteps:        config.MaxSteps,
		Temperature:     config.Temperature,
		MaxTokens:       config.MaxTokens,
		Timeout:         config.Timeout,
		VerboseLogging:  config.VerboseLogging,
	}

	systemPrompt := buildCodingExpertPrompt(config.SupportedLangs)

	agent := &CodingExpertAgent{
		BaseAgent:    NewBaseAgent(baseConfig, config.LLM),
		systemPrompt: systemPrompt,
	}

	// Add tools and memory
	if len(config.Tools) > 0 {
		agent.AddTools(config.Tools)
	}
	if config.Memory != nil {
		agent.SetMemory(config.Memory)
	}

	return agent
}

// Run executes the coding expert agent
func (a *CodingExpertAgent) Run(ctx context.Context, input types.AgentInput) (*types.AgentOutput, error) {
	output := types.NewAgentOutput()
	output.Metadata.AgentType = "CodingExpert"

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, a.config.Timeout)
	defer cancel()

	// Detect programming language from query
	language := detectLanguageFromQuery(input.Query)

	// Enhance query with context
	enhancedQuery := a.enhanceQuery(input.Query, language, input.Context)

	// Build the prompt
	prompt := a.buildPrompt(enhancedQuery, language)

	// Execute reasoning loop
	for step := 1; step <= a.config.MaxSteps; step++ {
		stepStart := time.Now()

		// Check context cancellation
		select {
		case <-ctx.Done():
			output.SetError("TIMEOUT", "Agent execution timeout", ctx.Err().Error())
			output.Finalize()
			return output, errors.NewAgentTimeoutError("Agent execution timeout")
		default:
		}

		// Get LLM response
		llmRequest := types.LLMRequest{
			Messages: []types.Message{
				{
					Role:    types.RoleSystem,
					Content: a.systemPrompt,
				},
				{
					Role:    types.RoleUser,
					Content: prompt,
				},
			},
			Temperature: a.config.Temperature,
			MaxTokens:   a.config.MaxTokens,
		}

		llmOutput, err := a.llm.Generate(ctx, llmRequest)
		if err != nil {
			output.SetError("LLM_ERROR", "Failed to get LLM response", err.Error())
			output.Finalize()
			return output, err
		}

		// Track token usage
		output.AddTokenUsage(llmOutput.Usage.PromptTokens, llmOutput.Usage.CompletionTokens)

		// Parse response
		response := llmOutput.Content

		// Check if we have a final answer
		if strings.Contains(strings.ToLower(response), "final answer:") {
			// Extract final answer
			parts := strings.Split(response, "Final Answer:")
			if len(parts) > 1 {
				output.Output = strings.TrimSpace(parts[1])
			} else {
				output.Output = response
			}

			// Add final step
			output.AddStep(types.AgentStep{
				StepNumber:  step,
				Thought:     "Providing final answer",
				Observation: output.Output,
				Timestamp:   time.Now(),
				Duration:    time.Since(stepStart),
			})

			break
		}

		// Parse thought and action
		thought, action, actionInput := a.parseResponse(response)

		// Add step
		agentStep := types.AgentStep{
			StepNumber:  step,
			Thought:     thought,
			Action:      action,
			ActionInput: actionInput,
			Timestamp:   time.Now(),
		}

		// Execute action if present
		if action != "" {
			observation, err := a.executeAction(ctx, action, actionInput)
			if err != nil {
				agentStep.Observation = fmt.Sprintf("Error: %v", err)
			} else {
				agentStep.Observation = observation
			}

			// Update prompt with observation
			prompt = fmt.Sprintf("%s\n\nObservation: %s\n\nThought:", prompt, agentStep.Observation)
		}

		agentStep.Duration = time.Since(stepStart)
		output.AddStep(agentStep)

		// If no action, we might have the answer
		if action == "" && thought != "" {
			output.Output = thought
			break
		}
	}

	// If no output yet, use last thought
	if output.Output == "" && len(output.Steps) > 0 {
		lastStep := output.Steps[len(output.Steps)-1]
		output.Output = lastStep.Thought
	}

	output.Finalize()
	return output, nil
}

// enhanceQuery adds context to the query
func (a *CodingExpertAgent) enhanceQuery(query, language string, context map[string]interface{}) string {
	enhanced := query

	if language != "" {
		enhanced = fmt.Sprintf("[Language: %s] %s", language, enhanced)
	}

	if len(context) > 0 {
		contextJSON, _ := json.Marshal(context)
		enhanced = fmt.Sprintf("%s\n\nContext: %s", enhanced, string(contextJSON))
	}

	return enhanced
}

// buildPrompt constructs the agent prompt
func (a *CodingExpertAgent) buildPrompt(query, language string) string {
	return fmt.Sprintf(`You are a coding expert assistant. Answer the following programming question:

%s

Think step by step:
1. Understand the question
2. Identify what information or tools you need
3. Use available tools if necessary
4. Provide a clear, accurate answer with code examples

Available tools: %s

Thought:`, query, a.getToolNames())
}

// parseResponse parses the LLM response
func (a *CodingExpertAgent) parseResponse(response string) (thought, action string, actionInput map[string]interface{}) {
	actionInput = make(map[string]interface{})

	// Simple parsing - in production, use more robust parsing
	lines := strings.Split(response, "\n")
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if strings.HasPrefix(trimmed, "Thought:") {
			thought = strings.TrimSpace(strings.TrimPrefix(trimmed, "Thought:"))
		} else if strings.HasPrefix(trimmed, "Action:") {
			action = strings.TrimSpace(strings.TrimPrefix(trimmed, "Action:"))
		} else if strings.HasPrefix(trimmed, "Action Input:") {
			// Try to parse JSON
			if i+1 < len(lines) {
				inputStr := strings.TrimSpace(strings.TrimPrefix(trimmed, "Action Input:"))
				json.Unmarshal([]byte(inputStr), &actionInput)
			}
		}
	}

	return thought, action, actionInput
}

// executeAction executes a tool action
func (a *CodingExpertAgent) executeAction(ctx context.Context, action string, input map[string]interface{}) (string, error) {
	// Find the tool
	var tool types.Tool
	for _, t := range a.tools {
		if t.Name() == action {
			tool = t
			break
		}
	}

	if tool == nil {
		return "", fmt.Errorf("tool not found: %s", action)
	}

	// Execute tool
	toolInput := types.NewToolInput(input)
	toolOutput, err := tool.Execute(ctx, toolInput)
	if err != nil {
		return "", err
	}

	if toolOutput.Error != nil {
		return "", fmt.Errorf("%s: %s", toolOutput.Error.Code, toolOutput.Error.Message)
	}

	// Convert result to string
	if str, ok := toolOutput.Result.(string); ok {
		return str, nil
	}

	resultJSON, _ := json.Marshal(toolOutput.Result)
	return string(resultJSON), nil
}

// getToolNames returns a comma-separated list of tool names
func (a *CodingExpertAgent) getToolNames() string {
	names := make([]string, len(a.tools))
	for i, tool := range a.tools {
		names[i] = tool.Name()
	}
	return strings.Join(names, ", ")
}

// buildCodingExpertPrompt creates the system prompt
func buildCodingExpertPrompt(supportedLangs []string) string {
	langs := "Go, Python, JavaScript, TypeScript, Rust, Java, C++, C"
	if len(supportedLangs) > 0 {
		langs = strings.Join(supportedLangs, ", ")
	}

	return fmt.Sprintf(`You are an expert programming assistant with deep knowledge of multiple programming languages including %s.

Your capabilities:
- Answer programming questions with accurate, well-explained solutions
- Debug code and identify issues
- Provide code examples and best practices
- Explain complex programming concepts clearly
- Suggest optimizations and improvements
- Reference official documentation and reliable sources

When answering:
1. Be precise and accurate
2. Provide working code examples
3. Explain your reasoning
4. Follow language-specific best practices
5. Consider performance and security
6. Reference documentation when helpful

Use available tools to:
- Analyze code for issues
- Execute code to verify solutions
- Search documentation
- Find relevant examples

Always provide clear, actionable answers with code examples when appropriate.`, langs)
}

// detectLanguageFromQuery attempts to detect the programming language
func detectLanguageFromQuery(query string) string {
	queryLower := strings.ToLower(query)

	// Check for explicit language mentions
	languages := map[string][]string{
		"go":         {"golang", "go ", " go"},
		"python":     {"python", "py "},
		"javascript": {"javascript", "js ", "node"},
		"typescript": {"typescript", "ts "},
		"rust":       {"rust"},
		"java":       {"java"},
		"cpp":        {"c++", "cpp"},
		"c":          {" c ", "c programming"},
	}

	for lang, keywords := range languages {
		for _, keyword := range keywords {
			if strings.Contains(queryLower, keyword) {
				return lang
			}
		}
	}

	return ""
}

