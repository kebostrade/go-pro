package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorRed    = "\033[31m"
)

type Chatbot struct {
	client   *openai.Client
	messages []openai.ChatCompletionMessage
	model    string
	stream   bool
}

func NewChatbot(apiKey string) *Chatbot {
	return &Chatbot{
		client: openai.NewClient(apiKey),
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful AI assistant. Be concise and clear in your responses.",
			},
		},
		model:  openai.GPT3Dot5Turbo,
		stream: true,
	}
}

func (c *Chatbot) Chat(ctx context.Context, userInput string) error {
	// Add user message to history
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	if c.stream {
		return c.chatStreaming(ctx)
	}
	return c.chatNonStreaming(ctx)
}

func (c *Chatbot) chatStreaming(ctx context.Context) error {
	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Messages:    c.messages,
		Temperature: 0.7,
		MaxTokens:   500,
		Stream:      true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	fmt.Print(colorBlue + "AI: " + colorReset)

	var fullResponse string

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			break
		}

		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}

		delta := response.Choices[0].Delta.Content
		fmt.Print(delta)
		fullResponse += delta
	}

	// Add assistant's response to history
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: fullResponse,
	})

	return nil
}

func (c *Chatbot) chatNonStreaming(ctx context.Context) error {
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.model,
		Messages:    c.messages,
		Temperature: 0.7,
		MaxTokens:   500,
	})

	if err != nil {
		return fmt.Errorf("API error: %w", err)
	}

	assistantMsg := resp.Choices[0].Message.Content

	// Add to conversation history
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: assistantMsg,
	})

	// Display response
	fmt.Printf("%sAI: %s%s\n", colorBlue, assistantMsg, colorReset)

	// Show token usage
	fmt.Printf("\n%s📊 Tokens: %d (prompt: %d, completion: %d)%s\n",
		colorYellow,
		resp.Usage.TotalTokens,
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		colorReset,
	)

	return nil
}

func (c *Chatbot) ClearHistory() {
	// Keep only system message
	c.messages = c.messages[:1]
	fmt.Println(colorGreen + "✓ Conversation history cleared" + colorReset)
}

func (c *Chatbot) SetModel(model string) {
	c.model = model
	fmt.Printf("%s✓ Model changed to: %s%s\n", colorGreen, model, colorReset)
}

func (c *Chatbot) ToggleStreaming() {
	c.stream = !c.stream
	status := "disabled"
	if c.stream {
		status = "enabled"
	}
	fmt.Printf("%s✓ Streaming %s%s\n", colorGreen, status, colorReset)
}

func (c *Chatbot) ShowStats() {
	messageCount := len(c.messages) - 1 // Exclude system message
	fmt.Printf("\n%s=== Conversation Stats ===%s\n", colorYellow, colorReset)
	fmt.Printf("Messages: %d\n", messageCount)
	fmt.Printf("Model: %s\n", c.model)
	fmt.Printf("Streaming: %v\n", c.stream)
	fmt.Println()
}

func printBanner() {
	fmt.Println(colorBlue + "╔═══════════════════════════════════════════════════════════╗" + colorReset)
	fmt.Println(colorBlue + "║                                                           ║" + colorReset)
	fmt.Println(colorBlue + "║              🤖 AI Chatbot CLI                            ║" + colorReset)
	fmt.Println(colorBlue + "║                                                           ║" + colorReset)
	fmt.Println(colorBlue + "╚═══════════════════════════════════════════════════════════╝" + colorReset)
	fmt.Println()
}

func printHelp() {
	fmt.Println(colorYellow + "\n=== Commands ===" + colorReset)
	fmt.Println("  /help      - Show this help message")
	fmt.Println("  /clear     - Clear conversation history")
	fmt.Println("  /stats     - Show conversation statistics")
	fmt.Println("  /stream    - Toggle streaming mode")
	fmt.Println("  /model     - Change AI model")
	fmt.Println("  /exit      - Exit the chatbot")
	fmt.Println()
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println(colorRed + "❌ Error: OPENAI_API_KEY not set" + colorReset)
		fmt.Println("Set it with: export OPENAI_API_KEY='your-key'")
		os.Exit(1)
	}

	// Create chatbot
	bot := NewChatbot(apiKey)

	// Print banner
	printBanner()
	fmt.Println("Type your message and press Enter.")
	fmt.Println("Type /help for commands or /exit to quit.")
	fmt.Println()

	// Create scanner for user input
	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	for {
		// Get user input
		fmt.Print(colorGreen + "You: " + colorReset)
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}

		// Handle commands
		switch userInput {
		case "/exit", "/quit":
			fmt.Println(colorYellow + "👋 Goodbye!" + colorReset)
			return

		case "/help":
			printHelp()
			continue

		case "/clear":
			bot.ClearHistory()
			continue

		case "/stats":
			bot.ShowStats()
			continue

		case "/stream":
			bot.ToggleStreaming()
			continue

		case "/model":
			fmt.Println("\nAvailable models:")
			fmt.Println("  1. gpt-3.5-turbo (fast, cheap)")
			fmt.Println("  2. gpt-4 (smart, expensive)")
			fmt.Println("  3. gpt-4-turbo-preview (balanced)")
			fmt.Print("\nEnter number: ")

			if scanner.Scan() {
				choice := strings.TrimSpace(scanner.Text())
				switch choice {
				case "1":
					bot.SetModel(openai.GPT3Dot5Turbo)
				case "2":
					bot.SetModel(openai.GPT4)
				case "3":
					bot.SetModel(openai.GPT4TurboPreview)
				default:
					fmt.Println(colorRed + "Invalid choice" + colorReset)
				}
			}
			continue
		}

		// Chat with AI
		if err := bot.Chat(ctx, userInput); err != nil {
			fmt.Printf("%s❌ Error: %v%s\n", colorRed, err, colorReset)
		}

		fmt.Println()
	}
}

