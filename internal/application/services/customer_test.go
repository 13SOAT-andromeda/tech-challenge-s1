package services

import (
	"context"
	"errors"
	"testing"

	customerModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	userModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestCustomerService(mockRepo *mocks.MockCustomerRepository, mockPersonRepo *mocks.MockPersonRepository, mockUserRepo *mocks.MockUserRepository) *customerService {
	return NewCustomerService(mockRepo, mockPersonRepo, mockUserRepo)
}

func makeTestPerson(name, email string) *domain.Person {
	return &domain.Person{
		Name:    name,
		Email:   email,
		Contact: "11999999999",
		Document: &domain.Document{
			Number: "293.034.620-50",
		},
		Address: &domain.Address{
			Address:       "Rua Teste",
			AddressNumber: "123",
			City:          "São Paulo",
			Neighborhood:  "Centro",
			Country:       "Brasil",
			ZipCode:       "01234-567",
		},
	}
}

func TestCustomerService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockHasher := new(encryption.MockHasher)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()
	inputCustomer := domain.Customer{
		Type:   "individual",
		Person: makeTestPerson("Gedan Magalhaes", "gedan@example.com"),
	}

	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)
	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	createdPerson := personModel.Model{}
	createdPerson.ID = 1
	createdPerson.Name = "Gedan Magalhaes"
	createdPerson.Email = "gedan@example.com"

	var mockModelCust customerModel.Model
	mockModelCust.PersonID = 1
	mockModelCust.Person = createdPerson

	createdUser := userModel.Model{}
	createdUser.ID = 1

	mockRepo.On("FindByDocument", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)
	mockUserRepo.On("GetByEmail", mock.Anything, "gedan@example.com").Return(nil, nil)
	mockPersonRepo.On("Create", mock.Anything, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*customer.Model")).Return(&mockModelCust, nil)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.Model")).Return(&createdUser, nil)

	result, err := service.Create(ctx, inputCustomer, password)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Person)
	assert.Equal(t, "Gedan Magalhaes", result.Person.Name)
	assert.Equal(t, "gedan@example.com", result.Person.Email)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestCustomerService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockHasher := new(encryption.MockHasher)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()
	inputCustomer := domain.Customer{
		Type:   "individual",
		Person: makeTestPerson("João Silva", "gedan@example.com"),
	}

	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)
	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	expectedError := errors.New("database connection error")

	createdPerson := personModel.Model{}
	createdPerson.ID = 1

	mockRepo.On("FindByDocument", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)
	mockUserRepo.On("GetByEmail", mock.Anything, "gedan@example.com").Return(nil, nil)
	mockPersonRepo.On("Create", mock.Anything, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*customer.Model")).Return(nil, expectedError)

	result, err := service.Create(ctx, inputCustomer, password)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCustomerService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()
	customerID := uint(1)

	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = "Gedan Magalhães"
	pm.Email = "gedan@example.com"

	var customerRepositoryResponse customerModel.Model
	customerRepositoryResponse.PersonID = 1
	customerRepositoryResponse.Person = pm

	mockRepo.On("FindByID", ctx, customerID).Return(&customerRepositoryResponse, nil)

	result, err := service.GetByID(ctx, customerID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Person)
	assert.Equal(t, "Gedan Magalhães", result.Person.Name)

	mockRepo.AssertExpectations(t)
}

func TestCustomerService_Search_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()

	pm1 := personModel.Model{}
	pm1.ID = 1
	pm1.Name = "Gedan Magalhaes"
	pm1.Email = "gedan@example.com"

	pm2 := personModel.Model{}
	pm2.ID = 2
	pm2.Name = "Elen Magalhaes"
	pm2.Email = "elen@example.com"

	expectedCustomers := []customerModel.Model{
		{PersonID: 1, Person: pm1},
		{PersonID: 2, Person: pm2},
	}

	mockRepo.On("Search", ctx, mock.Anything).Return(expectedCustomers, nil)

	result, err := service.Search(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Gedan Magalhaes", result[0].Person.Name)
	assert.Equal(t, "gedan@example.com", result[0].Person.Email)
	assert.Equal(t, "Elen Magalhaes", result[1].Person.Name)
	assert.Equal(t, "elen@example.com", result[1].Person.Email)

	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()
	customerID := uint(999)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("customer not found"))

	result, err := service.GetByID(ctx, customerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)
	ctx := context.Background()

	inputCustomer := domain.Customer{
		ID:   1,
		Type: "individual",
		Person: &domain.Person{
			Name:  "Gedan Magalhaes",
			Email: "gedan@example.com",
			Document: &domain.Document{
				Number: "293.034.620-50",
			},
		},
	}

	existingPm := personModel.Model{}
	existingPm.ID = 1
	existingPm.Name = "Gedan Magalhaes"

	var mockModel customerModel.Model
	mockModel.PersonID = 1
	mockModel.Person = existingPm

	existingPersonModel := personModel.Model{}
	existingPersonModel.ID = 1

	mockRepo.On("FindByID", ctx, uint(1)).Return(&mockModel, nil)
	mockRepo.On("FindByDocument", ctx, inputCustomer.Person.Document.GetDocumentNumber()).Return(nil, errors.New("record not found"))
	mockPersonRepo.On("FindByID", ctx, uint(1)).Return(&existingPersonModel, nil)
	mockPersonRepo.On("Update", ctx, mock.AnythingOfType("*person.Model")).Return(nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*customer.Model")).Return(nil)

	err := service.UpdateByID(ctx, 1, inputCustomer)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateByID_CustomerNotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)
	ctx := context.Background()

	inputCustomer := domain.Customer{Type: "individual"}

	mockRepo.On("FindByID", ctx, uint(2)).Return(nil, errors.New("record not found"))

	err := service.UpdateByID(ctx, 2, inputCustomer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Customer with Id 2 not found")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateByID_DocumentAlreadyInUse(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)
	ctx := context.Background()

	inputCustomer := domain.Customer{
		ID:   3,
		Type: "individual",
		Person: &domain.Person{
			Name: "Gedan",
			Document: &domain.Document{
				Number: "293.034.620-50",
			},
		},
	}

	pm := personModel.Model{}
	pm.ID = 1

	var mockModelCust customerModel.Model
	mockModelCust.ID = 3
	mockModelCust.PersonID = 1
	mockModelCust.Person = pm

	anotherModel := customerModel.Model{}
	anotherModel.ID = 4
	anotherModel.PersonID = 2

	mockRepo.On("FindByID", ctx, uint(3)).Return(&mockModelCust, nil)
	mockRepo.On("FindByDocument", ctx, inputCustomer.Person.Document.GetDocumentNumber()).Return(&anotherModel, nil)

	err := service.UpdateByID(ctx, 3, inputCustomer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Number is invalid or already in use")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_DeleteByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()
	customerID := uint(1)

	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = "John Doe"
	pm.Email = "john@example.com"

	var cm customerModel.Model
	cm.ID = customerID
	cm.PersonID = 1
	cm.Person = pm

	mockRepo.On("FindByID", ctx, customerID).Return(&cm, nil)
	mockRepo.On("Delete", ctx, customerID).Return(nil)

	result, err := service.DeleteByID(ctx, customerID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Person)
	assert.Equal(t, "John Doe", result.Person.Name)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_DeleteByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()
	customerID := uint(999)

	mockRepo.On("FindByID", ctx, customerID).Return(nil, errors.New("customer not found"))

	result, err := service.DeleteByID(ctx, customerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_DeleteByID_DeleteError(t *testing.T) {
	mockRepo := new(mocks.MockCustomerRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	service := newTestCustomerService(mockRepo, mockPersonRepo, mockUserRepo)

	ctx := context.Background()
	customerID := uint(1)

	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = "John Doe"

	var cm customerModel.Model
	cm.ID = customerID
	cm.PersonID = 1
	cm.Person = pm

	mockRepo.On("FindByID", ctx, customerID).Return(&cm, nil)
	mockRepo.On("Delete", ctx, customerID).Return(errors.New("delete error"))

	result, err := service.DeleteByID(ctx, customerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
