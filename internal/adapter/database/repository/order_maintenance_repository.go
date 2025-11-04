package repository

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type orderMaintenanceRepository struct {
	*BaseRepository[order_maintenance.Model]
}

func NewOrderMaintenanceRepository(db *gorm.DB) ports.OrderMaintenanceRepository {
	return &orderMaintenanceRepository{
		BaseRepository: NewBaseRepository[order_maintenance.Model](db),
	}
}
