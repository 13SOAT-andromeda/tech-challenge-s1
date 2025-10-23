package services

import (
	"context"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func float64Ptr(v float64) *float64 { return &v }

func TestMaintenanceService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	inputService := domain.Maintenance{
		Name:         "Maintenance Test",
		DefaultPrice: float64Ptr(150.0),
		CategoryId:   2,
		Number:       "SVC1001",
	}

	var mockModel maintenance.Model
	mockModel.FromDomain(&inputService)

	// Configurar o mock para esperar a chamada Create e retornar sucesso
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(&mockModel, nil)

	// Act
	result, err := service.Create(ctx, inputService)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Maintenance Test", result.Name)
	if result.DefaultPrice != nil {
		assert.Equal(t, 150.0, *result.DefaultPrice)
	} else {
		t.Fatalf("result.DefaultPrice is nil")
	}
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_Create_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	inputService := domain.Maintenance{
		Name:         "Maintenance Test",
		DefaultPrice: float64Ptr(150.0),
		CategoryId:   2,
		Number:       "SVC1001",
	}

	// Configurar o mock para esperar a chamada Create e retornar um erro
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	// Act
	result, err := service.Create(ctx, inputService)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_FindByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	expectedService := &maintenance.Model{
		Model:  gorm.Model{ID: serviceID},
		Name:   "Service Test",
		Number: "SVC1001",
	}

	// Configurar o mock para esperar a chamada FindByID e retornar sucesso
	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(expectedService, nil)

	// Act
	result, err := service.GetByID(ctx, serviceID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Maintenance Test", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_FindByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)

	// Configurar o mock para esperar a chamada FindByID e retornar um erro de não encontrado
	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.GetByID(ctx, serviceID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_UpdateByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	updateService := domain.Maintenance{
		ID:           serviceID,
		Name:         "Updated Maintenance",
		DefaultPrice: float64Ptr(200.0),
		CategoryId:   3,
		Number:       "SVC2002",
	}

	var mockModel maintenance.Model
	mockModel.FromDomain(&updateService)

	// Configurar o mock para esperar a chamada Update e retornar sucesso
	mockRepo.On("Update", mock.Anything, mock.Anything).
		Return(nil)

	// Act
	err := service.UpdateByID(ctx, serviceID, updateService)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_UpdateByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	updateService := domain.Maintenance{
		ID:           serviceID,
		Name:         "Updated Maintenance",
		DefaultPrice: float64Ptr(200.0),
		CategoryId:   3,
		Number:       "SVC2002",
	}

	var mockModel maintenance.Model
	mockModel.FromDomain(&updateService)

	// Configurar o mock para esperar a chamada Update e retornar um erro de não encontrado
	mockRepo.On("Update", mock.Anything, mock.Anything).
		Return(gorm.ErrRecordNotFound)

	// Act
	err := service.UpdateByID(ctx, serviceID, updateService)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_DeleteByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	existingService := &maintenance.Model{
		Model:  gorm.Model{ID: serviceID},
		Name:   "Service to Delete",
		Number: "SVC3003",
	}

	// Configurar o mock para esperar a chamada FindByID e retornar o serviço existente
	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(existingService, nil)

	// Configurar o mock para esperar a chamada Delete e retornar sucesso
	mockRepo.On("Delete", mock.Anything, serviceID).
		Return(nil)

	// Act
	result, err := service.DeleteByID(ctx, serviceID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Maintenance to Delete", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_DeleteByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaintenanceRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)

	// Configurar o mock para esperar a chamada FindByID e retornar um erro de não encontrado
	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.DeleteByID(ctx, serviceID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
