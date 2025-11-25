# 🚀 FinAgent Quick Start Guide

Get started with the Financial Services AI Agent Platform in 5 minutes!

## Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose (for full stack)
- OpenAI API key (or Anthropic API key)

## Option 1: Run Example (Fastest)

### 1. Set up environment

```bash
cd services/ai-agent-platform
cp .env.example .env
```

### 2. Add your OpenAI API key to `.env`

```bash
# Edit .env and add your key
OPENAI_API_KEY=sk-your-key-here
```

### 3. Install dependencies

```bash
make deps
```

### 4. Run the fraud detection example

```bash
make example-fraud
```

You should see output like:

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
   Tokens Used: 1250 (Prompt: 850, Completion: 400)

📊 Reasoning Steps:

   Step 1:
   💭 Thought: I need to first look up the transaction details
   🔧 Action: transaction_lookup
   📥 Input: {"transaction_id": "TXN_12345"}
   👁️  Observation: {"transaction_id":"TXN_12345","amount":1250.50,...}

   Step 2:
   💭 Thought: Now I should check this transaction for fraud
   🔧 Action: fraud_check
   📥 Input: {"transaction_id": "TXN_12345"}
   👁️  Observation: {"risk_score":0.15,"risk_level":"low",...}

   Step 3:
   💭 Thought: Based on the analysis, I can provide a final answer
   Final Answer: The transaction TXN_12345 appears to be legitimate...

✨ Final Answer:
   The transaction TXN_12345 appears to be legitimate with a low risk score of 0.15...
```

## Option 2: Run Full Stack with Docker

### 1. Set up environment

```bash
cd services/ai-agent-platform
cp .env.example .env
# Edit .env and add your API keys
```

### 2. Start all services

```bash
make docker-compose-up
```

This starts:
- PostgreSQL with pgvector (port 5432)
- Redis (port 6379)
- Qdrant vector database (port 6333)
- Jaeger tracing (port 16686)
- Prometheus (port 9090)
- Grafana (port 3000)
- AI Agent Platform (port 8080)

### 3. Access the services

- **Agent API**: http://localhost:8080
- **Jaeger UI**: http://localhost:16686
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090

### 4. Test the API

```bash
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Check transaction TXN_12345 for fraud",
    "agent_type": "react",
    "max_steps": 5
  }'
```

### 5. View traces in Jaeger

Open http://localhost:16686 and search for traces to see the complete execution flow.

## Option 3: Build from Source

### 1. Clone and setup

```bash
cd services/ai-agent-platform
make deps
```

### 2. Build the application

```bash
make build
```

### 3. Run the binary

```bash
export OPENAI_API_KEY=sk-your-key-here
./bin/finagent
```

## Quick Examples

### Example 1: Simple Fraud Check

```go
package main

import (
    "context"
    "fmt"
    "ai-agent-platform/internal/agent"
    "ai-agent-platform/internal/llm"
    "ai-agent-platform/internal/tools/financial"
    "ai-agent-platform/pkg/types"
)

func main() {
    // Create LLM provider
    llmProvider, _ := llm.NewOpenAIProvider(llm.OpenAIConfig{
        APIKey: "your-api-key",
        Model:  "gpt-4",
    })

    // Create agent with tools
    agent := agent.NewReActAgent(agent.ReActConfig{
        Name: "FraudAgent",
        LLM:  llmProvider,
        Tools: []types.Tool{
            financial.NewTransactionLookupTool(),
            financial.NewFraudCheckTool(),
        },
    })

    // Run agent
    result, _ := agent.Run(context.Background(), types.AgentInput{
        Query: "Is transaction TXN_12345 fraudulent?",
    })

    fmt.Println(result.Output)
}
```

### Example 2: Customer Support Agent

```go
// Create a customer support agent
supportAgent := agent.NewReActAgent(agent.ReActConfig{
    Name: "CustomerSupportAgent",
    LLM:  llmProvider,
    Tools: []types.Tool{
        financial.NewTransactionLookupTool(),
        financial.NewAccountInfoTool(),
        general.NewCalculatorTool(),
    },
    MaxSteps: 10,
})

// Handle customer query
result, _ := supportAgent.Run(ctx, types.AgentInput{
    Query: "I don't recognize a charge of $1,250.50 from yesterday",
    UserID: "user_123",
})
```

## Available Tools

### Financial Tools
- `transaction_lookup` - Look up transaction details
- `fraud_check` - Analyze transactions for fraud
- `account_info` - Get account information
- `risk_calculator` - Calculate risk scores

### General Tools
- `calculator` - Perform calculations
- `web_search` - Search the web
- `api_caller` - Call external APIs

## Development Commands

```bash
# Run tests
make test

# Run with hot reload
make dev

# Format code
make fmt

# Run linter
make lint

# Generate coverage report
make test-coverage

# Build Docker image
make docker-build

# View logs
make docker-compose-logs
```

## Monitoring & Observability

### View Traces
1. Open Jaeger UI: http://localhost:16686
2. Select service: `finagent`
3. Click "Find Traces"

### View Metrics
1. Open Grafana: http://localhost:3000
2. Login: admin/admin
3. Navigate to dashboards

### View Logs
```bash
# Application logs
docker-compose logs -f agent-platform

# All services
docker-compose logs -f
```

## Configuration

Key configuration options in `.env`:

```bash
# LLM Provider
DEFAULT_LLM_PROVIDER=openai
DEFAULT_MODEL=gpt-4

# Agent Settings
AGENT_MAX_STEPS=5
AGENT_TIMEOUT=5m
AGENT_TEMPERATURE=0.7

# Features
ENABLE_STREAMING=true
ENABLE_MEMORY=true
ENABLE_CACHING=true
```

## Troubleshooting

### "OPENAI_API_KEY not set"
```bash
export OPENAI_API_KEY=sk-your-key-here
```

### "Connection refused" errors
Make sure all services are running:
```bash
docker-compose ps
```

### High latency
- Check your LLM provider status
- Reduce `AGENT_MAX_STEPS`
- Enable caching: `ENABLE_CACHING=true`

## Next Steps

1. **Explore Examples**: Check out `examples/` directory
2. **Read Documentation**: See `docs/` for detailed guides
3. **Build Custom Agents**: Create your own agents and tools
4. **Deploy to Production**: Follow `docs/DEPLOYMENT.md`

## Support

- 📖 Documentation: `docs/`
- 💬 Issues: GitHub Issues
- 📧 Email: support@finagent.dev

Happy building! 🚀

