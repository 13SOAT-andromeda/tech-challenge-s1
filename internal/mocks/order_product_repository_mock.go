package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

type MockOrderProductRepository struct {
	mock.Mock
}

var _ ports.OrderProductRepository = (*MockOrderProductRepository)(nil)

func (m *MockOrderProductRepository) FindAll(ctx context.Context, includeDeleted bool) ([]order_product.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]order_product.Model), args.Error(1)
}

func (m *MockOrderProductRepository) FindByID(ctx context.Context, id uint) (*order_product.Model, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order_product.Model), args.Error(1)
}

func (m *MockOrderProductRepository) Create(ctx context.Context, entity *order_product.Model) (*order_product.Model, error) {
	args := m.Called(ctx, entity)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order_product.Model), args.Error(1)
}

func (m *MockOrderProductRepository) Update(ctx context.Context, entity *order_product.Model) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockOrderProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
