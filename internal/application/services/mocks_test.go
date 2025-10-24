package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/mock"
)

type MockCustomerRepository struct {
	mock.Mock
}

var _ ports.CustomerRepository = (*MockCustomerRepository)(nil)

func (m *MockCustomerRepository) FindByID(ctx context.Context, id uint) (*customer.Model, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Model), args.Error(1)
}

func (m *MockCustomerRepository) FindAll(ctx context.Context) ([]customer.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]customer.Model), args.Error(1)
}

func (m *MockCustomerRepository) Create(ctx context.Context, entity *customer.Model) (*customer.Model, error) {
	args := m.Called(ctx, entity)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*customer.Model), args.Error(1)
}

func (m *MockCustomerRepository) Update(ctx context.Context, entity *customer.Model) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerRepository) FindByEmail(ctx context.Context, email string) (*customer.Model, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Model), args.Error(1)
}

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

func (m *MockCompanyRepository) FindAll(ctx context.Context) ([]company.Model, error) {
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

type MockMaintenanceRepository struct {
	mock.Mock
}

var _ ports.MaintenanceRepository = (*MockMaintenanceRepository)(nil)

func (m *MockMaintenanceRepository) FindAll(ctx context.Context) ([]maintenance.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]maintenance.Model), args.Error(1)
}

func (m *MockMaintenanceRepository) FindByID(ctx context.Context, id uint) (*maintenance.Model, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*maintenance.Model), args.Error(1)
}

func (m *MockMaintenanceRepository) Create(ctx context.Context, entity *maintenance.Model) (*maintenance.Model, error) {
	args := m.Called(ctx, entity)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*maintenance.Model), args.Error(1)
}

func (m *MockMaintenanceRepository) Update(ctx context.Context, entity *maintenance.Model) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockMaintenanceRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

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

func (m *MockProductRepository) FindAll(ctx context.Context) ([]product.Model, error) {
	args := m.Called(ctx)
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
