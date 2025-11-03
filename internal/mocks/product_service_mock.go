package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
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

func (m *MockProductService) UpdateStock(ctx context.Context, products []domain.ProductItem, operation domain.StockOperation) error {
	args := m.Called(ctx, products, operation)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return nil
}

func (m *MockProductService) GetAll(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
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

func (m *MockProductService) CheckAvailability(ctx context.Context, productID uint, quantity uint) (bool, error) {
	args := m.Called(ctx, productID, quantity)
	return args.Bool(0), args.Error(1)
}
