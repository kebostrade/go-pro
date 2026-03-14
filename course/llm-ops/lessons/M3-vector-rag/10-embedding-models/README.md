# IO-10: Embedding Models

**Duration**: 2 hours
**Module**: 3 - Vector Databases & RAG

## Learning Objectives

- Understand how text embeddings are generated from transformer models
- Compare embedding models by dimension, performance, and use case
- Select appropriate embedding models for different AI applications
- Implement embedding generation in production pipelines
- Fine-tune embeddings for domain-specific tasks

## What Are Embedding Models?

Embedding models are specialized neural networks that convert text into dense vector representations. Unlike sparse representations (like TF-IDF), dense embeddings capture semantic meaning in a compact, fixed-dimensional space.

### How Text Embeddings Work

```
Input Text: "How do I calculate profit margin?"

┌─────────────────────────────────────────────────────────┐
│                  Transformer Model                       │
│  (BERT, SBERT, ada-002, etc.)                           │
├─────────────────────────────────────────────────────────┤
│  Tokenize → Encode → Pool → Dense Vector                │
└─────────────────────────────────────────────────────────┘
                  ↓
        Vector: [0.12, -0.34, 0.56, 0.89, -0.23, ...]
                  ↓ (1536 dims for ada-002)
```

### Key Components

1. **Tokenization**: Breaking text into tokens (words, subwords, characters)
2. **Encoding**: Passing tokens through transformer layers
3. **Pooling**: Combining token representations into single vector (CLS, mean, etc.)
4. **Projection**: Converting to final embedding dimensions

## Popular Embedding Models

### OpenAI Embeddings

| Model | Dimensions | MTEB Score | Cost | Best For |
|-------|------------|------------|------|----------|
| ada-002 | 1536 | 61.0 | Low | General purpose |
| text-embedding-3-small | 1536 | 62.3 | Very Low | Cost optimization |
| text-embedding-3-large | 3072 | 64.6 | High | Highest quality |

```python
# OpenAI Embeddings
from openai import OpenAI
client = OpenAI()

response = client.embeddings.create(
    model="text-embedding-3-small",
    input="Your text here"
)
vector = response.data[0].embedding
```

### Open-Source Models

| Model | Dimensions | MTEB Score | License | Host |
|-------|------------|------------|---------|------|
| sentence-transformers/all-MiniLM-L6-v2 | 384 | 69.2 | Apache 2 | Local |
| BAAI/bge-small-en-v1.5 | 384 | 71.6 | Apache 2 | Local |
| BAAI/bge-base-en-v1.5 | 768 | 72.3 | Apache 2 | Local |
| intfloat/e5-base-v2 | 768 | 72.7 | MIT | Local |

```python
# Sentence Transformers
from sentence_transformers import SentenceTransformer
model = SentenceTransformer('BAAI/bge-small-en-v1.5')
vector = model.encode("Your text here")
```

### Comparison by Use Case

| Use Case | Recommended Model | Reasoning |
|----------|-------------------|-----------|
| Chatbot FAQ | ada-002 or bge-small | Fast, good accuracy |
| Code search | bge-base or specialized | Code understanding |
| Legal docs | text-embedding-3-large | High precision needed |
| Multi-language | bge-m3 or multilingual | Cross-language support |
| Cost-sensitive | text-embedding-3-small | Best price/quality |

## Dimension Trade-offs

### Understanding Dimensions

The embedding dimension determines how much information can be captured:

```
Low dimension (64-128):    Fast, memory-efficient, but loses nuance
Medium (256-768):         Balance of speed and quality  
High (1024-4096):         Rich representation, higher cost
```

### Trade-off Analysis

| Dimension | Storage | Speed | Accuracy | Use Case |
|-----------|---------|-------|----------|----------|
| 384 | 1x | 1x | Baseline | Cost-sensitive |
| 768 | 2x | 0.8x | +5-10% | General purpose |
| 1536 | 4x | 0.6x | +10-15% | High precision |
| 3072 | 8x | 0.4x | +15-20% | Research/enterprise |

### Dimensionality Reduction

When working with high-dimensional embeddings, consider:

1. **PCA (Principal Component Analysis)**: Linear projection preserving variance
2. **UMAP**: Non-linear, preserves local and global structure
3. **Product Quantization**: Compress while maintaining similarity

## Model Selection Framework

### Decision Tree

```
                    ┌─────────────────────┐
                    │ What's your use?    │
                    └──────────┬──────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         ↓                     ↓                     ↓
   General Search          Code/Domain           Multi-language
         │                     │                     │
         ↓                     ↓                     ↓
   text-embedding-3-small   bge-base-en-v1.5     bge-m3
   or                       or                     or
   bge-small-en-v1.5       code专用模型           paraphrase-multilingual
```

### Key Factors to Consider

1. **MTEB (Massive Text Embedding Benchmark)**: Standardized evaluation
2. **Latency requirements**: Batch vs. real-time
3. **Cost**: API calls vs. infrastructure
4. **Data sensitivity**: On-prem vs. cloud
5. **Language**: English vs. multilingual

## Embedding Models in AI Agent Platform

The AI Agent Platform demonstrates embedding model implementation:

### Embedder Interface (services/ai-agent-platform/pkg/types/vector.go)

```go
// Embedder defines the interface for generating embeddings
type Embedder interface {
    // Embed generates an embedding for text
    Embed(text string) ([]float64, error)
    
    // EmbedBatch generates embeddings for multiple texts
    EmbedBatch(texts []string) ([][]float64, error)
    
    // Dimensions returns the embedding dimension
    Dimensions() int
}
```

### OpenAI Embedder Implementation (services/ai-agent-platform/internal/embeddings/openai.go)

```go
type OpenAIEmbedder struct {
    client *openai.Client
    model  string
}

func (e *OpenAIEmbedder) Embed(text string) ([]float64, error) {
    resp, err := e.client.Embeddings.Create(ctx, &openai.EmbeddingParams{
        Model: e.model,
        Input: text,
    })
    if err != nil {
        return nil, err
    }
    return resp.Data[0].Embedding, nil
}
```

### Mock Embedder for Testing

```go
// NewMockEmbedder creates a deterministic embedder for testing
func NewMockEmbedder(dimensions int) *MockEmbedder {
    return &MockEmbedder{
        dimensions: dimensions,
        cache:      make(map[string][]float64),
    }
}
```

## Fine-tuning Embeddings

### When to Fine-tune

Consider fine-tuning when:
- Domain-specific vocabulary (medical, legal, technical)
- Unique classification tasks
- General models underperform on your data
- You have labeled training data

### Approach: Contrastive Learning

```
Positive pairs: Similar documents → should be close in embedding space
Negative pairs: Different documents → should be far in embedding space

Loss = max(0, margin - similarity(pos) + similarity(neg))
```

### SentenceTransformers Fine-tuning Example

```python
from sentence_transformers import SentenceTransformer, InputExample, losses
from torch.utils.data import DataLoader

# Define training data
train_examples = [
    InputExample(texts=["query1", "positive_doc1"]),
    InputExample(texts=["query2", "positive_doc2"]),
]

# Configure training
model = SentenceTransformer('BAAI/bge-small-en-v1.5')
train_dataloader = DataLoader(train_examples, shuffle=True, batch_size=16)
train_loss = losses.CachedSoftmaxLoss(model)

# Train
model.fit(train_objectives=[(train_dataloader, train_loss)], epochs=3)
```

## Key Terminology

| Term | Definition |
|------|------------|
| **Embedding Model** | Neural network that converts text to vectors |
| **Dimension** | Size of output vector (e.g., 1536) |
| **MTEB** | Massive Text Embedding Benchmark - standard evaluation |
| **Pooling** | Combining token embeddings into single vector |
| **CLS Token** | First token used for sentence representation |
| **Fine-tuning** | Adapting pre-trained model to specific domain |
| **Contrastive Learning** | Training technique using positive/negative pairs |

## Exercise

### Exercise 10.1: Model Selection

A company building a legal document Q&A system needs:
- High accuracy for complex legal queries
- Can afford higher latency
- Budget: $1000/month for embedding API

Which model would you choose? Create a comparison with alternatives.

### Your Answer:

```
(Write your answer here)

```

---

### Exercise 10.2: Dimension Analysis

You're building a semantic search for a startup with 100k documents. Current model produces 1536-dim embeddings but memory is a concern.

1. What are your options to reduce memory?
2. What's the trade-off?
3. Implement a dimension reduction strategy in code

### Your Analysis:

```

```

---

### Exercise 10.3: Embedder Interface Implementation

Using the AI Agent Platform's embedder interface, implement a custom embedder that:
1. Uses a local model (sentence-transformers)
2. Caches embeddings to avoid recomputation
3. Handles batch processing efficiently

Reference: services/ai-agent-platform/internal/embeddings/

### Your Implementation:

```go
// Write your implementation

```

---

### Exercise 10.4: Cost Optimization

You have 1 million documents to embed. Compare the cost of:
1. OpenAI text-embedding-3-small ($0.02/1M tokens)
2. Self-hosted bge-small-en-v1.5 (assume $0.50/hour inference, 100 docs/sec)

Estimate total cost for initial embedding + monthly queries (100k searches).

### Cost Comparison:

| Component | OpenAI | Self-hosted |
|-----------|--------|-------------|
| Initial embedding | | |
| Monthly compute | | |
| Total per month | | |

---

## Key Takeaways

- ✅ Embedding models convert text to semantic vectors
- ✅ Model choice depends on accuracy, latency, and cost requirements
- ✅ Higher dimensions capture more nuance but cost more
- ✅ The AI Agent Platform provides embedder interface for flexibility
- ✅ Fine-tuning can significantly improve domain-specific performance

## Next Steps

→ [IO-11: RAG Architecture](../11-rag-architecture/README.md)

## Additional Resources

- [MTEB Leaderboard](https://huggingface.co/spaces/mteb/leaderboard)
- [OpenAI Embeddings Documentation](https://platform.openai.com/docs/guides/embeddings)
- [Sentence Transformers](https://sbert.net)
- [BGE Model Papers](https://github.com/FlagOpen/FlagEmbedding)