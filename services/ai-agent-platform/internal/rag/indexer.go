package rag

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// CodeIndexer indexes code files for semantic search
type CodeIndexer struct {
	pipeline *CodeRAGPipeline
	embedder types.Embedder
}

// NewCodeIndexer creates a new code indexer
func NewCodeIndexer(pipeline *CodeRAGPipeline, embedder types.Embedder) *CodeIndexer {
	return &CodeIndexer{
		pipeline: pipeline,
		embedder: embedder,
	}
}

// IndexDirectory indexes all code files in a directory
func (i *CodeIndexer) IndexDirectory(ctx context.Context, dirPath string, options IndexOptions) error {
	// Walk directory
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file should be indexed
		if !i.shouldIndex(path, options) {
			return nil
		}

		// Index file
		if err := i.IndexFile(ctx, path, options); err != nil {
			// Log error but continue
			fmt.Printf("Error indexing %s: %v\n", path, err)
		}

		return nil
	})
}

// IndexFile indexes a single code file
func (i *CodeIndexer) IndexFile(ctx context.Context, filePath string, options IndexOptions) error {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	code := string(content)

	// Detect language
	language := i.detectLanguage(filePath)

	// Generate ID
	id := i.generateID(filePath, code)

	// Generate embedding
	vector, err := i.embedder.Embed(code)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Create code embedding
	embedding := types.CodeEmbedding{
		ID:       id,
		Code:     code,
		Language: language,
		Vector:   vector,
		Metadata: types.CodeEmbeddingMetadata{
			FileName:   filepath.Base(filePath),
			FilePath:   filePath,
			Repository: options.Repository,
			Branch:     options.Branch,
			Commit:     options.Commit,
			LOC:        strings.Count(code, "\n") + 1,
		},
		CreatedAt: time.Now().Unix(),
	}

	// Add to pipeline
	return i.pipeline.AddCode(ctx, embedding)
}

// IndexOptions holds options for indexing
type IndexOptions struct {
	// Repository repository URL
	Repository string

	// Branch git branch
	Branch string

	// Commit git commit hash
	Commit string

	// IncludePatterns file patterns to include
	IncludePatterns []string

	// ExcludePatterns file patterns to exclude
	ExcludePatterns []string

	// MaxFileSize maximum file size in bytes
	MaxFileSize int64

	// ChunkSize chunk size for large files
	ChunkSize int
}

// shouldIndex checks if a file should be indexed
func (i *CodeIndexer) shouldIndex(path string, options IndexOptions) bool {
	// Check file size
	if options.MaxFileSize > 0 {
		info, err := os.Stat(path)
		if err != nil || info.Size() > options.MaxFileSize {
			return false
		}
	}

	// Check exclude patterns
	for _, pattern := range options.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return false
		}
	}

	// Check include patterns
	if len(options.IncludePatterns) > 0 {
		for _, pattern := range options.IncludePatterns {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				return true
			}
		}
		return false
	}

	// Check if it's a code file
	return i.isCodeFile(path)
}

// isCodeFile checks if a file is a code file
func (i *CodeIndexer) isCodeFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	codeExtensions := []string{
		".go", ".py", ".js", ".ts", ".jsx", ".tsx",
		".java", ".c", ".cpp", ".h", ".hpp",
		".rs", ".rb", ".php", ".swift", ".kt",
		".cs", ".scala", ".r", ".m", ".sh",
	}

	for _, codeExt := range codeExtensions {
		if ext == codeExt {
			return true
		}
	}

	return false
}

// detectLanguage detects the programming language from file extension
func (i *CodeIndexer) detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	
	languageMap := map[string]string{
		".go":   "go",
		".py":   "python",
		".js":   "javascript",
		".ts":   "typescript",
		".jsx":  "javascript",
		".tsx":  "typescript",
		".java": "java",
		".c":    "c",
		".cpp":  "cpp",
		".h":    "c",
		".hpp":  "cpp",
		".rs":   "rust",
		".rb":   "ruby",
		".php":  "php",
		".swift": "swift",
		".kt":   "kotlin",
		".cs":   "csharp",
		".scala": "scala",
		".r":    "r",
		".m":    "objective-c",
		".sh":   "bash",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}

	return "unknown"
}

// generateID generates a unique ID for a code file
func (i *CodeIndexer) generateID(path, content string) string {
	hash := sha256.Sum256([]byte(path + content))
	return fmt.Sprintf("%x", hash[:16])
}

// DocumentIndexer indexes documentation for semantic search
type DocumentIndexer struct {
	pipeline *DocumentRAGPipeline
	embedder types.Embedder
}

// NewDocumentIndexer creates a new documentation indexer
func NewDocumentIndexer(pipeline *DocumentRAGPipeline, embedder types.Embedder) *DocumentIndexer {
	return &DocumentIndexer{
		pipeline: pipeline,
		embedder: embedder,
	}
}

// IndexDocument indexes a documentation document
func (i *DocumentIndexer) IndexDocument(ctx context.Context, doc DocumentToIndex) error {
	// Generate embedding
	vector, err := i.embedder.Embed(doc.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Create document embedding
	embedding := types.DocumentEmbedding{
		ID:       doc.ID,
		Content:  doc.Content,
		Title:    doc.Title,
		Source:   doc.Source,
		Language: doc.Language,
		Vector:   vector,
		Metadata: types.DocumentEmbeddingMetadata{
			URL:         doc.URL,
			Section:     doc.Section,
			Category:    doc.Category,
			Tags:        doc.Tags,
			Version:     doc.Version,
			LastUpdated: time.Now().Unix(),
		},
		CreatedAt: time.Now().Unix(),
	}

	// Add to pipeline
	return i.pipeline.AddDocumentation(ctx, embedding)
}

// IndexDocuments indexes multiple documentation documents
func (i *DocumentIndexer) IndexDocuments(ctx context.Context, docs []DocumentToIndex) error {
	for _, doc := range docs {
		if err := i.IndexDocument(ctx, doc); err != nil {
			return fmt.Errorf("failed to index document %s: %w", doc.ID, err)
		}
	}
	return nil
}

// DocumentToIndex represents a document to be indexed
type DocumentToIndex struct {
	ID       string
	Content  string
	Title    string
	Source   string
	Language string
	URL      string
	Section  string
	Category string
	Tags     []string
	Version  string
}

// ChunkDocument splits a large document into chunks
func (i *DocumentIndexer) ChunkDocument(content string, chunkSize int, overlap int) []string {
	if len(content) <= chunkSize {
		return []string{content}
	}

	chunks := make([]string, 0)
	start := 0

	for start < len(content) {
		end := start + chunkSize
		if end > len(content) {
			end = len(content)
		}

		chunks = append(chunks, content[start:end])

		// Move start with overlap
		start = end - overlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}

// IndexDocumentWithChunking indexes a large document by chunking it
func (i *DocumentIndexer) IndexDocumentWithChunking(ctx context.Context, doc DocumentToIndex, chunkSize, overlap int) error {
	chunks := i.ChunkDocument(doc.Content, chunkSize, overlap)

	for idx, chunk := range chunks {
		chunkDoc := doc
		chunkDoc.ID = fmt.Sprintf("%s_chunk_%d", doc.ID, idx)
		chunkDoc.Content = chunk
		chunkDoc.Section = fmt.Sprintf("%s (Part %d/%d)", doc.Section, idx+1, len(chunks))

		if err := i.IndexDocument(ctx, chunkDoc); err != nil {
			return fmt.Errorf("failed to index chunk %d: %w", idx, err)
		}
	}

	return nil
}

// ReindexAll reindexes all documents
func (i *DocumentIndexer) ReindexAll(ctx context.Context, docs []DocumentToIndex) error {
	// Clear existing documents
	// Note: This would require a Clear method on the vector store

	// Index all documents
	return i.IndexDocuments(ctx, docs)
}

// GetIndexStats returns statistics about the index
func (i *DocumentIndexer) GetIndexStats(ctx context.Context) (map[string]interface{}, error) {
	count, err := i.pipeline.config.VectorStore.Count()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_documents": count,
		"last_updated":    time.Now().Unix(),
	}, nil
}

