package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type maintenanceRepository struct {
	*BaseRepository[maintenance.Model]
}

func NewMaintenanceRepository(db *gorm.DB) ports.MaintenanceRepository {
	return &maintenanceRepository{
		BaseRepository: NewBaseRepository[maintenance.Model](db),
	}
}

func (r *maintenanceRepository) FindByIDs(ctx context.Context, maintenanceIDs []uint) ([]maintenance.Model, error) {
	var maintenances []maintenance.Model

	err := r.db.WithContext(ctx).Where("id IN ?", maintenanceIDs).Find(&maintenances).Error
	if err != nil {
		return nil, err
	}

	return maintenances, nil
}
