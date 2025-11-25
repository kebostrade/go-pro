package rag

import (
	"context"
	"testing"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/embeddings"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/vectorstore"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

func TestNewRAGPipeline(t *testing.T) {
	embedder := embeddings.NewMockEmbedder(1536)
	vectorStore := vectorstore.NewMemoryVectorStore()

	config := types.RAGConfig{
		VectorStore: vectorStore,
		Embedder:    embedder,
	}

	pipeline, err := NewRAGPipeline(config)
	if err != nil {
		t.Fatalf("Failed to create pipeline: %v", err)
	}

	if pipeline == nil {
		t.Error("Pipeline is nil")
	}

	// Check defaults
	if pipeline.config.TopK != 5 {
		t.Errorf("Expected default TopK 5, got %d", pipeline.config.TopK)
	}

	if pipeline.config.MinScore != 0.7 {
		t.Errorf("Expected default MinScore 0.7, got %f", pipeline.config.MinScore)
	}
}

func TestNewRAGPipeline_NoVectorStore(t *testing.T) {
	embedder := embeddings.NewMockEmbedder(1536)

	config := types.RAGConfig{
		Embedder: embedder,
	}

	_, err := NewRAGPipeline(config)
	if err == nil {
		t.Error("Expected error for missing vector store")
	}
}

func TestNewRAGPipeline_NoEmbedder(t *testing.T) {
	vectorStore := vectorstore.NewMemoryVectorStore()

	config := types.RAGConfig{
		VectorStore: vectorStore,
	}

	_, err := NewRAGPipeline(config)
	if err == nil {
		t.Error("Expected error for missing embedder")
	}
}

func TestRAGPipeline_AddDocument(t *testing.T) {
	pipeline := createTestPipeline(t)
	ctx := context.Background()

	err := pipeline.AddDocument(ctx, "doc-1", "test content", map[string]interface{}{
		"language": "go",
	})

	if err != nil {
		t.Fatalf("AddDocument failed: %v", err)
	}

	// Verify document was added
	count, _ := pipeline.config.VectorStore.Count()
	if count != 1 {
		t.Errorf("Expected 1 document, got %d", count)
	}
}

func TestRAGPipeline_Retrieve(t *testing.T) {
	pipeline := createTestPipeline(t)
	ctx := context.Background()

	// Add test documents
	pipeline.AddDocument(ctx, "doc-1", "goroutines in Go", map[string]interface{}{
		"language": "go",
	})
	pipeline.AddDocument(ctx, "doc-2", "channels in Go", map[string]interface{}{
		"language": "go",
	})
	pipeline.AddDocument(ctx, "doc-3", "list comprehensions in Python", map[string]interface{}{
		"language": "python",
	})

	// Retrieve
	result, err := pipeline.Retrieve(ctx, "goroutines", nil)
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}

	if len(result.Results) == 0 {
		t.Error("Expected at least one result")
	}

	// Check metadata
	if result.Metadata.TotalResults == 0 {
		t.Error("Expected TotalResults > 0")
	}

	// Note: Timing may be 0ms for very fast operations, which is acceptable
}

func TestRAGPipeline_RetrieveWithFilters(t *testing.T) {
	pipeline := createTestPipeline(t)
	ctx := context.Background()

	// Add documents
	pipeline.AddDocument(ctx, "go-1", "Go code", map[string]interface{}{
		"language": "go",
	})
	pipeline.AddDocument(ctx, "py-1", "Python code", map[string]interface{}{
		"language": "python",
	})

	// Retrieve with filter
	filters := map[string]interface{}{"language": "go"}
	result, err := pipeline.Retrieve(ctx, "code", filters)
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}

	// All results should be Go
	for _, res := range result.Results {
		if res.Metadata["language"] != "go" {
			t.Errorf("Expected language 'go', got '%v'", res.Metadata["language"])
		}
	}
}

func TestRAGPipeline_UpdateDocument(t *testing.T) {
	pipeline := createTestPipeline(t)
	ctx := context.Background()

	// Add document
	pipeline.AddDocument(ctx, "doc-1", "original content", map[string]interface{}{
		"version": 1,
	})

	// Update document
	err := pipeline.UpdateDocument(ctx, "doc-1", "updated content", map[string]interface{}{
		"version": 2,
	})

	if err != nil {
		t.Fatalf("UpdateDocument failed: %v", err)
	}

	// Verify update
	result, _ := pipeline.config.VectorStore.Get("doc-1")
	if result.Metadata["version"] != 2 {
		t.Errorf("Expected version 2, got %v", result.Metadata["version"])
	}
}

func TestRAGPipeline_DeleteDocument(t *testing.T) {
	pipeline := createTestPipeline(t)
	ctx := context.Background()

	// Add document
	pipeline.AddDocument(ctx, "doc-1", "content", nil)

	// Delete document
	err := pipeline.DeleteDocument(ctx, "doc-1")
	if err != nil {
		t.Fatalf("DeleteDocument failed: %v", err)
	}

	// Verify deletion
	count, _ := pipeline.config.VectorStore.Count()
	if count != 0 {
		t.Errorf("Expected 0 documents, got %d", count)
	}
}

func TestCodeRAGPipeline_SearchCode(t *testing.T) {
	codePipeline := createTestCodePipeline(t)
	ctx := context.Background()

	// Add code snippets
	codePipeline.AddCode(ctx, types.CodeEmbedding{
		ID:       "go-1",
		Code:     "go func() {}",
		Language: "go",
		Vector:   make([]float64, 1536),
	})

	// Search
	request := types.CodeSearchRequest{
		Query:    "goroutine",
		Language: "go",
		Limit:    5,
	}

	response, err := codePipeline.SearchCode(ctx, request)
	if err != nil {
		t.Fatalf("SearchCode failed: %v", err)
	}

	if response.Query != "goroutine" {
		t.Errorf("Expected query 'goroutine', got '%s'", response.Query)
	}
}

func TestDocumentRAGPipeline_SearchDocumentation(t *testing.T) {
	docPipeline := createTestDocumentPipeline(t)
	ctx := context.Background()

	// Add documentation
	docPipeline.AddDocumentation(ctx, types.DocumentEmbedding{
		ID:       "doc-1",
		Content:  "Goroutines are lightweight threads",
		Title:    "Goroutines",
		Source:   "Go Docs",
		Language: "go",
		Vector:   make([]float64, 1536),
	})

	// Search
	request := types.DocumentSearchRequest{
		Query:    "goroutines",
		Language: "go",
		Limit:    5,
	}

	response, err := docPipeline.SearchDocumentation(ctx, request)
	if err != nil {
		t.Fatalf("SearchDocumentation failed: %v", err)
	}

	if response.Query != "goroutines" {
		t.Errorf("Expected query 'goroutines', got '%s'", response.Query)
	}
}

func TestRAGPipeline_FormatContext(t *testing.T) {
	pipeline := createTestPipeline(t)

	results := []types.VectorSearchResult{
		{
			ID:    "doc-1",
			Score: 0.9,
			Metadata: map[string]interface{}{
				"content": "Test content 1",
			},
		},
		{
			ID:    "doc-2",
			Score: 0.8,
			Metadata: map[string]interface{}{
				"content": "Test content 2",
			},
		},
	}

	context := pipeline.formatContext(results)

	if context == "" {
		t.Error("Expected non-empty context")
	}

	// Check if context contains content
	if !contains(context, "Test content 1") {
		t.Error("Context missing first result")
	}

	if !contains(context, "Test content 2") {
		t.Error("Context missing second result")
	}
}

func TestRAGPipeline_EstimateTokens(t *testing.T) {
	pipeline := createTestPipeline(t)

	tests := []struct {
		text     string
		expected int
	}{
		{"", 0},
		{"test", 1},
		{"hello world", 2},
		{"this is a test", 3},
	}

	for _, tt := range tests {
		tokens := pipeline.estimateTokens(tt.text)
		if tokens != tt.expected {
			t.Errorf("For text '%s', expected %d tokens, got %d", tt.text, tt.expected, tokens)
		}
	}
}

// Helper functions

func createTestPipeline(t *testing.T) *RAGPipeline {
	embedder := embeddings.NewMockEmbedder(1536)
	vectorStore := vectorstore.NewMemoryVectorStore()

	config := types.RAGConfig{
		VectorStore: vectorStore,
		Embedder:    embedder,
		TopK:        5,
		MinScore:    0.5,
	}

	pipeline, err := NewRAGPipeline(config)
	if err != nil {
		t.Fatalf("Failed to create test pipeline: %v", err)
	}

	return pipeline
}

func createTestCodePipeline(t *testing.T) *CodeRAGPipeline {
	embedder := embeddings.NewMockEmbedder(1536)
	vectorStore := vectorstore.NewMemoryVectorStore()

	config := types.RAGConfig{
		VectorStore: vectorStore,
		Embedder:    embedder,
		TopK:        5,
		MinScore:    0.5,
	}

	pipeline, err := NewCodeRAGPipeline(config)
	if err != nil {
		t.Fatalf("Failed to create test code pipeline: %v", err)
	}

	return pipeline
}

func createTestDocumentPipeline(t *testing.T) *DocumentRAGPipeline {
	embedder := embeddings.NewMockEmbedder(1536)
	vectorStore := vectorstore.NewMemoryVectorStore()

	config := types.RAGConfig{
		VectorStore: vectorStore,
		Embedder:    embedder,
		TopK:        5,
		MinScore:    0.5,
	}

	pipeline, err := NewDocumentRAGPipeline(config)
	if err != nil {
		t.Fatalf("Failed to create test document pipeline: %v", err)
	}

	return pipeline
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
