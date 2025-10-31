package service

import (
	"context"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/model"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/repository"
)

// ProductService handles product business logic
type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// Create creates a new product
func (s *ProductService) Create(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetByID retrieves a product by ID
func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves a list of products
func (s *ProductService) List(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	return s.repo.List(ctx, limit, offset)
}

// ListByCategory retrieves products by category
func (s *ProductService) ListByCategory(ctx context.Context, category string, limit, offset int) ([]*model.Product, error) {
	return s.repo.ListByCategory(ctx, category, limit, offset)
}

// Update updates a product
func (s *ProductService) Update(ctx context.Context, id int64, req *model.UpdateProductRequest) (*model.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.Active != nil {
		product.Active = *req.Active
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete deletes a product
func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// Count returns the total number of products
func (s *ProductService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

