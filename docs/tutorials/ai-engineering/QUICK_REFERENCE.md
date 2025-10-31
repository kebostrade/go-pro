# AI Engineering Quick Reference Guide

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║                    🚀 AI ENGINEERING QUICK REFERENCE                         ║
║                                                                              ║
║              Common Patterns, Code Snippets, and Best Practices              ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 📚 Table of Contents

1. [LLM Integration](#llm-integration)
2. [Prompt Engineering](#prompt-engineering)
3. [Embeddings & Vectors](#embeddings--vectors)
4. [RAG Systems](#rag-systems)
5. [AI Agents](#ai-agents)
6. [Error Handling](#error-handling)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)

---

## 🤖 LLM Integration

### Basic Chat Completion

```go
import openai "github.com/sashabaranov/go-openai"

client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

resp, err := client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleSystem,
                Content: "You are a helpful assistant.",
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: "Hello!",
            },
        },
        Temperature: 0.7,
        MaxTokens:   500,
    },
)

if err != nil {
    log.Fatal(err)
}

answer := resp.Choices[0].Message.Content
```

### Streaming Responses

```go
stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
    Model:    openai.GPT3Dot5Turbo,
    Messages: messages,
    Stream:   true,
})
defer stream.Close()

for {
    response, err := stream.Recv()
    if errors.Is(err, io.EOF) {
        break
    }
    if err != nil {
        return err
    }
    
    fmt.Print(response.Choices[0].Delta.Content)
}
```

### Function Calling

```go
functions := []openai.FunctionDefinition{
    {
        Name:        "get_weather",
        Description: "Get current weather for a location",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "location": map[string]interface{}{
                    "type":        "string",
                    "description": "City name",
                },
            },
            "required": []string{"location"},
        },
    },
}

resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
    Model:     openai.GPT3Dot5Turbo,
    Messages:  messages,
    Functions: functions,
})

// Check if function was called
if resp.Choices[0].Message.FunctionCall != nil {
    funcName := resp.Choices[0].Message.FunctionCall.Name
    funcArgs := resp.Choices[0].Message.FunctionCall.Arguments
    // Execute function...
}
```

---

## 📝 Prompt Engineering

### System Message Templates

```go
// Code Assistant
systemMsg := `You are an expert Go programmer. Write clean, idiomatic code 
with proper error handling. Include comments and follow Go best practices.`

// Technical Writer
systemMsg := `You are a technical writer. Explain complex topics clearly 
using simple language, examples, and analogies.`

// Code Reviewer
systemMsg := `You are a senior code reviewer. Analyze code for bugs, 
security issues, and performance problems. Be constructive and specific.`
```

### Few-Shot Learning

```go
messages := []openai.ChatCompletionMessage{
    {Role: "system", Content: "You extract structured data from text."},
    
    // Example 1
    {Role: "user", Content: "John Doe, john@example.com, 555-1234"},
    {Role: "assistant", Content: `{"name": "John Doe", "email": "john@example.com", "phone": "555-1234"}`},
    
    // Example 2
    {Role: "user", Content: "Jane Smith, jane@test.com, 555-5678"},
    {Role: "assistant", Content: `{"name": "Jane Smith", "email": "jane@test.com", "phone": "555-5678"}`},
    
    // Actual request
    {Role: "user", Content: "Bob Johnson, bob@demo.com, 555-9999"},
}
```

### Chain-of-Thought

```go
prompt := `Analyze this code for thread safety. Think step by step:

1. Identify shared state
2. Check for synchronization mechanisms
3. Look for potential race conditions
4. Provide verdict with explanation

Code:
%s`
```

---

## 🔢 Embeddings & Vectors

### Generate Embeddings

```go
// Using OpenAI
resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
    Model: openai.AdaEmbeddingV2,
    Input: []string{"Your text here"},
})

embedding := resp.Data[0].Embedding // []float32
```

### Cosine Similarity

```go
func cosineSimilarity(a, b []float32) float32 {
    var dotProduct, normA, normB float32
    
    for i := range a {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    
    return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
```

### Vector Search (Simple)

```go
type Document struct {
    ID        string
    Content   string
    Embedding []float32
}

func findSimilar(query []float32, docs []Document, topK int) []Document {
    type scored struct {
        doc   Document
        score float32
    }
    
    scores := make([]scored, len(docs))
    for i, doc := range docs {
        scores[i] = scored{
            doc:   doc,
            score: cosineSimilarity(query, doc.Embedding),
        }
    }
    
    // Sort by score descending
    sort.Slice(scores, func(i, j int) bool {
        return scores[i].score > scores[j].score
    })
    
    results := make([]Document, min(topK, len(scores)))
    for i := range results {
        results[i] = scores[i].doc
    }
    
    return results
}
```

---

## 📚 RAG Systems

### Basic RAG Pipeline

```go
type RAGPipeline struct {
    client      *openai.Client
    vectorStore VectorStore
}

func (r *RAGPipeline) Query(ctx context.Context, question string) (string, error) {
    // 1. Generate query embedding
    queryEmb, err := r.generateEmbedding(ctx, question)
    if err != nil {
        return "", err
    }
    
    // 2. Retrieve relevant documents
    docs := r.vectorStore.Search(queryEmb, 5)
    
    // 3. Build context from documents
    context := r.buildContext(docs)
    
    // 4. Generate answer with context
    prompt := fmt.Sprintf(`Answer the question based on this context:

Context:
%s

Question: %s

Answer:`, context, question)
    
    resp, err := r.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {Role: "system", Content: "You are a helpful assistant. Answer based on the provided context."},
            {Role: "user", Content: prompt},
        },
    })
    
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}
```

### Document Chunking

```go
func chunkText(text string, chunkSize, overlap int) []string {
    words := strings.Fields(text)
    var chunks []string
    
    for i := 0; i < len(words); i += chunkSize - overlap {
        end := i + chunkSize
        if end > len(words) {
            end = len(words)
        }
        
        chunk := strings.Join(words[i:end], " ")
        chunks = append(chunks, chunk)
        
        if end == len(words) {
            break
        }
    }
    
    return chunks
}
```

---

## 🤖 AI Agents

### ReAct Agent Loop

```go
type Agent struct {
    llm   LLMProvider
    tools map[string]Tool
}

func (a *Agent) Run(ctx context.Context, task string) (string, error) {
    messages := []Message{
        {Role: "system", Content: a.systemPrompt()},
        {Role: "user", Content: task},
    }
    
    for i := 0; i < maxIterations; i++ {
        // 1. Think (LLM decides next action)
        resp, err := a.llm.Generate(ctx, messages)
        if err != nil {
            return "", err
        }
        
        // 2. Check if done
        if resp.IsFinal {
            return resp.Content, nil
        }
        
        // 3. Act (execute tool)
        if resp.FunctionCall != nil {
            result, err := a.executeTool(ctx, resp.FunctionCall)
            if err != nil {
                return "", err
            }
            
            // 4. Observe (add result to context)
            messages = append(messages, Message{
                Role:    "function",
                Content: result,
                Name:    resp.FunctionCall.Name,
            })
        }
    }
    
    return "", errors.New("max iterations reached")
}
```

### Tool Definition

```go
type Tool struct {
    Name        string
    Description string
    Parameters  map[string]interface{}
    Execute     func(ctx context.Context, args map[string]interface{}) (string, error)
}

// Example: Weather tool
var WeatherTool = Tool{
    Name:        "get_weather",
    Description: "Get current weather for a location",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "location": map[string]interface{}{
                "type":        "string",
                "description": "City name",
            },
        },
        "required": []string{"location"},
    },
    Execute: func(ctx context.Context, args map[string]interface{}) (string, error) {
        location := args["location"].(string)
        // Call weather API...
        return fmt.Sprintf("Weather in %s: 72°F, Sunny", location), nil
    },
}
```

---

## ⚠️ Error Handling

### Retry Logic

```go
func withRetry(ctx context.Context, maxRetries int, fn func() error) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        
        // Exponential backoff
        wait := time.Duration(math.Pow(2, float64(i))) * time.Second
        select {
        case <-time.After(wait):
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return err
}
```

### Rate Limiting

```go
type RateLimiter struct {
    tokens chan struct{}
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, requestsPerSecond),
    }
    
    go func() {
        ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
        defer ticker.Stop()
        
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
            }
        }
    }()
    
    return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

## ✅ Best Practices

### 1. Always Use Context

```go
// ✅ Good
func query(ctx context.Context, prompt string) (string, error) {
    resp, err := client.CreateChatCompletion(ctx, request)
    // ...
}

// ❌ Bad
func query(prompt string) (string, error) {
    resp, err := client.CreateChatCompletion(context.Background(), request)
    // ...
}
```

### 2. Handle Errors Properly

```go
// ✅ Good
resp, err := client.CreateChatCompletion(ctx, request)
if err != nil {
    var apiErr *openai.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.HTTPStatusCode {
        case 429:
            return "", errors.New("rate limit exceeded")
        case 401:
            return "", errors.New("invalid API key")
        default:
            return "", fmt.Errorf("API error: %w", err)
        }
    }
    return "", err
}
```

### 3. Cache Responses

```go
type Cache struct {
    store map[string]CachedResponse
    mu    sync.RWMutex
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    cached, ok := c.store[key]
    if !ok || time.Now().After(cached.ExpiresAt) {
        return "", false
    }
    
    return cached.Response, true
}
```

### 4. Monitor Token Usage

```go
type TokenTracker struct {
    totalTokens      int64
    totalCost        float64
    mu               sync.Mutex
}

func (t *TokenTracker) Track(usage openai.Usage, model string) {
    t.mu.Lock()
    defer t.mu.Unlock()
    
    t.totalTokens += int64(usage.TotalTokens)
    t.totalCost += calculateCost(usage, model)
}
```

---

## 🔧 Troubleshooting

### Common Issues

**Issue**: "Invalid API key"
```bash
# Solution: Check environment variable
echo $OPENAI_API_KEY

# Set it properly
export OPENAI_API_KEY="sk-..."
```

**Issue**: "Rate limit exceeded"
```go
// Solution: Implement exponential backoff
err := withRetry(ctx, 3, func() error {
    resp, err := client.CreateChatCompletion(ctx, request)
    return err
})
```

**Issue**: "Context length exceeded"
```go
// Solution: Trim conversation history
func trimMessages(messages []Message, maxTokens int) []Message {
    // Keep system message + recent messages
    if len(messages) <= 1 {
        return messages
    }
    
    system := messages[0]
    recent := messages[max(1, len(messages)-10):]
    
    return append([]Message{system}, recent...)
}
```

---

## 📊 Model Comparison

| Model | Context | Speed | Cost | Best For |
|-------|---------|-------|------|----------|
| GPT-3.5-turbo | 4K | ⚡⚡⚡ | $ | Development, simple tasks |
| GPT-3.5-turbo-16k | 16K | ⚡⚡ | $$ | Longer contexts |
| GPT-4 | 8K | ⚡ | $$$$ | Complex reasoning |
| GPT-4-turbo | 128K | ⚡⚡ | $$$ | Long documents, production |

---

## 💰 Cost Optimization

```go
// 1. Use cheaper models for simple tasks
model := openai.GPT3Dot5Turbo
if complexity == "high" {
    model = openai.GPT4
}

// 2. Cache responses
cacheKey := hashPrompt(prompt)
if cached, ok := cache.Get(cacheKey); ok {
    return cached, nil
}

// 3. Limit max tokens
MaxTokens: 500  // Don't use more than needed

// 4. Batch requests when possible
// Process multiple items in one request
```

---

**💡 Bookmark this page for quick reference while building AI applications!**

