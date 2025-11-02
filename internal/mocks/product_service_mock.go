package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

var _ ports.ProductService = (*MockProductService)(nil)

func (m *MockProductService) Update(ctx context.Context, p domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) GetById(ctx context.Context, productID uint) (*domain.Product, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) GetByIds(ctx context.Context, productIDs []uint) ([]domain.Product, error) {
	args := m.Called(ctx, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductService) UpdateStock(ctx context.Context, p domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) GetAll(ctx context.Context, productFilter *filter.ProductFilter) ([]domain.Product, error) {
	args := m.Called(ctx, productFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductService) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) Delete(ctx context.Context, productID uint) (*domain.Product, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) ManageStockItem(ctx context.Context, productID uint, quantity uint, operation string) (*domain.Product, error) {
	args := m.Called(ctx, productID, quantity, operation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) AddStockItem(ctx context.Context, productID uint, quantity uint) (*domain.Product, error) {
	args := m.Called(ctx, productID, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) RemoveStockItem(ctx context.Context, productID uint, quantity uint) (*domain.Product, error) {
	args := m.Called(ctx, productID, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}
