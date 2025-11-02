package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockMaintenanceService struct {
	mock.Mock
}

var _ ports.MaintenanceService = (*MockMaintenanceService)(nil)

func (m *MockMaintenanceService) Create(ctx context.Context, c domain.Maintenance) (*domain.Maintenance, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Maintenance), args.Error(1)
}

func (m *MockMaintenanceService) GetByIDs(ctx context.Context, ids []uint) ([]domain.Maintenance, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Maintenance), args.Error(1)
}

func (m *MockMaintenanceService) GetByID(ctx context.Context, id uint) (*domain.Maintenance, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Maintenance), args.Error(1)
}

func (m *MockMaintenanceService) GetAll(ctx context.Context) ([]domain.Maintenance, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Maintenance), args.Error(1)
}

func (m *MockMaintenanceService) Update(ctx context.Context, c domain.Maintenance) (*domain.Maintenance, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Maintenance), args.Error(1)
}

func (m *MockMaintenanceService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMaintenanceService) UpdateByID(ctx context.Context, id uint, c domain.Maintenance) error {
	args := m.Called(ctx, id, c)
	return args.Error(0)
}

func (m *MockMaintenanceService) DeleteByID(ctx context.Context, id uint) (*domain.Maintenance, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Maintenance), args.Error(1)
}
