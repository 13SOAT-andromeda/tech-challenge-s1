package repository

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type maintenanceRepository struct {
	*BaseRepository[maintenance.Model]
}

func NewMaintenenceRepository(db *gorm.DB) ports.MaintenanceRepository {
	return &maintenanceRepository{
		BaseRepository: NewBaseRepository[maintenance.Model](db),
	}
}
