package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
)

type MockCustomerVehicleRepository struct {
	MockGenericRepository[customer_vehicle.Model]
}

var _ ports.CustomerVehicleRepository = (*MockCustomerVehicleRepository)(nil)

func (m *MockCustomerVehicleRepository) FindByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) (*customer_vehicle.Model, error) {
	args := m.Called(ctx, customerID, vehicleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer_vehicle.Model), args.Error(1)
}

func (m *MockCustomerVehicleRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]customer_vehicle.Model, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]customer_vehicle.Model), args.Error(1)
}

func (m *MockCustomerVehicleRepository) DeleteByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) error {
	args := m.Called(ctx, customerID, vehicleID)
	return args.Error(0)
}
