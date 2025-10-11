package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do CustomerRepository
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) FindByID(ctx context.Context, id uint) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindAll(ctx context.Context) ([]domain.Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Create(ctx context.Context, entity *domain.Customer) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockCustomerRepository) Update(ctx context.Context, entity *domain.Customer) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

// Testes
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

	// Configurar o mock para esperar a chamada Create e retornar sucesso
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

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
		//Address: &domain.Address{
		//	Address:       "Rua Teste",
		//	AddressNumber: "123",
		//	City:          "São Paulo",
		//	Neighborhood:  "Centro",
		//	Country:       "Brasil",
		//	ZipCode:       "01234-567",
		//},
	}

	expectedError := errors.New("database connection error")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Customer")).Return(expectedError)

	// Act
	result, err := service.Create(ctx, inputCustomer)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_FindByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	ctx := context.Background()
	customerID := uint(1)
	expectedCustomer := &domain.Customer{
		Name:  "João Silva",
		Email: "joao@example.com",
	}

	mockRepo.On("FindByID", ctx, customerID).Return(expectedCustomer, nil)

	// Act
	result, err := service.GetByID(ctx, customerID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCustomer.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_FindByID_NotFound(t *testing.T) {
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
