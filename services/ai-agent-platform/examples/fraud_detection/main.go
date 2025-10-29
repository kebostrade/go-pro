package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"ai-agent-platform/internal/agent"
	"ai-agent-platform/internal/llm"
	"ai-agent-platform/internal/tools/financial"
	"ai-agent-platform/internal/tools/general"
	"ai-agent-platform/pkg/types"
)

func main() {
	fmt.Println("🤖 FinAgent - Fraud Detection Example")
	fmt.Println("=====================================")

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create LLM provider
	llmProvider, err := llm.NewOpenAIProvider(llm.OpenAIConfig{
		APIKey:  apiKey,
		Model:   "gpt-4",
		Timeout: 60 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create LLM provider: %v", err)
	}

	fmt.Println("✅ LLM Provider initialized (OpenAI GPT-4)")

	// Create tools
	tools := []types.Tool{
		financial.NewTransactionLookupTool(),
		financial.NewFraudCheckTool(),
		general.NewCalculatorTool(),
	}

	fmt.Printf("✅ Loaded %d tools\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("   - %s: %s\n", tool.Name(), tool.Description())
	}

	// Create ReAct agent
	reactAgent := agent.NewReActAgent(agent.ReActConfig{
		Name:           "FraudDetectionAgent",
		Description:    "An AI agent specialized in detecting fraudulent transactions",
		LLM:            llmProvider,
		Tools:          tools,
		MaxSteps:       5,
		Temperature:    0.7,
		VerboseLogging: true,
	})

	fmt.Println("\n✅ ReAct Agent created")
	fmt.Println("\n" + strings.Repeat("=", 60))

	// Example queries
	queries := []string{
		"Look up transaction TXN_12345 and check if it's fraudulent",
		"What is the risk score for transaction TXN_67890?",
		"Analyze transaction TXN_ABC123 for fraud and tell me if I should approve it",
	}

	// Run each query
	for i, query := range queries {
		fmt.Printf("\n\n📝 Query %d: %s\n", i+1, query)
		fmt.Println(strings.Repeat("-", 60))

		ctx := context.Background()
		input := types.NewAgentInput(query)

		// Run agent
		startTime := time.Now()
		output, err := reactAgent.Run(ctx, input)
		duration := time.Since(startTime)

		if err != nil {
			log.Printf("❌ Error: %v", err)
			continue
		}

		// Display results
		fmt.Println("\n🔍 Agent Execution:")
		fmt.Printf("   Duration: %v\n", duration)
		fmt.Printf("   Steps: %d\n", len(output.Steps))
		fmt.Printf("   Tokens Used: %d (Prompt: %d, Completion: %d)\n",
			output.Metadata.TokensUsed.TotalTokens,
			output.Metadata.TokensUsed.PromptTokens,
			output.Metadata.TokensUsed.CompletionTokens)

		// Display steps
		fmt.Println("\n📊 Reasoning Steps:")
		for _, step := range output.Steps {
			fmt.Printf("\n   Step %d:\n", step.StepNumber)
			fmt.Printf("   💭 Thought: %s\n", step.Thought)
			if step.Action != "" {
				fmt.Printf("   🔧 Action: %s\n", step.Action)
				if len(step.ActionInput) > 0 {
					inputJSON, _ := json.MarshalIndent(step.ActionInput, "      ", "  ")
					fmt.Printf("   📥 Input: %s\n", string(inputJSON))
				}
			}
			if step.Observation != "" {
				fmt.Printf("   👁️  Observation: %s\n", truncate(step.Observation, 200))
			}
		}

		// Display tool calls
		if len(output.ToolCalls) > 0 {
			fmt.Println("\n🛠️  Tool Calls:")
			for _, toolCall := range output.ToolCalls {
				fmt.Printf("   - %s (ID: %s)\n", toolCall.ToolName, toolCall.ID[:8])
				if toolCall.Output != nil {
					resultJSON, _ := json.MarshalIndent(toolCall.Output.Result, "     ", "  ")
					fmt.Printf("     Result: %s\n", truncate(string(resultJSON), 300))
				}
			}
		}

		// Display final answer
		fmt.Println("\n✨ Final Answer:")
		fmt.Printf("   %s\n", output.Output)

		// Display metadata
		if output.Error != nil {
			fmt.Printf("\n⚠️  Error: %s - %s\n", output.Error.Code, output.Error.Message)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("\n✅ Fraud Detection Example Complete!")
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
