# Tutorial 1: LLM Basics - Build Your First AI Chatbot

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟢 BEGINNER                                    ⏱️  30 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: CLI Chatbot with OpenAI                                    │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ OpenAI API integration in Go                                     │
│     ✓ Streaming vs non-streaming responses                             │
│     ✓ Token management and counting                                    │
│     ✓ Temperature and model parameters                                 │
│     ✓ Conversation history management                                  │
│                                                                          │
│  🛠️ TECH STACK: OpenAI API, Go standard library                         │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 🎓 Understanding LLMs

### What is a Large Language Model?

```
┌─────────────────────────────────────────────────────────────────┐
│  LLM: A neural network trained on massive amounts of text       │
│                                                                  │
│  Input (Prompt):                                                │
│    "Explain quantum computing in simple terms"                  │
│                                                                  │
│  LLM Processing:                                                │
│    • Tokenizes input                                            │
│    • Predicts next tokens based on patterns                     │
│    • Generates coherent response                                │
│                                                                  │
│  Output (Completion):                                           │
│    "Quantum computing uses quantum mechanics principles..."     │
└─────────────────────────────────────────────────────────────────┘
```

### Key Concepts

**1. Tokens**
```
Text: "Hello, world!"
Tokens: ["Hello", ",", " world", "!"]  ← ~4 tokens

Rule of thumb: 1 token ≈ 4 characters (English)
```

**2. Context Window**
```
GPT-3.5-turbo: 4,096 tokens
GPT-4: 8,192 tokens
GPT-4-turbo: 128,000 tokens

Context = Input + Output combined
```

**3. Temperature**
```
Temperature 0.0: Deterministic, focused
Temperature 0.7: Balanced creativity
Temperature 1.0+: More random, creative
```

---

## 🚀 Project: Build a CLI Chatbot

### Step 1: Setup

```bash
# Create project directory
mkdir -p basic/projects/ai-engineering/chatbot-cli
cd basic/projects/ai-engineering/chatbot-cli

# Initialize Go module
go mod init chatbot-cli

# Install OpenAI SDK
go get github.com/sashabaranov/go-openai
```

### Step 2: Get Your API Key

```
┌──────────────────────────────────────────────────────────────────┐
│  1. Visit https://platform.openai.com/api-keys                   │
│  2. Sign up or log in                                            │
│  3. Click "Create new secret key"                                │
│  4. Copy the key (starts with sk-...)                            │
│  5. Set environment variable:                                    │
│     export OPENAI_API_KEY="sk-..."                               │
└──────────────────────────────────────────────────────────────────┘
```

### Step 3: Create the Chatbot

Create `main.go`:

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: OPENAI_API_KEY not set")
		fmt.Println("Set it with: export OPENAI_API_KEY='your-key'")
		os.Exit(1)
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Store conversation history
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful AI assistant.",
		},
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                           ║")
	fmt.Println("║              🤖 AI Chatbot CLI                            ║")
	fmt.Println("║                                                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Type your message and press Enter. Type 'exit' to quit.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Get user input
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}

		if userInput == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		// Add user message to history
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		// Create chat completion request
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				Messages:    messages,
				Temperature: 0.7,
				MaxTokens:   500,
			},
		)

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		// Get assistant's response
		assistantMsg := resp.Choices[0].Message.Content

		// Add to conversation history
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantMsg,
		})

		// Display response
		fmt.Printf("\nAI: %s\n\n", assistantMsg)

		// Show token usage
		fmt.Printf("📊 Tokens used: %d (prompt: %d, completion: %d)\n\n",
			resp.Usage.TotalTokens,
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens,
		)
	}
}
```

### Step 4: Run Your Chatbot

```bash
# Make sure API key is set
export OPENAI_API_KEY="sk-..."

# Run the chatbot
go run main.go
```

**📤 Expected Output:**
```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║              🤖 AI Chatbot CLI                            ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

Type your message and press Enter. Type 'exit' to quit.

You: Hello! What can you help me with?

AI: Hello! I'm an AI assistant and I can help you with many things:
- Answer questions on various topics
- Help with coding and technical problems
- Provide explanations and tutorials
- Assist with writing and editing
- And much more!

What would you like help with today?

📊 Tokens used: 87 (prompt: 25, completion: 62)

You: Explain Go channels in simple terms

AI: Go channels are like pipes that allow goroutines (lightweight threads)
to communicate with each other safely. Think of it as a mailbox:
- One goroutine can send a message into the channel
- Another goroutine can receive that message from the channel
- This prevents race conditions and makes concurrent programming safer

Example:
ch := make(chan int)  // Create a channel
go func() { ch <- 42 }()  // Send value
value := <-ch  // Receive value

📊 Tokens used: 156 (prompt: 89, completion: 67)
```

---

## 🔧 Understanding the Code

### 1. Client Initialization

```go
client := openai.NewClient(apiKey)
```

Creates a client to communicate with OpenAI's API.

### 2. Message Structure

```go
type ChatCompletionMessage struct {
    Role    string  // "system", "user", or "assistant"
    Content string  // The actual message
}
```

**Roles:**
- `system`: Sets behavior/personality
- `user`: Your input
- `assistant`: AI's responses

### 3. Conversation History

```go
messages := []openai.ChatCompletionMessage{
    {Role: "system", Content: "You are helpful"},
    {Role: "user", Content: "Hello"},
    {Role: "assistant", Content: "Hi there!"},
    {Role: "user", Content: "How are you?"},
}
```

The LLM sees the entire conversation for context.

### 4. API Request

```go
resp, err := client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest{
        Model:       openai.GPT3Dot5Turbo,
        Messages:    messages,
        Temperature: 0.7,
        MaxTokens:   500,
    },
)
```

**Parameters:**
- `Model`: Which LLM to use
- `Messages`: Conversation history
- `Temperature`: Creativity (0-2)
- `MaxTokens`: Max response length

---

## ⚡ Adding Streaming Responses

Streaming shows responses as they're generated (like ChatGPT).

Create `streaming.go`:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func streamingChat(client *openai.Client, messages []openai.ChatCompletionMessage) (string, error) {
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   500,
		Stream:      true, // Enable streaming
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	fmt.Print("AI: ")

	var fullResponse string

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println() // New line after response
			break
		}

		if err != nil {
			return "", err
		}

		// Get the delta (new text chunk)
		delta := response.Choices[0].Delta.Content

		// Print it immediately
		fmt.Print(delta)

		// Accumulate full response
		fullResponse += delta
	}

	return fullResponse, nil
}
```

**Benefits of Streaming:**
- ✅ Better user experience (see response immediately)
- ✅ Feels more interactive
- ✅ Can cancel long responses early

---

## 📊 Token Management

### Counting Tokens

```go
func estimateTokens(text string) int {
	// Rough estimate: 1 token ≈ 4 characters
	return len(text) / 4
}

func checkTokenLimit(messages []openai.ChatCompletionMessage, limit int) bool {
	total := 0
	for _, msg := range messages {
		total += estimateTokens(msg.Content)
	}
	return total < limit
}
```

### Managing Context Window

```go
func trimConversation(messages []openai.ChatCompletionMessage, maxTokens int) []openai.ChatCompletionMessage {
	// Always keep system message
	if len(messages) <= 1 {
		return messages
	}

	// Keep system message + recent messages
	systemMsg := messages[0]
	recentMsgs := messages[1:]

	// Estimate and trim if needed
	for estimateTokens(formatMessages(recentMsgs)) > maxTokens {
		// Remove oldest user-assistant pair
		if len(recentMsgs) >= 2 {
			recentMsgs = recentMsgs[2:]
		} else {
			break
		}
	}

	return append([]openai.ChatCompletionMessage{systemMsg}, recentMsgs...)
}
```

---

## 🎛️ Model Parameters

### Temperature

```go
// Deterministic (good for factual answers)
Temperature: 0.0

// Balanced (default)
Temperature: 0.7

// Creative (good for brainstorming)
Temperature: 1.5
```

### Max Tokens

```go
// Short responses
MaxTokens: 100

// Medium responses
MaxTokens: 500

// Long responses
MaxTokens: 2000
```

### Top P (Nucleus Sampling)

```go
// More focused
TopP: 0.1

// Balanced
TopP: 0.9

// More diverse
TopP: 1.0
```

---

## 🎯 Challenges

### Challenge 1: Add Commands
Add special commands to your chatbot:
- `/clear` - Clear conversation history
- `/tokens` - Show total tokens used
- `/model <name>` - Switch models

### Challenge 2: Save Conversations
Save conversation history to a file:
```go
func saveConversation(messages []openai.ChatCompletionMessage, filename string) error {
	// Implement JSON serialization
}
```

### Challenge 3: Add Streaming
Modify your chatbot to use streaming responses.

### Challenge 4: Multi-turn Context
Implement a sliding window to keep only the last N messages.

---

## ✅ What You Learned

```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ OpenAI API integration in Go                                  │
│  ✓ Chat completion requests                                      │
│  ✓ Conversation history management                               │
│  ✓ Token counting and limits                                     │
│  ✓ Streaming responses                                           │
│  ✓ Model parameters (temperature, max_tokens)                    │
│  ✓ Error handling for API calls                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Next Steps

**Immediate:**
1. Build and run the chatbot
2. Experiment with different temperatures
3. Try different models (GPT-4, GPT-3.5-turbo-16k)
4. Implement the challenges

**Next Tutorial:**
[Tutorial 2: Prompt Engineering →](02_PROMPT_ENGINEERING.md)

Learn how to design effective prompts to get better AI responses!

---

## 📚 Additional Resources

- [OpenAI API Documentation](https://platform.openai.com/docs)
- [go-openai SDK](https://github.com/sashabaranov/go-openai)
- [Token Counting](https://platform.openai.com/tokenizer)
- [Model Pricing](https://openai.com/pricing)

---

**💡 Pro Tip**: Start with GPT-3.5-turbo for development (cheaper, faster), then upgrade to GPT-4 for production if needed.

