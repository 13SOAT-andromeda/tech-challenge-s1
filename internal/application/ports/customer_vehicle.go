package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
)

type CustomerVehicleRepository interface {
	Repository[customer_vehicle.Model]
	FindByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) (*customer_vehicle.Model, error)
	FindByCustomerID(ctx context.Context, customerID uint) ([]customer_vehicle.Model, error)
	DeleteByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) error
}
