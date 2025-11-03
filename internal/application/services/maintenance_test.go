package services

import (
	"context"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestMaintenanceService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	input := domain.Maintenance{
		Name:  "Maintenance Test",
		Price: 150,
	}

	var mockModel maintenance.Model
	mockModel.FromDomain(&input)

	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(&mockModel, nil)

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Maintenance Test", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	input := domain.Maintenance{
		Name:  "Maintenance Test",
		Price: 150,
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	expected := &maintenance.Model{
		Model:      gorm.Model{ID: serviceID},
		Name:       "Maintenance Test",
		Price:      100,
		CategoryId: "standard",
	}

	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(expected, nil)

	result, err := service.GetByID(ctx, serviceID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Maintenance Test", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)

	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetByID(ctx, serviceID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_GetByIDs_Success(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceIDs := []uint{1, 2}
	expected := []maintenance.Model{
		{Model: gorm.Model{ID: 1}, Name: "Maintenance 1", Price: 100, CategoryId: "standard"},
		{Model: gorm.Model{ID: 2}, Name: "Maintenance 2", Price: 150, CategoryId: "premium"},
	}

	mockRepo.On("FindByIDs", mock.Anything, serviceIDs).
		Return(expected, nil)

	result, err := service.GetByIDs(ctx, serviceIDs)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Maintenance 1", result[0].Name)
	assert.Equal(t, uint(2), result[1].ID)
	assert.Equal(t, "Maintenance 2", result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_GetByIDs_EmptyInput(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceIDs := []uint{}

	result, err := service.GetByIDs(ctx, serviceIDs)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	mockRepo.AssertNotCalled(t, "FindByIDs")
}

func TestMaintenanceService_GetByIDs_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceIDs := []uint{1, 2}

	mockRepo.On("FindByIDs", mock.Anything, serviceIDs).
		Return([]maintenance.Model{}, assert.AnError)

	result, err := service.GetByIDs(ctx, serviceIDs)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_UpdateByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	update := domain.Maintenance{
		ID:    serviceID,
		Name:  "Updated Maintenance",
		Price: 200,
	}

	var mockModel maintenance.Model
	mockModel.FromDomain(&update)

	mockRepo.On("Update", mock.Anything, mock.Anything).
		Return(nil)

	err := service.UpdateByID(ctx, serviceID, update)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_UpdateByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	update := domain.Maintenance{
		ID:    serviceID,
		Name:  "Updated Maintenance",
		Price: 200,
	}

	var mockModel maintenance.Model
	mockModel.FromDomain(&update)

	mockRepo.On("Update", mock.Anything, mock.Anything).
		Return(gorm.ErrRecordNotFound)

	err := service.UpdateByID(ctx, serviceID, update)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_UpdateByID_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	update := domain.Maintenance{
		ID:    serviceID,
		Name:  "Updated Maintenance",
		Price: 200,
	}

	mockRepo.On("Update", mock.Anything, mock.Anything).
		Return(assert.AnError)

	err := service.UpdateByID(ctx, serviceID, update)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_DeleteByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	existing := &maintenance.Model{
		Model: gorm.Model{ID: serviceID},
		Name:  "Maintenance to Delete",
	}

	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(existing, nil)
	mockRepo.On("Delete", mock.Anything, serviceID).
		Return(nil)

	result, err := service.DeleteByID(ctx, serviceID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Maintenance to Delete", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_DeleteByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)

	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(nil, gorm.ErrRecordNotFound)

	result, err := service.DeleteByID(ctx, serviceID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_DeleteByID_DeleteError(t *testing.T) {
	mockRepo := new(mocks.MockMaintenanceRepository)
	service := NewMaintenanceService(mockRepo)

	ctx := context.Background()
	serviceID := uint(1)
	existing := &maintenance.Model{
		Model: gorm.Model{ID: serviceID},
		Name:  "Maintenance to Delete",
	}

	mockRepo.On("FindByID", mock.Anything, serviceID).
		Return(existing, nil)
	mockRepo.On("Delete", mock.Anything, serviceID).
		Return(assert.AnError)

	result, err := service.DeleteByID(ctx, serviceID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
