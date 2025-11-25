# 🚀 Coding Expert AI Agents - Quick Start Guide

Get up and running with the Coding Expert AI Agent system in 5 minutes!

## 📋 Prerequisites

Before you begin, ensure you have:

1. **Go 1.22 or higher**
   ```bash
   go version
   # Should show: go version go1.22.0 or higher
   ```

2. **OpenAI API Key**
   - Get one from: https://platform.openai.com/api-keys
   - Or use another LLM provider (Anthropic Claude, Ollama, etc.)

3. **Git** (to clone the repository)

## ⚡ 5-Minute Quick Start

### Step 1: Navigate to the Platform

```bash
cd services/ai-agent-platform
```

### Step 2: Set Your API Key

```bash
# On Linux/macOS
export OPENAI_API_KEY="sk-your-api-key-here"

# On Windows (PowerShell)
$env:OPENAI_API_KEY="sk-your-api-key-here"

# On Windows (CMD)
set OPENAI_API_KEY=sk-your-api-key-here
```

### Step 3: Install Dependencies

```bash
go mod download
```

### Step 4: Run the Example

```bash
go run examples/coding_qa/main.go
```

That's it! 🎉 You should see the agent answering programming questions.

## 📝 What You'll See

The example demonstrates:

1. **Agent Initialization**
   ```
   🤖 Coding Expert AI Agent - Programming Q&A System
   ============================================================
   
   🔧 Initializing LLM provider...
   ✅ LLM provider initialized
   🔧 Setting up language support...
   ✅ Registered Go language support
   🔧 Creating programming tools...
   ✅ Created 5 programming tools
   🔧 Initializing Coding Expert Agent...
   ✅ Coding Expert Agent ready
   ```

2. **Example Questions**
   - How to use goroutines in Go
   - Code analysis for issues
   - Best practices for error handling

3. **Interactive Mode**
   - Ask your own programming questions
   - Get instant answers with code examples

## 🎯 Try These Examples

### Example 1: Ask a Programming Question

```
You: How do I create a REST API in Go?
```

The agent will:
- Search documentation
- Find code examples
- Provide a complete answer with working code

### Example 2: Analyze Code

```
You: Analyze this code: package main

func main() {
    var x int
    if x == 0 {
    }
}
```

The agent will:
- Parse the code using AST
- Identify issues (empty if statement)
- Suggest improvements

### Example 3: Debug Code

```
You: Why does this code panic? 
var s []int
s[0] = 1
```

The agent will:
- Identify the issue (nil slice)
- Explain the problem
- Provide the correct solution

## 🛠️ Customization

### Use a Different Model

Edit `examples/coding_qa/main.go`:

```go
llmProvider, err := llm.NewOpenAIProvider(llm.OpenAIConfig{
    APIKey:      apiKey,
    Model:       "gpt-3.5-turbo", // Change this
    Temperature: 0.7,
    MaxTokens:   2000,
})
```

### Add More Languages

```go
// Add Python support (when implemented)
pythonProvider := python.NewProvider()
languageRegistry.Register(pythonProvider)
```

### Customize Agent Behavior

```go
codingAgent := agent.NewCodingExpertAgent(agent.CodingExpertConfig{
    Name:           "MyCustomAgent",
    LLM:            llmProvider,
    Tools:          tools,
    MaxSteps:       10,              // More reasoning steps
    Temperature:    0.5,             // More deterministic
    VerboseLogging: true,            // Detailed logs
    SupportedLangs: []string{"Go"},  // Specific languages
})
```

## 📚 Next Steps

### 1. Explore the Code

```bash
# View the agent implementation
cat internal/agent/coding_expert.go

# View the Go language analyzer
cat internal/languages/golang/analyzer.go

# View the tools
ls internal/tools/programming/
```

### 2. Read the Documentation

- `CODING_AGENTS_README.md` - Complete user guide
- `CODING_AGENTS_IMPLEMENTATION.md` - Implementation details
- `pkg/types/coding.go` - Type definitions

### 3. Build Your Own Agent

Create a new file `my_agent.go`:

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "ai-agent-platform/internal/agent"
    "ai-agent-platform/internal/languages/common"
    "ai-agent-platform/internal/languages/golang"
    "ai-agent-platform/internal/llm"
    "ai-agent-platform/internal/tools/programming"
    "ai-agent-platform/pkg/types"
)

func main() {
    // Setup
    llmProvider, _ := llm.NewOpenAIProvider(llm.OpenAIConfig{
        APIKey: os.Getenv("OPENAI_API_KEY"),
        Model:  "gpt-4",
    })
    
    languageRegistry := common.NewLanguageRegistry()
    languageRegistry.Register(golang.NewProvider())
    
    tools := []types.Tool{
        programming.NewCodeAnalysisTool(languageRegistry),
        programming.NewCodeExecutionTool(languageRegistry),
    }
    
    agent := agent.NewCodingExpertAgent(agent.CodingExpertConfig{
        LLM:   llmProvider,
        Tools: tools,
    })
    
    // Ask a question
    result, _ := agent.Run(context.Background(), types.AgentInput{
        Query: "How do I use channels in Go?",
    })
    
    fmt.Println(result.Output)
}
```

Run it:
```bash
go run my_agent.go
```

## 🔧 Troubleshooting

### Issue: "OPENAI_API_KEY not set"

**Solution**: Make sure you've exported the environment variable:
```bash
export OPENAI_API_KEY="your-key-here"
```

### Issue: "Module not found"

**Solution**: Run `go mod download` to install dependencies:
```bash
go mod download
```

### Issue: "Connection refused"

**Solution**: Check your internet connection and API key validity.

### Issue: Agent is slow

**Solutions**:
1. Use a faster model: `gpt-3.5-turbo` instead of `gpt-4`
2. Reduce max tokens: `MaxTokens: 1000`
3. Enable caching (if implemented)

## 💡 Tips & Tricks

### 1. Better Questions Get Better Answers

❌ Bad: "How to use Go?"
✅ Good: "How do I implement a concurrent worker pool in Go with graceful shutdown?"

### 2. Provide Context

```go
agent.Run(ctx, types.AgentInput{
    Query: "Why is this slow?",
    Context: map[string]interface{}{
        "code": "your code here",
        "performance_metrics": "...",
    },
})
```

### 3. Use Specific Languages

```go
agent.Run(ctx, types.AgentInput{
    Query: "[Go] How do I handle errors?",
})
```

### 4. Request Code Examples

```
You: Show me an example of using goroutines with channels
```

### 5. Ask for Best Practices

```
You: What are the best practices for structuring a Go project?
```

## 📊 Understanding the Output

### Agent Response Structure

```go
type AgentOutput struct {
    Output    string        // The final answer
    Steps     []AgentStep   // Reasoning steps
    ToolCalls []ToolCall    // Tools that were used
    Metadata  AgentMetadata // Execution info
}
```

### Metadata Includes

- **Execution Time**: How long it took
- **Steps Taken**: Number of reasoning steps
- **Tokens Used**: LLM token consumption
- **Tools Used**: Which tools were called

### Example Output

```
🎯 Answer:
To use goroutines in Go, you use the `go` keyword...
[detailed answer with code examples]

📊 Metadata:
   - Execution Time: 2.5s
   - Steps Taken: 3
   - Tokens Used: 1250 (Prompt: 500, Completion: 750)

🔍 Reasoning Steps:
   Step 1: Understanding the question about goroutines
   Action: documentation_search
   Observation: Found official Go documentation...
   
   Step 2: Providing code examples
   Action: code_execution
   Observation: Code executed successfully...
```

## 🎓 Learning Path

### Beginner (Week 1)
1. ✅ Run the quick start example
2. ✅ Ask simple programming questions
3. ✅ Understand the agent output
4. ✅ Read the README

### Intermediate (Week 2)
1. ⏳ Customize agent configuration
2. ⏳ Use different tools
3. ⏳ Analyze your own code
4. ⏳ Explore the codebase

### Advanced (Week 3+)
1. ⏳ Add new language support
2. ⏳ Create custom tools
3. ⏳ Build specialized agents
4. ⏳ Integrate with your applications

## 🚀 Production Deployment

### Docker (Coming Soon)

```bash
docker build -t coding-agent .
docker run -e OPENAI_API_KEY=your-key coding-agent
```

### Kubernetes (Coming Soon)

```bash
kubectl apply -f k8s/deployment.yaml
```

### API Server (Coming Soon)

```bash
go run cmd/coding-agent-server/main.go
```

## 📞 Get Help

- **Documentation**: Check the `docs/` folder
- **Examples**: See `examples/` directory
- **Code**: Read the source in `internal/` and `pkg/`
- **Issues**: Open a GitHub issue

## 🎉 You're Ready!

You now have a working Coding Expert AI Agent system. Start asking programming questions and explore the capabilities!

**Happy Coding!** 🚀

---

**Next**: Read [CODING_AGENTS_README.md](CODING_AGENTS_README.md) for complete documentation

