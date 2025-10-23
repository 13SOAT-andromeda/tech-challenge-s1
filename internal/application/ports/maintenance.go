package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type MaintenanceRepository interface {
	Repository[maintenance.Model]
}

type MaintenanceService interface {
	Create(ctx context.Context, c domain.Maintenance) (*domain.Maintenance, error)
	GetByID(ctx context.Context, id uint) (*domain.Maintenance, error)
	UpdateByID(ctx context.Context, id uint, c domain.Maintenance) error
	DeleteByID(ctx context.Context, id uint) (*domain.Maintenance, error)
}
