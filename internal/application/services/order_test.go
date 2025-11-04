package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestOrderModel(id uint) *order.Model {
	d := time.Now()
	om := &order.Model{
		DateIn:            d,
		Status:            string(domain.OrderStatuses.RECEIVED),
		VehicleKilometers: 100,
		UserID:            1,
		CustomerVehicleID: 1,
		CompanyID:         1,
	}
	om.ID = id
	return om
}

func TestOrderService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	domainOrder := domain.Order{
		DateIn:            time.Now(),
		Status:            domain.OrderStatuses.RECEIVED,
		VehicleKilometers: 100,
		User:              domain.User{ID: 1},
		CustomerVehicle:   domain.CustomerVehicle{ID: 1},
		Company:           domain.Company{ID: 1},
	}

	createdModel := createTestOrderModel(1)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*order.Model")).Return(createdModel, nil)

	result, err := service.Create(ctx, domainOrder)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.OrderStatuses.RECEIVED, result.Status)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	domainOrder := domain.Order{Status: domain.OrderStatuses.FINISHED}
	expectedErr := errors.New("create failed")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*order.Model")).Return(nil, expectedErr)

	result, err := service.Create(ctx, domainOrder)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	om := createTestOrderModel(1)
	mockRepo.On("FindOrderByID", ctx, uint(1)).Return(om, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindOrderByID", ctx, uint(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	models := []order.Model{*createTestOrderModel(1), *createTestOrderModel(2)}
	mockRepo.On("Search", mock.Anything, mock.Anything).Return(models, nil)

	result, err := service.GetAll(ctx, map[string]interface{}{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Search", ctx, mock.Anything).Return(nil, errors.New("search error"))

	result, err := service.GetAll(ctx, map[string]interface{}{})

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := service.Delete(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_Delete_DeleteError(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("delete error"))

	err := service.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrOrderDelete, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	domainOrder := domain.Order{ID: 1, VehicleKilometers: 200}
	mockRepo.On("Update", ctx, mock.AnythingOfType("*order.Model")).Return(nil)

	err := service.Update(ctx, domainOrder)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_Update_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	domainOrder := domain.Order{ID: 1}
	mockRepo.On("Update", ctx, mock.AnythingOfType("*order.Model")).Return(errors.New("update failed"))

	err := service.Update(ctx, domainOrder)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_WithStatusFilter(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	models := []order.Model{*createTestOrderModel(1)}
	mockRepo.On("Search", mock.Anything, mock.Anything).Return(models, nil)

	params := map[string]interface{}{
		"status": "completed",
	}
	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_WithEnabledFalse(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	expectedSearch := ports.OrderSearch{
		Status:   "",
		Enabled:  false,
		OrderBy:  "",
		SortDesc: false,
	}

	models := []order.Model{*createTestOrderModel(1)}
	mockRepo.On("Search", ctx, expectedSearch).Return(models, nil)

	params := map[string]interface{}{
		"enabled": "false",
	}
	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_WithOrderBy(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	expectedSearch := ports.OrderSearch{
		Status:   "",
		Enabled:  true,
		OrderBy:  "created_at",
		SortDesc: false,
	}

	models := []order.Model{*createTestOrderModel(1), *createTestOrderModel(2)}
	mockRepo.On("Search", ctx, expectedSearch).Return(models, nil)

	params := map[string]interface{}{
		"orderby": "created_at",
	}
	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_WithSortDescTrue(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	expectedSearch := ports.OrderSearch{
		Status:   "",
		Enabled:  true,
		OrderBy:  "created_at",
		SortDesc: true,
	}

	models := []order.Model{*createTestOrderModel(2), *createTestOrderModel(1)}
	mockRepo.On("Search", ctx, expectedSearch).Return(models, nil)

	params := map[string]interface{}{
		"orderby":  "created_at",
		"sortdesc": "true",
	}
	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_WithSortDescOne(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	expectedSearch := ports.OrderSearch{
		Status:   "",
		Enabled:  true,
		OrderBy:  "",
		SortDesc: true,
	}

	models := []order.Model{*createTestOrderModel(1)}
	mockRepo.On("Search", ctx, expectedSearch).Return(models, nil)

	params := map[string]interface{}{
		"sortdesc": "1",
	}
	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_WithAllParameters(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	expectedSearch := ports.OrderSearch{
		Status:   "pending",
		Enabled:  false,
		OrderBy:  "id",
		SortDesc: true,
	}

	models := []order.Model{*createTestOrderModel(1)}
	mockRepo.On("Search", ctx, expectedSearch).Return(models, nil)

	params := map[string]interface{}{
		"status":   "pending",
		"enabled":  "false",
		"orderby":  "id",
		"sortdesc": "true",
	}
	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_EmptyResult(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	models := []order.Model{}
	mockRepo.On("Search", ctx, mock.Anything).Return(models, nil)

	result, err := service.GetAll(ctx, map[string]interface{}{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 0)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetAll_WithInvalidSortDesc(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	// sortdesc com valor inválido deve manter o padrão false
	expectedSearch := ports.OrderSearch{
		Status:   "",
		Enabled:  true,
		OrderBy:  "",
		SortDesc: false,
	}

	models := []order.Model{*createTestOrderModel(1)}
	mockRepo.On("Search", ctx, expectedSearch).Return(models, nil)

	params := map[string]interface{}{
		"sortdesc": "invalid",
	}
	result, err := service.GetAll(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}
