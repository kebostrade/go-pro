# 🌐 Coding Agent API Documentation

## Overview

The Coding Agent API provides RESTful endpoints for programming assistance, code analysis, execution, and debugging.

**Base URL**: `http://localhost:8080/api/v1`

## 🚀 Quick Start

### Start the Server

```bash
cd services/ai-agent-platform
export OPENAI_API_KEY="your-api-key"
go run cmd/coding-agent-server/main.go
```

The server will start on port 8080 (configurable via `PORT` environment variable).

### Test the API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# List supported languages
curl http://localhost:8080/api/v1/languages
```

## 📡 API Endpoints

### 1. Ask Programming Question

Ask the coding agent a programming question and get an answer with code examples.

**Endpoint**: `POST /api/v1/coding/ask`

**Request Body**:
```json
{
  "query": "How do I create a goroutine in Go?",
  "context": {
    "difficulty": "beginner",
    "include_examples": true
  }
}
```

**Response**:
```json
{
  "answer": "To create a goroutine in Go, use the `go` keyword...",
  "steps": 3,
  "metadata": {
    "execution_time_ms": 2500,
    "tokens_used": 1250
  }
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/coding/ask \
  -H "Content-Type: application/json" \
  -d '{
    "query": "How do I handle errors in Go?"
  }'
```

---

### 2. Analyze Code

Analyze code for issues, security vulnerabilities, and best practices.

**Endpoint**: `POST /api/v1/coding/analyze`

**Request Body**:
```json
{
  "code": "package main\n\nfunc main() {\n    var x int\n    if x == 0 {\n    }\n}",
  "language": "go"
}
```

**Response**:
```json
{
  "language": "go",
  "valid_syntax": true,
  "issues": [
    {
      "type": "best_practice",
      "severity": "low",
      "line": 4,
      "message": "Empty if statement",
      "suggestion": "Remove empty if block or add implementation"
    }
  ],
  "metrics": {
    "lines_of_code": 6,
    "complexity": 1,
    "maintainability_index": 85
  },
  "security_issues": [],
  "performance_issues": []
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/coding/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello\")\n}",
    "language": "go"
  }'
```

---

### 3. Execute Code

Execute code in a secure sandbox and get the output.

**Endpoint**: `POST /api/v1/coding/execute`

**Request Body**:
```json
{
  "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}",
  "language": "go",
  "input": "",
  "timeout": 30
}
```

**Response**:
```json
{
  "success": true,
  "output": "Hello, World!\n",
  "error": "",
  "exit_code": 0,
  "execution_time": 1250,
  "memory_used": 45000000,
  "cpu_time": 800
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/coding/execute \
  -H "Content-Type: application/json" \
  -d '{
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello\")\n}",
    "language": "go"
  }'
```

---

### 4. Debug Code

Get help debugging code with error messages.

**Endpoint**: `POST /api/v1/coding/debug`

**Request Body**:
```json
{
  "code": "var s []int\ns[0] = 1",
  "language": "go",
  "error": "panic: runtime error: index out of range [0] with length 0"
}
```

**Response**:
```json
{
  "explanation": "The error occurs because you're trying to access index 0 of a nil slice...",
  "suggestions": [
    "Initialize the slice: s := make([]int, 1)",
    "Use append: s = append(s, 1)"
  ]
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/coding/debug \
  -H "Content-Type: application/json" \
  -d '{
    "code": "var s []int\ns[0] = 1",
    "language": "go",
    "error": "panic: index out of range"
  }'
```

---

### 5. Health Check

Check if the API server is running.

**Endpoint**: `GET /api/v1/health`

**Response**:
```json
{
  "status": "healthy",
  "timestamp": 1704067200,
  "version": "1.0.0"
}
```

**cURL Example**:
```bash
curl http://localhost:8080/api/v1/health
```

---

### 6. List Supported Languages

Get a list of supported programming languages.

**Endpoint**: `GET /api/v1/languages`

**Response**:
```json
{
  "languages": ["go", "python", "javascript", "typescript"],
  "count": 4
}
```

**cURL Example**:
```bash
curl http://localhost:8080/api/v1/languages
```

## 🔐 Authentication

Currently, the API does not require authentication. In production, implement:

- API keys
- JWT tokens
- Rate limiting
- IP whitelisting

## ⚠️ Error Responses

All errors follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

**Common HTTP Status Codes**:

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Bad Request (invalid input) |
| 405 | Method Not Allowed |
| 500 | Internal Server Error |

## 📊 Rate Limits

**Current Limits** (to be implemented):
- 100 requests per minute per IP
- 1000 requests per hour per IP
- Maximum request size: 1MB

## 🧪 Testing

### Using cURL

```bash
# Test ask endpoint
curl -X POST http://localhost:8080/api/v1/coding/ask \
  -H "Content-Type: application/json" \
  -d '{"query": "How to use channels in Go?"}'

# Test analyze endpoint
curl -X POST http://localhost:8080/api/v1/coding/analyze \
  -H "Content-Type: application/json" \
  -d '{"code": "package main\n\nfunc main() {}", "language": "go"}'
```

### Using JavaScript (fetch)

```javascript
// Ask a question
fetch('http://localhost:8080/api/v1/coding/ask', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    query: 'How do I create a REST API in Go?'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

### Using Python (requests)

```python
import requests

# Ask a question
response = requests.post(
    'http://localhost:8080/api/v1/coding/ask',
    json={'query': 'How to handle errors in Go?'}
)
print(response.json())

# Analyze code
response = requests.post(
    'http://localhost:8080/api/v1/coding/analyze',
    json={
        'code': 'package main\n\nfunc main() {}',
        'language': 'go'
    }
)
print(response.json())
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `OPENAI_API_KEY` | OpenAI API key | Required |
| `LLM_MODEL` | LLM model to use | gpt-4 |

### Example

```bash
export PORT=3000
export OPENAI_API_KEY="sk-..."
export LLM_MODEL="gpt-3.5-turbo"
go run cmd/coding-agent-server/main.go
```

## 📈 Performance

### Typical Response Times

| Endpoint | Average | P95 | P99 |
|----------|---------|-----|-----|
| /ask | 2.5s | 5s | 8s |
| /analyze | 500ms | 1s | 2s |
| /execute | 1.5s | 3s | 5s |
| /debug | 3s | 6s | 10s |
| /health | 5ms | 10ms | 20ms |

## 🚀 Production Deployment

### Docker

```bash
# Build image
docker build -t coding-agent-api .

# Run container
docker run -p 8080:8080 \
  -e OPENAI_API_KEY="your-key" \
  coding-agent-api
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: coding-agent-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: coding-agent-api
  template:
    metadata:
      labels:
        app: coding-agent-api
    spec:
      containers:
      - name: api
        image: coding-agent-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: openai-key
```

## 📞 Support

For issues or questions:
- Check the health endpoint
- Review error messages
- Enable verbose logging
- Check server logs

---

**Built for developers, by developers** 🚀

