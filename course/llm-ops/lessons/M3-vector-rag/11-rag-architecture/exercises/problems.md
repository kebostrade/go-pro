# Exercise: RAG Architecture

## Problem 1: Document Processing Pipeline Design

You're building a RAG system for a company with the following document sources:

| Source | Format | Volume | Update Frequency |
|--------|--------|--------|-------------------|
| Product documentation | PDF | 500 files | Weekly |
| API reference | OpenAPI (JSON) | 200 files | Daily |
| Support articles | HTML | 10,000 files | Monthly |
| Code comments | Markdown | 50,000 files | On-commit |

Design a unified document processing pipeline that handles all these sources.

### Pipeline Design:

| Stage | Description | Implementation |
|-------|-------------|----------------|
| 1. Ingestion | | |
| 2. Extraction | | |
| 3. Cleaning | | |
| 4. Chunking | | |
| 5. Embedding | | |
| 6. Storage | | |

---

## Problem 2: Chunking Strategy Analysis

Compare different chunking strategies for a 10,000-word technical article about distributed systems:

1. Fixed-size (512 tokens, 50 overlap)
2. Semantic (by headers and paragraphs)
3. Recursive (by \n\n, \n, sentence)

For each strategy, estimate:
- Number of chunks produced
- Coherence (how well chunks maintain context)
- Implementation complexity

### Comparison:

| Strategy | Est. Chunks | Coherence | Complexity |
|----------|-------------|-----------|------------|
| | | | |
| | | | |
| | | | |

---

## Problem 3: Token Budget Optimization

Your LLM has a 4096 token context window. You want to retrieve up to 10 documents but must respect token limits.

Given:
- System prompt: 150 tokens
- User query: 50 tokens
- Each retrieved chunk: ~200 tokens average
- Retrieved results: 10 chunks = 2000 tokens

**Problem**: You can't fit all 10 chunks. What's your strategy?

### Strategy Options:

| Strategy | Description | Pros | Cons |
|----------|-------------|------|------|
| A. Truncate | | | |
| B. Summarize | | | |
| C. Select top-N | | | |
| D. Prority + truncate | | | |

Recommended strategy: ____________

---

## Problem 4: RAG Pipeline Implementation

Write Go code for a simple RAG pipeline using the AI Agent Platform's types:

```go
// Using these interfaces from pkg/types/vector.go
type VectorStore interface {
    Store(id string, vector []float64, metadata map[string]interface{}) error
    Search(ctx context.Context, queryVector []float64, limit int) ([]SearchResult, error)
    Delete(id string) error
}

type Embedder interface {
    Embed(text string) ([]float64, error)
    EmbedBatch(texts []string) ([][]float64, error)
    Dimensions() int
}

// Implement RAG pipeline
type SimpleRAG struct {
    vectorStore VectorStore
    embedder    Embedder
    topK        int
}

func NewSimpleRAG(vs VectorStore, e Embedder, k int) *SimpleRAG {
    // Your code here
}

func (r *SimpleRAG) Index(ctx context.Context, id, content string) error {
    // Your code here
}

func (r *SimpleRAG) Search(ctx context.Context, query string) ([]string, error) {
    // Your code here
}
```

---

## Problem 5: Metadata Schema Design

Design metadata for a multi-tenant documentation system where:
- Each tenant has their own documents
- Documents have version and language
- Some documents require access control

### Metadata Schema:

```json
{
  "document": {
    "id": "...",
    "metadata": {
      // Your fields here
    }
  }
}
```

List at least 8 metadata fields with their types and use cases.

---

## Problem 6: Context Assembly

Given these retrieved document chunks:

```
Chunk 1: "The profit margin is calculated as (Revenue - Cost) / Revenue..."
Chunk 2: "Profit margin = Net Income / Total Revenue..."
Chunk 3: "This formula applies only to companies in the tech sector..."
Chunk 4: "For non-tech companies, see the financial accounting standards..."
Chunk 5: "Historical data shows margins between 15-25% for industry leaders..."
```

Query: "How do I calculate profit margin?"

1. Which chunks are most relevant? Order them.
2. If token limit allows only 2 chunks, which do you choose?
3. Write the assembled prompt.

### Relevance Ranking:

| Rank | Chunk | Reason |
|------|-------|--------|
| 1 | | |
| 2 | | |
| 3 | | |
| 4 | | |
| 5 | | |

### Assembled Prompt:

```
(Write your prompt here)

```

---

## Problem 7: Error Handling in RAG

What could go wrong at each stage of the RAG pipeline? List failures and recovery strategies:

| Stage | Potential Failures | Recovery Strategy |
|-------|--------------------|--------------------|
| Ingestion | | |
| Embedding | | |
| Vector search | | |
| Context assembly | | |
| Generation | | |

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.