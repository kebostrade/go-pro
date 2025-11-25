package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/embeddings"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/rag"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/vectorstore"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

func main() {
	fmt.Println("🔍 RAG (Retrieval-Augmented Generation) Demo")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	// Check for API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	useRealEmbeddings := apiKey != ""

	var embedder types.Embedder
	var err error

	if useRealEmbeddings {
		fmt.Println("✅ Using OpenAI embeddings")
		embedder, err = embeddings.NewOpenAIEmbedder(embeddings.OpenAIEmbedderConfig{
			APIKey: apiKey,
			Model:  "text-embedding-3-small",
		})
		if err != nil {
			log.Fatalf("Failed to create embedder: %v", err)
		}
	} else {
		fmt.Println("⚠️  OPENAI_API_KEY not set, using mock embeddings")
		embedder = embeddings.NewMockEmbedder(1536)
	}

	// Create vector store
	vectorStore := vectorstore.NewMemoryVectorStore()
	fmt.Println("✅ Vector store created")

	// Create RAG pipeline
	ragConfig := types.RAGConfig{
		VectorStore:     vectorStore,
		Embedder:        embedder,
		TopK:            3,
		MinScore:        0.5,
		MaxTokens:       4000,
		IncludeMetadata: true,
	}

	pipeline, err := rag.NewRAGPipeline(ragConfig)
	if err != nil {
		log.Fatalf("Failed to create RAG pipeline: %v", err)
	}
	fmt.Println("✅ RAG pipeline created")
	fmt.Println()

	// Demo 1: Code Search
	fmt.Println("📝 Demo 1: Code Search")
	fmt.Println("---")
	demoCodeSearch(pipeline)
	fmt.Println()

	// Demo 2: Documentation Search
	fmt.Println("📚 Demo 2: Documentation Search")
	fmt.Println("---")
	demoDocumentationSearch(pipeline)
	fmt.Println()

	// Demo 3: Code Indexing
	fmt.Println("🗂️  Demo 3: Code Indexing")
	fmt.Println("---")
	demoCodeIndexing(embedder, vectorStore)
	fmt.Println()

	fmt.Println("✅ RAG Demo Complete!")
}

func demoCodeSearch(pipeline *rag.RAGPipeline) {
	ctx := context.Background()

	// Add sample code snippets
	codeSnippets := []struct {
		ID       string
		Code     string
		Language string
		Desc     string
	}{
		{
			ID:       "go-goroutine-1",
			Code:     "go func() {\n    fmt.Println(\"Hello from goroutine\")\n}()",
			Language: "go",
			Desc:     "Basic goroutine example",
		},
		{
			ID:       "go-channel-1",
			Code:     "ch := make(chan int)\ngo func() {\n    ch <- 42\n}()\nvalue := <-ch",
			Language: "go",
			Desc:     "Channel communication example",
		},
		{
			ID:       "python-list-comp-1",
			Code:     "squares = [x**2 for x in range(10)]",
			Language: "python",
			Desc:     "List comprehension example",
		},
		{
			ID:       "js-async-1",
			Code:     "async function fetchData() {\n    const response = await fetch(url);\n    return response.json();\n}",
			Language: "javascript",
			Desc:     "Async/await example",
		},
	}

	fmt.Println("Adding code snippets to vector store...")
	for _, snippet := range codeSnippets {
		metadata := map[string]interface{}{
			"content":     snippet.Code,
			"language":    snippet.Language,
			"description": snippet.Desc,
		}

		if err := pipeline.AddDocument(ctx, snippet.ID, snippet.Code, metadata); err != nil {
			log.Printf("Error adding snippet %s: %v", snippet.ID, err)
		}
	}
	fmt.Println("✅ Added 4 code snippets")
	fmt.Println()

	// Search for code
	queries := []string{
		"How to use goroutines in Go?",
		"Python list operations",
		"Asynchronous JavaScript",
	}

	for _, query := range queries {
		fmt.Printf("Query: %s\n", query)

		result, err := pipeline.Retrieve(ctx, query, nil)
		if err != nil {
			log.Printf("Error searching: %v", err)
			continue
		}

		fmt.Printf("Found %d results:\n", result.Metadata.TotalResults)
		for i, res := range result.Results {
			lang := res.Metadata["language"]
			desc := res.Metadata["description"]
			fmt.Printf("  %d. [%s] %s (Score: %.2f)\n", i+1, lang, desc, res.Score)
		}
		fmt.Println()
	}
}

func demoDocumentationSearch(pipeline *rag.RAGPipeline) {
	ctx := context.Background()

	// Add sample documentation
	docs := []struct {
		ID      string
		Content string
		Title   string
		Source  string
	}{
		{
			ID:      "go-doc-1",
			Content: "Goroutines are lightweight threads managed by the Go runtime. They are created using the 'go' keyword followed by a function call.",
			Title:   "Goroutines",
			Source:  "Go Documentation",
		},
		{
			ID:      "go-doc-2",
			Content: "Channels are typed conduits through which you can send and receive values with the channel operator <-. By default, sends and receives block until the other side is ready.",
			Title:   "Channels",
			Source:  "Go Documentation",
		},
		{
			ID:      "python-doc-1",
			Content: "List comprehensions provide a concise way to create lists. Common applications are to make new lists where each element is the result of some operations applied to each member of another sequence.",
			Title:   "List Comprehensions",
			Source:  "Python Documentation",
		},
	}

	fmt.Println("Adding documentation to vector store...")
	for _, doc := range docs {
		metadata := map[string]interface{}{
			"content": doc.Content,
			"title":   doc.Title,
			"source":  doc.Source,
		}

		if err := pipeline.AddDocument(ctx, doc.ID, doc.Content, metadata); err != nil {
			log.Printf("Error adding doc %s: %v", doc.ID, err)
		}
	}
	fmt.Println("✅ Added 3 documentation entries")
	fmt.Println()

	// Search documentation
	queries := []string{
		"What are goroutines?",
		"How do channels work?",
		"Python list operations",
	}

	for _, query := range queries {
		fmt.Printf("Query: %s\n", query)

		result, err := pipeline.Retrieve(ctx, query, nil)
		if err != nil {
			log.Printf("Error searching: %v", err)
			continue
		}

		fmt.Printf("Found %d results:\n", result.Metadata.TotalResults)
		for i, res := range result.Results {
			title := res.Metadata["title"]
			source := res.Metadata["source"]
			fmt.Printf("  %d. %s - %s (Score: %.2f)\n", i+1, title, source, res.Score)
		}
		fmt.Println()
	}
}

func demoCodeIndexing(embedder types.Embedder, vectorStore types.VectorStore) {
	ctx := context.Background()

	// Create code RAG pipeline
	codeRAGConfig := types.RAGConfig{
		VectorStore: vectorStore,
		Embedder:    embedder,
		TopK:        5,
		MinScore:    0.6,
	}

	codePipeline, err := rag.NewCodeRAGPipeline(codeRAGConfig)
	if err != nil {
		log.Fatalf("Failed to create code pipeline: %v", err)
	}

	// Create indexer
	_ = rag.NewCodeIndexer(codePipeline, embedder) // Will be used for actual indexing

	// Index options
	options := rag.IndexOptions{
		Repository:      "github.com/example/repo",
		Branch:          "main",
		Commit:          "abc123",
		IncludePatterns: []string{"*.go", "*.py", "*.js"},
		MaxFileSize:     1024 * 1024, // 1MB
	}

	fmt.Println("Indexing configuration:")
	fmt.Printf("  Repository: %s\n", options.Repository)
	fmt.Printf("  Branch: %s\n", options.Branch)
	fmt.Printf("  Patterns: %v\n", options.IncludePatterns)
	fmt.Println()

	// Note: In a real scenario, you would index actual files
	fmt.Println("ℹ️  To index a directory, use:")
	fmt.Println("  indexer.IndexDirectory(ctx, \"/path/to/code\", options)")
	fmt.Println()

	// Search code
	searchRequest := types.CodeSearchRequest{
		Query:    "goroutine example",
		Language: "go",
		Limit:    3,
		MinScore: 0.5,
	}

	fmt.Printf("Searching for: %s\n", searchRequest.Query)
	response, err := codePipeline.SearchCode(ctx, searchRequest)
	if err != nil {
		log.Printf("Error searching code: %v", err)
		return
	}

	fmt.Printf("Found %d code snippets:\n", response.TotalResults)
	for i, result := range response.Results {
		fmt.Printf("  %d. [%s] Score: %.2f\n", i+1, result.Language, result.Score)
		if result.Description != "" {
			fmt.Printf("     %s\n", result.Description)
		}
	}
}

