package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type VehicleSearch struct {
	Plate  string
	Name   string
	Year   int
	Brand  string
	Color  string
	Status bool
}

type VehicleRepository interface {
	Repository[vehicle.Model]
	Search(ctx context.Context, params VehicleSearch) []vehicle.Model
	GetByPlate(ctx context.Context, plate string) (*vehicle.Model, error)
}

type VehicleService interface {
	Create(ctx context.Context, u domain.Vehicle) (*domain.Vehicle, error)
	GetAll(ctx context.Context, params map[string]interface{}) (*[]domain.Vehicle, error)
	GetByID(ctx context.Context, id uint) (*domain.Vehicle, error)
	GetByPlate(ctx context.Context, plate string) (*domain.Vehicle, error)
	Update(ctx context.Context, u domain.Vehicle) (*domain.Vehicle, error)
	Delete(ctx context.Context, id uint) error
}
