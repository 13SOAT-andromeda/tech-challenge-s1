package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CompanyRepository interface {
	Save(ctx context.Context, c domain.Company) (*domain.Company, error)
	FindByID(ctx context.Context, id uint) (*domain.Company, error)
	UpdateByID(ctx context.Context, id uint, c domain.Company) (*domain.Company, error)
	DeleteByID(ctx context.Context, id uint) (*domain.Company, error)
}

type CompanyService interface {
	Create(ctx context.Context, c domain.Company) (*domain.Company, error)
	GetByID(ctx context.Context, id uint) (*domain.Company, error)
	UpdateByID(ctx context.Context, id uint, c domain.Company) (*domain.Company, error)
	DeleteByID(ctx context.Context, id uint) (*domain.Company, error)
}
