package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type OrderRepository struct {
	*BaseRepository[order.Model]
}

func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &OrderRepository{
		BaseRepository: NewBaseRepository[order.Model](db),
	}
}

func (r *OrderRepository) FindOrderByID(ctx context.Context, id uint) (*order.Model, error) {
	var models order.Model
	err := r.db.WithContext(ctx).
		Preload("CustomerVehicle").
		Preload("Company").
		Preload("User").
		Where("id = ?", id).
		Find(&models).Error

	return &models, err
}

func (r *OrderRepository) Search(ctx context.Context, params ports.OrderSearch) ([]order.Model, error) {
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
