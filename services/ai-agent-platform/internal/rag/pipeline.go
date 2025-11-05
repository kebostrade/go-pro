package rag

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// RAGPipeline implements Retrieval-Augmented Generation
type RAGPipeline struct {
	config types.RAGConfig
}

// NewRAGPipeline creates a new RAG pipeline
func NewRAGPipeline(config types.RAGConfig) (*RAGPipeline, error) {
	if config.VectorStore == nil {
		return nil, fmt.Errorf("vector store is required")
	}
	if config.Embedder == nil {
		return nil, fmt.Errorf("embedder is required")
	}

	// Set defaults
	if config.TopK == 0 {
		config.TopK = 5
	}
	if config.MinScore == 0 {
		config.MinScore = 0.7
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4000
	}

	return &RAGPipeline{
		config: config,
	}, nil
}

// Retrieve retrieves relevant context for a query
func (r *RAGPipeline) Retrieve(ctx context.Context, query string, filters map[string]interface{}) (*types.RAGResult, error) {
	// Generate query embedding
	embeddingStart := time.Now()
	queryVector, err := r.config.Embedder.Embed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}
	embeddingTime := time.Since(embeddingStart).Milliseconds()

	// Search vector store
	retrievalStart := time.Now()
	results, err := r.config.VectorStore.Search(queryVector, r.config.TopK, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}
	retrievalTime := time.Since(retrievalStart).Milliseconds()

	// Filter by minimum score
	filteredResults := make([]types.VectorSearchResult, 0)
	for _, result := range results {
		if result.Score >= r.config.MinScore {
			filteredResults = append(filteredResults, result)
		}
	}

	// Format context
	context := r.formatContext(filteredResults)

	// Estimate tokens
	contextTokens := r.estimateTokens(context)

	return &types.RAGResult{
		Query:   query,
		Results: filteredResults,
		Context: context,
		Metadata: types.RAGMetadata{
			TotalResults:  len(filteredResults),
			RetrievalTime: retrievalTime,
			EmbeddingTime: embeddingTime,
			ContextTokens: contextTokens,
			Filters:       filters,
		},
	}, nil
}

// formatContext formats search results into context for LLM
func (r *RAGPipeline) formatContext(results []types.VectorSearchResult) string {
	if len(results) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("Relevant Context:\n\n")

	for i, result := range results {
		builder.WriteString(fmt.Sprintf("--- Result %d (Score: %.2f) ---\n", i+1, result.Score))

		// Extract content from metadata
		if content, ok := result.Metadata["content"].(string); ok {
			builder.WriteString(content)
			builder.WriteString("\n\n")
		}

		// Add metadata if enabled
		if r.config.IncludeMetadata {
			builder.WriteString("Metadata:\n")
			for key, value := range result.Metadata {
				if key != "content" {
					builder.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
				}
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// estimateTokens estimates the number of tokens in text
func (r *RAGPipeline) estimateTokens(text string) int {
	// Rough estimate: ~4 characters per token
	return len(text) / 4
}

// AddDocument adds a document to the vector store
func (r *RAGPipeline) AddDocument(ctx context.Context, id, content string, metadata map[string]interface{}) error {
	// Generate embedding
	vector, err := r.config.Embedder.Embed(content)
	if err != nil {
		return fmt.Errorf("failed to embed document: %w", err)
	}

	// Add content to metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["content"] = content

	// Store in vector store
	if err := r.config.VectorStore.Store(id, vector, metadata); err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	return nil
}

// AddDocuments adds multiple documents to the vector store
func (r *RAGPipeline) AddDocuments(ctx context.Context, documents []struct {
	ID       string
	Content  string
	Metadata map[string]interface{}
}) error {
	for _, doc := range documents {
		if err := r.AddDocument(ctx, doc.ID, doc.Content, doc.Metadata); err != nil {
			return fmt.Errorf("failed to add document %s: %w", doc.ID, err)
		}
	}
	return nil
}

// DeleteDocument deletes a document from the vector store
func (r *RAGPipeline) DeleteDocument(ctx context.Context, id string) error {
	return r.config.VectorStore.Delete(id)
}

// UpdateDocument updates a document in the vector store
func (r *RAGPipeline) UpdateDocument(ctx context.Context, id, content string, metadata map[string]interface{}) error {
	// Generate new embedding
	vector, err := r.config.Embedder.Embed(content)
	if err != nil {
		return fmt.Errorf("failed to embed document: %w", err)
	}

	// Add content to metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["content"] = content

	// Update in vector store
	if err := r.config.VectorStore.Update(id, vector, metadata); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// CodeRAGPipeline extends RAGPipeline for code-specific operations
type CodeRAGPipeline struct {
	*RAGPipeline
}

// NewCodeRAGPipeline creates a new code RAG pipeline
func NewCodeRAGPipeline(config types.RAGConfig) (*CodeRAGPipeline, error) {
	pipeline, err := NewRAGPipeline(config)
	if err != nil {
		return nil, err
	}

	return &CodeRAGPipeline{
		RAGPipeline: pipeline,
	}, nil
}

// SearchCode searches for code snippets
func (c *CodeRAGPipeline) SearchCode(ctx context.Context, request types.CodeSearchRequest) (*types.CodeSearchResponse, error) {
	// Build filters
	filters := make(map[string]interface{})
	if request.Language != "" {
		filters["language"] = request.Language
	}
	if request.Repository != "" {
		filters["repository"] = request.Repository
	}

	// Set defaults
	if request.Limit == 0 {
		request.Limit = 5
	}
	if request.MinScore == 0 {
		request.MinScore = 0.7
	}

	// Override config for this search
	originalTopK := c.config.TopK
	originalMinScore := c.config.MinScore
	c.config.TopK = request.Limit
	c.config.MinScore = request.MinScore

	// Retrieve results
	ragResult, err := c.Retrieve(ctx, request.Query, filters)
	if err != nil {
		return nil, err
	}

	// Restore original config
	c.config.TopK = originalTopK
	c.config.MinScore = originalMinScore

	// Convert to code search results
	results := make([]types.CodeSearchResult, 0, len(ragResult.Results))
	for _, result := range ragResult.Results {
		codeResult := types.CodeSearchResult{
			Score: result.Score,
		}

		// Extract code from metadata
		if code, ok := result.Metadata["content"].(string); ok {
			codeResult.Code = code
		}
		if lang, ok := result.Metadata["language"].(string); ok {
			codeResult.Language = lang
		}
		if desc, ok := result.Metadata["description"].(string); ok {
			codeResult.Description = desc
		}

		results = append(results, codeResult)
	}

	return &types.CodeSearchResponse{
		Results:      results,
		TotalResults: len(results),
		Query:        request.Query,
		Metadata: map[string]interface{}{
			"retrieval_time_ms": ragResult.Metadata.RetrievalTime,
			"embedding_time_ms": ragResult.Metadata.EmbeddingTime,
		},
	}, nil
}

// AddCode adds a code snippet to the vector store
func (c *CodeRAGPipeline) AddCode(ctx context.Context, embedding types.CodeEmbedding) error {
	metadata := map[string]interface{}{
		"content":     embedding.Code,
		"language":    embedding.Language,
		"description": embedding.Description,
		"file_name":   embedding.Metadata.FileName,
		"file_path":   embedding.Metadata.FilePath,
		"repository":  embedding.Metadata.Repository,
	}

	return c.config.VectorStore.Store(embedding.ID, embedding.Vector, metadata)
}

// DocumentRAGPipeline extends RAGPipeline for documentation-specific operations
type DocumentRAGPipeline struct {
	*RAGPipeline
}

// NewDocumentRAGPipeline creates a new documentation RAG pipeline
func NewDocumentRAGPipeline(config types.RAGConfig) (*DocumentRAGPipeline, error) {
	pipeline, err := NewRAGPipeline(config)
	if err != nil {
		return nil, err
	}

	return &DocumentRAGPipeline{
		RAGPipeline: pipeline,
	}, nil
}

// SearchDocumentation searches for documentation
func (d *DocumentRAGPipeline) SearchDocumentation(ctx context.Context, request types.DocumentSearchRequest) (*types.DocumentSearchResponse, error) {
	// Build filters
	filters := make(map[string]interface{})
	if request.Language != "" {
		filters["language"] = request.Language
	}
	if request.Source != "" {
		filters["source"] = request.Source
	}
	if request.Category != "" {
		filters["category"] = request.Category
	}

	// Set defaults
	if request.Limit == 0 {
		request.Limit = 5
	}
	if request.MinScore == 0 {
		request.MinScore = 0.7
	}

	// Override config
	originalTopK := d.config.TopK
	originalMinScore := d.config.MinScore
	d.config.TopK = request.Limit
	d.config.MinScore = request.MinScore

	// Retrieve results
	ragResult, err := d.Retrieve(ctx, request.Query, filters)
	if err != nil {
		return nil, err
	}

	// Restore config
	d.config.TopK = originalTopK
	d.config.MinScore = originalMinScore

	// Convert to documentation search results
	results := make([]types.DocumentSearchResult, 0, len(ragResult.Results))
	for _, result := range ragResult.Results {
		docResult := types.DocumentSearchResult{
			Score: result.Score,
		}

		// Extract from metadata
		if content, ok := result.Metadata["content"].(string); ok {
			docResult.Content = content
		}
		if title, ok := result.Metadata["title"].(string); ok {
			docResult.Title = title
		}
		if source, ok := result.Metadata["source"].(string); ok {
			docResult.Source = source
		}
		if url, ok := result.Metadata["url"].(string); ok {
			docResult.URL = url
		}

		results = append(results, docResult)
	}

	return &types.DocumentSearchResponse{
		Results:      results,
		TotalResults: len(results),
		Query:        request.Query,
		Metadata: map[string]interface{}{
			"retrieval_time_ms": ragResult.Metadata.RetrievalTime,
			"embedding_time_ms": ragResult.Metadata.EmbeddingTime,
		},
	}, nil
}

// AddDocumentation adds documentation to the vector store
func (d *DocumentRAGPipeline) AddDocumentation(ctx context.Context, embedding types.DocumentEmbedding) error {
	metadata := map[string]interface{}{
		"content":  embedding.Content,
		"title":    embedding.Title,
		"source":   embedding.Source,
		"language": embedding.Language,
		"url":      embedding.Metadata.URL,
		"category": embedding.Metadata.Category,
		"version":  embedding.Metadata.Version,
	}

	return d.config.VectorStore.Store(embedding.ID, embedding.Vector, metadata)
}
