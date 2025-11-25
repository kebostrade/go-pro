# 🔍 RAG (Retrieval-Augmented Generation) Guide

## Overview

The RAG system provides semantic search capabilities for code and documentation using vector embeddings. It enables the coding agents to retrieve relevant context from a knowledge base to improve answer quality.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                  RAG Pipeline                        │
├─────────────────────────────────────────────────────┤
│                                                      │
│  ┌──────────────┐    ┌──────────────┐              │
│  │   Embedder   │    │ Vector Store │              │
│  │  (OpenAI)    │    │  (Memory)    │              │
│  └──────┬───────┘    └──────┬───────┘              │
│         │                    │                       │
│         ▼                    ▼                       │
│  ┌─────────────────────────────────┐                │
│  │      Similarity Search          │                │
│  │   (Cosine Similarity)           │                │
│  └─────────────┬───────────────────┘                │
│                │                                     │
│                ▼                                     │
│  ┌─────────────────────────────────┐                │
│  │    Context Formatting           │                │
│  │  (For LLM Consumption)          │                │
│  └─────────────────────────────────┘                │
│                                                      │
└─────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### 1. Basic Setup

```go
import (
    "ai-agent-platform/internal/embeddings"
    "ai-agent-platform/internal/rag"
    "ai-agent-platform/internal/vectorstore"
    "ai-agent-platform/pkg/types"
)

// Create embedder
embedder, err := embeddings.NewOpenAIEmbedder(embeddings.OpenAIEmbedderConfig{
    APIKey: os.Getenv("OPENAI_API_KEY"),
    Model:  "text-embedding-3-small",
})

// Create vector store
vectorStore := vectorstore.NewMemoryVectorStore()

// Create RAG pipeline
pipeline, err := rag.NewRAGPipeline(types.RAGConfig{
    VectorStore:     vectorStore,
    Embedder:        embedder,
    TopK:            5,
    MinScore:        0.7,
    MaxTokens:       4000,
    IncludeMetadata: true,
})
```

### 2. Add Documents

```go
ctx := context.Background()

// Add a code snippet
metadata := map[string]interface{}{
    "language":    "go",
    "description": "Goroutine example",
}

err := pipeline.AddDocument(ctx, "doc-1", codeSnippet, metadata)
```

### 3. Search

```go
// Search for relevant documents
result, err := pipeline.Retrieve(ctx, "How to use goroutines?", nil)

fmt.Printf("Found %d results\n", result.Metadata.TotalResults)
for _, res := range result.Results {
    fmt.Printf("Score: %.2f - %s\n", res.Score, res.Metadata["description"])
}
```

## 📊 Components

### 1. Vector Store

**Interface**: `types.VectorStore`

**Implementations**:
- `MemoryVectorStore` - In-memory storage (development)
- Future: Redis, PostgreSQL with pgvector, Pinecone, etc.

**Operations**:
```go
// Store a vector
vectorStore.Store(id, vector, metadata)

// Search for similar vectors
results, err := vectorStore.Search(queryVector, limit, filters)

// Delete a vector
vectorStore.Delete(id)

// Update a vector
vectorStore.Update(id, newVector, newMetadata)

// Get count
count, err := vectorStore.Count()
```

### 2. Embedder

**Interface**: `types.Embedder`

**Implementations**:
- `OpenAIEmbedder` - Uses OpenAI's embedding API
- `MockEmbedder` - For testing without API calls

**Operations**:
```go
// Embed single text
vector, err := embedder.Embed("Hello, world!")

// Embed batch
vectors, err := embedder.EmbedBatch([]string{"text1", "text2"})

// Get dimensions
dims := embedder.Dimensions() // 1536 for text-embedding-3-small
```

### 3. RAG Pipeline

**Core Pipeline**: `rag.RAGPipeline`

**Specialized Pipelines**:
- `CodeRAGPipeline` - For code search
- `DocumentRAGPipeline` - For documentation search

**Operations**:
```go
// Retrieve relevant context
result, err := pipeline.Retrieve(ctx, query, filters)

// Add document
pipeline.AddDocument(ctx, id, content, metadata)

// Update document
pipeline.UpdateDocument(ctx, id, newContent, newMetadata)

// Delete document
pipeline.DeleteDocument(ctx, id)
```

## 🔧 Configuration

### RAG Config

```go
type RAGConfig struct {
    VectorStore     VectorStore  // Required
    Embedder        Embedder     // Required
    TopK            int          // Number of results (default: 5)
    MinScore        float64      // Minimum similarity (default: 0.7)
    MaxTokens       int          // Max context tokens (default: 4000)
    IncludeMetadata bool         // Include metadata in results
    RerankerEnabled bool         // Enable reranking (future)
}
```

### Embedder Config

```go
type OpenAIEmbedderConfig struct {
    APIKey     string  // Required
    Model      string  // Default: "text-embedding-3-small"
    Dimensions int     // Optional: custom dimensions
}
```

## 📝 Code Search

### Setup

```go
// Create code RAG pipeline
codePipeline, err := rag.NewCodeRAGPipeline(ragConfig)

// Create indexer
indexer := rag.NewCodeIndexer(codePipeline, embedder)
```

### Index Code

```go
// Index a directory
options := rag.IndexOptions{
    Repository:      "github.com/user/repo",
    Branch:          "main",
    IncludePatterns: []string{"*.go", "*.py"},
    MaxFileSize:     1024 * 1024, // 1MB
}

err := indexer.IndexDirectory(ctx, "/path/to/code", options)
```

### Search Code

```go
request := types.CodeSearchRequest{
    Query:    "goroutine example",
    Language: "go",
    Limit:    5,
    MinScore: 0.7,
}

response, err := codePipeline.SearchCode(ctx, request)

for _, result := range response.Results {
    fmt.Printf("Code: %s\n", result.Code)
    fmt.Printf("Language: %s\n", result.Language)
    fmt.Printf("Score: %.2f\n", result.Score)
}
```

## 📚 Documentation Search

### Setup

```go
// Create documentation RAG pipeline
docPipeline, err := rag.NewDocumentRAGPipeline(ragConfig)

// Create indexer
indexer := rag.NewDocumentIndexer(docPipeline, embedder)
```

### Index Documentation

```go
doc := rag.DocumentToIndex{
    ID:       "go-doc-1",
    Content:  "Goroutines are lightweight threads...",
    Title:    "Goroutines",
    Source:   "Go Documentation",
    Language: "go",
    URL:      "https://go.dev/doc/goroutines",
    Category: "concurrency",
    Tags:     []string{"goroutines", "concurrency"},
}

err := indexer.IndexDocument(ctx, doc)
```

### Search Documentation

```go
request := types.DocumentSearchRequest{
    Query:    "What are goroutines?",
    Language: "go",
    Source:   "Go Documentation",
    Limit:    3,
}

response, err := docPipeline.SearchDocumentation(ctx, request)

for _, result := range response.Results {
    fmt.Printf("Title: %s\n", result.Title)
    fmt.Printf("Content: %s\n", result.Content)
    fmt.Printf("Score: %.2f\n", result.Score)
}
```

## 🎯 Use Cases

### 1. Code Examples

```go
// User asks: "How to use channels in Go?"

// RAG retrieves relevant code examples
result, _ := codePipeline.SearchCode(ctx, types.CodeSearchRequest{
    Query:    "channel communication in Go",
    Language: "go",
    Limit:    3,
})

// Agent uses retrieved examples to generate answer
```

### 2. Documentation Lookup

```go
// User asks: "What is the syntax for list comprehensions?"

// RAG retrieves relevant documentation
result, _ := docPipeline.SearchDocumentation(ctx, types.DocumentSearchRequest{
    Query:    "list comprehension syntax",
    Language: "python",
    Limit:    2,
})

// Agent uses documentation to provide accurate answer
```

### 3. Similar Code Search

```go
// User provides code and asks: "Find similar examples"

// Embed user's code
userVector, _ := embedder.Embed(userCode)

// Search for similar code
results, _ := vectorStore.Search(userVector, 5, map[string]interface{}{
    "language": "go",
})
```

## 📈 Performance

### Embedding Models

| Model | Dimensions | Speed | Cost |
|-------|-----------|-------|------|
| text-embedding-3-small | 1536 | Fast | Low |
| text-embedding-3-large | 3072 | Medium | Medium |
| text-embedding-ada-002 | 1536 | Fast | Low |

### Similarity Metrics

- **Cosine Similarity**: Default, works well for most cases
- **Euclidean Distance**: Alternative metric
- **Dot Product**: For normalized vectors

### Optimization Tips

1. **Batch Embeddings**: Use `EmbedBatch()` for multiple texts
2. **Cache Embeddings**: Store embeddings to avoid recomputation
3. **Filter Early**: Use metadata filters to reduce search space
4. **Chunk Large Documents**: Split into smaller, focused chunks
5. **Adjust TopK**: Balance between quality and speed

## 🔒 Security

### API Key Management

```go
// Never hardcode API keys
apiKey := os.Getenv("OPENAI_API_KEY")

// Use secrets management in production
```

### Data Privacy

- Embeddings are sent to OpenAI API
- Consider using local models for sensitive data
- Implement access controls on vector store

## 🧪 Testing

### Mock Embedder

```go
// Use mock embedder for testing
embedder := embeddings.NewMockEmbedder(1536)

// No API calls, deterministic results
vector, _ := embedder.Embed("test")
```

### Example Test

```go
func TestRAGPipeline(t *testing.T) {
    embedder := embeddings.NewMockEmbedder(1536)
    vectorStore := vectorstore.NewMemoryVectorStore()
    
    pipeline, _ := rag.NewRAGPipeline(types.RAGConfig{
        VectorStore: vectorStore,
        Embedder:    embedder,
        TopK:        3,
    })
    
    // Add test documents
    pipeline.AddDocument(ctx, "doc1", "content1", nil)
    
    // Search
    result, _ := pipeline.Retrieve(ctx, "query", nil)
    
    assert.Equal(t, 1, result.Metadata.TotalResults)
}
```

## 📊 Monitoring

### Metrics to Track

```go
// Retrieval time
fmt.Printf("Retrieval: %dms\n", result.Metadata.RetrievalTime)

// Embedding time
fmt.Printf("Embedding: %dms\n", result.Metadata.EmbeddingTime)

// Context tokens
fmt.Printf("Tokens: %d\n", result.Metadata.ContextTokens)

// Results count
fmt.Printf("Results: %d\n", result.Metadata.TotalResults)
```

## 🚀 Production Deployment

### Vector Store Options

1. **PostgreSQL with pgvector**
   - Persistent storage
   - SQL queries
   - Good for moderate scale

2. **Redis with RediSearch**
   - Fast in-memory search
   - Good for caching
   - Moderate scale

3. **Pinecone/Weaviate**
   - Managed vector databases
   - High scale
   - Advanced features

### Best Practices

1. **Index Management**
   - Regular reindexing
   - Version control for embeddings
   - Backup vector store

2. **Monitoring**
   - Track search latency
   - Monitor embedding costs
   - Alert on failures

3. **Scaling**
   - Horizontal scaling for vector store
   - Batch processing for indexing
   - Caching for frequent queries

## 📞 Support

For issues or questions:
- Check the examples in `examples/rag_demo/`
- Review the code in `internal/rag/`
- Test with mock embedder first

---

**Built for semantic search and knowledge retrieval** 🔍

