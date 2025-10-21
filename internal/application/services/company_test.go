package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompanyService_Create_Success(t *testing.T) {
	// Arrange (Preparar)
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()
	inputCompany := domain.Company{
		Name:     "Company Test",
		Email:    "company_test@example.com",
		Document: "12345678900",
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

	var mockModel company.Model
	mockModel.FromDomain(&inputCompany)

	// Configurar o mock para esperar a chamada Create e retornar sucesso
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType(reflect.TypeOf(&company.Model{}).String())).
		Return(&mockModel, nil)
	// Act (Agir)
	result, err := service.Create(ctx, inputCompany)

	// Assert (Verificar)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Company Test", result.Name)
	assert.Equal(t, "company_test@example.com", result.Email)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_Create_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()

	inputCompany := domain.Company{
		Name:  "Company Teste",
		Email: "company_test@example.com",
	}

	expectedError := errors.New("database connection error")

	mockRepo.On("Create", ctx, mock.AnythingOfType("*company.Model")).Return(nil, expectedError)

	// Act
	result, err := service.Create(ctx, inputCompany)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_FindById_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()
	CompanyID := uint(1)
	expectedCompany := &domain.Company{
		Name:  "Company Test",
		Email: "company_test@example.com",
	}

	var CompanyRepositoryResponse company.Model
	CompanyRepositoryResponse.FromDomain(expectedCompany)

	mockRepo.On("FindByID", ctx, CompanyID).Return(&CompanyRepositoryResponse, nil)

	// Act
	result, err := service.GetByID(ctx, CompanyID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCompany.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_FindById_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()
	CompanyID := uint(999)

	mockRepo.On("FindByID", ctx, CompanyID).Return(nil, errors.New("Company not found"))

	// Act
	result, err := service.GetByID(ctx, CompanyID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tests for UpdateById
func TestCompanyService_UpdateById_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()
	CompanyID := uint(1)
	inputCompany := domain.Company{
		ID:    CompanyID,
		Name:  "Updated Company",
		Email: "updated@example.com",
	}

	// Expect repository Update to be called and return nil (success)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *company.Model) bool {
		return m != nil && m.ID == CompanyID
	})).Return(nil)

	// Act
	err := service.UpdateByID(ctx, CompanyID, inputCompany)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_UpdateById_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()
	CompanyID := uint(999)
	inputCompany := domain.Company{
		ID:    CompanyID,
		Name:  "Does Not Exist",
		Email: "noone@example.com",
	}

	// Simulate repository returning an error when updating a non-existent entity
	repoErr := errors.New("company not found")
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*company.Model")).Return(repoErr)

	// Act
	err := service.UpdateByID(ctx, CompanyID, inputCompany)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, repoErr, err)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_DeleteById_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()
	companyID := uint(1)
	expectedCompany := &domain.Company{
		ID:    companyID,
		Name:  "Company To Delete",
		Email: "delete@example.com",
	}

	var companyModel company.Model
	companyModel.FromDomain(expectedCompany)

	mockRepo.On("FindByID", ctx, companyID).Return(&companyModel, nil)
	mockRepo.On("Delete", ctx, companyID).Return(nil)

	// Act
	result, err := service.DeleteByID(ctx, companyID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCompany.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCompanyService_DeleteById_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCompanyRepository)
	service := NewCompanyService(mockRepo)

	ctx := context.Background()
	companyID := uint(999)

	mockRepo.On("FindByID", ctx, companyID).Return(nil, errors.New("Company not found"))

	// Act
	result, err := service.DeleteByID(ctx, companyID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
