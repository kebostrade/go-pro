# Tutorial 3: Embeddings & Vector Search

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  35 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Semantic Search Engine                                     │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Text embeddings and vector representations                       │
│     ✓ Cosine similarity and distance metrics                           │
│     ✓ Vector databases (Qdrant)                                        │
│     ✓ Semantic search implementation                                   │
│     ✓ Document indexing and retrieval                                  │
│                                                                          │
│  🛠️ TECH STACK: OpenAI Embeddings, Qdrant, Go                          │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 🧠 Understanding Embeddings

### What are Embeddings?

**Embeddings** convert text into numerical vectors that capture semantic meaning.

```
┌─────────────────────────────────────────────────────────────────┐
│  Text: "The cat sat on the mat"                                 │
│                                                                  │
│  Embedding: [0.023, -0.145, 0.892, ..., 0.234]                 │
│             ↑                                                    │
│             1536 dimensions (for OpenAI ada-002)                │
│                                                                  │
│  Similar texts have similar vectors!                            │
└─────────────────────────────────────────────────────────────────┘
```

### Why Embeddings?

```
Traditional Search (Keyword):
  Query: "automobile"
  Document: "car"
  Match: ❌ No match (different words)

Semantic Search (Embeddings):
  Query: "automobile" → [0.1, 0.8, ...]
  Document: "car" → [0.12, 0.79, ...]
  Similarity: ✅ 0.95 (very similar!)
```

---

## 📐 Vector Similarity

### Cosine Similarity

Measures the angle between two vectors (0 to 1, higher = more similar).

```go
func cosineSimilarity(a, b []float32) float32 {
    var dotProduct, normA, normB float32
    
    for i := range a {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    
    if normA == 0 || normB == 0 {
        return 0
    }
    
    return dotProduct / (float32(math.Sqrt(float64(normA))) * 
                         float32(math.Sqrt(float64(normB))))
}
```

### Distance Metrics

```
┌──────────────────────────────────────────────────────────────────┐
│  Cosine Similarity:                                              │
│    • Range: 0 to 1 (1 = identical)                               │
│    • Best for: Text, semantic similarity                         │
│    • Use when: Direction matters more than magnitude             │
│                                                                   │
│  Euclidean Distance:                                             │
│    • Range: 0 to ∞ (0 = identical)                               │
│    • Best for: Spatial data, exact matches                       │
│    • Use when: Magnitude matters                                 │
│                                                                   │
│  Dot Product:                                                    │
│    • Range: -∞ to ∞                                              │
│    • Best for: Fast approximate search                           │
│    • Use when: Speed is critical                                 │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Project: Semantic Search Engine

### Step 1: Generate Embeddings

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    openai "github.com/sashabaranov/go-openai"
)

type Document struct {
    ID        string
    Content   string
    Embedding []float32
}

func generateEmbedding(client *openai.Client, text string) ([]float32, error) {
    resp, err := client.CreateEmbeddings(
        context.Background(),
        openai.EmbeddingRequest{
            Model: openai.AdaEmbeddingV2,
            Input: []string{text},
        },
    )
    
    if err != nil {
        return nil, err
    }
    
    return resp.Data[0].Embedding, nil
}

func main() {
    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
    
    // Generate embedding for a document
    text := "Go is a statically typed, compiled programming language"
    embedding, err := generateEmbedding(client, text)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Generated embedding with %d dimensions\n", len(embedding))
    fmt.Printf("First 5 values: %v\n", embedding[:5])
}
```

### Step 2: Build In-Memory Vector Store

```go
type VectorStore struct {
    documents []Document
    mu        sync.RWMutex
}

func NewVectorStore() *VectorStore {
    return &VectorStore{
        documents: make([]Document, 0),
    }
}

func (vs *VectorStore) Add(doc Document) {
    vs.mu.Lock()
    defer vs.mu.Unlock()
    vs.documents = append(vs.documents, doc)
}

func (vs *VectorStore) Search(queryEmbedding []float32, topK int) []Document {
    vs.mu.RLock()
    defer vs.mu.RUnlock()
    
    type scored struct {
        doc   Document
        score float32
    }
    
    scores := make([]scored, len(vs.documents))
    for i, doc := range vs.documents {
        scores[i] = scored{
            doc:   doc,
            score: cosineSimilarity(queryEmbedding, doc.Embedding),
        }
    }
    
    // Sort by score descending
    sort.Slice(scores, func(i, j int) bool {
        return scores[i].score > scores[j].score
    })
    
    // Return top K
    k := min(topK, len(scores))
    results := make([]Document, k)
    for i := 0; i < k; i++ {
        results[i] = scores[i].doc
    }
    
    return results
}
```

### Step 3: Complete Search Engine

```go
type SearchEngine struct {
    client *openai.Client
    store  *VectorStore
}

func NewSearchEngine(apiKey string) *SearchEngine {
    return &SearchEngine{
        client: openai.NewClient(apiKey),
        store:  NewVectorStore(),
    }
}

func (se *SearchEngine) IndexDocument(id, content string) error {
    embedding, err := generateEmbedding(se.client, content)
    if err != nil {
        return fmt.Errorf("failed to generate embedding: %w", err)
    }
    
    se.store.Add(Document{
        ID:        id,
        Content:   content,
        Embedding: embedding,
    })
    
    return nil
}

func (se *SearchEngine) Search(query string, topK int) ([]Document, error) {
    // Generate query embedding
    queryEmbedding, err := generateEmbedding(se.client, query)
    if err != nil {
        return nil, fmt.Errorf("failed to generate query embedding: %w", err)
    }
    
    // Search vector store
    results := se.store.Search(queryEmbedding, topK)
    
    return results, nil
}
```

### Step 4: Use the Search Engine

```go
func main() {
    engine := NewSearchEngine(os.Getenv("OPENAI_API_KEY"))
    
    // Index documents
    docs := []struct {
        id      string
        content string
    }{
        {"1", "Go is a statically typed, compiled programming language"},
        {"2", "Python is a high-level, interpreted programming language"},
        {"3", "JavaScript is a dynamic, prototype-based language"},
        {"4", "Rust is a systems programming language focused on safety"},
        {"5", "Java is an object-oriented programming language"},
    }
    
    fmt.Println("Indexing documents...")
    for _, doc := range docs {
        if err := engine.IndexDocument(doc.id, doc.content); err != nil {
            panic(err)
        }
    }
    
    // Search
    query := "compiled programming language"
    fmt.Printf("\nSearching for: %s\n\n", query)
    
    results, err := engine.Search(query, 3)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Top 3 Results:")
    for i, result := range results {
        fmt.Printf("%d. [%s] %s\n", i+1, result.ID, result.Content)
    }
}
```

**📤 Expected Output:**
```
Indexing documents...

Searching for: compiled programming language

Top 3 Results:
1. [1] Go is a statically typed, compiled programming language
2. [4] Rust is a systems programming language focused on safety
3. [5] Java is an object-oriented programming language
```

---

## 🗄️ Using Qdrant Vector Database

### Setup Qdrant

```bash
# Run Qdrant with Docker
docker run -p 6333:6333 qdrant/qdrant
```

### Integrate Qdrant

```go
import (
    "github.com/qdrant/go-client/qdrant"
)

type QdrantStore struct {
    client     *qdrant.Client
    collection string
}

func NewQdrantStore(url, collection string) (*QdrantStore, error) {
    client, err := qdrant.NewClient(&qdrant.Config{
        Host: url,
    })
    if err != nil {
        return nil, err
    }
    
    // Create collection
    err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
        CollectionName: collection,
        VectorsConfig: qdrant.VectorsConfig{
            Params: &qdrant.VectorParams{
                Size:     1536, // OpenAI ada-002 dimension
                Distance: qdrant.Distance_Cosine,
            },
        },
    })
    
    return &QdrantStore{
        client:     client,
        collection: collection,
    }, nil
}

func (qs *QdrantStore) Upsert(id string, embedding []float32, payload map[string]interface{}) error {
    _, err := qs.client.Upsert(context.Background(), &qdrant.UpsertPoints{
        CollectionName: qs.collection,
        Points: []*qdrant.PointStruct{
            {
                Id:      qdrant.NewIDNum(hashID(id)),
                Vectors: qdrant.NewVectors(embedding...),
                Payload: qdrant.NewValueMap(payload),
            },
        },
    })
    
    return err
}

func (qs *QdrantStore) Search(embedding []float32, limit int) ([]*qdrant.ScoredPoint, error) {
    results, err := qs.client.Search(context.Background(), &qdrant.SearchPoints{
        CollectionName: qs.collection,
        Vector:         embedding,
        Limit:          uint64(limit),
        WithPayload:    qdrant.NewWithPayload(true),
    })
    
    if err != nil {
        return nil, err
    }
    
    return results, nil
}
```

---

## 📊 Advanced Techniques

### 1. Chunking Long Documents

```go
func chunkText(text string, chunkSize, overlap int) []string {
    words := strings.Fields(text)
    var chunks []string
    
    for i := 0; i < len(words); i += chunkSize - overlap {
        end := min(i+chunkSize, len(words))
        chunk := strings.Join(words[i:end], " ")
        chunks = append(chunks, chunk)
        
        if end == len(words) {
            break
        }
    }
    
    return chunks
}
```

### 2. Batch Embedding Generation

```go
func generateEmbeddingsBatch(client *openai.Client, texts []string) ([][]float32, error) {
    resp, err := client.CreateEmbeddings(
        context.Background(),
        openai.EmbeddingRequest{
            Model: openai.AdaEmbeddingV2,
            Input: texts,
        },
    )
    
    if err != nil {
        return nil, err
    }
    
    embeddings := make([][]float32, len(resp.Data))
    for i, data := range resp.Data {
        embeddings[i] = data.Embedding
    }
    
    return embeddings, nil
}
```

### 3. Hybrid Search (Keyword + Semantic)

```go
func hybridSearch(query string, keywordResults, semanticResults []Document) []Document {
    // Combine and re-rank results
    scoreMap := make(map[string]float32)
    
    // Keyword scores
    for i, doc := range keywordResults {
        scoreMap[doc.ID] += float32(len(keywordResults)-i) * 0.3
    }
    
    // Semantic scores
    for i, doc := range semanticResults {
        scoreMap[doc.ID] += float32(len(semanticResults)-i) * 0.7
    }
    
    // Sort by combined score
    // ... implementation
    
    return rankedResults
}
```

---

## 🎯 Challenges

### Challenge 1: Multi-Language Search
Build a search engine that works across multiple languages.

### Challenge 2: Image Search
Use CLIP embeddings to search images with text queries.

### Challenge 3: Recommendation System
Build a content recommendation system using embeddings.

### Challenge 4: Duplicate Detection
Find duplicate or near-duplicate documents using embeddings.

---

## ✅ What You Learned

```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ Text embeddings and vector representations                    │
│  ✓ Cosine similarity and distance metrics                        │
│  ✓ In-memory vector store implementation                         │
│  ✓ Qdrant vector database integration                            │
│  ✓ Semantic search implementation                                │
│  ✓ Document chunking strategies                                  │
│  ✓ Batch embedding generation                                    │
│  ✓ Hybrid search techniques                                      │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Next Steps

**Practice:**
1. Build the semantic search engine
2. Index 100+ documents
3. Try different chunk sizes
4. Experiment with Qdrant

**Next Tutorial:**
[Tutorial 4: RAG Systems →](04_RAG_SYSTEMS.md)

Learn how to build Retrieval-Augmented Generation systems!

---

## 📚 Resources

- [OpenAI Embeddings Guide](https://platform.openai.com/docs/guides/embeddings)
- [Qdrant Documentation](https://qdrant.tech/documentation/)
- [Vector Search Explained](https://www.pinecone.io/learn/vector-search/)

---

**💡 Pro Tip**: Cache embeddings to save API costs - they don't change for the same text!

