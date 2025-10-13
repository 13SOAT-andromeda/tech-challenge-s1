package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCustomerRepository_FindByEmail_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	email := "test@example.com"

	expectedCustomer := &domain.Customer{
		ID:       1,
		Name:     "John Doe",
		Email:    "test@example.com",
		Document: "12345678900",
	}

	mockRepo.On("FindByEmail", ctx, email).Return(expectedCustomer, nil)

	// Act
	result, err := mockRepo.FindByEmail(ctx, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCustomer.ID, result.ID)
	assert.Equal(t, expectedCustomer.Name, result.Name)
	assert.Equal(t, expectedCustomer.Email, result.Email)
	assert.Equal(t, expectedCustomer.Document, result.Document)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_FindByEmail_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	email := "notfound@example.com"

	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("record not found"))

	// Act
	result, err := mockRepo.FindByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_FindByEmail_EmptyEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	email := ""

	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("record not found"))

	// Act
	result, err := mockRepo.FindByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_FindByEmail_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	email := "error@example.com"
	expectedError := errors.New("database connection error")

	mockRepo.On("FindByEmail", ctx, email).Return(nil, expectedError)

	// Act
	result, err := mockRepo.FindByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// ==========================================
// TESTES DOS MÉTODOS DO BaseRepository
// ==========================================

func TestCustomerRepository_FindByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	customerID := uint(1)

	expectedCustomer := &domain.Customer{
		ID:    1,
		Name:  "John Doe",
		Email: "test@example.com",
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)

	// Act
	result, err := mockRepo.FindByID(ctx, customerID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCustomer.ID, result.ID)
	assert.Equal(t, expectedCustomer.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_FindByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	customerID := uint(999)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("record not found"))

	// Act
	result, err := mockRepo.FindByID(ctx, customerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()

	newCustomer := &domain.Customer{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Document: "98765432100",
		Type:     "individual",
	}

	mockRepo.On("Create", ctx, newCustomer).Return(nil)

	// Act
	err := mockRepo.Create(ctx, newCustomer)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_Create_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	expectedError := errors.New("duplicate email")

	newCustomer := &domain.Customer{
		Name:  "Jane Doe",
		Email: "existing@example.com",
	}

	mockRepo.On("Create", ctx, newCustomer).Return(expectedError)

	// Act
	err := mockRepo.Create(ctx, newCustomer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_Update_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()

	customer := &domain.Customer{
		ID:    1,
		Name:  "John Doe Updated",
		Email: "john.updated@example.com",
	}

	mockRepo.On("Update", ctx, customer).Return(nil)

	// Act
	err := mockRepo.Update(ctx, customer)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()
	customerID := uint(1)

	mockRepo.On("Delete", ctx, customerID).Return(nil)

	// Act
	err := mockRepo.Delete(ctx, customerID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_FindAll_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()

	expectedCustomers := []domain.Customer{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Type: "individual"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com", Type: "individual"},
	}

	mockRepo.On("FindAll", ctx).Return(expectedCustomers, nil)

	// Act
	result, err := mockRepo.FindAll(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "John Doe", result[0].Name)
	assert.Equal(t, "Jane Doe", result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestCustomerRepository_FindAll_Empty(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return([]domain.Customer{}, nil)

	// Act
	result, err := mockRepo.FindAll(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}
