package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CustomerRepository interface {
	Repository[domain.Customer]
	FindByEmail(ctx context.Context, email string) (*domain.Customer, error)
}

type CustomerService interface {
	Create(ctx context.Context, c domain.Customer) (*domain.Customer, error)
	GetAll(ctx context.Context) ([]domain.Customer, error)
	GetByID(ctx context.Context, id uint) (*domain.Customer, error)
}
