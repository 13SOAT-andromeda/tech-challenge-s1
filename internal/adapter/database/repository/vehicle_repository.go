package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type vehicleRepository struct {
	*BaseRepository[vehicle.Model]
}

func NewVehicleRepository(db *gorm.DB) ports.VehicleRepository {
	return &vehicleRepository{
		BaseRepository: NewBaseRepository[vehicle.Model](db),
	}
}

func (v *vehicleRepository) Search(ctx context.Context, params ports.VehicleSearch) []vehicle.Model {
	vehicles := []vehicle.Model{}
	q := v.db.Model(&vehicles)

	if !params.Status {
		q = v.db.Unscoped()
	}

	if params.Plate != "" {
		q.Where("lower(plate) LIKE ?", "%"+strings.ToLower(params.Plate)+"%")
	}
	if params.Brand != "" {
		q.Where("lower(brand) LIKE ?", "%"+strings.ToLower(params.Brand)+"%")
	}
	if params.Color != "" {
		q.Where("lower(color) LIKE ?", "%"+strings.ToLower(params.Color)+"%")
	}
	if params.Name != "" {
		q.Where("lower(name) LIKE ?", "%"+strings.ToLower(params.Name)+"%")
	}
	if params.Year != 0 {
		q.Where("year = ?", params.Year)
	}

	q.Find(&vehicles)

	return vehicles
}

func (v *vehicleRepository) GetByPlate(ctx context.Context, plate string) (*vehicle.Model, error) {
	vehicle := vehicle.Model{}
	err := v.db.Where("lower(plate) = ?", strings.ToLower(plate)).First(&vehicle).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &vehicle, nil
}
