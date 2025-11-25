package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/model"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

// MemoryProductRepository implements ProductRepository using in-memory storage
type MemoryProductRepository struct {
	mu       sync.RWMutex
	products map[int64]*model.Product
	nextID   int64
}

// NewMemoryProductRepository creates a new in-memory product repository
func NewMemoryProductRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[int64]*model.Product),
		nextID:   1,
	}
}

// Create creates a new product
func (r *MemoryProductRepository) Create(ctx context.Context, product *model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product.ID = r.nextID
	r.nextID++
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.Active = true

	r.products[product.ID] = product

	return nil
}

// GetByID retrieves a product by ID
func (r *MemoryProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, ErrProductNotFound
	}

	return product, nil
}

// List retrieves a list of products with pagination
func (r *MemoryProductRepository) List(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*model.Product, 0, len(r.products))
	for _, product := range r.products {
		if product.Active {
			products = append(products, product)
		}
	}

	// Apply pagination
	start := offset
	if start > len(products) {
		start = len(products)
	}

	end := start + limit
	if end > len(products) {
		end = len(products)
	}

	return products[start:end], nil
}

// ListByCategory retrieves products by category with pagination
func (r *MemoryProductRepository) ListByCategory(ctx context.Context, category string, limit, offset int) ([]*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*model.Product, 0)
	for _, product := range r.products {
		if product.Active && product.Category == category {
			products = append(products, product)
		}
	}

	// Apply pagination
	start := offset
	if start > len(products) {
		start = len(products)
	}

	end := start + limit
	if end > len(products) {
		end = len(products)
	}

	return products[start:end], nil
}

// Update updates a product
func (r *MemoryProductRepository) Update(ctx context.Context, product *model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return ErrProductNotFound
	}

	product.UpdatedAt = time.Now()
	r.products[product.ID] = product

	return nil
}

// Delete deletes a product (soft delete)
func (r *MemoryProductRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, exists := r.products[id]
	if !exists {
		return ErrProductNotFound
	}

	product.Active = false
	product.UpdatedAt = time.Now()

	return nil
}

// Count returns the total number of active products
func (r *MemoryProductRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, product := range r.products {
		if product.Active {
			count++
		}
	}

	return count, nil
}

