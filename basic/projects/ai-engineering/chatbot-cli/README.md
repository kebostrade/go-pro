# 🤖 AI Chatbot CLI

A simple command-line chatbot powered by OpenAI's GPT models, built with Go.

## ✨ Features

- 💬 Interactive chat with AI
- 🌊 Streaming responses (real-time output)
- 📝 Conversation history management
- 🎛️ Multiple model support (GPT-3.5, GPT-4)
- 📊 Token usage tracking
- 🎨 Colorful terminal output
- ⚡ Fast and lightweight

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- OpenAI API key ([Get one here](https://platform.openai.com/api-keys))

### Installation

```bash
# Clone or navigate to the project
cd basic/projects/ai-engineering/chatbot-cli

# Install dependencies
go mod init chatbot-cli
go get github.com/sashabaranov/go-openai

# Set your API key
export OPENAI_API_KEY="sk-..."

# Run the chatbot
go run main.go
```

## 📖 Usage

### Basic Chat

```
You: Hello! What can you help me with?

AI: Hello! I'm an AI assistant and I can help you with many things:
- Answer questions on various topics
- Help with coding and technical problems
- Provide explanations and tutorials
- Assist with writing and editing
- And much more!

What would you like help with today?
```

### Commands

| Command | Description |
|---------|-------------|
| `/help` | Show available commands |
| `/clear` | Clear conversation history |
| `/stats` | Show conversation statistics |
| `/stream` | Toggle streaming mode on/off |
| `/model` | Change AI model (GPT-3.5, GPT-4) |
| `/exit` | Exit the chatbot |

### Example Session

```bash
$ go run main.go

╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║              🤖 AI Chatbot CLI                            ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

Type your message and press Enter.
Type /help for commands or /exit to quit.

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

You: /stats

=== Conversation Stats ===
Messages: 2
Model: gpt-3.5-turbo
Streaming: true

You: /exit
👋 Goodbye!
```

## 🎛️ Configuration

### Change Model

```
You: /model

Available models:
  1. gpt-3.5-turbo (fast, cheap)
  2. gpt-4 (smart, expensive)
  3. gpt-4-turbo-preview (balanced)

Enter number: 2
✓ Model changed to: gpt-4
```

### Toggle Streaming

```
You: /stream
✓ Streaming disabled
```

## 🏗️ Project Structure

```
chatbot-cli/
├── main.go           # Main application code
├── go.mod            # Go module file
└── README.md         # This file
```

## 🔧 Code Overview

### Main Components

**Chatbot Struct**
```go
type Chatbot struct {
    client   *openai.Client
    messages []openai.ChatCompletionMessage
    model    string
    stream   bool
}
```

**Key Methods**
- `Chat()` - Send message and get response
- `chatStreaming()` - Handle streaming responses
- `chatNonStreaming()` - Handle non-streaming responses
- `ClearHistory()` - Reset conversation
- `SetModel()` - Change AI model
- `ToggleStreaming()` - Enable/disable streaming

## 📊 Token Usage

The chatbot tracks token usage for each request:

```
📊 Tokens: 156 (prompt: 89, completion: 67)
```

**Understanding Tokens:**
- Prompt tokens: Your input + conversation history
- Completion tokens: AI's response
- Total tokens: Sum of both (affects cost)

**Approximate Costs (as of 2024):**
- GPT-3.5-turbo: $0.0015 per 1K tokens
- GPT-4: $0.03 per 1K tokens
- GPT-4-turbo: $0.01 per 1K tokens

## 🎨 Customization

### Change System Prompt

Edit the system message in `NewChatbot()`:

```go
messages: []openai.ChatCompletionMessage{
    {
        Role:    openai.ChatMessageRoleSystem,
        Content: "You are a helpful coding assistant specializing in Go.",
    },
},
```

### Adjust Temperature

Modify the `Temperature` parameter in `Chat()`:

```go
Temperature: 0.7,  // 0.0 = focused, 2.0 = creative
```

### Change Max Tokens

Adjust `MaxTokens` to control response length:

```go
MaxTokens: 500,  // Increase for longer responses
```

## 🐛 Troubleshooting

### "Invalid API key" Error

```bash
# Make sure your API key is set correctly
echo $OPENAI_API_KEY

# If empty, set it:
export OPENAI_API_KEY="sk-..."
```

### "Rate limit exceeded" Error

You're making too many requests. Wait a moment and try again, or upgrade your OpenAI plan.

### "Context length exceeded" Error

Your conversation is too long. Use `/clear` to reset the history.

## 🚀 Next Steps

### Enhancements to Try

1. **Save Conversations**
   - Add ability to save chat history to file
   - Load previous conversations

2. **Custom Commands**
   - Add `/save <filename>` command
   - Add `/load <filename>` command

3. **Better Error Handling**
   - Retry failed requests
   - Handle network errors gracefully

4. **Advanced Features**
   - Add function calling support
   - Implement RAG for document Q&A
   - Add multi-turn context management

## 📚 Learn More

- [Tutorial 1: LLM Basics](../../../docs/tutorials/ai-engineering/01_LLM_BASICS.md)
- [Tutorial 2: Prompt Engineering](../../../docs/tutorials/ai-engineering/02_PROMPT_ENGINEERING.md)
- [OpenAI API Documentation](https://platform.openai.com/docs)
- [go-openai SDK](https://github.com/sashabaranov/go-openai)

## 📝 License

This project is part of the Go Learn AI Engineering tutorial series.

---

**💡 Pro Tip**: Start with GPT-3.5-turbo for development (it's faster and cheaper), then switch to GPT-4 when you need more advanced reasoning!

