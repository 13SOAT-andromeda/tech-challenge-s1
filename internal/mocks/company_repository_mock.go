package mocks

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

type MockCompanyRepository struct {
	mock.Mock
}

var _ ports.CompanyRepository = (*MockCompanyRepository)(nil)

func (m *MockCompanyRepository) FindByID(ctx context.Context, id uint) (*company.Model, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Model), args.Error(1)
}

func (m *MockCompanyRepository) FindAll(ctx context.Context, includeDeleted bool) ([]company.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]company.Model), args.Error(1)
}

func (m *MockCompanyRepository) Create(ctx context.Context, entity *company.Model) (*company.Model, error) {
	args := m.Called(ctx, entity)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*company.Model), args.Error(1)
}

func (m *MockCompanyRepository) Update(ctx context.Context, entity *company.Model) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockCompanyRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
