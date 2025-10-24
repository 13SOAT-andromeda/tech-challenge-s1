package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CustomerRepository interface {
	Repository[customer.Model]
	FindByEmail(ctx context.Context, email string) (*customer.Model, error)
	FindByDocument(ctx context.Context, document string) (*customer.Model, error)
}

type CustomerService interface {
	Create(ctx context.Context, c domain.Customer) (*domain.Customer, error)
	GetAll(ctx context.Context) ([]domain.Customer, error)
	GetByID(ctx context.Context, id uint) (*domain.Customer, error)
	UpdateByID(ctx context.Context, id uint, c domain.Customer) error
	DeleteByID(ctx context.Context, id uint) (*domain.Customer, error)
}

type CustomerUseCase interface {
}
