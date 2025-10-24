package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCustomerService_Create_Success(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	inputCustomer := domain.Customer{
		Name:     "Gedan Magalhaes",
		Email:    "gedan@example.com",
		Document: "293.034.620-50",
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

	var mockModel customer.Model
	mockModel.FromDomain(&inputCustomer)

	mockRepo.On("FindByDocument", mock.Anything, mock.AnythingOfType("string")).
		Return(nil, nil)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*customer.Model")).
		Return(&mockModel, nil)

	result, err := service.Create(ctx, inputCustomer)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Gedan Magalhaes", result.Name)
	assert.Equal(t, "gedan@example.com", result.Email)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_Create_Documento_Is_Invalid(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	inputCustomer := domain.Customer{
		Name:     "Gedan Magalhaes",
		Email:    "gedan@example.com",
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

	var mockModel customer.Model
	mockModel.FromDomain(&inputCustomer)

	result, err := service.Create(ctx, inputCustomer)

	assert.Nil(t, result)
	assert.EqualError(t, err, "Document is invalid")
}

func TestCustomerService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()

	inputCustomer := domain.Customer{
		Name:     "João Silva",
		Email:    "gedan@example.com",
		Document: "293.034.620-50",
	}

	expectedError := errors.New("database connection error")

	mockRepo.On("FindByDocument", mock.Anything, mock.AnythingOfType("string")).
		Return(nil, nil)

	mockRepo.On("Create", ctx, mock.AnythingOfType("*customer.Model")).Return(nil, expectedError)

	result, err := service.Create(ctx, inputCustomer)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	customerID := uint(1)

	expectedCustomer := domain.Customer{
		Name:  "Gedan Magalhães",
		Email: "gedan@example.com",
	}

	var customerRepositoryResponse customer.Model
	customerRepositoryResponse.FromDomain(&expectedCustomer)

	mockRepo.On("FindByID", ctx, customerID).Return(&customerRepositoryResponse, nil)

	result, err := service.GetByID(ctx, customerID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCustomer.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetAll_Success(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()

	expectedCustomers := []customer.Model{
		{
			Name:  "Gedan Magalhaes",
			Email: "gedan@example.com",
		},
		{
			Name:  "Elen Magalhaes",
			Email: "elen@example.com",
		},
	}

	mockRepo.On("FindAll", ctx).Return(expectedCustomers, nil)

	result, err := service.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Gedan Magalhaes", result[0].Name)
	assert.Equal(t, "gedan@example.com", result[0].Email)
	assert.Equal(t, "Elen Magalhaes", result[1].Name)
	assert.Equal(t, "elen@example.com", result[1].Email)

	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	customerID := uint(999)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("customer not found"))

	result, err := service.GetByID(ctx, customerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateByID_Success(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	inputCustomer := domain.Customer{
		ID:       1,
		Name:     "Gedan Magalhaes",
		Email:    "gedan@example.com",
		Document: "293.034.620-50",
	}
	var mockModel customer.Model

	mockModel.FromDomain(&inputCustomer)

	mockRepo.On("FindByID", ctx, uint(1)).Return(&mockModel, nil)

	mockRepo.On("FindByDocument", ctx, inputCustomer.Document).Return(nil, errors.New("record not found"))

	mockRepo.On("Update", ctx, mock.AnythingOfType("*customer.Model")).Return(nil)

	err := service.UpdateByID(ctx, 1, inputCustomer)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateByID_CustomerNotFound(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	inputCustomer := domain.Customer{Name: "Gedan"}

	mockRepo.On("FindByID", ctx, uint(2)).Return(nil, errors.New("record not found"))

	err := service.UpdateByID(ctx, 2, inputCustomer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Customer with Id 2 not found")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateByID_DocumentAlreadyInUse(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	inputCustomer := domain.Customer{
		ID:       3,
		Name:     "Gedan",
		Document: "293.034.620-50",
	}
	var mockModel customer.Model
	mockModel.FromDomain(&inputCustomer)

	mockRepo.On("FindByID", ctx, uint(3)).Return(&mockModel, nil)

	anotherInputCustomer := domain.Customer{
		ID:       4,
		Name:     "Gedan",
		Document: "293.034.620-50",
	}
	mockModel.FromDomain(&anotherInputCustomer)

	mockRepo.On("FindByDocument", ctx, inputCustomer.Document).Return(&mockModel, nil)

	err := service.UpdateByID(ctx, 3, inputCustomer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Document is invalid or already in use")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateByID_UpdateFails(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	inputCustomer := domain.Customer{
		ID:       5,
		Name:     "Gedan",
		Document: "293.034.620-50",
	}

	var mockModel customer.Model

	mockModel.FromDomain(&inputCustomer)

	mockRepo.On("FindByID", ctx, uint(5)).Return(&mockModel, nil)

	mockRepo.On("FindByDocument", ctx, inputCustomer.Document).Return(nil, errors.New("record not found"))

	mockRepo.On("Update", ctx, mock.AnythingOfType("*customer.Model")).Return(errors.New("db error"))

	err := service.UpdateByID(ctx, 5, inputCustomer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to update customer")
	mockRepo.AssertExpectations(t)
}
