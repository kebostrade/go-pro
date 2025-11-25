package types

import (
	"context"
	"time"
)

// Memory represents a memory system for agents
type Memory interface {
	// Add adds a message to memory
	Add(ctx context.Context, message Message) error

	// Get retrieves messages from memory
	Get(ctx context.Context, limit int) ([]Message, error)

	// Clear clears all messages from memory
	Clear(ctx context.Context) error

	// GetBySessionID retrieves messages for a specific session
	GetBySessionID(ctx context.Context, sessionID string, limit int) ([]Message, error)

	// Search searches memory for relevant messages
	Search(ctx context.Context, query string, limit int) ([]Message, error)
}

// VectorMemory extends Memory with vector search capabilities
type VectorMemory interface {
	Memory

	// AddWithEmbedding adds a message with its embedding
	AddWithEmbedding(ctx context.Context, message Message, embedding []float32) error

	// SearchSimilar finds similar messages using vector similarity
	SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]Message, error)

	// GetEmbedding retrieves the embedding for a message
	GetEmbedding(ctx context.Context, messageID string) ([]float32, error)
}

// BufferMemory stores recent messages in a buffer
type BufferMemory interface {
	Memory

	// GetRecent retrieves the N most recent messages
	GetRecent(ctx context.Context, n int) ([]Message, error)

	// GetWindow retrieves messages within a time window
	GetWindow(ctx context.Context, start, end time.Time) ([]Message, error)
}

// SummaryMemory stores summarized conversation history
type SummaryMemory interface {
	Memory

	// Summarize creates a summary of the conversation
	Summarize(ctx context.Context) (string, error)

	// GetSummary retrieves the current summary
	GetSummary(ctx context.Context) (string, error)

	// UpdateSummary updates the summary with new messages
	UpdateSummary(ctx context.Context, newMessages []Message) error
}

// EntityMemory tracks entities mentioned in conversations
type EntityMemory interface {
	Memory

	// AddEntity adds or updates an entity
	AddEntity(ctx context.Context, entity Entity) error

	// GetEntity retrieves an entity by name
	GetEntity(ctx context.Context, name string) (*Entity, error)

	// GetEntities retrieves all tracked entities
	GetEntities(ctx context.Context) ([]Entity, error)

	// UpdateEntity updates an entity's information
	UpdateEntity(ctx context.Context, entity Entity) error
}

// Entity represents an entity tracked in memory
type Entity struct {
	// Name of the entity
	Name string `json:"name"`

	// Type of entity (person, organization, location, etc.)
	Type string `json:"type"`

	// Attributes associated with the entity
	Attributes map[string]interface{} `json:"attributes"`

	// FirstMentioned when the entity was first mentioned
	FirstMentioned time.Time `json:"first_mentioned"`

	// LastMentioned when the entity was last mentioned
	LastMentioned time.Time `json:"last_mentioned"`

	// MentionCount how many times the entity was mentioned
	MentionCount int `json:"mention_count"`
}

// MemoryConfig holds configuration for memory systems
type MemoryConfig struct {
	// Type of memory (buffer, summary, vector, entity)
	Type MemoryType `json:"type"`

	// MaxMessages maximum number of messages to store
	MaxMessages int `json:"max_messages"`

	// TTL time-to-live for messages
	TTL time.Duration `json:"ttl"`

	// SessionID for session-specific memory
	SessionID string `json:"session_id,omitempty"`

	// VectorStoreConfig for vector memory
	VectorStoreConfig *VectorStoreConfig `json:"vector_store_config,omitempty"`

	// SummaryConfig for summary memory
	SummaryConfig *SummaryConfig `json:"summary_config,omitempty"`
}

// MemoryType defines the type of memory
type MemoryType string

const (
	MemoryTypeBuffer  MemoryType = "buffer"
	MemoryTypeSummary MemoryType = "summary"
	MemoryTypeVector  MemoryType = "vector"
	MemoryTypeEntity  MemoryType = "entity"
)

// VectorStoreConfig holds configuration for vector stores
type VectorStoreConfig struct {
	// Provider (pgvector, qdrant, redis)
	Provider string `json:"provider"`

	// ConnectionString for the vector store
	ConnectionString string `json:"connection_string"`

	// CollectionName for the vector collection
	CollectionName string `json:"collection_name"`

	// EmbeddingDimension size of embeddings
	EmbeddingDimension int `json:"embedding_dimension"`

	// SimilarityMetric (cosine, euclidean, dot_product)
	SimilarityMetric string `json:"similarity_metric"`
}

// SummaryConfig holds configuration for summary memory
type SummaryConfig struct {
	// MaxTokens maximum tokens in summary
	MaxTokens int `json:"max_tokens"`

	// SummaryInterval how often to update summary
	SummaryInterval int `json:"summary_interval"`

	// LLMProvider for generating summaries
	LLMProvider string `json:"llm_provider"`

	// PromptTemplate for summary generation
	PromptTemplate string `json:"prompt_template,omitempty"`
}

// ConversationHistory represents a conversation history
type ConversationHistory struct {
	// SessionID identifies the conversation
	SessionID string `json:"session_id"`

	// Messages in the conversation
	Messages []Message `json:"messages"`

	// Summary of the conversation
	Summary string `json:"summary,omitempty"`

	// Entities mentioned in the conversation
	Entities []Entity `json:"entities,omitempty"`

	// StartTime when the conversation started
	StartTime time.Time `json:"start_time"`

	// LastActivity when the last message was added
	LastActivity time.Time `json:"last_activity"`

	// Metadata additional information
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewMemoryConfig creates a new MemoryConfig with defaults
func NewMemoryConfig(memoryType MemoryType) MemoryConfig {
	return MemoryConfig{
		Type:        memoryType,
		MaxMessages: 100,
		TTL:         24 * time.Hour,
	}
}

// NewEntity creates a new Entity
func NewEntity(name, entityType string) Entity {
	return Entity{
		Name:           name,
		Type:           entityType,
		Attributes:     make(map[string]interface{}),
		FirstMentioned: time.Now(),
		LastMentioned:  time.Now(),
		MentionCount:   1,
	}
}

