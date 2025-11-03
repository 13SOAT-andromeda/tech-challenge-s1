package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"github.com/stretchr/testify/mock"
)

type MockCustomerService struct {
	mock.Mock
}

var _ ports.CustomerService = (*MockCustomerService)(nil)

func (m *MockCustomerService) Create(ctx context.Context, p domain.Customer) (*domain.Customer, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerService) Search(ctx context.Context, filter *filter.CustomerFilter) ([]domain.Customer, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Customer), args.Error(1)
}

func (m *MockCustomerService) GetByID(ctx context.Context, id uint) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerService) UpdateByID(ctx context.Context, id uint, c domain.Customer) error {
	args := m.Called(ctx, id, c)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *MockCustomerService) DeleteByID(ctx context.Context, id uint) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}
