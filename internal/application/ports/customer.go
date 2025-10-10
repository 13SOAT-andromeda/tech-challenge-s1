package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CustomerRepository interface {
	Save(ctx context.Context, c domain.Customer) (*domain.Customer, error)
	FindAll(ctx context.Context) ([]domain.Customer, error)
	FindByID(ctx context.Context, id uint) (*domain.Customer, error)
}

type CustomerService interface {
	Create(ctx context.Context, c domain.Customer) (*domain.Customer, error)
	GetAll(ctx context.Context) ([]domain.Customer, error)
	GetByID(ctx context.Context, id uint) (*domain.Customer, error)
}
