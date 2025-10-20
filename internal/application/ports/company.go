package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CompanyRepository interface {
	Repository[company.Model]
}

type CompanyService interface {
	Create(ctx context.Context, c domain.Company) (*domain.Company, error)
	GetByID(ctx context.Context, id uint) (*domain.Company, error)
	UpdateByID(ctx context.Context, id uint, c domain.Company) error
	DeleteByID(ctx context.Context, id uint) (*domain.Company, error)
}
