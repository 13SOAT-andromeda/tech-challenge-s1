package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
)

type MockVehicleRepository struct {
	MockGenericRepository[vehicle.Model]
}

var _ ports.VehicleRepository = (*MockVehicleRepository)(nil)

func (m *MockVehicleRepository) Search(ctx context.Context, params ports.VehicleSearch) []vehicle.Model {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]vehicle.Model)
}

func (m *MockVehicleRepository) GetByPlate(ctx context.Context, plate string) (*vehicle.Model, error) {
	args := m.Called(ctx, plate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vehicle.Model), args.Error(1)
}
