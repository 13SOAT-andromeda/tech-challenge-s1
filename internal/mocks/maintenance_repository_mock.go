package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

type MockMaintenanceRepository struct {
	mock.Mock
}

var _ ports.MaintenanceRepository = (*MockMaintenanceRepository)(nil)

func (m *MockMaintenanceRepository) FindAll(ctx context.Context, includeDeleted bool) ([]maintenance.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]maintenance.Model), args.Error(1)
}

func (m *MockMaintenanceRepository) FindByID(ctx context.Context, id uint) (*maintenance.Model, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*maintenance.Model), args.Error(1)
}

func (m *MockMaintenanceRepository) Create(ctx context.Context, entity *maintenance.Model) (*maintenance.Model, error) {
	args := m.Called(ctx, entity)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*maintenance.Model), args.Error(1)
}

func (m *MockMaintenanceRepository) Update(ctx context.Context, entity *maintenance.Model) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockMaintenanceRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
