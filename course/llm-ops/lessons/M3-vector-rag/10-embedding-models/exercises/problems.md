# Exercise: Embedding Models

## Problem 1: Model Benchmark Analysis

Analyze the MTEB (Massive Text Embedding Benchmark) leaderboard data for the following models and complete the table:

| Model | MTEB Score | Dimensions | Speed (relative) | Best Use Case |
|-------|------------|------------|-----------------|---------------|
| ada-002 | 61.0 | | | |
| bge-small-en-v1.5 | | 384 | | |
| bge-base-en-v1.5 | | | | |
| text-embedding-3-large | | | Very slow | |
| all-MiniLM-L6-v2 | | | Fast | |

Research and fill in the missing values.

---

## Problem 2: Embedding Model Implementation

Write a Python script that:
1. Loads the `BAAI/bge-small-en-v1.5` model
2. Generates embeddings for a list of documents
3. Computes pairwise cosine similarity
4. Returns the most similar document pairs

```python
# Write your solution here
import torch
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

# Solution

```

---

## Problem 3: Batch Processing Optimization

You need to embed 10,000 documents. The naive approach processes one at a time, taking 50ms each.

1. How long would naive processing take?
2. If you batch 32 documents per API call, how long would it take (assuming 100ms per batch)?
3. What is the speedup ratio?

### Your Analysis:

| Approach | Documents | Time per Doc | Total Time |
|----------|-----------|--------------|------------|
| Naive | 10,000 | 50ms | |
| Batched (32) | 10,000 | ~3ms/batch | |

Speedup: ____________

---

## Problem 4: Dimension Reduction Trade-offs

You're considering reducing 1536-dimensional embeddings to 384 dimensions for storage efficiency.

1. What techniques could you use?
2. What information might be lost?
3. How would you measure the impact on search quality?

### Your Answers:

```
(Write your answers here)

```

---

## Problem 5: Cost Analysis for Production

Calculate the monthly embedding costs for a RAG system with:
- 50,000 documents indexed once
- 100,000 user queries per month
- Average document length: 500 tokens
- Average query length: 50 tokens

Compare:
1. OpenAI text-embedding-3-small ($0.02/million tokens)
2. Self-hosted on GPU (assume A100 GPU, $1/hour, 500 docs/minute indexing)

### Cost Analysis:

| Cost Category | OpenAI | Self-hosted |
|---------------|--------|-------------|
| Indexing (50k docs × 500 tokens) | | |
| Queries (100k × 50 tokens) | | |
| GPU hours | N/A | |
| Total/month | | |

---

## Problem 6: Fine-tuning Decision

A healthcare company wants to use embeddings for medical document search.

Questions to evaluate if they need fine-tuning:
1. How many medical-specific terms exist that general models don't understand?
2. Do they have labeled training data (query → relevant document pairs)?
3. What's the current accuracy? Is it acceptable?

For each scenario, recommend whether to fine-tune:
- 70% accuracy with general model, limited training data
- 65% accuracy with general model, 10k labeled pairs
- 85% accuracy with general model, no training data

### Recommendations:

| Scenario | Fine-tune? | Justification |
|----------|-----------|---------------|
| | | |
| | | |
| | | |

---

## Problem 7: Embedder Interface Usage

Using the AI Agent Platform's embedder interface from `pkg/types/vector.go`, write Go code that:

1. Creates an embedder with the correct dimensions
2. Embeds a batch of texts
3. Handles errors appropriately

```go
// Reference interface
type Embedder interface {
    Embed(text string) ([]float64, error)
    EmbedBatch(texts []string) ([][]float64, error)
    Dimensions() int
}

// Your implementation

```

---

## Submission

Save your answers in a file and be prepared to discuss them in the next lesson.