# 🎯 Getting Started with FinAgent

Welcome to FinAgent - the Financial Services AI Agent Platform built in Go!

## 🚀 Quick Start (5 Minutes)

### Step 1: Prerequisites

Make sure you have:
- ✅ Go 1.22+ installed
- ✅ OpenAI API key (get one at https://platform.openai.com/)
- ✅ Git installed

### Step 2: Clone and Setup

```bash
# Navigate to the project
cd services/ai-agent-platform

# Copy environment template
cp .env.example .env

# Edit .env and add your OpenAI API key
# OPENAI_API_KEY=sk-your-key-here
```

### Step 3: Install Dependencies

```bash
make deps
```

This will download all required Go modules.

### Step 4: Run Your First Agent!

```bash
make example-fraud
```

You should see output like this:

```
🤖 FinAgent - Fraud Detection Example
=====================================

✅ LLM Provider initialized (OpenAI GPT-4)
✅ Loaded 3 tools
   - transaction_lookup: Look up details of a financial transaction...
   - fraud_check: Analyze a transaction for fraud indicators...
   - calculator: Perform mathematical calculations...

✅ ReAct Agent created

📝 Query 1: Look up transaction TXN_12345 and check if it's fraudulent
------------------------------------------------------------

🔍 Agent Execution:
   Duration: 3.2s
   Steps: 3
   Tokens Used: 1250

📊 Reasoning Steps:
   Step 1:
   💭 Thought: I need to first look up the transaction details
   🔧 Action: transaction_lookup
   ...

✨ Final Answer:
   The transaction TXN_12345 appears to be legitimate...
```

Congratulations! 🎉 You just ran your first AI agent!

---

## 📖 Understanding What Just Happened

### The ReAct Pattern

The agent you just ran uses the **ReAct (Reasoning + Acting)** pattern:

1. **Think** - The agent reasons about what to do
2. **Act** - The agent uses tools to gather information
3. **Observe** - The agent sees the results
4. **Repeat** - Until it has enough information to answer

### Example Flow

```
User Query: "Is transaction TXN_12345 fraudulent?"
    ↓
Step 1: Thought: "I need transaction details"
        Action: transaction_lookup(TXN_12345)
        Observation: {amount: $1250.50, merchant: "Amazon"...}
    ↓
Step 2: Thought: "Now I should check for fraud"
        Action: fraud_check(TXN_12345)
        Observation: {risk_score: 0.15, risk_level: "low"...}
    ↓
Step 3: Thought: "I have enough information"
        Final Answer: "Transaction appears legitimate..."
```

---

## 🛠️ Building Your Own Agent

### Example 1: Simple Agent

Create a file `my_agent.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "ai-agent-platform/internal/agent"
    "ai-agent-platform/internal/llm"
    "ai-agent-platform/internal/tools/general"
    "ai-agent-platform/pkg/types"
)

func main() {
    // 1. Create LLM provider
    llmProvider, err := llm.NewOpenAIProvider(llm.OpenAIConfig{
        APIKey: os.Getenv("OPENAI_API_KEY"),
        Model:  "gpt-4",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 2. Create tools
    tools := []types.Tool{
        general.NewCalculatorTool(),
    }

    // 3. Create agent
    myAgent := agent.NewReActAgent(agent.ReActConfig{
        Name:        "MathAgent",
        Description: "An agent that helps with math",
        LLM:         llmProvider,
        Tools:       tools,
        MaxSteps:    3,
    })

    // 4. Run agent
    result, err := myAgent.Run(context.Background(), types.AgentInput{
        Query: "What is 25 * 4 + 10?",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 5. Print result
    fmt.Println("Answer:", result.Output)
}
```

Run it:
```bash
go run my_agent.go
```

### Example 2: Custom Tool

Create your own tool:

```go
package main

import (
    "context"
    "fmt"
    "ai-agent-platform/pkg/types"
)

type WeatherTool struct{}

func (t *WeatherTool) Name() string {
    return "get_weather"
}

func (t *WeatherTool) Description() string {
    return "Get current weather for a city"
}

func (t *WeatherTool) Execute(ctx context.Context, input types.ToolInput) (*types.ToolOutput, error) {
    city, _ := input.GetString("city")
    
    // In production, call a real weather API
    weather := map[string]interface{}{
        "city":        city,
        "temperature": 72,
        "condition":   "sunny",
    }
    
    return types.NewToolOutput(weather), nil
}

func (t *WeatherTool) GetSchema() types.ToolSchema {
    return types.ToolSchema{
        Type: "object",
        Properties: map[string]types.PropertySchema{
            "city": {
                Type:        "string",
                Description: "Name of the city",
            },
        },
        Required: []string{"city"},
    }
}

func (t *WeatherTool) Validate(input types.ToolInput) error {
    if _, ok := input.GetString("city"); !ok {
        return fmt.Errorf("city is required")
    }
    return nil
}
```

Use it in your agent:

```go
tools := []types.Tool{
    &WeatherTool{},
}

agent := agent.NewReActAgent(agent.ReActConfig{
    Name:  "WeatherAgent",
    LLM:   llmProvider,
    Tools: tools,
})

result, _ := agent.Run(ctx, types.AgentInput{
    Query: "What's the weather in San Francisco?",
})
```

---

## 🎓 Next Steps

### 1. Explore Examples

```bash
# Fraud detection
make example-fraud

# Customer support (coming soon)
make example-customer-support

# Risk assessment (coming soon)
make example-risk-assessment
```

### 2. Read Documentation

- **Architecture**: See `README.md`
- **Implementation Plan**: See `IMPLEMENTATION_GUIDE.md`
- **Project Summary**: See `PROJECT_SUMMARY.md`

### 3. Run with Docker

```bash
# Start full stack (PostgreSQL, Redis, Jaeger, etc.)
make docker-compose-up

# View logs
make docker-compose-logs

# Stop services
make docker-compose-down
```

### 4. Access Monitoring

Once Docker Compose is running:

- **Jaeger (Tracing)**: http://localhost:16686
- **Grafana (Dashboards)**: http://localhost:3000 (admin/admin)
- **Prometheus (Metrics)**: http://localhost:9090

### 5. Build Your Own Agent

Start with the examples and modify them:

```bash
# Copy an example
cp examples/fraud_detection/main.go my_custom_agent.go

# Edit and customize
# Run it
go run my_custom_agent.go
```

---

## 🔧 Common Tasks

### Add a New Tool

1. Create file in `internal/tools/`
2. Implement the `Tool` interface
3. Register it with your agent

### Change LLM Model

```go
llmProvider, _ := llm.NewOpenAIProvider(llm.OpenAIConfig{
    APIKey: os.Getenv("OPENAI_API_KEY"),
    Model:  "gpt-3.5-turbo", // Faster and cheaper
})
```

### Adjust Agent Behavior

```go
agent := agent.NewReActAgent(agent.ReActConfig{
    MaxSteps:    10,        // More reasoning steps
    Temperature: 0.3,       // More deterministic
    MaxTokens:   4000,      // Longer responses
})
```

### Enable Verbose Logging

```go
agent := agent.NewReActAgent(agent.ReActConfig{
    VerboseLogging: true,  // See detailed execution
})
```

---

## 🐛 Troubleshooting

### "OPENAI_API_KEY not set"

```bash
# Set in .env file
echo "OPENAI_API_KEY=sk-your-key" >> .env

# Or export directly
export OPENAI_API_KEY=sk-your-key
```

### "Module not found"

```bash
# Download dependencies
make deps

# Or manually
go mod download
go mod tidy
```

### "Connection refused" (Docker)

```bash
# Check services are running
docker-compose ps

# Restart services
make docker-compose-down
make docker-compose-up
```

### Agent is slow

- Use `gpt-3.5-turbo` instead of `gpt-4`
- Reduce `MaxSteps`
- Enable caching: `ENABLE_CACHING=true`

---

## 📚 Learning Resources

### Go Resources
- [Go Tour](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

### AI Agent Concepts
- [ReAct Paper](https://arxiv.org/abs/2210.03629)
- [Langchain Docs](https://python.langchain.com/docs/)
- [OpenAI Function Calling](https://platform.openai.com/docs/guides/function-calling)

### Financial Services
- [Fraud Detection Techniques](https://www.kaggle.com/mlg-ulb/creditcardfraud)
- [Risk Assessment Models](https://www.investopedia.com/terms/r/risk-assessment.asp)

---

## 💬 Get Help

- **Documentation**: Check `docs/` folder
- **Examples**: See `examples/` folder
- **Issues**: Open a GitHub issue
- **Questions**: Ask in discussions

---

## 🎯 What's Next?

Now that you've got the basics, you can:

1. ✅ Build custom agents for your use case
2. ✅ Create domain-specific tools
3. ✅ Integrate with your existing systems
4. ✅ Deploy to production

**Happy building!** 🚀

---

**FinAgent - Built with ❤️ in Go for production financial services**

