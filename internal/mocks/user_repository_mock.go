package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
)

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
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.Model), args.Error(1)
}
