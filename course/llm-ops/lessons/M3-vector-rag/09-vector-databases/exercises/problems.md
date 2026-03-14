# Exercise: Vector Database Fundamentals

## Problem 1: Vector Database Comparison

Complete the comparison table for the following vector databases:

| Database | Deployment Model | Best Use Case | Max Vectors | Key Limitation |
|----------|------------------|---------------|-------------|----------------|
| pgvector | | | | |
| Qdrant | | | | |
| Weaviate | | | | |
| Pinecone | | | | |
| Milvus | | | | |

---

## Problem 2: Index Selection

A healthcare company needs to store 5 million medical records as embeddings. They require:
- 99% accuracy
- < 50ms query latency
- Ability to filter by patient demographics

What index type would you recommend? Explain your reasoning.

### Your Answer:

```
(Write your answer here)

```

---

## Problem 3: Distance Metric Selection

For each use case, recommend the best distance metric (cosine, euclidean, dot product):

1. **Customer service chatbot**: Looking for similar past conversations
2. **Recommendation system**: Finding items similar to user preferences
3. **Document classification**: Categorizing documents by topic

### Your Answers:

| Use Case | Metric | Justification |
|----------|--------|---------------|
| | | |
| | | |
| | | |

---

## Problem 4: Implement Vector Search

Write Go code that:
1. Initializes a vector store
2. Stores 3 documents with embeddings (using mock vectors)
3. Searches for the most similar document

Use the AI Agent Platform's vector store interface as reference:

```go
// Reference from services/ai-agent-platform/pkg/types/vector.go
type VectorStore interface {
    Store(id string, vector []float64, metadata map[string]interface{}) error
    Search(ctx context.Context, queryVector []float64, limit int) ([]SearchResult, error)
    Delete(id string) error
    Get(id string) (VectorEntry, error)
    Count() int
}
```

### Your Implementation:

```go
// Write your code here

```

---

## Problem 5: Database Migration

A company currently uses Elasticsearch for keyword search. They want to add semantic search with vectors.

Option A: Add pgvector to existing PostgreSQL
Option B: Deploy Qdrant as separate service
Option C: Migrate everything to Pinecone

Compare the three approaches considering:
- Implementation effort
- Maintenance overhead
- Cost implications
- Performance

### Your Analysis:

| Factor | Option A (pgvector) | Option B (Qdrant) | Option C (Pinecone) |
|--------|---------------------|-------------------|---------------------|
| Implementation | | | |
| Maintenance | | | |
| Cost | | | |
| Performance | | | |

Recommended: ____________

---

## Problem 6: Scaling Analysis

The AI Agent Platform's vector store currently handles 10,000 documents. What happens to each component when scaling to:

1. **1 million documents** - What changes needed?
2. **100 million documents** - Architecture evolution?

### Scenario 1 - 1M Documents:

```

```

### Scenario 2 - 100M Documents:

```

```

---

## Problem 7: Debugging Vector Search

Users report that semantic search is returning irrelevant results for queries with negations (e.g., "find documents about risks but NOT compliance issues").

What could be causing this? How would you diagnose and fix it?

### Diagnosis:

```

```

### Fix:

```

```

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.