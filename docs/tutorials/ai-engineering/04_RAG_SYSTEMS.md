# Tutorial 4: RAG Systems - Retrieval-Augmented Generation

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  45 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Document Q&A System with RAG                               │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ RAG architecture and pipeline                                    │
│     ✓ Document ingestion and processing                                │
│     ✓ Chunking strategies                                              │
│     ✓ Retrieval and ranking                                            │
│     ✓ Context-aware answer generation                                  │
│                                                                          │
│  🛠️ TECH STACK: OpenAI, Qdrant, Go                                     │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 🧠 Understanding RAG

### What is RAG?

**RAG (Retrieval-Augmented Generation)** combines document retrieval with LLM generation to answer questions based on your own data.

```
┌─────────────────────────────────────────────────────────────────┐
│  Traditional LLM:                                                │
│    Question → LLM → Answer (from training data only)            │
│                                                                  │
│  RAG System:                                                     │
│    Question → Retrieve Docs → LLM + Context → Answer            │
│                                                                  │
│  Benefits:                                                       │
│    ✓ Use your own data                                          │
│    ✓ Up-to-date information                                     │
│    ✓ Cite sources                                               │
│    ✓ Reduce hallucinations                                      │
└─────────────────────────────────────────────────────────────────┘
```

### RAG Pipeline

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                   │
│  1. INDEXING (Offline)                                           │
│     Documents → Chunk → Embed → Store in Vector DB              │
│                                                                   │
│  2. RETRIEVAL (Online)                                           │
│     Query → Embed → Search Vector DB → Top K Documents          │
│                                                                   │
│  3. GENERATION (Online)                                          │
│     Query + Retrieved Docs → LLM → Answer                       │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Project: Document Q&A System

### Step 1: Document Ingestion

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strings"
    
    openai "github.com/sashabaranov/go-openai"
)

type Document struct {
    ID       string
    Content  string
    Metadata map[string]string
}

type Chunk struct {
    ID        string
    Content   string
    DocID     string
    ChunkNum  int
    Embedding []float32
}

// ChunkDocument splits a document into smaller chunks
func ChunkDocument(doc Document, chunkSize, overlap int) []Chunk {
    words := strings.Fields(doc.Content)
    chunks := make([]Chunk, 0)
    chunkNum := 0
    
    for i := 0; i < len(words); i += chunkSize - overlap {
        end := min(i+chunkSize, len(words))
        content := strings.Join(words[i:end], " ")
        
        chunks = append(chunks, Chunk{
            ID:       fmt.Sprintf("%s_chunk_%d", doc.ID, chunkNum),
            Content:  content,
            DocID:    doc.ID,
            ChunkNum: chunkNum,
        })
        
        chunkNum++
        
        if end == len(words) {
            break
        }
    }
    
    return chunks
}
```

### Step 2: RAG Pipeline

```go
type RAGPipeline struct {
    llmClient     *openai.Client
    vectorStore   VectorStore
    chunkSize     int
    chunkOverlap  int
    topK          int
}

func NewRAGPipeline(apiKey string, vectorStore VectorStore) *RAGPipeline {
    return &RAGPipeline{
        llmClient:    openai.NewClient(apiKey),
        vectorStore:  vectorStore,
        chunkSize:    500,  // words
        chunkOverlap: 50,   // words
        topK:         5,    // retrieve top 5 chunks
    }
}

// IndexDocument processes and stores a document
func (r *RAGPipeline) IndexDocument(ctx context.Context, doc Document) error {
    // 1. Chunk the document
    chunks := ChunkDocument(doc, r.chunkSize, r.chunkOverlap)
    
    // 2. Generate embeddings for each chunk
    for i := range chunks {
        embedding, err := r.generateEmbedding(ctx, chunks[i].Content)
        if err != nil {
            return fmt.Errorf("failed to generate embedding: %w", err)
        }
        chunks[i].Embedding = embedding
        
        // 3. Store in vector database
        if err := r.vectorStore.Upsert(ctx, chunks[i]); err != nil {
            return fmt.Errorf("failed to store chunk: %w", err)
        }
    }
    
    return nil
}

func (r *RAGPipeline) generateEmbedding(ctx context.Context, text string) ([]float32, error) {
    resp, err := r.llmClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{
        Model: openai.AdaEmbeddingV2,
        Input: []string{text},
    })
    
    if err != nil {
        return nil, err
    }
    
    return resp.Data[0].Embedding, nil
}
```

### Step 3: Query and Retrieval

```go
// Query answers a question using RAG
func (r *RAGPipeline) Query(ctx context.Context, question string) (string, error) {
    // 1. Generate query embedding
    queryEmbedding, err := r.generateEmbedding(ctx, question)
    if err != nil {
        return "", fmt.Errorf("failed to generate query embedding: %w", err)
    }
    
    // 2. Retrieve relevant chunks
    chunks, err := r.vectorStore.Search(ctx, queryEmbedding, r.topK)
    if err != nil {
        return "", fmt.Errorf("failed to search: %w", err)
    }
    
    // 3. Build context from chunks
    context := r.buildContext(chunks)
    
    // 4. Generate answer with LLM
    answer, err := r.generateAnswer(ctx, question, context)
    if err != nil {
        return "", fmt.Errorf("failed to generate answer: %w", err)
    }
    
    return answer, nil
}

func (r *RAGPipeline) buildContext(chunks []Chunk) string {
    var builder strings.Builder
    
    for i, chunk := range chunks {
        builder.WriteString(fmt.Sprintf("Document %d:\n%s\n\n", i+1, chunk.Content))
    }
    
    return builder.String()
}

func (r *RAGPipeline) generateAnswer(ctx context.Context, question, context string) (string, error) {
    prompt := fmt.Sprintf(`Answer the question based on the context below. If the answer is not in the context, say "I don't have enough information to answer that."

Context:
%s

Question: %s

Answer:`, context, question)
    
    resp, err := r.llmClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleSystem,
                Content: "You are a helpful assistant that answers questions based on the provided context. Always cite which document you're referencing.",
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: prompt,
            },
        },
        Temperature: 0.3, // Lower temperature for factual answers
        MaxTokens:   500,
    })
    
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}
```

### Step 4: Complete Example

```go
func main() {
    ctx := context.Background()
    
    // Initialize RAG pipeline
    vectorStore := NewInMemoryVectorStore()
    rag := NewRAGPipeline(os.Getenv("OPENAI_API_KEY"), vectorStore)
    
    // Index documents
    docs := []Document{
        {
            ID: "go-basics",
            Content: `Go is a statically typed, compiled programming language designed at Google. 
                     It is syntactically similar to C, but with memory safety, garbage collection, 
                     structural typing, and CSP-style concurrency. Go was created to improve 
                     programming productivity in an era of multicore processors and large codebases.`,
        },
        {
            ID: "go-concurrency",
            Content: `Go has built-in concurrency features. Goroutines are lightweight threads 
                     managed by the Go runtime. Channels are typed conduits through which you can 
                     send and receive values with the channel operator <-. The select statement 
                     lets a goroutine wait on multiple communication operations.`,
        },
        {
            ID: "go-packages",
            Content: `Go programs are organized into packages. A package is a collection of source 
                     files in the same directory that are compiled together. The main package is 
                     special - it defines a standalone executable program, not a library. The 
                     import statement is used to use code from other packages.`,
        },
    }
    
    fmt.Println("Indexing documents...")
    for _, doc := range docs {
        if err := rag.IndexDocument(ctx, doc); err != nil {
            panic(err)
        }
    }
    
    // Query the system
    questions := []string{
        "What is Go?",
        "How does concurrency work in Go?",
        "What are packages in Go?",
        "What is Python?", // Should say "I don't have enough information"
    }
    
    for _, question := range questions {
        fmt.Printf("\nQ: %s\n", question)
        answer, err := rag.Query(ctx, question)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }
        fmt.Printf("A: %s\n", answer)
    }
}
```

**📤 Expected Output:**
```
Indexing documents...

Q: What is Go?
A: Based on Document 1, Go is a statically typed, compiled programming language 
designed at Google. It is syntactically similar to C, but includes memory safety, 
garbage collection, structural typing, and CSP-style concurrency features.

Q: How does concurrency work in Go?
A: According to Document 2, Go has built-in concurrency features including goroutines 
(lightweight threads managed by the Go runtime) and channels (typed conduits for 
sending and receiving values). The select statement allows goroutines to wait on 
multiple communication operations.

Q: What are packages in Go?
A: Document 3 explains that Go programs are organized into packages, which are 
collections of source files in the same directory compiled together. The main package 
is special as it defines a standalone executable program.

Q: What is Python?
A: I don't have enough information to answer that.
```

---

## 🎨 Advanced RAG Techniques

### 1. Re-ranking Retrieved Documents

```go
func (r *RAGPipeline) rerank(ctx context.Context, query string, chunks []Chunk) []Chunk {
    // Use cross-encoder or LLM to re-rank
    type scored struct {
        chunk Chunk
        score float32
    }
    
    scores := make([]scored, len(chunks))
    
    for i, chunk := range chunks {
        // Calculate relevance score
        score := r.calculateRelevance(ctx, query, chunk.Content)
        scores[i] = scored{chunk: chunk, score: score}
    }
    
    // Sort by score
    sort.Slice(scores, func(i, j int) bool {
        return scores[i].score > scores[j].score
    })
    
    reranked := make([]Chunk, len(scores))
    for i, s := range scores {
        reranked[i] = s.chunk
    }
    
    return reranked
}
```

### 2. Hybrid Search (Keyword + Semantic)

```go
func (r *RAGPipeline) hybridSearch(ctx context.Context, query string, topK int) ([]Chunk, error) {
    // 1. Semantic search
    semanticResults, err := r.semanticSearch(ctx, query, topK*2)
    if err != nil {
        return nil, err
    }
    
    // 2. Keyword search (BM25)
    keywordResults := r.keywordSearch(query, topK*2)
    
    // 3. Combine and re-rank
    combined := r.combineResults(semanticResults, keywordResults, topK)
    
    return combined, nil
}
```

### 3. Query Expansion

```go
func (r *RAGPipeline) expandQuery(ctx context.Context, query string) ([]string, error) {
    prompt := fmt.Sprintf(`Generate 3 alternative phrasings of this question:
"%s"

Return only the questions, one per line.`, query)
    
    resp, err := r.llmClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {Role: openai.ChatMessageRoleUser, Content: prompt},
        },
    })
    
    if err != nil {
        return nil, err
    }
    
    expanded := strings.Split(resp.Choices[0].Message.Content, "\n")
    return append([]string{query}, expanded...), nil
}
```

### 4. Citation and Source Tracking

```go
type Answer struct {
    Content string
    Sources []Source
}

type Source struct {
    DocID    string
    ChunkNum int
    Content  string
    Score    float32
}

func (r *RAGPipeline) QueryWithSources(ctx context.Context, question string) (Answer, error) {
    // Retrieve chunks
    chunks, err := r.retrieve(ctx, question)
    if err != nil {
        return Answer{}, err
    }
    
    // Generate answer
    content, err := r.generateAnswer(ctx, question, r.buildContext(chunks))
    if err != nil {
        return Answer{}, err
    }
    
    // Build sources
    sources := make([]Source, len(chunks))
    for i, chunk := range chunks {
        sources[i] = Source{
            DocID:    chunk.DocID,
            ChunkNum: chunk.ChunkNum,
            Content:  chunk.Content,
            Score:    chunk.Score,
        }
    }
    
    return Answer{
        Content: content,
        Sources: sources,
    }, nil
}
```

---

## 📊 Chunking Strategies

### Fixed-Size Chunking

```go
// Simple, predictable, but may split sentences
func fixedSizeChunk(text string, size int) []string {
    words := strings.Fields(text)
    chunks := make([]string, 0)
    
    for i := 0; i < len(words); i += size {
        end := min(i+size, len(words))
        chunks = append(chunks, strings.Join(words[i:end], " "))
    }
    
    return chunks
}
```

### Sentence-Based Chunking

```go
// Respects sentence boundaries
func sentenceChunk(text string, maxSentences int) []string {
    sentences := splitSentences(text)
    chunks := make([]string, 0)
    
    for i := 0; i < len(sentences); i += maxSentences {
        end := min(i+maxSentences, len(sentences))
        chunks = append(chunks, strings.Join(sentences[i:end], " "))
    }
    
    return chunks
}
```

### Semantic Chunking

```go
// Groups semantically similar sentences
func semanticChunk(text string, threshold float32) []string {
    sentences := splitSentences(text)
    embeddings := generateEmbeddings(sentences)
    
    chunks := make([]string, 0)
    currentChunk := []string{sentences[0]}
    
    for i := 1; i < len(sentences); i++ {
        similarity := cosineSimilarity(embeddings[i-1], embeddings[i])
        
        if similarity > threshold {
            currentChunk = append(currentChunk, sentences[i])
        } else {
            chunks = append(chunks, strings.Join(currentChunk, " "))
            currentChunk = []string{sentences[i]}
        }
    }
    
    chunks = append(chunks, strings.Join(currentChunk, " "))
    return chunks
}
```

---

## 🎯 Challenges

### Challenge 1: Multi-Document Q&A
Build a system that can answer questions across multiple PDF documents.

### Challenge 2: Conversational RAG
Add conversation history to maintain context across multiple questions.

### Challenge 3: Streaming Answers
Implement streaming responses for better UX.

### Challenge 4: Evaluation System
Build a system to evaluate RAG answer quality.

---

## ✅ What You Learned

```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ RAG architecture and pipeline                                 │
│  ✓ Document ingestion and chunking                               │
│  ✓ Embedding generation and storage                              │
│  ✓ Semantic retrieval                                            │
│  ✓ Context-aware answer generation                               │
│  ✓ Re-ranking and hybrid search                                  │
│  ✓ Query expansion techniques                                    │
│  ✓ Citation and source tracking                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Next Steps

**Next Tutorial:**
[Tutorial 5: AI Agents →](05_AI_AGENTS.md)

Learn how to build autonomous AI agents with the ReAct pattern!

---

**💡 Pro Tip**: Experiment with chunk sizes - smaller chunks are more precise, larger chunks provide more context!

