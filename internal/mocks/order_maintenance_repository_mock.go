package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

type MockOrderMaintenanceRepository struct {
	mock.Mock
}

var _ ports.OrderMaintenanceRepository = (*MockOrderMaintenanceRepository)(nil)

func (m *MockOrderMaintenanceRepository) FindAll(ctx context.Context, includeDeleted bool) ([]order_maintenance.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]order_maintenance.Model), args.Error(1)
}

func (m *MockOrderMaintenanceRepository) FindByID(ctx context.Context, id uint) (*order_maintenance.Model, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order_maintenance.Model), args.Error(1)
}

func (m *MockOrderMaintenanceRepository) Create(ctx context.Context, entity *order_maintenance.Model) (*order_maintenance.Model, error) {
	args := m.Called(ctx, entity)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order_maintenance.Model), args.Error(1)
}

func (m *MockOrderMaintenanceRepository) Update(ctx context.Context, entity *order_maintenance.Model) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockOrderMaintenanceRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
