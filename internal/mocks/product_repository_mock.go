package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

var _ ports.ProductRepository = (*MockProductRepository)(nil)

func (m *MockProductRepository) FindByName(ctx context.Context, name string) (*product.Model, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Model), args.Error(1)
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*product.Model, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Model), args.Error(1)
}

func (m *MockProductRepository) FindAll(ctx context.Context, includeDeleted bool) ([]product.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]product.Model), args.Error(1)
}

func (m *MockProductRepository) Search(ctx context.Context, filter filter.ProductFilter) ([]product.Model, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]product.Model), args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, entity *product.Model) (*product.Model, error) {
	args := m.Called(ctx, entity)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*product.Model), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, entity *product.Model) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id uint, quantity int, operation domain.StockOperation) error {
	args := m.Called(ctx, id, quantity, operation)
	return args.Error(0)
}

func (m *MockProductRepository) FindByIDs(ctx context.Context, productIDs []uint) ([]product.Model, error) {
	args := m.Called(ctx, productIDs)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]product.Model), args.Error(1)
}
