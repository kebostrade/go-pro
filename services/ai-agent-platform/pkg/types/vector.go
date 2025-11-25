package types

// VectorStore defines the interface for vector storage and retrieval
type VectorStore interface {
	// Store stores a vector with metadata
	Store(id string, vector []float64, metadata map[string]interface{}) error

	// Search performs similarity search
	Search(query []float64, limit int, filters map[string]interface{}) ([]VectorSearchResult, error)

	// Delete deletes a vector by ID
	Delete(id string) error

	// Update updates a vector and its metadata
	Update(id string, vector []float64, metadata map[string]interface{}) error

	// Get retrieves a vector by ID
	Get(id string) (*VectorSearchResult, error)

	// Count returns the total number of vectors
	Count() (int64, error)

	// Clear removes all vectors
	Clear() error
}

// VectorSearchResult represents a search result
type VectorSearchResult struct {
	// ID unique identifier
	ID string `json:"id"`

	// Score similarity score (0-1, higher is more similar)
	Score float64 `json:"score"`

	// Vector the embedding vector
	Vector []float64 `json:"vector,omitempty"`

	// Metadata associated metadata
	Metadata map[string]interface{} `json:"metadata"`

	// Content the original content (if stored in metadata)
	Content string `json:"content,omitempty"`
}

// Embedder defines the interface for generating embeddings
type Embedder interface {
	// Embed generates an embedding for text
	Embed(text string) ([]float64, error)

	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(texts []string) ([][]float64, error)

	// Dimensions returns the embedding dimension
	Dimensions() int
}

// CodeEmbedding represents an embedded code snippet
type CodeEmbedding struct {
	// ID unique identifier
	ID string `json:"id"`

	// Code the code snippet
	Code string `json:"code"`

	// Language programming language
	Language string `json:"language"`

	// Description code description
	Description string `json:"description,omitempty"`

	// Vector embedding vector
	Vector []float64 `json:"vector"`

	// Metadata additional metadata
	Metadata CodeEmbeddingMetadata `json:"metadata"`

	// CreatedAt timestamp
	CreatedAt int64 `json:"created_at"`
}

// CodeEmbeddingMetadata holds metadata for code embeddings
type CodeEmbeddingMetadata struct {
	// FileName source file name
	FileName string `json:"file_name,omitempty"`

	// FilePath source file path
	FilePath string `json:"file_path,omitempty"`

	// LineStart starting line number
	LineStart int `json:"line_start,omitempty"`

	// LineEnd ending line number
	LineEnd int `json:"line_end,omitempty"`

	// FunctionName function or method name
	FunctionName string `json:"function_name,omitempty"`

	// ClassName class name
	ClassName string `json:"class_name,omitempty"`

	// PackageName package or module name
	PackageName string `json:"package_name,omitempty"`

	// Tags custom tags
	Tags []string `json:"tags,omitempty"`

	// Complexity code complexity score
	Complexity int `json:"complexity,omitempty"`

	// LOC lines of code
	LOC int `json:"loc,omitempty"`

	// Repository repository URL
	Repository string `json:"repository,omitempty"`

	// Branch git branch
	Branch string `json:"branch,omitempty"`

	// Commit git commit hash
	Commit string `json:"commit,omitempty"`
}

// DocumentEmbedding represents an embedded documentation chunk
type DocumentEmbedding struct {
	// ID unique identifier
	ID string `json:"id"`

	// Content documentation content
	Content string `json:"content"`

	// Title document title
	Title string `json:"title,omitempty"`

	// Source documentation source
	Source string `json:"source"`

	// Language programming language
	Language string `json:"language,omitempty"`

	// Vector embedding vector
	Vector []float64 `json:"vector"`

	// Metadata additional metadata
	Metadata DocumentEmbeddingMetadata `json:"metadata"`

	// CreatedAt timestamp
	CreatedAt int64 `json:"created_at"`
}

// DocumentEmbeddingMetadata holds metadata for documentation embeddings
type DocumentEmbeddingMetadata struct {
	// URL source URL
	URL string `json:"url,omitempty"`

	// Section document section
	Section string `json:"section,omitempty"`

	// Subsection document subsection
	Subsection string `json:"subsection,omitempty"`

	// Category documentation category
	Category string `json:"category,omitempty"`

	// Tags custom tags
	Tags []string `json:"tags,omitempty"`

	// Version documentation version
	Version string `json:"version,omitempty"`

	// LastUpdated last update timestamp
	LastUpdated int64 `json:"last_updated,omitempty"`
}

// RAGConfig holds configuration for RAG pipeline
type RAGConfig struct {
	// VectorStore the vector store to use
	VectorStore VectorStore

	// Embedder the embedder to use
	Embedder Embedder

	// TopK number of results to retrieve
	TopK int

	// MinScore minimum similarity score
	MinScore float64

	// MaxTokens maximum tokens for context
	MaxTokens int

	// IncludeMetadata whether to include metadata in results
	IncludeMetadata bool

	// RerankerEnabled whether to use reranking
	RerankerEnabled bool
}

// RAGResult represents a RAG retrieval result
type RAGResult struct {
	// Query the original query
	Query string `json:"query"`

	// Results retrieved results
	Results []VectorSearchResult `json:"results"`

	// Context formatted context for LLM
	Context string `json:"context"`

	// Metadata retrieval metadata
	Metadata RAGMetadata `json:"metadata"`
}

// RAGMetadata holds metadata about RAG retrieval
type RAGMetadata struct {
	// TotalResults total number of results found
	TotalResults int `json:"total_results"`

	// RetrievalTime time taken for retrieval (ms)
	RetrievalTime int64 `json:"retrieval_time_ms"`

	// EmbeddingTime time taken for embedding (ms)
	EmbeddingTime int64 `json:"embedding_time_ms"`

	// ContextTokens estimated tokens in context
	ContextTokens int `json:"context_tokens"`

	// Filters applied filters
	Filters map[string]interface{} `json:"filters,omitempty"`
}

// CodeSearchRequest represents a code search request
type CodeSearchRequest struct {
	// Query search query
	Query string `json:"query"`

	// Language filter by language
	Language string `json:"language,omitempty"`

	// Repository filter by repository
	Repository string `json:"repository,omitempty"`

	// Tags filter by tags
	Tags []string `json:"tags,omitempty"`

	// MinComplexity minimum complexity
	MinComplexity int `json:"min_complexity,omitempty"`

	// MaxComplexity maximum complexity
	MaxComplexity int `json:"max_complexity,omitempty"`

	// Limit number of results
	Limit int `json:"limit,omitempty"`

	// MinScore minimum similarity score
	MinScore float64 `json:"min_score,omitempty"`
}

// CodeSearchResponse represents a code search response
type CodeSearchResponse struct {
	// Results search results
	Results []CodeSearchResult `json:"results"`

	// TotalResults total number of results
	TotalResults int `json:"total_results"`

	// Query the original query
	Query string `json:"query"`

	// Metadata search metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CodeSearchResult represents a single code search result
type CodeSearchResult struct {
	// Code the code snippet
	Code string `json:"code"`

	// Language programming language
	Language string `json:"language"`

	// Description code description
	Description string `json:"description,omitempty"`

	// Score similarity score
	Score float64 `json:"score"`

	// Metadata code metadata
	Metadata CodeEmbeddingMetadata `json:"metadata"`

	// Highlights highlighted portions
	Highlights []string `json:"highlights,omitempty"`
}

// DocumentSearchRequest represents a documentation search request
type DocumentSearchRequest struct {
	// Query search query
	Query string `json:"query"`

	// Language filter by language
	Language string `json:"language,omitempty"`

	// Source filter by source
	Source string `json:"source,omitempty"`

	// Category filter by category
	Category string `json:"category,omitempty"`

	// Version filter by version
	Version string `json:"version,omitempty"`

	// Limit number of results
	Limit int `json:"limit,omitempty"`

	// MinScore minimum similarity score
	MinScore float64 `json:"min_score,omitempty"`
}

// DocumentSearchResponse represents a documentation search response
type DocumentSearchResponse struct {
	// Results search results
	Results []DocumentSearchResult `json:"results"`

	// TotalResults total number of results
	TotalResults int `json:"total_results"`

	// Query the original query
	Query string `json:"query"`

	// Metadata search metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// DocumentSearchResult represents a single documentation search result
type DocumentSearchResult struct {
	// Content documentation content
	Content string `json:"content"`

	// Title document title
	Title string `json:"title,omitempty"`

	// Source documentation source
	Source string `json:"source"`

	// Score similarity score
	Score float64 `json:"score"`

	// Metadata document metadata
	Metadata DocumentEmbeddingMetadata `json:"metadata"`

	// URL source URL
	URL string `json:"url,omitempty"`
}

