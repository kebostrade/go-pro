package vectorstore

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// MemoryVectorStore is an in-memory vector store implementation
type MemoryVectorStore struct {
	mu      sync.RWMutex
	vectors map[string]*vectorEntry
}

// vectorEntry represents a stored vector with metadata
type vectorEntry struct {
	ID       string
	Vector   []float64
	Metadata map[string]interface{}
}

// NewMemoryVectorStore creates a new in-memory vector store
func NewMemoryVectorStore() *MemoryVectorStore {
	return &MemoryVectorStore{
		vectors: make(map[string]*vectorEntry),
	}
}

// Store stores a vector with metadata
func (m *MemoryVectorStore) Store(id string, vector []float64, metadata map[string]interface{}) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if len(vector) == 0 {
		return fmt.Errorf("vector cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.vectors[id] = &vectorEntry{
		ID:       id,
		Vector:   vector,
		Metadata: metadata,
	}

	return nil
}

// Search performs similarity search using cosine similarity
func (m *MemoryVectorStore) Search(query []float64, limit int, filters map[string]interface{}) ([]types.VectorSearchResult, error) {
	if len(query) == 0 {
		return nil, fmt.Errorf("query vector cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Calculate similarities for all vectors
	results := make([]types.VectorSearchResult, 0)

	for _, entry := range m.vectors {
		// Apply filters
		if !m.matchesFilters(entry.Metadata, filters) {
			continue
		}

		// Calculate cosine similarity
		similarity := cosineSimilarity(query, entry.Vector)

		results = append(results, types.VectorSearchResult{
			ID:       entry.ID,
			Score:    similarity,
			Vector:   entry.Vector,
			Metadata: entry.Metadata,
		})
	}

	// Sort by score (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit results
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// Delete deletes a vector by ID
func (m *MemoryVectorStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.vectors[id]; !exists {
		return fmt.Errorf("vector with id %s not found", id)
	}

	delete(m.vectors, id)
	return nil
}

// Update updates a vector and its metadata
func (m *MemoryVectorStore) Update(id string, vector []float64, metadata map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.vectors[id]; !exists {
		return fmt.Errorf("vector with id %s not found", id)
	}

	m.vectors[id] = &vectorEntry{
		ID:       id,
		Vector:   vector,
		Metadata: metadata,
	}

	return nil
}

// Get retrieves a vector by ID
func (m *MemoryVectorStore) Get(id string) (*types.VectorSearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.vectors[id]
	if !exists {
		return nil, fmt.Errorf("vector with id %s not found", id)
	}

	return &types.VectorSearchResult{
		ID:       entry.ID,
		Score:    1.0, // Perfect match
		Vector:   entry.Vector,
		Metadata: entry.Metadata,
	}, nil
}

// Count returns the total number of vectors
func (m *MemoryVectorStore) Count() (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return int64(len(m.vectors)), nil
}

// Clear removes all vectors
func (m *MemoryVectorStore) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.vectors = make(map[string]*vectorEntry)
	return nil
}

// matchesFilters checks if metadata matches filters
func (m *MemoryVectorStore) matchesFilters(metadata, filters map[string]interface{}) bool {
	if len(filters) == 0 {
		return true
	}

	for key, filterValue := range filters {
		metadataValue, exists := metadata[key]
		if !exists {
			return false
		}

		// Simple equality check
		if metadataValue != filterValue {
			return false
		}
	}

	return true
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// euclideanDistance calculates Euclidean distance between two vectors
func euclideanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.MaxFloat64
	}

	var sum float64
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

// dotProduct calculates dot product of two vectors
func dotProduct(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var sum float64
	for i := 0; i < len(a); i++ {
		sum += a[i] * b[i]
	}

	return sum
}

// normalize normalizes a vector to unit length
func normalize(v []float64) []float64 {
	var norm float64
	for _, val := range v {
		norm += val * val
	}
	norm = math.Sqrt(norm)

	if norm == 0 {
		return v
	}

	normalized := make([]float64, len(v))
	for i, val := range v {
		normalized[i] = val / norm
	}

	return normalized
}

// BatchStore stores multiple vectors at once
func (m *MemoryVectorStore) BatchStore(entries []struct {
	ID       string
	Vector   []float64
	Metadata map[string]interface{}
}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range entries {
		if entry.ID == "" {
			return fmt.Errorf("id cannot be empty")
		}
		if len(entry.Vector) == 0 {
			return fmt.Errorf("vector cannot be empty for id %s", entry.ID)
		}

		m.vectors[entry.ID] = &vectorEntry{
			ID:       entry.ID,
			Vector:   entry.Vector,
			Metadata: entry.Metadata,
		}
	}

	return nil
}

// BatchDelete deletes multiple vectors at once
func (m *MemoryVectorStore) BatchDelete(ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, id := range ids {
		delete(m.vectors, id)
	}

	return nil
}

// GetAll returns all vectors (use with caution for large datasets)
func (m *MemoryVectorStore) GetAll() ([]types.VectorSearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]types.VectorSearchResult, 0, len(m.vectors))

	for _, entry := range m.vectors {
		results = append(results, types.VectorSearchResult{
			ID:       entry.ID,
			Score:    1.0,
			Vector:   entry.Vector,
			Metadata: entry.Metadata,
		})
	}

	return results, nil
}

// SearchByMetadata searches vectors by metadata only (no vector similarity)
func (m *MemoryVectorStore) SearchByMetadata(filters map[string]interface{}, limit int) ([]types.VectorSearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]types.VectorSearchResult, 0)

	for _, entry := range m.vectors {
		if m.matchesFilters(entry.Metadata, filters) {
			results = append(results, types.VectorSearchResult{
				ID:       entry.ID,
				Score:    1.0,
				Vector:   entry.Vector,
				Metadata: entry.Metadata,
			})
		}
	}

	// Limit results
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

