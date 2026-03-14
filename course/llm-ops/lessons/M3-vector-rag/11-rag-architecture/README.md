# IO-11: RAG Architecture

**Duration**: 3 hours
**Module**: 3 - Vector Databases & RAG

## Learning Objectives

- Understand the end-to-end RAG pipeline architecture
- Implement document processing and chunking strategies
- Build retrieval systems that fetch relevant context
- Assemble context for LLM generation
- Connect to AI Agent Platform's RAG implementation

## What is RAG?

Retrieval-Augmented Generation (RAG) combines the power of information retrieval with LLM generation. Instead of relying solely on LLM knowledge, RAG retrieves relevant context from a knowledge base and includes it in the prompt.

```
┌─────────────────────────────────────────────────────────────────┐
│                    RAG Pipeline                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Query: "How do I calculate profit margin?"                     │
│                    ↓                                             │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    RETRIEVAL                             │   │
│  │  1. Generate query embedding                            │   │
│  │  2. Search vector database                              │   │
│  │  3. Fetch top-K relevant documents                      │   │
│  └─────────────────────────────────────────────────────────┘   │
│                    ↓                                             │
│  Retrieved Context:                                             │
│  - Document 1: "Revenue minus costs = profit..."               │
│  - Document 2: "Profit margin formula: (Revenue - Cost)/..."    │
│                    ↓                                             │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                 AUGMENTATION                             │   │
│  │  Assemble context into prompt                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                    ↓                                             │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                 GENERATION                               │   │
│  │  LLM generates answer using retrieved context            │   │
│  └─────────────────────────────────────────────────────────┘   │
│                    ↓                                             │
│  Answer: "Profit margin is calculated by subtracting..."        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Document Processing Pipeline

### Step 1: Document Ingestion

Documents come in various formats and need standardized processing:

```go
// Document types supported
type Document struct {
    ID        string
    Content   string
    Metadata  map[string]interface{}
    Source    string      // pdf, html, md, txt, etc.
}
```

### Step 2: Text Extraction

Extract clean text from various sources:

| Format | Tools | Challenges |
|--------|-------|------------|
| PDF | pdfplumber, PyMuPDF | Tables, images, layout |
| HTML | BeautifulSoup, jsoup | Tags, scripts |
| Markdown | Custom parsers | Headers, code blocks |
| Office | python-docx, openpyxl | Formatting, embedded content |

### Step 3: Cleaning & Normalization

Clean extracted text:

1. Remove extra whitespace
2. Normalize unicode characters
3. Fix encoding issues
4. Remove boilerplate (headers, footers)
5. Handle special characters

```python
def clean_text(text: str) -> str:
    # Remove extra whitespace
    text = re.sub(r'\s+', ' ', text)
    # Remove non-printable characters
    text = ''.join(c for c in text if c.isprintable())
    # Normalize quotes
    text = text.replace('"', '"').replace('"', '"')
    return text.strip()
```

## Chunking Strategies

### Why Chunking Matters

Chunking determines how document content is split into searchable units. Too small = lose context. Too large = dilute relevance.

### Fixed-Size Chunking

Simple approach with overlaps:

```python
def chunk_by_size(text: str, chunk_size: int = 512, overlap: int = 50) -> list[str]:
    tokens = text.split()
    chunks = []
    for i in range(0, len(tokens), chunk_size - overlap):
        chunk = ' '.join(tokens[i:i + chunk_size])
        chunks.append(chunk)
    return chunks
```

**Pros**: Simple, consistent size
**Cons**: May split sentences, lose semantic coherence

### Semantic Chunking

Split at natural boundaries:

```python
def semantic_chunks(text: str) -> list[str]:
    # Split by paragraphs, headers, sentences
    splits = re.split(r'\n\n+|\n(?=#)|(?<=[.!?])\s+', text)
    # Merge small chunks
    return merge_small_chunks(splits, min_size=100)
```

**Pros**: Respects content structure
**Cons**: Variable sizes, more complex

### Recursive Chunking

Hierarchical splitting based on delimiters:

```python
def recursive_chunk(text: str, delimiters: list[str] = ['\n\n', '\n', '. ']) -> list[str]:
    if not delimiters:
        return [text]
    
    delimiter = delimiters[0]
    chunks = text.split(delimiter)
    
    result = []
    for chunk in chunks:
        if len(chunk) > 2000:
            result.extend(recursive_chunk(chunk, delimiters[1:]))
        else:
            result.append(chunk)
    return result
```

### Chunking Comparison

| Strategy | Best For | Limitations |
|----------|----------|-------------|
| Fixed-size | Uniform content | May break语义 |
| Semantic | Well-structured docs | Requires good delimiters |
| Recursive | Mixed content | Complex tuning |
| Markdown-aware | Technical docs | Only for MD |
| Table-aware | Data-heavy docs | Specialized parsing |

### Optimal Chunk Size by Use Case

| Use Case | Recommended Size | Overlap |
|----------|------------------|---------|
| FAQ | 128-256 tokens | 20-50 |
| Long-form articles | 512-1024 tokens | 50-100 |
| Code repositories | 256-512 tokens | 50-100 |
| Legal documents | 512-1024 tokens | 100-150 |
| Scientific papers | 768-1024 tokens | 50-100 |

## Retrieval System Design

### Query Processing

```go
type RetrievalConfig struct {
    Embedder    types.Embedder
    TopK        int
    MaxTokens   int // Context window constraint
    Filters     map[string]interface{}
}
```

### Query Understanding

1. **Query expansion**: Add synonyms
2. **Query rewriting**: Reformulate for better retrieval
3. **hyde (Hypothetical Document Embeddings)**: Generate hypothetical answer

### Vector Search

```go
func (r *DocumentRAGPipeline) Search(ctx context.Context, query string) (*RAGResult, error) {
    // 1. Generate query embedding
    queryVector, err := r.config.Embedder.Embed(query)
    if err != nil {
        return nil, err
    }
    
    // 2. Search vector store
    results, err := r.vectorStore.Search(ctx, queryVector, r.config.TopK)
    if err != nil {
        return nil, err
    }
    
    // 3. Apply filters if specified
    filtered := applyFilters(results, r.config.Filters)
    
    // 4. Return with assembled context
    return r.assembleContext(filtered), nil
}
```

### Re-ranking

Post-retrieval re-ranking can improve quality:

```python
from sentence_transformers import CrossEncoder

cross_encoder = CrossEncoder('cross-encoder/ms-marco-MiniLM-L-6-v2')

# Re-rank initial results
scores = cross_encoder.predict([(query, doc) for doc in initial_docs])
ranked = sorted(zip(initial_docs, scores), key=lambda x: x[1], reverse=True)
```

## Context Assembly

### Prompt Engineering for RAG

```python
RAG_PROMPT = """You are a helpful assistant. Use the following context to answer the user's question.

Context:
{context}

Question: {question}

Instructions:
- Only use information from the provided context
- If the context doesn't contain enough information, say so
- Provide specific details and citations when possible
- Be concise but thorough

Answer:"""

def assemble_prompt(query: str, retrieved_docs: list[str]) -> str:
    context = "\n\n---\n\n".join(retrieved_docs)
    return RAG_PROMPT.format(context=context, question=query)
```

### Context Window Management

Token budgets require careful management:

```
LLM Context Window: 8192 tokens (e.g., GPT-3.5)
├── System prompt: ~500 tokens
├── Retrieved context: ~7000 tokens  
└── User query: ~100 tokens
```

Strategies:
1. **Truncate**: Cut off oldest/least relevant content
2. **Summarize**: Use smaller LLM to summarize chunks
3. **Select**: Choose top-K by relevance scores

### Metadata Filtering

Filter retrieval results by metadata:

```python
# Filter by source, date, category
results = vector_store.search(
    query_embedding,
    filter={"source": "documentation", "version": "2.0"},
    top_k=10
)
```

## RAG in AI Agent Platform

The AI Agent Platform provides production RAG implementation:

### RAG Pipeline (services/ai-agent-platform/internal/rag/pipeline.go)

```go
type DocumentRAGPipeline struct {
    config        *Config
    vectorStore   types.VectorStore
    embedder      types.Embedder
}

type Config struct {
    Embedder   types.Embedder
    VectorStore types.VectorStore
    TopK       int
    MaxContextTokens int
}
```

### Code RAG Pipeline

```go
// Specialized for code search
type CodeRAGPipeline struct {
    config      *Config
    vectorStore types.VectorStore
    embedder    types.Embedder
}
```

### Indexer Implementation

```go
// DocumentIndexer indexes documents into vector store
type DocumentIndexer struct {
    pipeline *DocumentRAGPipeline
    embedder types.Embedder
}

func (i *DocumentIndexer) IndexDocument(ctx context.Context, doc Document) error {
    // Generate embedding
    vector, err := i.embedder.Embed(doc.Content)
    if err != nil {
        return err
    }
    
    // Store in vector store
    embedding := types.DocumentEmbedding{
        Vector: vector,
        Metadata: types.DocumentEmbeddingMetadata{
            Source:     doc.Source,
            Title:      doc.Title,
            Category:   doc.Category,
            FilePath:   doc.FilePath,
            Language:   doc.Language,
        },
    }
    
    return i.pipeline.AddDocumentation(ctx, embedding)
}
```

### Complete RAG Flow

```go
func (r *DocumentRAGPipeline) Search(ctx context.Context, query string) (*RAGResult, error) {
    start := time.Now()
    
    // Generate embedding for query
    queryVector, err := r.config.Embedder.Embed(query)
    if err != nil {
        return nil, err
    }
    
    // Search vector store
    results, err := r.vectorStore.Search(ctx, queryVector, r.config.TopK)
    if err != nil {
        return nil, err
    }
    
    // Calculate timing
    embeddingTime := time.Since(start).Milliseconds()
    
    return &RAGResult{
        Query:         query,
        Results:       results,
        Context:       assembleContext(results),
        Metadata: RAGMetadata{
            RetrievedDocs: len(results),
            EmbeddingTime: embeddingTime,
            TotalTime:     time.Since(start).Milliseconds(),
        },
    }, nil
}
```

## Key Terminology

| Term | Definition |
|------|------------|
| **Chunking** | Splitting documents into smaller, searchable units |
| **Overlap** | Shared content between adjacent chunks for context |
| **Indexing** | Process of storing documents with embeddings |
| **Retrieval** | Finding relevant documents for a query |
| **Re-ranking** | Post-processing to improve result ordering |
| **Context assembly** | Combining retrieved docs into prompt |
| **Token budget** | Limits on context window size |

## Exercise

### Exercise 11.1: Document Processing Pipeline

Design a document processing pipeline for a company with:
- PDF contracts (100MB total)
- Markdown technical docs (50MB)
- HTML help articles (20MB)

1. What extraction tools would you use for each format?
2. How would you handle the different formats in a unified pipeline?
3. What metadata would you capture?

### Your Design:

| Format | Extractor | Metadata to Capture |
|--------|-----------|---------------------|
| PDF | | |
| Markdown | | |
| HTML | | |

---

### Exercise 11.2: Chunking Strategy Selection

Choose the best chunking strategy for each scenario:

| Scenario | Content Type | Recommended Strategy | Chunk Size |
|----------|--------------|---------------------|------------|
| Stack Overflow Q&A | Text | | |
| GitHub repositories | Code | | |
| Legal contracts | Mixed | | |
| Product catalogs | Tables | | |

Justify your choices.

---

### Exercise 11.3: Implement RAG Pipeline

Using the AI Agent Platform's RAG pipeline structure, implement a simple RAG system that:
1. Indexes 5 documents
2. Searches with a query
3. Returns assembled context

Reference: services/ai-agent-platform/internal/rag/pipeline.go

### Your Implementation:

```go
// Write your implementation

```

---

### Exercise 11.4: Token Budget Calculation

You have an LLM with 4096 token context window.

Given:
- System prompt: 200 tokens
- User query: 100 tokens
- Retrieved chunks (each 512 tokens): 10 chunks

1. How many tokens of retrieved content can you include?
2. How would you select which chunks to include?
3. What's the strategy?

### Your Analysis:

```
(Write your analysis)

```

---

### Exercise 11.5: Metadata Filtering

Design a metadata schema for a documentation search system with:
- Multiple product versions (v1.0, v2.0, v3.0)
- Different content types (api-ref, tutorials, guides)
- Multiple languages (en, es, ja)

How would you filter to get "English API reference for v2.0"?

### Schema Design:

```json
{
  "metadata_fields": {
    "version": [...],
    "type": [...],
    "language": [...]
  }
}
```

---

## Key Takeaways

- ✅ RAG combines retrieval with generation for grounded responses
- ✅ Document processing is critical for quality RAG systems
- ✅ Chunking strategy affects retrieval relevance
- ✅ Token budgets require careful context management
- ✅ AI Agent Platform demonstrates production RAG implementation

## Next Steps

→ [IO-12: RAG Optimization](../12-rag-optimization/README.md)

## Additional Resources

- [LangChain RAG Tutorial](https://python.langchain.com/docs/modules/data_connection/)
- [LlamaIndex Documentation](https://docs.llamaindex.ai)
- [RAG Survey Paper](https://arxiv.org/abs/2312.10997)
- [Chunking Strategies](https://github.com/Run-llama/llama_index/blob/main/docs/docs/module_guides/loading/node_parsers.md)