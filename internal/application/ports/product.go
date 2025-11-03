package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type ProductRepository interface {
	Repository[product.Model]
	UpdateStock(ctx context.Context, productID uint, quantity uint, operation domain.StockOperation) error
	FindByIDs(ctx context.Context, productIDs []uint) ([]product.Model, error)
}

type ProductService interface {
	Create(ctx context.Context, p domain.Product) (*domain.Product, error)
	Update(ctx context.Context, p domain.Product) (*domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetById(ctx context.Context, productID uint) (*domain.Product, error)
	GetByIds(ctx context.Context, productIDs []uint) ([]domain.Product, error)
	Delete(ctx context.Context, productID uint) (*domain.Product, error)
	CheckAvailability(ctx context.Context, productID uint, quantity uint) (bool, error)
	UpdateStock(ctx context.Context, products []domain.StockItem) error
}
