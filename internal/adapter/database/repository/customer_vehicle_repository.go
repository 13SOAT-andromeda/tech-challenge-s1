package repository

import (
	"context"
	"errors"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type customerVehicleRepository struct {
	*BaseRepository[customer_vehicle.Model]
}

func NewCustomerVehicleRepository(db *gorm.DB) ports.CustomerVehicleRepository {
	return &customerVehicleRepository{
		BaseRepository: NewBaseRepository[customer_vehicle.Model](db),
	}
}

func (r *customerVehicleRepository) FindByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) (*customer_vehicle.Model, error) {
	var model customer_vehicle.Model

	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("customer_id = ? AND vehicle_id = ?", customerID, vehicleID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (r *customerVehicleRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]customer_vehicle.Model, error) {
	var models []customer_vehicle.Model
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Customer").
		Where("customer_id = ?", customerID).
		Find(&models).Error

	return models, err
}

func (r *customerVehicleRepository) DeleteByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) error {
	result := r.db.WithContext(ctx).
		Where("customer_id = ? AND vehicle_id = ?", customerID, vehicleID).
		Delete(&customer_vehicle.Model{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("customer vehicle association not found")
	}

	return nil
}
