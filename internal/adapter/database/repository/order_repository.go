package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type orderRepository struct {
	*BaseRepository[order.Model]
}

func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &orderRepository{
		BaseRepository: NewBaseRepository[order.Model](db),
	}
}

func (r *orderRepository) Search(ctx context.Context, params ports.OrderSearch) ([]order.Model, error) {
	var model []order.Model

	db := r.db.WithContext(ctx)

	if !params.Enabled {
		db = db.Unscoped()
	}

	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}

	err := db.Find(&model).Error

	return model, err

}
