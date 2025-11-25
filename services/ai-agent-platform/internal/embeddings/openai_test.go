package embeddings

import (
	"testing"
)

func TestMockEmbedder_Embed(t *testing.T) {
	embedder := NewMockEmbedder(1536)

	vector, err := embedder.Embed("test text")
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vector) != 1536 {
		t.Errorf("Expected vector length 1536, got %d", len(vector))
	}

	// Check all values are between 0 and 1
	for i, v := range vector {
		if v < 0 || v > 1 {
			t.Errorf("Vector value at index %d out of range: %f", i, v)
		}
	}
}

func TestMockEmbedder_EmbedBatch(t *testing.T) {
	embedder := NewMockEmbedder(1536)

	texts := []string{"text1", "text2", "text3"}
	vectors, err := embedder.EmbedBatch(texts)
	if err != nil {
		t.Fatalf("EmbedBatch failed: %v", err)
	}

	if len(vectors) != 3 {
		t.Errorf("Expected 3 vectors, got %d", len(vectors))
	}

	for i, vector := range vectors {
		if len(vector) != 1536 {
			t.Errorf("Vector %d has wrong length: %d", i, len(vector))
		}
	}
}

func TestMockEmbedder_Dimensions(t *testing.T) {
	embedder := NewMockEmbedder(768)

	dims := embedder.Dimensions()
	if dims != 768 {
		t.Errorf("Expected dimensions 768, got %d", dims)
	}
}

func TestMockEmbedder_ConsistentEmbeddings(t *testing.T) {
	embedder := NewMockEmbedder(1536)

	// Same text should produce same embedding
	vector1, _ := embedder.Embed("test")
	vector2, _ := embedder.Embed("test")

	if len(vector1) != len(vector2) {
		t.Error("Vectors have different lengths")
	}

	for i := range vector1 {
		if vector1[i] != vector2[i] {
			t.Error("Same text produced different embeddings")
			break
		}
	}
}

func TestMockEmbedder_DifferentTexts(t *testing.T) {
	embedder := NewMockEmbedder(1536)

	vector1, _ := embedder.Embed("text1")
	vector2, _ := embedder.Embed("text2")

	// Different texts should produce different embeddings
	same := true
	for i := range vector1 {
		if vector1[i] != vector2[i] {
			same = false
			break
		}
	}

	if same {
		t.Error("Different texts produced identical embeddings")
	}
}

func TestCodeEmbedder_Embed(t *testing.T) {
	mockEmbedder := NewMockEmbedder(1536)
	codeEmbedder := NewCodeEmbedder(mockEmbedder)

	code := `package main

import "fmt"

func main() {
    fmt.Println("Hello")
}`

	vector, err := codeEmbedder.Embed(code)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vector) != 1536 {
		t.Errorf("Expected vector length 1536, got %d", len(vector))
	}
}

func TestCodeEmbedder_EmbedBatch(t *testing.T) {
	mockEmbedder := NewMockEmbedder(1536)
	codeEmbedder := NewCodeEmbedder(mockEmbedder)

	codes := []string{
		"func foo() {}",
		"def bar(): pass",
		"function baz() {}",
	}

	vectors, err := codeEmbedder.EmbedBatch(codes)
	if err != nil {
		t.Fatalf("EmbedBatch failed: %v", err)
	}

	if len(vectors) != 3 {
		t.Errorf("Expected 3 vectors, got %d", len(vectors))
	}
}

func TestCodeEmbedder_Dimensions(t *testing.T) {
	mockEmbedder := NewMockEmbedder(768)
	codeEmbedder := NewCodeEmbedder(mockEmbedder)

	dims := codeEmbedder.Dimensions()
	if dims != 768 {
		t.Errorf("Expected dimensions 768, got %d", dims)
	}
}

func TestOpenAIEmbedder_EstimateTokens(t *testing.T) {
	embedder := &OpenAIEmbedder{}

	tests := []struct {
		text     string
		expected int
	}{
		{"", 0},
		{"test", 1},
		{"hello world", 2},
		{"this is a longer text", 5},
	}

	for _, tt := range tests {
		tokens := embedder.EstimateTokens(tt.text)
		if tokens != tt.expected {
			t.Errorf("For text '%s', expected %d tokens, got %d", tt.text, tt.expected, tokens)
		}
	}
}

func TestOpenAIEmbedder_ChunkText(t *testing.T) {
	embedder := &OpenAIEmbedder{}

	text := "This is a test text that needs to be chunked into smaller pieces"
	chunks := embedder.ChunkText(text, 5) // 5 tokens = ~20 chars

	if len(chunks) == 0 {
		t.Error("Expected at least one chunk")
	}

	// Verify all chunks together equal original text
	combined := ""
	for _, chunk := range chunks {
		combined += chunk
	}

	if combined != text {
		t.Error("Chunks don't combine to original text")
	}
}

func TestOpenAIEmbedder_ChunkText_SmallText(t *testing.T) {
	embedder := &OpenAIEmbedder{}

	text := "short"
	chunks := embedder.ChunkText(text, 100)

	if len(chunks) != 1 {
		t.Errorf("Expected 1 chunk for small text, got %d", len(chunks))
	}

	if chunks[0] != text {
		t.Error("Single chunk doesn't match original text")
	}
}

func TestNewOpenAIEmbedder_NoAPIKey(t *testing.T) {
	_, err := NewOpenAIEmbedder(OpenAIEmbedderConfig{
		APIKey: "",
	})

	if err == nil {
		t.Error("Expected error for empty API key")
	}
}

func TestNewOpenAIEmbedder_DefaultModel(t *testing.T) {
	embedder, err := NewOpenAIEmbedder(OpenAIEmbedderConfig{
		APIKey: "test-key",
	})

	if err != nil {
		t.Fatalf("Failed to create embedder: %v", err)
	}

	if embedder.model != "text-embedding-3-small" {
		t.Errorf("Expected default model 'text-embedding-3-small', got '%s'", embedder.model)
	}

	if embedder.dimensions != 1536 {
		t.Errorf("Expected default dimensions 1536, got %d", embedder.dimensions)
	}
}

func TestNewOpenAIEmbedder_CustomModel(t *testing.T) {
	embedder, err := NewOpenAIEmbedder(OpenAIEmbedderConfig{
		APIKey:     "test-key",
		Model:      "text-embedding-3-large",
		Dimensions: 3072,
	})

	if err != nil {
		t.Fatalf("Failed to create embedder: %v", err)
	}

	if embedder.model != "text-embedding-3-large" {
		t.Errorf("Expected model 'text-embedding-3-large', got '%s'", embedder.model)
	}

	if embedder.dimensions != 3072 {
		t.Errorf("Expected dimensions 3072, got %d", embedder.dimensions)
	}
}

func TestNewMockEmbedder_DefaultDimensions(t *testing.T) {
	embedder := NewMockEmbedder(0)

	if embedder.dimensions != 1536 {
		t.Errorf("Expected default dimensions 1536, got %d", embedder.dimensions)
	}
}

func TestMockEmbedder_EmptyText(t *testing.T) {
	embedder := NewMockEmbedder(1536)

	vector, err := embedder.Embed("")
	if err != nil {
		t.Fatalf("Embed failed for empty text: %v", err)
	}

	if len(vector) != 1536 {
		t.Errorf("Expected vector length 1536, got %d", len(vector))
	}
}

func TestMockEmbedder_LongText(t *testing.T) {
	embedder := NewMockEmbedder(1536)

	// Create a long text
	longText := ""
	for i := 0; i < 1000; i++ {
		longText += "word "
	}

	vector, err := embedder.Embed(longText)
	if err != nil {
		t.Fatalf("Embed failed for long text: %v", err)
	}

	if len(vector) != 1536 {
		t.Errorf("Expected vector length 1536, got %d", len(vector))
	}
}

func BenchmarkMockEmbedder_Embed(b *testing.B) {
	embedder := NewMockEmbedder(1536)
	text := "This is a test text for benchmarking"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = embedder.Embed(text)
	}
}

func BenchmarkMockEmbedder_EmbedBatch(b *testing.B) {
	embedder := NewMockEmbedder(1536)
	texts := []string{
		"text1", "text2", "text3", "text4", "text5",
		"text6", "text7", "text8", "text9", "text10",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = embedder.EmbedBatch(texts)
	}
}

