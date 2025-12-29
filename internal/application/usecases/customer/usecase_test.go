package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddVehicleToCustomer_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(1)

	expectedCustomer := &customer.Model{
		Name:  "Gedan Magalhaes",
		Email: "gedan@example.com",
	}
	expectedCustomer.ID = customerID

	expectedVehicle := &domain.Vehicle{
		ID:    vehicleID,
		Brand: "Toyota",
		Name:  "Corolla 2020",
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(expectedVehicle, nil)
	mockCustomerVehicleRepo.On("FindByCustomerAndVehicle", ctx, customerID, vehicleID).Return(nil, nil)
	mockCustomerVehicleRepo.On("Create", ctx, mock.AnythingOfType("*customer_vehicle.Model")).Return(&customer_vehicle.Model{}, nil)

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

// TestAddVehicleToCustomer_CustomerNotFound testa quando o cliente não é encontrado
func TestAddVehicleToCustomer_CustomerNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(999)
	vehicleID := uint(1)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("customer not found"))

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer not found")
	mockRepo.AssertExpectations(t)
}

// TestAddVehicleToCustomer_CustomerNil testa quando o cliente retorna nil
func TestAddVehicleToCustomer_CustomerNil(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(999)
	vehicleID := uint(1)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, nil)

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrCustomerNotFound.Message)
	mockRepo.AssertExpectations(t)
}

// TestAddVehicleToCustomer_VehicleNotFound testa quando o veículo não é encontrado
func TestAddVehicleToCustomer_VehicleNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(999)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(nil, errors.New("vehicle not found"))

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vehicle not found")
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
}

// TestAddVehicleToCustomer_VehicleNil testa quando o veículo retorna nil
func TestAddVehicleToCustomer_VehicleNil(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(999)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(nil, nil)

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrVehicleNotFound.Message)
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
}

// TestAddVehicleToCustomer_AlreadyAssociated testa quando a associação já existe
// TestAddVehicleToCustomer_AlreadyAssociated testa quando a associação já existe
func TestAddVehicleToCustomer_AlreadyAssociated(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(1)

	expectedCustomer := &customer.Model{
		Name:  "Gedan Magalhaes",
		Email: "gedan@example.com",
	}
	expectedCustomer.ID = customerID

	expectedVehicle := &domain.Vehicle{
		ID:    vehicleID,
		Brand: "Toyota",
		Name:  "Corolla 2020",
	}

	existingAssociation := &customer_vehicle.Model{
		CustomerID: customerID,
		VehicleID:  vehicleID,
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(expectedVehicle, nil)
	mockCustomerVehicleRepo.On("FindByCustomerAndVehicle", ctx, customerID, vehicleID).Return(existingAssociation, nil)

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vehicle is already associated with this customer")
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

// TestAddVehicleToCustomer_ErrorCheckingExisting testa erro ao verificar associação existente
func TestAddVehicleToCustomer_ErrorCheckingExisting(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(1)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	expectedVehicle := &domain.Vehicle{
		ID:    vehicleID,
		Brand: "Toyota",
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(expectedVehicle, nil)
	mockCustomerVehicleRepo.On("FindByCustomerAndVehicle", ctx, customerID, vehicleID).Return(nil, errors.New("db error"))

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking existing association")
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

// TestAddVehicleToCustomer_CreateError testa erro ao criar a associação
func TestAddVehicleToCustomer_CreateError(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(1)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	expectedVehicle := &domain.Vehicle{
		ID:    vehicleID,
		Brand: "Toyota",
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(expectedVehicle, nil)
	mockCustomerVehicleRepo.On("FindByCustomerAndVehicle", ctx, customerID, vehicleID).Return(nil, nil)
	mockCustomerVehicleRepo.On("Create", ctx, mock.AnythingOfType("*customer_vehicle.Model")).Return(nil, errors.New("create error"))

	err := useCase.AddVehicleToCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating customer-vehicle association")
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

// TestRemoveVehicleFromCustomer_Success testa a remoção bem-sucedida de um veículo de um cliente
func TestRemoveVehicleFromCustomer_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(1)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	expectedVehicle := &domain.Vehicle{
		ID:    vehicleID,
		Brand: "Toyota",
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(expectedVehicle, nil)
	mockCustomerVehicleRepo.On("DeleteByCustomerAndVehicle", ctx, customerID, vehicleID).Return(nil)

	err := useCase.RemoveVehicleFromCustomer(ctx, customerID, vehicleID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

// TestRemoveVehicleFromCustomer_CustomerNotFound testa quando o cliente não é encontrado
func TestRemoveVehicleFromCustomer_CustomerNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(999)
	vehicleID := uint(1)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("customer not found"))

	err := useCase.RemoveVehicleFromCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer not found")
	mockRepo.AssertExpectations(t)
}

// TestRemoveVehicleFromCustomer_CustomerNil testa quando o cliente retorna nil
func TestRemoveVehicleFromCustomer_CustomerNil(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(999)
	vehicleID := uint(1)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, nil)

	err := useCase.RemoveVehicleFromCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrCustomerNotFound.Message)
	mockRepo.AssertExpectations(t)
}

// TestRemoveVehicleFromCustomer_VehicleNotFound testa quando o veículo não é encontrado
func TestRemoveVehicleFromCustomer_VehicleNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(999)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(nil, errors.New("vehicle not found"))

	err := useCase.RemoveVehicleFromCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vehicle not found")
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
}

func TestRemoveVehicleFromCustomer_VehicleNil(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(999)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(nil, nil)

	err := useCase.RemoveVehicleFromCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrVehicleNotFound.Message)
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
}

func TestRemoveVehicleFromCustomer_DeleteError(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)
	vehicleID := uint(1)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	expectedVehicle := &domain.Vehicle{
		ID:    vehicleID,
		Brand: "Toyota",
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockVehicleService.On("GetByID", ctx, vehicleID).Return(expectedVehicle, nil)
	mockCustomerVehicleRepo.On("DeleteByCustomerAndVehicle", ctx, customerID, vehicleID).Return(errors.New("delete error"))

	err := useCase.RemoveVehicleFromCustomer(ctx, customerID, vehicleID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error removing customer-vehicle association")
	mockRepo.AssertExpectations(t)
	mockVehicleService.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

func TestGetCustomerVehicles_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)

	expectedCustomer := &customer.Model{
		Name:  "Gedan Magalhaes",
		Email: "gedan@example.com",
	}
	expectedCustomer.ID = customerID

	vehicleModel1 := vehicle.Model{
		Plate: "ABC1234",
		Name:  "Corolla 2020",
		Year:  2020,
		Brand: "Toyota",
		Color: "Prata",
	}
	vehicleModel1.ID = 1

	vehicleModel2 := vehicle.Model{
		Plate: "XYZ5678",
		Name:  "Civic 2021",
		Year:  2021,
		Brand: "Honda",
		Color: "Preto",
	}
	vehicleModel2.ID = 2

	customerVehicles := []customer_vehicle.Model{
		{
			CustomerID: customerID,
			VehicleID:  1,
			Vehicle:    vehicleModel1,
		},
		{
			CustomerID: customerID,
			VehicleID:  2,
			Vehicle:    vehicleModel2,
		},
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockCustomerVehicleRepo.On("FindByCustomerID", ctx, customerID).Return(customerVehicles, nil)

	result, err := useCase.GetCustomerVehicles(ctx, customerID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, uint(1), result[0].VehicleId)
	assert.Equal(t, customerID, result[0].CustomerId)
	assert.Equal(t, "Toyota", result[0].Vehicle.Brand)
	assert.Equal(t, "Corolla 2020", result[0].Vehicle.Name)
	assert.Equal(t, uint(2), result[1].VehicleId)
	assert.Equal(t, customerID, result[1].CustomerId)
	assert.Equal(t, "Honda", result[1].Vehicle.Brand)
	assert.Equal(t, "Civic 2021", result[1].Vehicle.Name)
	mockRepo.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

// TestGetCustomerVehicles_CustomerNotFound testa quando o cliente não é encontrado
func TestGetCustomerVehicles_CustomerNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(999)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("customer not found"))

	result, err := useCase.GetCustomerVehicles(ctx, customerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "customer not found")
	mockRepo.AssertExpectations(t)
}

// TestGetCustomerVehicles_CustomerNil testa quando o cliente retorna nil
func TestGetCustomerVehicles_CustomerNil(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(999)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, nil)

	result, err := useCase.GetCustomerVehicles(ctx, customerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), ErrCustomerNotFound.Message)
	mockRepo.AssertExpectations(t)
}

// TestGetCustomerVehicles_FetchError testa erro ao buscar veículos
func TestGetCustomerVehicles_FetchError(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockCustomerVehicleRepo.On("FindByCustomerID", ctx, customerID).Return(nil, errors.New("db error"))

	result, err := useCase.GetCustomerVehicles(ctx, customerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error fetching customer vehicles")
	mockRepo.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

// TestGetCustomerVehicles_EmptyList testa quando não há veículos associados
func TestGetCustomerVehicles_EmptyList(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockCustomerVehicleRepo.On("FindByCustomerID", ctx, customerID).Return([]customer_vehicle.Model{}, nil)

	result, err := useCase.GetCustomerVehicles(ctx, customerID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}

func TestGetCustomerVehicles_SkipsInvalidVehicles(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockCustomerVehicleRepo := new(mocks.MockCustomerVehicleRepository)
	mockVehicleService := new(mocks.MockVehicleService)
	useCase := NewCustomerUseCase(mockRepo, mockCustomerVehicleRepo, mockVehicleService)

	ctx := context.Background()
	customerID := uint(1)

	expectedCustomer := &customer.Model{
		Name: "Gedan Magalhaes",
	}
	expectedCustomer.ID = customerID

	validVehicle := vehicle.Model{
		Plate: "ABC1234",
		Name:  "Corolla 2020",
		Year:  2020,
		Brand: "Toyota",
		Color: "Prata",
	}
	validVehicle.ID = 1

	invalidVehicle := vehicle.Model{
		Plate: "INVALID",
		Name:  "Invalid Name",
		Year:  0,
		Brand: "Invalid",
		Color: "Invalid",
	}
	invalidVehicle.ID = 0

	customerVehicles := []customer_vehicle.Model{
		{
			CustomerID: customerID,
			VehicleID:  1,
			Vehicle:    validVehicle,
		},
		{
			CustomerID: customerID,
			VehicleID:  0,
			Vehicle:    invalidVehicle,
		},
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)
	mockCustomerVehicleRepo.On("FindByCustomerID", ctx, customerID).Return(customerVehicles, nil)

	result, err := useCase.GetCustomerVehicles(ctx, customerID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, uint(1), result[0].VehicleId)
	assert.Equal(t, customerID, result[0].CustomerId)
	assert.Equal(t, "Toyota", result[0].Vehicle.Brand)
	assert.Equal(t, "Corolla 2020", result[0].Vehicle.Name)
	assert.Equal(t, uint(0), result[1].VehicleId)
	assert.Equal(t, customerID, result[1].CustomerId)
	mockRepo.AssertExpectations(t)
	mockCustomerVehicleRepo.AssertExpectations(t)
}
