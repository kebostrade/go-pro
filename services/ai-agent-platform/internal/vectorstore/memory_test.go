package vectorstore

import (
	"testing"
)

func TestMemoryVectorStore_Store(t *testing.T) {
	store := NewMemoryVectorStore()

	vector := []float64{1.0, 2.0, 3.0}
	metadata := map[string]interface{}{
		"language": "go",
		"content":  "test code",
	}

	err := store.Store("test-1", vector, metadata)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Verify count
	count, _ := store.Count()
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestMemoryVectorStore_StoreEmptyID(t *testing.T) {
	store := NewMemoryVectorStore()

	err := store.Store("", []float64{1.0}, nil)
	if err == nil {
		t.Error("Expected error for empty ID")
	}
}

func TestMemoryVectorStore_StoreEmptyVector(t *testing.T) {
	store := NewMemoryVectorStore()

	err := store.Store("test-1", []float64{}, nil)
	if err == nil {
		t.Error("Expected error for empty vector")
	}
}

func TestMemoryVectorStore_Get(t *testing.T) {
	store := NewMemoryVectorStore()

	vector := []float64{1.0, 2.0, 3.0}
	metadata := map[string]interface{}{"key": "value"}

	store.Store("test-1", vector, metadata)

	result, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if result.ID != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", result.ID)
	}

	if result.Score != 1.0 {
		t.Errorf("Expected score 1.0, got %f", result.Score)
	}
}

func TestMemoryVectorStore_GetNotFound(t *testing.T) {
	store := NewMemoryVectorStore()

	_, err := store.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent ID")
	}
}

func TestMemoryVectorStore_Delete(t *testing.T) {
	store := NewMemoryVectorStore()

	store.Store("test-1", []float64{1.0}, nil)

	err := store.Delete("test-1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	count, _ := store.Count()
	if count != 0 {
		t.Errorf("Expected count 0 after delete, got %d", count)
	}
}

func TestMemoryVectorStore_Update(t *testing.T) {
	store := NewMemoryVectorStore()

	store.Store("test-1", []float64{1.0}, map[string]interface{}{"version": 1})

	newVector := []float64{2.0, 3.0}
	newMetadata := map[string]interface{}{"version": 2}

	err := store.Update("test-1", newVector, newMetadata)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	result, _ := store.Get("test-1")
	if len(result.Vector) != 2 {
		t.Errorf("Expected vector length 2, got %d", len(result.Vector))
	}

	if result.Metadata["version"] != 2 {
		t.Errorf("Expected version 2, got %v", result.Metadata["version"])
	}
}

func TestMemoryVectorStore_Search(t *testing.T) {
	store := NewMemoryVectorStore()

	// Add test vectors
	store.Store("vec-1", []float64{1.0, 0.0, 0.0}, map[string]interface{}{"name": "vec1"})
	store.Store("vec-2", []float64{0.9, 0.1, 0.0}, map[string]interface{}{"name": "vec2"})
	store.Store("vec-3", []float64{0.0, 1.0, 0.0}, map[string]interface{}{"name": "vec3"})

	// Search for similar to vec-1
	query := []float64{1.0, 0.0, 0.0}
	results, err := store.Search(query, 2, nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// First result should be vec-1 (perfect match)
	if results[0].ID != "vec-1" {
		t.Errorf("Expected first result to be vec-1, got %s", results[0].ID)
	}

	// Results should be sorted by score (descending)
	if results[0].Score < results[1].Score {
		t.Error("Results not sorted by score")
	}
}

func TestMemoryVectorStore_SearchWithFilters(t *testing.T) {
	store := NewMemoryVectorStore()

	store.Store("go-1", []float64{1.0, 0.0}, map[string]interface{}{"language": "go"})
	store.Store("py-1", []float64{0.9, 0.1}, map[string]interface{}{"language": "python"})
	store.Store("go-2", []float64{0.8, 0.2}, map[string]interface{}{"language": "go"})

	// Search with language filter
	filters := map[string]interface{}{"language": "go"}
	results, err := store.Search([]float64{1.0, 0.0}, 10, filters)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 Go results, got %d", len(results))
	}

	for _, result := range results {
		if result.Metadata["language"] != "go" {
			t.Errorf("Expected language 'go', got '%v'", result.Metadata["language"])
		}
	}
}

func TestMemoryVectorStore_Clear(t *testing.T) {
	store := NewMemoryVectorStore()

	store.Store("test-1", []float64{1.0}, nil)
	store.Store("test-2", []float64{2.0}, nil)

	err := store.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	count, _ := store.Count()
	if count != 0 {
		t.Errorf("Expected count 0 after clear, got %d", count)
	}
}

func TestMemoryVectorStore_BatchStore(t *testing.T) {
	store := NewMemoryVectorStore()

	entries := []struct {
		ID       string
		Vector   []float64
		Metadata map[string]interface{}
	}{
		{"batch-1", []float64{1.0}, map[string]interface{}{"index": 1}},
		{"batch-2", []float64{2.0}, map[string]interface{}{"index": 2}},
		{"batch-3", []float64{3.0}, map[string]interface{}{"index": 3}},
	}

	err := store.BatchStore(entries)
	if err != nil {
		t.Fatalf("BatchStore failed: %v", err)
	}

	count, _ := store.Count()
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestMemoryVectorStore_BatchDelete(t *testing.T) {
	store := NewMemoryVectorStore()

	store.Store("test-1", []float64{1.0}, nil)
	store.Store("test-2", []float64{2.0}, nil)
	store.Store("test-3", []float64{3.0}, nil)

	err := store.BatchDelete([]string{"test-1", "test-3"})
	if err != nil {
		t.Fatalf("BatchDelete failed: %v", err)
	}

	count, _ := store.Count()
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Verify test-2 still exists
	_, err = store.Get("test-2")
	if err != nil {
		t.Error("test-2 should still exist")
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1.0, 0.0, 0.0},
			b:        []float64{1.0, 0.0, 0.0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float64{1.0, 0.0},
			b:        []float64{0.0, 1.0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			a:        []float64{1.0, 0.0},
			b:        []float64{-1.0, 0.0},
			expected: -1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			if abs(result-tt.expected) > 0.0001 {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestCosineSimilarity_DifferentLengths(t *testing.T) {
	result := cosineSimilarity([]float64{1.0}, []float64{1.0, 2.0})
	if result != 0.0 {
		t.Errorf("Expected 0.0 for different length vectors, got %f", result)
	}
}

func TestNormalize(t *testing.T) {
	vector := []float64{3.0, 4.0}
	normalized := normalize(vector)

	// Length should be 1
	length := 0.0
	for _, v := range normalized {
		length += v * v
	}
	length = sqrt(length)

	if abs(length-1.0) > 0.0001 {
		t.Errorf("Expected normalized length 1.0, got %f", length)
	}
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func sqrt(x float64) float64 {
	// Simple Newton's method
	if x == 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}
