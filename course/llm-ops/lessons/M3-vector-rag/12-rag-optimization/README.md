# IO-12: RAG Optimization

**Duration**: 3 hours
**Module**: 3 - Vector Databases & RAG

## Learning Objectives

- Implement hybrid search combining keyword and vector search
- Apply re-ranking to improve retrieval quality
- Use query transformation techniques
- Evaluate RAG systems with proper metrics
- Optimize RAG performance for production workloads

## Beyond Basic Retrieval

Basic RAG (vector search only) has limitations:
- Poor with rare terms or specific keywords
- Doesn't understand exact terminology
- Can miss exact matches

Optimization techniques address these gaps.

## Hybrid Search

### Concept

Hybrid search combines:
1. **Semantic search**: Vector similarity for meaning
2. **Keyword search**: BM25/TF-IDF for exact matching

```
┌─────────────────────────────────────────────────────────────────┐
│                    Hybrid Search Architecture                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Query: "API authentication error 401"                          │
│                    ↓                                             │
│  ┌──────────────────┐    ┌──────────────────┐                  │
│  │  Vector Search   │    │  Keyword Search  │                  │
│  │  (semantic)       │    │  (BM25)           │                  │
│  │                   │    │                   │                  │
│  │  Results:         │    │  Results:         │                  │
│  │  - /auth/errors   │    │  - 401 unauthorized│                │
│  │  - token validation│    │  - API error codes │               │
│  │  - OAuth flow     │    │  - auth endpoints │                │
│  └──────────────────┘    └──────────────────┘                  │
│                    ↓                                             │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Fusion/Reciprocal Rank                │   │
│  │  Combine and re-rank results                             │   │
│  └─────────────────────────────────────────────────────────┘   │
│                    ↓                                             │
│        Final Results (semantic + keyword)                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Implementation with RRF

Reciprocal Rank Fusion combines rankings:

```python
def reciprocal_rank_fusion(results_list: list[list], k: int = 60) -> list:
    """
    Combine multiple result lists using RRF
    
    RRF score = sum(1 / (k + rank)) for each list
    """
    scores = defaultdict(float)
    
    for results in results_list:
        for rank, doc in enumerate(results, 1):
            scores[doc['id']] += 1 / (k + rank)
    
    # Sort by fused score
    return sorted(scores.items(), key=lambda x: x[1], reverse=True)

# Example
vector_results = [
    {'id': 'doc1', 'score': 0.9},
    {'id': 'doc2', 'score': 0.8},
    {'id': 'doc3', 'score': 0.7},
]

bm25_results = [
    {'id': 'doc3', 'score': 0.95},
    {'id': 'doc1', 'score': 0.6},
    {'id': 'doc4', 'score': 0.5},
]

fused = reciprocal_rank_fusion([vector_results, bm25_results])
```

### Hybrid Search with Elasticsearch/Opensearch

```json
{
  "query": {
    "hybrid": {
      "queries": [
        {
          "match": {
            "content": "authentication"
          }
        },
        {
          "knn": {
            "field": "embedding",
            "query_vector": [0.1, 0.3, ...],
            "k": 10,
            "num_candidates": 100
          }
        }
      ]
    }
  }
}
```

## Re-ranking

### Why Re-rank?

Initial retrieval prioritizes speed over accuracy. Re-ranking applies more sophisticated scoring.

```
Stage 1: Fast vector search (Top-100 candidates)
    ↓
Stage 2: Cross-encoder re-rank (Top-10 final)
```

### Cross-encoder Re-ranking

```python
from sentence_transformers import CrossEncoder

# Cross-encoder for re-ranking
cross_encoder = CrossEncoder('cross-encoder/ms-marco-MiniLM-L-6-v2')

def rerank(query: str, candidates: list[str], top_k: int = 5):
    # Score all candidate-document pairs
    scores = cross_encoder.predict([(query, doc) for doc in candidates])
    
    # Sort by score
    ranked = sorted(zip(candidates, scores), key=lambda x: x[1], reverse=True)
    return ranked[:top_k]
```

### Re-ranking Models

| Model | Speed | Quality | Use Case |
|-------|-------|---------|----------|
| ms-marco-MiniLM-L-6-v2 | Fast | Good | General |
| ms-marco-MiniLM-L-12-v2 | Medium | Better | Production |
| ms-marco-electra-base | Slow | Best | High precision |
| bge-reranker-base | Medium | Better | Multi-language |

## Query Transformation

### Query Expansion

Add related terms to improve recall:

```python
def expand_query(query: str) -> list[str]:
    expansions = [
        # Synonyms
        ("API", ["API", "endpoint", "interface"]),
        ("auth", ["auth", "authentication", "authorization"]),
        ("error", ["error", "exception", "failure"]),
    ]
    
    expanded = [query]
    for term, synonyms in expansions:
        if term.lower() in query.lower():
            for syn in synonyms:
                expanded.append(query.replace(term, syn))
    
    return expanded
```

### Query Rewriting

Use LLM to rewrite query for better retrieval:

```python
QUERY_REWRITE_PROMPT = """Given this user query, rewrite it to be more effective for retrieving relevant documents from a knowledge base.

Original: {query}

Rewrite to be:
- More specific and clear
- Use technical terms from the domain
- Include key concepts
- Remove ambiguity

Rewritten:"""

def rewrite_query(query: str, llm) -> str:
    prompt = QUERY_REWRITE_PROMPT.format(query=query)
    return llm.generate(prompt)
```

### HyDE (Hypothetical Document Embeddings)

Generate a hypothetical answer and embed it:

```python
def hyde_search(query: str, llm, vector_store, embedder) -> list:
    # 1. Generate hypothetical answer
    hypothetical = llm.generate(
        f"Give a brief answer to: {query}"
    )
    
    # 2. Embed both query and hypothetical answer
    query_vec = embedder.embed(query)
    hypo_vec = embedder.embed(hypothetical)
    
    # 3. Combine embeddings (or use hypothetical alone)
    combined = average([query_vec, hypo_vec])
    
    # 4. Search with combined embedding
    return vector_store.search(combined)
```

## RAG Evaluation Metrics

### Retrieval Metrics

| Metric | Description | Formula |
|--------|-------------|---------|
| Precision@K | Fraction of top-K that are relevant | #relevant_in_K / K |
| Recall@K | Fraction of relevant found in top-K | #relevant_in_K / #total_relevant |
| MRR | Mean reciprocal rank of first relevant | avg(1/rank_1st_relevant) |
| NDCG | Normalized discounted cumulative gain | See standard formula |

### Generation Metrics

| Metric | Description |
|--------|-------------|
| Faithfulness | Does answer match retrieved context? |
| Answer relevance | Does answer address the question? |
| Context precision | Is retrieved context high quality? |

### RAGAS Score

RAGAS (RAG Assessment) combines retrieval and generation:

```python
from ragas import evaluate
from ragas.metrics import faithfulness, answer_relevancy, context_precision

eval_dataset = [
    {
        "question": "How do I calculate profit margin?",
        "answer": "Profit margin is calculated as (Revenue - Cost) / Revenue.",
        "contexts": ["Profit margin = Net Income / Revenue..."],
        "ground_truth": "Profit margin = (Revenue - Cost) / Revenue"
    }
]

results = evaluate(eval_dataset, metrics=[faithfulness, answer_relevancy, context_precision])
```

## Query Optimization Strategies

### Query Understanding

```python
# Analyze query to choose retrieval strategy
def select_retrieval_strategy(query: str) -> str:
    keywords = extract_keywords(query)
    is_specific = len(keywords) > 3
    
    if is_specific:
        return "hybrid"  # Use both semantic + keyword
    else:
        return "semantic"  # Pure vector search
```

### Adaptive Top-K

```python
def adaptive_top_k(query: str, base_k: int = 5) -> int:
    # Longer, complex queries need more context
    if len(query.split()) > 20:
        return base_k * 2
    # Technical queries benefit from more candidates
    elif any(term in query.lower() for term in ['error', 'debug', 'issue']):
        return base_k + 3
    else:
        return base_k
```

## Optimization in AI Agent Platform

The AI Agent Platform's RAG system implements several optimization patterns:

### Configurable Retrieval

```go
type Config struct {
    Embedder          types.Embedder
    VectorStore       types.VectorStore
    TopK              int
    MaxContextTokens  int
    EnableHybrid      bool  // Toggle hybrid search
    EnableRerank      bool  // Toggle reranking
}
```

### Metrics Collection

```go
// Track retrieval performance
func (r *RAGResult) GetMetrics() RAGMetrics {
    return RAGMetrics{
        RetrievedDocs:  len(r.Results),
        EmbeddingTime:  r.Metadata.EmbeddingTime,
        SearchTime:     r.Metadata.SearchTime,
        TotalTime:     r.Metadata.TotalTime,
        AvgScore:      calculateAvgScore(r.Results),
    }
}
```

## Key Terminology

| Term | Definition |
|------|------------|
| **Hybrid Search** | Combining vector and keyword search |
| **BM25** | Classic keyword ranking algorithm |
| **RRF** | Reciprocal Rank Fusion - combining rankings |
| **Re-ranking** | Post-retrieval scoring with better model |
| **Cross-encoder** | Scorer that processes query-doc pairs together |
| **Query expansion** | Adding synonyms to improve recall |
| **Query rewriting** | LLM-based query improvement |
| **HyDE** | Using hypothetical documents for retrieval |
| **RAGAS** | RAG Assessment - evaluation framework |

## Exercise

### Exercise 12.1: Hybrid Search Implementation

Implement a hybrid search function that:
1. Performs vector search (Top-50)
2. Performs BM25 search (Top-50)
3. Combines results using RRF (k=60)
4. Returns Top-10 final results

```python
# Write your implementation

def hybrid_search(query: str, vector_store, bm25_index, top_k: int = 10):
    # Your code here
    pass

```

---

### Exercise 12.2: Re-ranking Pipeline

Given these initial search results:
```
1. "Error 401: The request was rejected because..."
2. "API Authentication: Best practices guide..."
3. "OAuth 2.0 implementation tutorial..."
4. "HTTP status codes cheat sheet..."
5. "Debugging authentication issues in production..."
```

Query: "how to fix 401 authentication error in REST API"

1. Which documents should be re-ranked to the top?
2. Implement a cross-encoder re-ranking function

### Your Ranking:

| Original Rank | After Re-rank | Reason |
|---------------|---------------|--------|
| | | |
| | | |
| | | |
| | | |
| | | |

---

### Exercise 12.3: Evaluation Metrics

You have a test set of 100 queries with ground truth. Calculate metrics:

| Query | Retrieved (Top-5) | Relevant | Precision@5 | Is Top-1 Relevant? |
|-------|-------------------|----------|--------------|--------------------|
| 1 | [A,B,C,D,E] | [A,B,C] | | Yes |
| 2 | [F,G,H,I,J] | [F,G] | | No |
| 3 | [K,L,M,N,O] | [K,L,M,N] | | Yes |
| ... | ... | ... | | |

Calculate:
- Average Precision@5
- Recall@5
- MRR

### Metrics:

| Metric | Value |
|--------|-------|
| Precision@5 | |
| Recall@5 | |
| MRR | |

---

### Exercise 12.4: Query Transformation

Transform these queries using query rewriting:

| Original Query | Expanded Query | Rationale |
|----------------|-----------------|------------|
| "how to setup auth" | | |
| "api failing" | | |
| "deployment error" | | |

---

### Exercise 12.5: Performance Optimization

For a production RAG system:
- 1M documents
- 10k queries/day
- P99 latency target: 200ms

What optimizations would you implement?

| Optimization | Expected Impact | Implementation Complexity |
|--------------|-----------------|--------------------------|
| | | |
| | | |
| | | |
| | | |

---

### Exercise 12.6: Debug Retrieval Issues

Users report: "The search is returning irrelevant results for technical queries"

What diagnostics would you run? What's the likely cause and fix?

### Diagnostic Steps:

1. 
2. 
3. 

### Likely Cause:

```

```

### Fix:

```

```

---

## Key Takeaways

- ✅ Hybrid search combines semantic and keyword matching for better results
- ✅ Re-ranking improves precision but adds latency
- ✅ Query transformation (rewrite, expand) improves recall
- ✅ RAGAS provides comprehensive evaluation metrics
- ✅ Production RAG requires careful optimization

## Module Summary

This module covered:
- Vector database fundamentals and selection
- Embedding models and their trade-offs
- End-to-end RAG architecture
- Optimization techniques for production

## Additional Resources

- [Hybrid Search Blog](https://www.elastic.co/blog/hybrid-search-llm)
- [RAGAS Documentation](https://docs.ragas.io)
- [HyDE Paper](https://arxiv.org/abs/2212.10496)
- [BGE Reranker](https://github.com/FlagOpen/FlagEmbedding)