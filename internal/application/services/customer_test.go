package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCustomerService_Create_Success(t *testing.T) {
	// Arrange (Preparar)
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	inputCustomer := domain.Customer{
		Name:     "João Silva",
		Email:    "joao@example.com",
		Document: "12345678900",
		Type:     "individual",
		Contact:  "11999999999",
		Address: &domain.Address{
			Address:       "Rua Teste",
			AddressNumber: "123",
			City:          "São Paulo",
			Neighborhood:  "Centro",
			Country:       "Brasil",
			ZipCode:       "01234-567",
		},
	}

	mockModel := model.FromDomain(inputCustomer)

	// Configurar o mock para esperar a chamada Create e retornar sucesso
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.CustomerModel")).
		Return(&mockModel, nil)
	// Act (Agir)
	result, err := service.Create(ctx, inputCustomer)

	// Assert (Verificar)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "João Silva", result.Name)
	assert.Equal(t, "joao@example.com", result.Email)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_Create_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()

	inputCustomer := domain.Customer{
		Name:  "João Silva",
		Email: "joao@example.com",
	}

	expectedError := errors.New("database connection error")

	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.CustomerModel")).Return(nil, expectedError)

	// Act
	result, err := service.Create(ctx, inputCustomer)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	customerID := uint(1)
	expectedCustomer := &domain.Customer{
		Name:  "João Silva",
		Email: "joao@example.com",
	}

	customerRepositoryResponse := model.FromDomain(*expectedCustomer)

	mockRepo.On("FindByID", ctx, customerID).Return(&customerRepositoryResponse, nil)

	// Act
	result, err := service.GetByID(ctx, customerID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCustomer.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetAll_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()

	expectedCustomers := []model.CustomerModel{
		{
			Name:  "João Silva",
			Email: "joao@example.com",
		},
		{
			Name:  "Maria Santos",
			Email: "maria@example.com",
		},
	}

	mockRepo.On("FindAll", ctx).Return(expectedCustomers, nil)

	// Act
	result, err := service.GetAll(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "João Silva", result[0].Name)
	assert.Equal(t, "joao@example.com", result[0].Email)
	assert.Equal(t, "Maria Santos", result[1].Name)
	assert.Equal(t, "maria@example.com", result[1].Email)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	customerID := uint(999)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("customer not found"))

	// Act
	result, err := service.GetByID(ctx, customerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
