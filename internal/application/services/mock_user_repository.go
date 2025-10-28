package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

// MockGenericRepository genérico
type MockGenericRepository[T any] struct {
	mock.Mock
}

func (m *MockGenericRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockGenericRepository[T]) FindAll(ctx context.Context, includeDeleted bool) ([]T, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockGenericRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockGenericRepository[T]) Update(ctx context.Context, entity *T) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockGenericRepository[T]) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserRepository struct {
	MockGenericRepository[user.Model]
}

var _ ports.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Search(ctx context.Context, params ports.UserSearch) []user.Model {
	args := m.Called(ctx, params)
	return args.Get(0).([]user.Model)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*user.Model, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(0)
	}
	return args.Get(0).(*user.Model), args.Error(1)
}
