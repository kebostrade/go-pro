package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// OpenAIEmbedder implements the Embedder interface using OpenAI's API
type OpenAIEmbedder struct {
	apiKey     string
	model      string
	dimensions int
	client     *http.Client
}

// OpenAIEmbedderConfig holds configuration for OpenAI embedder
type OpenAIEmbedderConfig struct {
	APIKey     string
	Model      string // e.g., "text-embedding-3-small", "text-embedding-ada-002"
	Dimensions int    // Optional: for models that support custom dimensions
}

// NewOpenAIEmbedder creates a new OpenAI embedder
func NewOpenAIEmbedder(config OpenAIEmbedderConfig) (*OpenAIEmbedder, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	model := config.Model
	if model == "" {
		model = "text-embedding-3-small" // Default model
	}

	dimensions := config.Dimensions
	if dimensions == 0 {
		// Default dimensions for common models
		switch model {
		case "text-embedding-3-small":
			dimensions = 1536
		case "text-embedding-3-large":
			dimensions = 3072
		case "text-embedding-ada-002":
			dimensions = 1536
		default:
			dimensions = 1536
		}
	}

	return &OpenAIEmbedder{
		apiKey:     config.APIKey,
		model:      model,
		dimensions: dimensions,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Embed generates an embedding for text
func (e *OpenAIEmbedder) Embed(text string) ([]float64, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	embeddings, err := e.EmbedBatch([]string{text})
	if err != nil {
		return nil, err
	}

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return embeddings[0], nil
}

// EmbedBatch generates embeddings for multiple texts
func (e *OpenAIEmbedder) EmbedBatch(texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts cannot be empty")
	}

	// Prepare request
	reqBody := map[string]interface{}{
		"input": texts,
		"model": e.model,
	}

	// Add dimensions if supported
	if e.model == "text-embedding-3-small" || e.model == "text-embedding-3-large" {
		reqBody["dimensions"] = e.dimensions
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	// Send request
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result openAIEmbeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract embeddings
	embeddings := make([][]float64, len(result.Data))
	for i, data := range result.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}

// Dimensions returns the embedding dimension
func (e *OpenAIEmbedder) Dimensions() int {
	return e.dimensions
}

// openAIEmbeddingResponse represents the API response
type openAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// MockEmbedder is a mock embedder for testing
type MockEmbedder struct {
	dimensions int
}

// NewMockEmbedder creates a new mock embedder
func NewMockEmbedder(dimensions int) *MockEmbedder {
	if dimensions == 0 {
		dimensions = 1536
	}
	return &MockEmbedder{
		dimensions: dimensions,
	}
}

// Embed generates a mock embedding
func (m *MockEmbedder) Embed(text string) ([]float64, error) {
	// Generate a simple hash-based embedding
	embedding := make([]float64, m.dimensions)
	
	// Simple hash function
	hash := 0
	for _, c := range text {
		hash = (hash*31 + int(c)) % 1000000
	}

	// Fill embedding with pseudo-random values based on hash
	for i := 0; i < m.dimensions; i++ {
		hash = (hash*1103515245 + 12345) % 1000000
		embedding[i] = float64(hash) / 1000000.0
	}

	return embedding, nil
}

// EmbedBatch generates mock embeddings for multiple texts
func (m *MockEmbedder) EmbedBatch(texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))
	for i, text := range texts {
		embedding, err := m.Embed(text)
		if err != nil {
			return nil, err
		}
		embeddings[i] = embedding
	}
	return embeddings, nil
}

// Dimensions returns the embedding dimension
func (m *MockEmbedder) Dimensions() int {
	return m.dimensions
}

// CodeEmbedder wraps an embedder with code-specific preprocessing
type CodeEmbedder struct {
	embedder types.Embedder
}

// NewCodeEmbedder creates a new code embedder
func NewCodeEmbedder(embedder types.Embedder) *CodeEmbedder {
	return &CodeEmbedder{
		embedder: embedder,
	}
}

// Embed generates an embedding for code
func (c *CodeEmbedder) Embed(code string) ([]float64, error) {
	// Preprocess code (remove comments, normalize whitespace, etc.)
	processed := c.preprocessCode(code)
	return c.embedder.Embed(processed)
}

// EmbedBatch generates embeddings for multiple code snippets
func (c *CodeEmbedder) EmbedBatch(codes []string) ([][]float64, error) {
	processed := make([]string, len(codes))
	for i, code := range codes {
		processed[i] = c.preprocessCode(code)
	}
	return c.embedder.EmbedBatch(processed)
}

// Dimensions returns the embedding dimension
func (c *CodeEmbedder) Dimensions() int {
	return c.embedder.Dimensions()
}

// preprocessCode preprocesses code for embedding
func (c *CodeEmbedder) preprocessCode(code string) string {
	// TODO: Implement code preprocessing
	// - Remove comments
	// - Normalize whitespace
	// - Extract function signatures
	// - etc.
	return code
}

// EmbedWithContext generates an embedding with additional context
func (e *OpenAIEmbedder) EmbedWithContext(text, context string) ([]float64, error) {
	// Combine text with context
	combined := fmt.Sprintf("Context: %s\n\nText: %s", context, text)
	return e.Embed(combined)
}

// EmbedCode generates an embedding specifically for code
func (e *OpenAIEmbedder) EmbedCode(code, language string) ([]float64, error) {
	// Add language context
	text := fmt.Sprintf("Language: %s\n\nCode:\n%s", language, code)
	return e.Embed(text)
}

// EmbedDocumentation generates an embedding for documentation
func (e *OpenAIEmbedder) EmbedDocumentation(content, title, source string) ([]float64, error) {
	// Combine documentation elements
	text := fmt.Sprintf("Title: %s\nSource: %s\n\nContent:\n%s", title, source, content)
	return e.Embed(text)
}

// EstimateTokens estimates the number of tokens in text
func (e *OpenAIEmbedder) EstimateTokens(text string) int {
	// Rough estimate: ~4 characters per token
	return len(text) / 4
}

// ChunkText splits text into chunks suitable for embedding
func (e *OpenAIEmbedder) ChunkText(text string, maxTokens int) []string {
	// Simple chunking by character count
	maxChars := maxTokens * 4
	
	if len(text) <= maxChars {
		return []string{text}
	}

	chunks := make([]string, 0)
	for i := 0; i < len(text); i += maxChars {
		end := i + maxChars
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	return chunks
}

