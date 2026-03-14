# IO-09: Vector Database Fundamentals

**Duration**: 3 hours
**Module**: 3 - Vector Databases & RAG

## Learning Objectives

- Understand vector embeddings and their role in AI applications
- Compare vector databases with traditional databases
- Implement similarity search with various algorithms
- Evaluate and select vector databases for production use cases
- Connect vector databases to AI Agent Platform's RAG implementation

## What Are Vector Databases?

Vector databases are specialized database systems designed to store, index, and query high-dimensional vector embeddings. Unlike traditional databases that handle structured data (rows and columns), vector databases excel at similarity search over unstructured data.

### The Vector Embedding Concept

When we convert text, images, or any data into numerical vectors, we capture semantic meaning in a format that computers can compare mathematically:

```
Text: "The stock price increased significantly"
  → Vector: [0.12, -0.34, 0.56, 0.89, -0.23, ...]  (1536 dimensions for ada-002)

Text: "Revenue grew a lot"
  → Vector: [0.11, -0.31, 0.58, 0.85, -0.25, ...]  (similar!)

Text: "The weather is sunny"
  → Vector: [-0.45, 0.23, -0.12, 0.78, 0.34, ...]  (different!)
```

### Why Traditional Databases Don't Work

| Feature | Traditional DB | Vector Database |
|---------|---------------|-----------------|
| Query Type | Exact match, range | Similarity search |
| Indexing | B-tree, hash | HNSW, IVF, PQ |
| Scalability | Linear with data | Sub-linear (approximate) |
| Distance Metrics | SQL operators | Cosine, Euclidean, dot product |

## Vector Similarity Search

### Distance Metrics

Vector databases use different similarity metrics depending on the embedding model and use case:

1. **Cosine Similarity** - Measures angle between vectors (most common for text)
   ```
   cosine(A, B) = (A · B) / (||A|| × ||B||)
   ```

2. **Euclidean Distance** - Straight-line distance between points
   ```
   euclidean(A, B) = √(Σ(Ai - Bi)²)
   ```

3. **Dot Product** - Projection of one vector onto another
   ```
   dot(A, B) = Σ(Ai × Bi)
   ```

### Indexing Strategies

Searching through millions of vectors requires efficient indexing:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Vector Search Process                         │
├─────────────────────────────────────────────────────────────────┤
│  Query Vector → Approximate Nearest Neighbor Search → Top-K   │
│                      ↓                                          │
│              Index (HNSW/IVF)                                   │
│                      ↓                                          │
│         [Search time: O(1) vs O(n)]                            │
└─────────────────────────────────────────────────────────────────┘
```

**HNSW (Hierarchical Navigable Small World)**:
- Builds a navigable graph structure
- Excellent for high accuracy with low latency
- Memory-intensive but fast
- Used by Qdrant, Weaviate, pgvector

**IVF (Inverted File Index)**:
- Clusters vectors into partitions
- Searches subset of clusters
- Balance between speed and accuracy
- Good for very large datasets

**Product Quantization (PQ)**:
- Compresses vectors into smaller codes
- Reduces memory footprint dramatically
- Enables billion-scale search
- Often combined with HNSW

## Vector Databases Overview

### pgvector

Open-source extension for PostgreSQL, perfect for teams already using PostgreSQL:

```sql
-- Create vector column
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    content TEXT,
    embedding vector(1536)
);

-- Create HNSW index
CREATE INDEX ON documents 
USING hnsw (embedding vector_cosine_ops);

-- Similarity search
SELECT id, content, 
       1 - (embedding <=> query_embedding) as similarity
FROM documents
ORDER BY embedding <=> query_embedding
LIMIT 5;
```

**Pros**: PostgreSQL ecosystem, ACID compliance, SQL interface
**Cons**: Not specialized, limited advanced features
**Use Cases**: Small-medium datasets, existing PostgreSQL users

### Qdrant

Purpose-built cloud-native vector database with excellent performance:

```go
// Qdrant client usage
client := qdrant.NewClient("localhost:6333")

// Search
results, err := client.Search(ctx, &qdrant.SearchParams{
    CollectionName: "documents",
    Vector:        queryVector,
    Limit:         5,
    WithPayload:   true,
})
```

**Pros**: High performance, cloud-native, REST + gRPC APIs
**Cons**: Newer ecosystem, fewer integrations
**Use Cases**: Production AI apps, real-time search

### Weaviate

Multi-model vector database with graph connections:

```python
import weaviate

client = weaviate.Client("http://localhost:8080")

# Search with filters
result = client.query.get(
    "Document",
    ["title", "content", "category"]
).with_near_text({
    "concepts": ["financial report"]
}).with_limit(5).do()
```

**Pros**: Multi-model support, GraphQL, rich filters
**Cons**: Resource-intensive
**Use Cases**: Multi-modal search, knowledge graphs

### Pinecone

Managed vector database with excellent scalability:

**Pros**: Fully managed, excellent scaling, hybrid search
**Cons**: Vendor lock-in, costs at scale
**Use Cases**: Enterprise applications, large-scale RAG

### Milvus

Open-source, cloud-native with powerful features:

**Pros**: Open-source, rich features, strong community
**Cons**: Complex setup
**Use Cases**: Large-scale AI applications, multi-modal

## Comparison Matrix

| Feature | pgvector | Qdrant | Weaviate | Pinecone | Milvus |
|---------|----------|--------|----------|----------|--------|
| **Type** | Extension | Native | Native | Managed | Native |
| **Deployment** | Self-hosted | Self/Cloud | Self/Cloud | Cloud | Self/Cloud |
| **Index Types** | HNSW, IVF | HNSW, PQ | HNSW | HNSW | HNSW, IVF, PQ |
| **Scale** | 10M vectors | 1B+ vectors | 1B+ vectors | 1B+ vectors | 1B+ vectors |
| **Latency** | Low | Very Low | Low | Low | Low |
| **Cloud** | No | Yes | Yes | Yes | Yes |
| **Cost** | Low | Medium | Medium | High | Medium |

## Vector Database in AI Agent Platform

The AI Agent Platform demonstrates production vector database usage:

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    AI Agent Platform                             │
├─────────────────────────────────────────────────────────────────┤
│  Vector Store                                                    │
│  ┌─────────────┬─────────────┬─────────────┐                   │
│  │   pgvector  │   Qdrant    │   Memory    │                   │
│  │  (Storage)  │ (Optional)  │ (Dev/Test)  │                   │
│  └─────────────┴─────────────┴─────────────┘                   │
│                                                                  │
│  RAG Pipelines                                                   │
│  ┌──────────────────┬──────────────────┐                         │
│  │  Code RAG        │  Document RAG   │                         │
│  │  - Code indexing │  - Doc indexing │                         │
│  │  - Snippet search│  - Semantic search│                        │
│  └──────────────────┴──────────────────┘                         │
└─────────────────────────────────────────────────────────────────┘
```

### Implementation (services/ai-agent-platform/internal/rag/pipeline.go)

```go
// RAG Pipeline uses vector store for semantic search
type DocumentRAGPipeline struct {
    config        *Config
    vectorStore   types.VectorStore
    embedder      types.Embedder
}

// Search retrieves relevant documents
func (r *DocumentRAGPipeline) Search(ctx context.Context, query string) (*RAGResult, error) {
    // Generate query embedding
    queryVector, err := r.config.Embedder.Embed(query)
    if err != nil {
        return nil, err
    }
    
    // Search vector store
    results, err := r.vectorStore.Search(ctx, queryVector, r.config.TopK)
    if err != nil {
        return nil, err
    }
    
    return r.assembleContext(results), nil
}
```

### Docker Setup (services/ai-agent-platform/docker-compose.yml)

```yaml
# PostgreSQL with pgvector
postgres:
  image: ankane/pgvector:latest
  environment:
    POSTGRES_DB: finagent

# Qdrant (optional, for production)
qdrant:
  image: qdrant/qdrant:latest
  ports:
    - "6333:6333"
```

## Key Terminology

| Term | Definition |
|------|------------|
| **Vector Embedding** | Numerical representation of data capturing semantic meaning |
| **Similarity Search** | Finding vectors most similar to a query vector |
| **HNSW** | Hierarchical Navigable Small World - graph-based indexing |
| **IVF** | Inverted File Index - clustering-based indexing |
| **Product Quantization** | Vector compression technique for large-scale search |
| **Top-K Retrieval** | Finding K most similar vectors |
| **Vector Index** | Data structure enabling efficient similarity search |

## Exercise

### Exercise 9.1: Vector Database Selection

A company is building a legal document search system with the following requirements:
- 500,000 legal documents
- 100 concurrent users
- Must support semantic search ("find contracts about liability clauses")
- Budget: $500/month
- Team has PostgreSQL expertise

Which vector database would you recommend and why? Create a comparison table.

### Exercise 9.2: Implement Vector Operations

Write pseudocode for the following operations using a vector database:
1. Store a document with embedding
2. Search for similar documents
3. Update a document's embedding
4. Delete a document

### Exercise 9.3: Analyze Index Performance

Given the following scenarios, recommend the best indexing strategy:

| Dataset Size | Latency Requirement | Accuracy Requirement |
|--------------|---------------------|----------------------|
| 10,000 vectors | < 10ms | 95% |
| 10M vectors | < 100ms | 90% |
| 100M vectors | < 200ms | 85% |

For each scenario, explain: index type, why, and trade-offs.

### Exercise 9.4: Connect to AI Agent Platform

The AI Agent Platform has a `MemoryVectorStore` implementation. Research how it works and write code to:
1. Initialize a vector store with 1536-dimensional embeddings
2. Store a document with metadata
3. Search for similar documents using cosine similarity

Use the implementation at: `services/ai-agent-platform/internal/vectorstore/memory.go`

## Key Takeaways

- ✅ Vector databases enable semantic search over embeddings
- ✅ HNSW and IVF are primary indexing strategies with different trade-offs
- ✅ pgvector is great for existing PostgreSQL users
- ✅ Qdrant offers excellent performance for production workloads
- ✅ The AI Agent Platform demonstrates real-world vector store implementation

## Next Steps

→ [IO-10: Embedding Models](../10-embedding-models/README.md)

## Additional Resources

- [pgvector Documentation](https://github.com/pgvector/pgvector)
- [Qdrant Documentation](https://qdrant.tech/documentation/)
- [Weaviate Documentation](https://weaviate.io/developers/weaviate)
- [Vector Search Algorithms](https://arxiv.org/abs/2110.08499)