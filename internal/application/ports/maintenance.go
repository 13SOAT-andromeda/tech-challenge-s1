package ports

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
)

type MaintenanceRepository interface {
	Repository[maintenance.Model]
}

type MaintenanceService interface {
	Create(maintenance maintenance.Model) (*maintenance.Model, error)
	GetByID(id uint) (*maintenance.Model, error)
	UpdateByID(id uint, service maintenance.Model) error
	DeleteByID(id uint) (*maintenance.Model, error)
}
