package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"github.com/stretchr/testify/mock"
)

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
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.Model), args.Error(1)
}

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

func (m *MockCustomerRepository) FindAll(ctx context.Context, includeDeleted bool) ([]customer.Model, error) {
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

func (m *MockCustomerRepository) FindByDocument(ctx context.Context, document string) (*customer.Model, error) {
	args := m.Called(ctx, document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Model), args.Error(1)
}

func (m *MockCustomerRepository) Search(ctx context.Context, filters filter.CustomerFilter) ([]customer.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]customer.Model), args.Error(1)
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

type MockMaintenanceRepository struct {
	mock.Mock
}

var _ ports.MaintenanceRepository = (*MockMaintenanceRepository)(nil)

func (m *MockMaintenanceRepository) FindAll(ctx context.Context, includeDeleted bool) ([]maintenance.Model, error) {
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

func (m *MockProductRepository) FindAll(ctx context.Context, includeDeleted bool) ([]product.Model, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]product.Model), args.Error(1)
}

func (m *MockProductRepository) Search(ctx context.Context, filters filter.ProductFilter) ([]product.Model, error) {
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

func (m *MockProductRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockProductRepository) FindByIDs(ctx context.Context, productIDs []uint) ([]product.Model, error) {
	args := m.Called(ctx, productIDs)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]product.Model), args.Error(1)
}

type MockVehicleRepository struct {
	MockGenericRepository[vehicle.Model]
}

var _ ports.VehicleRepository = (*MockVehicleRepository)(nil)

func (m *MockVehicleRepository) Search(ctx context.Context, params ports.VehicleSearch) []vehicle.Model {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]vehicle.Model)
}

func (m *MockVehicleRepository) GetByPlate(ctx context.Context, plate string) (*vehicle.Model, error) {
	args := m.Called(ctx, plate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vehicle.Model), args.Error(1)
}

type MockCustomerVehicleRepository struct {
	MockGenericRepository[customer_vehicle.Model]
}

var _ ports.CustomerVehicleRepository = (*MockCustomerVehicleRepository)(nil)

func (m *MockCustomerVehicleRepository) FindByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) (*customer_vehicle.Model, error) {
	args := m.Called(ctx, customerID, vehicleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer_vehicle.Model), args.Error(1)
}

func (m *MockCustomerVehicleRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]customer_vehicle.Model, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]customer_vehicle.Model), args.Error(1)
}

func (m *MockCustomerVehicleRepository) DeleteByCustomerAndVehicle(ctx context.Context, customerID, vehicleID uint) error {
	args := m.Called(ctx, customerID, vehicleID)
	return args.Error(0)
}

type MockVehicleService struct {
	mock.Mock
}

var _ ports.VehicleService = (*MockVehicleService)(nil)

func (m *MockVehicleService) Create(ctx context.Context, v domain.Vehicle) (*domain.Vehicle, error) {
	args := m.Called(ctx, v)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) GetAll(ctx context.Context, params map[string]interface{}) (*[]domain.Vehicle, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) GetByID(ctx context.Context, id uint) (*domain.Vehicle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) GetByPlate(ctx context.Context, plate string) (*domain.Vehicle, error) {
	args := m.Called(ctx, plate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) Update(ctx context.Context, v domain.Vehicle) (*domain.Vehicle, error) {
	args := m.Called(ctx, v)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
