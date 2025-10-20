package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

// MockRepository genérico
type MockRepository[T any] struct {
	mock.Mock
}

func (m *MockRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockRepository[T]) Update(ctx context.Context, entity *T) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockRepository[T]) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockCustomerRepository struct {
	MockRepository[model.CustomerModel]
}

var _ ports.CustomerRepository = (*MockCustomerRepository)(nil)

func (m *MockCustomerRepository) FindByEmail(ctx context.Context, email string) (*model.CustomerModel, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CustomerModel), args.Error(1)
}

type MockUserRepository struct {
	MockRepository[model.UserModel]
}

var _ ports.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Search(ctx context.Context, params ports.UserSearch) []model.UserModel {
	args := m.Called(ctx, params)
	return args.Get(0).([]model.UserModel)
}

func (m *MockUserRepository) Exists(ctx context.Context, id uint, email string) (bool, error) {
	args := m.Called(ctx, id, email)
	return args.Bool(0), args.Error(1)
}
