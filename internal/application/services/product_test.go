package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestProductService_GetById_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productID := uint(1)

	expectedProduct := &domain.Product{
		ID:    productID,
		Name:  "Test Product",
		Price: 10.0,
		Stock: 100,
	}

	var productModel product.Model
	productModel.FromDomain(expectedProduct)

	mockRepo.On("FindByID", ctx, productID).Return(&productModel, nil)

	result, err := service.GetById(ctx, productID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProduct, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetById_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productID := uint(1)

	mockRepo.On("FindByID", ctx, productID).Return(nil, ErrProductNotFound)

	result, err := service.GetById(ctx, productID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetByIds_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productIDs := []uint{1, 2}

	expectedProducts := []domain.Product{
		{ID: 1, Name: "Product 1"},
		{ID: 2, Name: "Product 2"},
	}

	var productModels []product.Model
	for _, p := range expectedProducts {
		var model product.Model
		model.FromDomain(&p)
		productModels = append(productModels, model)
	}

	mockRepo.On("FindByIDs", ctx, productIDs).Return(productModels, nil)

	result, err := service.GetByIds(ctx, productIDs)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProducts, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()

	expectedProducts := []domain.Product{
		{ID: 1, Name: "Product 1"},
		{ID: 2, Name: "Product 2"},
	}

	var productModels []product.Model
	for _, p := range expectedProducts {
		var model product.Model
		model.FromDomain(&p)
		productModels = append(productModels, model)
	}

	mockRepo.On("FindAll", ctx).Return(productModels, nil)

	result, err := service.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProducts, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productToUpdate := domain.Product{
		ID:   1,
		Name: "Updated Product",
	}

	var model product.Model
	model.FromDomain(&productToUpdate)

	mockRepo.On("Update", ctx, &model).Return(nil)
	mockRepo.On("FindByID", ctx, productToUpdate.ID).Return(&model, nil)

	result, err := service.Update(ctx, productToUpdate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &productToUpdate, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_Fail(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productToUpdate := domain.Product{
		ID:   1,
		Name: "Updated Product",
	}

	var model product.Model
	model.FromDomain(&productToUpdate)

	dbErr := errors.New("db error")
	mockRepo.On("Update", ctx, &model).Return(dbErr)

	result, err := service.Update(ctx, productToUpdate)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update product")
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productToCreate := domain.Product{
		Name: "New Product",
	}

	var model product.Model
	model.FromDomain(&productToCreate)

	createdProduct := domain.Product{
		ID:   1,
		Name: "New Product",
	}
	var createdModel product.Model
	createdModel.FromDomain(&createdProduct)

	mockRepo.On("Create", ctx, &model).Return(&createdModel, nil)

	result, err := service.Create(ctx, productToCreate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &createdProduct, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create_Fail(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productToCreate := domain.Product{
		Name: "New Product",
	}

	var model product.Model
	model.FromDomain(&productToCreate)

	dbErr := errors.New("db error")
	mockRepo.On("Create", ctx, &model).Return(nil, dbErr)

	result, err := service.Create(ctx, productToCreate)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbErr, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productID := uint(1)

	deletedProduct := &domain.Product{ID: productID, Name: "Product to delete"}
	var deletedModel product.Model
	deletedModel.FromDomain(deletedProduct)

	mockRepo.On("FindByID", ctx, productID).Return(&deletedModel, nil)
	mockRepo.On("Delete", ctx, productID).Return(nil)

	result, err := service.Delete(ctx, productID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, deletedProduct, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productID := uint(1)

	mockRepo.On("FindByID", ctx, productID).Return(nil, ErrProductNotFound)

	result, err := service.Delete(ctx, productID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete_Fail(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productID := uint(1)

	productToDelete := &domain.Product{ID: productID, Name: "Product to delete"}
	var model product.Model
	model.FromDomain(productToDelete)

	dbErr := errors.New("db error")
	mockRepo.On("FindByID", ctx, productID).Return(&model, nil)
	mockRepo.On("Delete", ctx, productID).Return(dbErr)

	result, err := service.Delete(ctx, productID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbErr, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CheckAvailability_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productID := uint(1)
	quantity := uint(5)

	availableProduct := &domain.Product{ID: productID, Stock: 10}
	var model product.Model
	model.FromDomain(availableProduct)

	mockRepo.On("FindByID", ctx, productID).Return(&model, nil)

	available, err := service.CheckAvailability(ctx, productID, quantity)

	assert.NoError(t, err)
	assert.True(t, available)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CheckAvailability_InsufficientStock(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	productID := uint(1)
	quantity := uint(15)

	availableProduct := &domain.Product{ID: productID, Stock: 10}
	var model product.Model
	model.FromDomain(availableProduct)

	mockRepo.On("FindByID", ctx, productID).Return(&model, nil)

	available, err := service.CheckAvailability(ctx, productID, quantity)

	assert.Error(t, err)
	assert.False(t, available)
	assert.Equal(t, domain.ErrInsufficientStock, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateStock_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	operation := domain.StockOperationIncrease
	products := []domain.StockItem{
		{ID: 1, Quantity: 2, Operation: &operation},
		{ID: 2, Quantity: 3, Operation: &operation},
	}

	mockRepo.On("UpdateStock", ctx, uint(1), 2).Return(nil)
	mockRepo.On("UpdateStock", ctx, uint(2), 3).Return(nil)

	err := service.UpdateStock(ctx, products)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateStock_Fail(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()

	operation := domain.StockOperationIncrease

	products := []domain.StockItem{
		{ID: 1, Quantity: 2, Operation: &operation},
		{ID: 2, Quantity: 3, Operation: &operation},
	}

	dbErr := errors.New("db error")
	mockRepo.On("UpdateStock", ctx, uint(1), 2).Return(dbErr)

	err := service.UpdateStock(ctx, products)

	assert.Error(t, err)
	assert.Equal(t, dbErr, err)
	mockRepo.AssertExpectations(t)
}
