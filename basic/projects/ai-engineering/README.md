# 🤖 AI Engineering Projects

Welcome to the AI Engineering projects! This directory contains hands-on projects to help you master building production-ready AI applications with Go.

---

## 📚 Projects Overview

### 🟢 Beginner Projects

#### 1. CLI Chatbot
**Duration**: 1-2 hours | **Difficulty**: Beginner

Build an interactive command-line chatbot powered by OpenAI's GPT models.

**What You'll Learn**:
- OpenAI API integration
- Streaming responses
- Conversation history management
- Token counting and optimization
- Model parameter tuning

**Tech Stack**: OpenAI API, Go standard library

**Directory**: [`chatbot-cli/`](chatbot-cli/)

---

### 🟡 Intermediate Projects

#### 2. Semantic Search Engine
**Duration**: 3-4 hours | **Difficulty**: Intermediate

Build a semantic search engine that understands meaning, not just keywords.

**What You'll Learn**:
- Text embeddings generation
- Vector similarity calculations
- In-memory vector store
- Qdrant integration
- Document indexing

**Tech Stack**: OpenAI Embeddings, Qdrant, Go

**Directory**: [`semantic-search/`](semantic-search/)

**Status**: 🚧 Template (implement based on Tutorial 3)

---

#### 3. RAG Q&A System
**Duration**: 4-5 hours | **Difficulty**: Intermediate

Build a Retrieval-Augmented Generation system for document Q&A.

**What You'll Learn**:
- RAG pipeline architecture
- Document chunking strategies
- Retrieval and ranking
- Context-aware generation
- Source citation

**Tech Stack**: OpenAI, Qdrant, Go

**Directory**: [`rag-qa-system/`](rag-qa-system/)

**Status**: 🚧 Template (implement based on Tutorial 4)

---

### 🔴 Advanced Projects

#### 4. Coding Assistant Agent
**Duration**: 6-8 hours | **Difficulty**: Advanced

Build an autonomous AI agent that can write, run, and analyze code.

**What You'll Learn**:
- ReAct pattern implementation
- Tool calling and execution
- Multi-step reasoning
- Agent state management
- Error recovery

**Tech Stack**: OpenAI Function Calling, Go

**Directory**: [`coding-assistant/`](coding-assistant/)

**Status**: 🚧 Template (implement based on Tutorial 5)

---

#### 5. Production AI Service
**Duration**: 10-12 hours | **Difficulty**: Advanced

Build a production-ready AI service with monitoring, caching, and rate limiting.

**What You'll Learn**:
- Production architecture
- Caching strategies
- Rate limiting
- Monitoring and logging
- Cost optimization
- Error handling at scale

**Tech Stack**: OpenAI, Redis, Prometheus, Grafana, Go

**Status**: 🎯 Coming Soon

---

## 🚀 Getting Started

### Prerequisites

1. **Go 1.21+** installed
2. **OpenAI API Key** ([Get one here](https://platform.openai.com/api-keys))
3. **Docker** (for Qdrant and other services)

### Quick Start

```bash
# 1. Clone or navigate to the project
cd basic/projects/ai-engineering

# 2. Choose a project
cd chatbot-cli

# 3. Set up environment
export OPENAI_API_KEY="sk-..."

# 4. Install dependencies
go mod tidy

# 5. Run the project
go run main.go
```

---

## 📖 Learning Path

### Recommended Order

```
1. CLI Chatbot (Week 1-2)
   ↓
2. Semantic Search (Week 3-4)
   ↓
3. RAG Q&A System (Week 5-7)
   ↓
4. Coding Assistant (Week 8-10)
   ↓
5. Production AI Service (Week 11-14)
```

### By Skill Level

**Beginner** (Start Here):
- CLI Chatbot

**Intermediate** (After Tutorials 1-4):
- Semantic Search Engine
- RAG Q&A System

**Advanced** (After Tutorials 1-5):
- Coding Assistant Agent
- Production AI Service

---

## 🎯 Project Comparison

| Project | Duration | Difficulty | Key Concepts | Prerequisites |
|---------|----------|------------|--------------|---------------|
| CLI Chatbot | 1-2h | 🟢 Beginner | LLM API, Streaming | Tutorial 1 |
| Semantic Search | 3-4h | 🟡 Intermediate | Embeddings, Vectors | Tutorial 3 |
| RAG Q&A | 4-5h | 🟡 Intermediate | RAG, Retrieval | Tutorial 4 |
| Coding Assistant | 6-8h | 🔴 Advanced | Agents, ReAct | Tutorial 5 |
| Production AI | 10-12h | 🔴 Advanced | Deployment, Scale | All Tutorials |

---

## 🛠️ Common Setup

### Environment Variables

Create a `.env` file in each project:

```bash
# OpenAI Configuration
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-3.5-turbo

# Vector Database (for projects 2-5)
QDRANT_URL=http://localhost:6333

# Redis (for projects 4-5)
REDIS_URL=localhost:6379

# Monitoring (for project 5)
PROMETHEUS_PORT=9090
```

### Docker Services

Start required services:

```bash
# Qdrant (Vector Database)
docker run -p 6333:6333 qdrant/qdrant

# Redis (Caching)
docker run -p 6379:6379 redis:alpine

# PostgreSQL (Optional)
docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:15
```

---

## 📊 Cost Estimation

### OpenAI API Costs (Approximate)

| Project | Estimated Cost | Notes |
|---------|---------------|-------|
| CLI Chatbot | $0.10 - $0.50 | Depends on usage |
| Semantic Search | $1 - $5 | Embedding generation |
| RAG Q&A | $2 - $10 | Embeddings + completions |
| Coding Assistant | $5 - $20 | Multiple tool calls |
| Production AI | Variable | Depends on traffic |

**💡 Cost Saving Tips**:
- Use GPT-3.5-turbo for development
- Cache embeddings (they don't change)
- Implement rate limiting
- Monitor token usage

---

## 🧪 Testing Your Projects

### Unit Tests

```bash
# Run tests for a project
cd chatbot-cli
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Integration Tests

```bash
# Make sure services are running
docker-compose up -d

# Run integration tests
go test -tags=integration ./...
```

---

## 🎓 Learning Resources

### Tutorials
- [Tutorial 0: AI Engineering Overview](../../../docs/tutorials/ai-engineering/00_AI_ENGINEERING_OVERVIEW.md)
- [Tutorial 1: LLM Basics](../../../docs/tutorials/ai-engineering/01_LLM_BASICS.md)
- [Tutorial 2: Prompt Engineering](../../../docs/tutorials/ai-engineering/02_PROMPT_ENGINEERING.md)
- [Tutorial 3: Embeddings & Vectors](../../../docs/tutorials/ai-engineering/03_EMBEDDINGS_VECTORS.md)
- [Tutorial 4: RAG Systems](../../../docs/tutorials/ai-engineering/04_RAG_SYSTEMS.md)
- [Tutorial 5: AI Agents](../../../docs/tutorials/ai-engineering/05_AI_AGENTS.md)

### Quick Reference
- [AI Engineering Quick Reference](../../../docs/tutorials/ai-engineering/QUICK_REFERENCE.md)

### External Resources
- [OpenAI Documentation](https://platform.openai.com/docs)
- [Qdrant Documentation](https://qdrant.tech/documentation/)
- [LangChain Concepts](https://python.langchain.com/docs/get_started/introduction)

---

## 🤝 Contributing

### Adding Your Own Project

1. Create a new directory: `basic/projects/ai-engineering/your-project/`
2. Add `README.md`, `main.go`, and `go.mod`
3. Follow the existing project structure
4. Update this README with your project

### Improving Existing Projects

- Add features
- Improve error handling
- Add tests
- Optimize performance
- Update documentation

---

## 🐛 Troubleshooting

### Common Issues

**"Invalid API key"**
```bash
# Check if key is set
echo $OPENAI_API_KEY

# Set it properly
export OPENAI_API_KEY="sk-..."
```

**"Connection refused" (Qdrant)**
```bash
# Make sure Qdrant is running
docker ps | grep qdrant

# Start Qdrant
docker run -p 6333:6333 qdrant/qdrant
```

**"Rate limit exceeded"**
- Wait a moment and retry
- Implement exponential backoff
- Upgrade your OpenAI plan

**"Context length exceeded"**
- Reduce conversation history
- Use smaller chunks
- Implement sliding window

---

## 📈 Next Steps

After completing these projects:

1. **Explore the AI Agent Platform**: Check out `services/ai-agent-platform/` for production examples
2. **Build Your Own**: Create a custom AI application
3. **Deploy to Production**: Use Docker, Kubernetes, or cloud platforms
4. **Contribute**: Share your improvements and new projects

---

## 📝 License

These projects are part of the Go-Pro learning platform.

---

**🚀 Ready to build AI applications? Start with the [CLI Chatbot](chatbot-cli/)!**

