package repository

import (
	"context"
	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/internal/model"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	List(ctx context.Context, limit, offset int) ([]*model.Product, error)
	ListByCategory(ctx context.Context, category string, limit, offset int) ([]*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}

