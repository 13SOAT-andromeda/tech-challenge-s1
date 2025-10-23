package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type ProductRepository interface {
	Repository[product.Model]
	FindByName(ctx context.Context, email string) (*product.Model, error)
}

type ProductService interface {
	Create(ctx context.Context, p domain.Product) (*domain.Product, error)
	// Update(ctx context.Context, p domain.Product) (*domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetById(ctx context.Context, id uint) (*domain.Product, error)
	// GetByName(ctx context.Context, p domain.ProductName) (*domain.Product, error)
	// UpdateStock(ctx context.Context, id uint, quantity int) error
	Delete(ctx context.Context, id uint) error
}
