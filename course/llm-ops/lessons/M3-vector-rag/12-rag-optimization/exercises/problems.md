# Exercise: RAG Optimization

## Problem 1: RRF Implementation

Implement the Reciprocal Rank Fusion algorithm:

```python
from collections import defaultdict

def reciprocal_rank_fusion(results_list: list[list], k: int = 60) -> list[tuple]:
    """
    Combine multiple ranked result lists using RRF
    
    Args:
        results_list: List of result lists, each containing document IDs in ranking order
        k: Constant to prevent large rankings from dominating
    
    Returns:
        List of (doc_id, score) tuples sorted by fused score
    """
    # Your implementation here
    pass

# Test with sample data
results_a = ['doc1', 'doc2', 'doc3', 'doc4']
results_b = ['doc3', 'doc1', 'doc4', 'doc2']
results_c = ['doc2', 'doc3', 'doc1', 'doc4']

fused = reciprocal_rank_fusion([results_a, results_b, results_c])
print(fused)  # Expected: [('doc1', ...), ('doc3', ...), ('doc2', ...), ('doc4', ...)]
```

---

## Problem 2: Re-ranking Pipeline Analysis

You have two retrieval stages:
- Stage 1: Vector search (returns 50 candidates in 50ms)
- Stage 2: Cross-encoder re-ranking (scores 50 docs in 500ms)

If you want P99 latency under 600ms, what strategies could you use?

| Strategy | Description | Impact |
|----------|-------------|--------|
| | | |
| | | |
| | | |

---

## Problem 3: Query Expansion

Create a query expansion function for technical documentation:

```python
EXPANSIONS = {
    'api': ['API', 'endpoint', 'REST', 'interface'],
    'auth': ['authentication', 'authorization', 'OAuth', 'JWT'],
    'error': ['error', 'exception', 'failure', 'issue', 'bug'],
    'deploy': ['deployment', 'deploy', 'release', 'push'],
    'config': ['configuration', 'settings', 'setup', 'options'],
}

def expand_query(query: str) -> list[str]:
    """
    Expand query with synonyms and related terms
    """
    # Your implementation here
    pass

# Test
print(expand_query("api auth error"))
# Expected: ["api auth error", "API auth error", "endpoint authentication error", ...]
```

---

## Problem 4: Evaluation Metrics

Calculate metrics for a test set:

Given:
- 20 queries tested
- 100 total relevant documents across all queries
- Sum of (relevant_found / total_relevant) = 45
- First relevant position sum = 8

### Calculate:

| Metric | Formula | Value |
|--------|---------|-------|
| Average Recall@10 | sum(relevant_found / total_relevant) / queries | |
| MRR | sum(1 / first_relevant_position) / queries | |

---

## Problem 5: Hybrid Search Design

Design a hybrid search for a system with:
- 500k documents
- Mix of structured (metadata) and unstructured (content) data
- Need for exact keyword matches (part numbers, error codes)

### Architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Query                                     │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────┬─────────────────────────┐
│    Vector Search        │    BM25/Keyword         │
│                         │                         │
│    - Semantic matching  │    - Exact matching     │
│    - Top-100            │    - Top-100            │
└─────────────────────────┴─────────────────────────┘
                              ↓
                    ┌─────────────────┐
                    │  RRF Fusion     │
                    │  k=60           │
                    └─────────────────┘
                              ↓
                    ┌─────────────────┐
                    │  Re-ranker      │
                    │  (optional)     │
                    └─────────────────┘
                              ↓
                    ┌─────────────────┐
                    │  Top-10 Results │
                    └─────────────────┘
```

Explain each component's role and why this architecture works.

---

## Problem 6: Performance Tuning

A RAG system is too slow:
- Vector search: 100ms
- Re-ranking 50 docs: 400ms
- Context assembly: 50ms
- Total: 550ms (target: 300ms)

### Optimizations:

| Component | Current | Optimization | Expected |
|-----------|---------|--------------|----------|
| Vector search | 100ms | | |
| Re-ranking | 400ms | | |
| Context assembly | 50ms | | |
| **Total** | **550ms** | | **300ms** |

---

## Problem 7: Debugging Retrieval

Users report that search for specific error codes returns wrong results:
- Query: "error 503"
- Returns: "500 Internal Server Error" instead of "503 Service Unavailable"

What could be wrong? How would you fix it?

### Root Cause Analysis:

```

```

### Fix:

```

```

---

## Problem 8: RAGAS Integration

Write code to calculate RAGAS metrics:

```python
from ragas import evaluate
from ragas.metrics import faithfulness, answer_relevancy, context_precision

# Sample eval data
eval_data = [
    {
        "question": "What is profit margin?",
        "answer": "Profit margin is revenue minus costs divided by revenue.",
        "contexts": ["Profit margin = (Revenue - Cost) / Revenue"],
        "ground_truth": "(Revenue - Cost) / Revenue"
    },
    {
        "question": "How does OAuth work?",
        "answer": "OAuth uses access tokens to authorize third-party access.",
        "contexts": ["OAuth 2.0 uses access tokens..."],
        "ground_truth": "Authorization framework using access tokens"
    }
]

# Your code here
def evaluate_rag(eval_data):
    # Calculate faithfulness, answer_relevancy, context_precision
    pass
```

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.