package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function para criar um produto de teste
func createTestProductModel(id uint, stock uint, price int64) *product.Model {
	p := &product.Model{
		Name:  "Test Product",
		Price: price,
		Stock: stock,
	}
	p.ID = id // Importante: setar o ID do gorm.Model
	return p
}

func TestProductService_ConfirmOrderProducts_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("UpdateStock", ctx, uint(1), 5).Return(nil)

	err := service.ConfirmOrderProducts(ctx, 1, 5)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_ConfirmOrderProducts_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	expectedError := errors.New("update failed")
	mockRepo.On("UpdateStock", ctx, uint(1), 5).Return(expectedError)

	err := service.ConfirmOrderProducts(ctx, 1, 5)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	existingProduct := createTestProductModel(1, 10, 10000)
	domainProduct := domain.Product{
		ID:    1,
		Name:  "Updated Product",
		Price: 15000,
		Stock: 15,
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(nil)

	result, err := service.Update(ctx, domainProduct)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Product", result.Name)
	assert.Equal(t, int64(15000), result.Price)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_ProductNotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	domainProduct := domain.Product{ID: 999}

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("not found"))

	result, err := service.Update(ctx, domainProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found or disabled")
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_UpdateFails(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	existingProduct := createTestProductModel(1, 10, 10000)
	domainProduct := domain.Product{ID: 1, Name: "Updated"}

	mockRepo.On("FindByID", ctx, uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(errors.New("update error"))

	result, err := service.Update(ctx, domainProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update product")
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateStock_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	existingProduct := createTestProductModel(1, 10, 10000) // Stock atual = 10
	domainProduct := domain.Product{
		ID:    1,
		Name:  "Updated Product",
		Price: 12000,
		Stock: 20, // Novo estoque
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(m *product.Model) bool {
		return m.Stock == 20 // Esperamos que o estoque seja atualizado
	})).Return(nil)

	result, err := service.UpdateStock(ctx, domainProduct)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(20), result.Stock)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateStock_ProductNotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	domainProduct := domain.Product{ID: 999, Stock: 5}

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("not found"))

	result, err := service.UpdateStock(ctx, domainProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found or disabled")
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateStock_UpdateFails(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	existingProduct := createTestProductModel(1, 10, 10000)
	domainProduct := domain.Product{
		ID:    1,
		Name:  "Updated Product",
		Price: 12000,
		Stock: 50,
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(errors.New("update error"))

	result, err := service.UpdateStock(ctx, domainProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update product")
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_DoesNotChangeStock(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	existingProduct := createTestProductModel(1, 10, 10000) // Stock atual = 10
	domainProduct := domain.Product{
		ID:    1,
		Name:  "Updated Product",
		Price: 15000,
		Stock: 0, // Mesmo que venha 0, não deve sobrescrever
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(m *product.Model) bool {
		return m.Stock == 10 // Esperamos que continue com o valor antigo
	})).Return(nil)

	result, err := service.Update(ctx, domainProduct)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(10), existingProduct.Stock)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetById_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)

	result, err := service.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Test Product", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetById_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("not found"))

	result, err := service.GetById(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetByIds_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	products := []product.Model{
		*createTestProductModel(1, 10, 10000),
		*createTestProductModel(2, 20, 20000),
	}

	mockRepo.On("FindByIDs", ctx, []uint{1, 2}).Return(products, nil)

	result, err := service.GetByIds(ctx, []uint{1, 2})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, uint(2), result[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetByIds_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByIDs", ctx, []uint{1, 2}).Return(nil, errors.New("database error"))

	result, err := service.GetByIds(ctx, []uint{1, 2})

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	products := []product.Model{
		*createTestProductModel(1, 10, 10000),
		*createTestProductModel(2, 20, 20000),
	}

	mockRepo.On("Search", ctx, mock.Anything).Return(products, nil)

	result, err := service.GetAll(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Test Product", result[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetAll_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Search", ctx, mock.AnythingOfType("filter.ProductFilter")).Return(nil, errors.New("database error"))

	result, err := service.GetAll(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	domainProduct := domain.Product{
		Name:  "New Product",
		Price: 10000,
		Stock: 10,
	}

	createdProduct := createTestProductModel(1, 10, 10000)
	createdProduct.Name = "New Product"
	mockRepo.On("Create", ctx, mock.AnythingOfType("*product.Model")).Return(createdProduct, nil)

	result, err := service.Create(ctx, domainProduct)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Product", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	domainProduct := domain.Product{Name: "New Product"}
	expectedError := errors.New("create error")

	mockRepo.On("Create", ctx, mock.AnythingOfType("*product.Model")).Return(nil, expectedError)

	result, err := service.Create(ctx, domainProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	result, err := service.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete_ProductNotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("not found"))

	result, err := service.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete_DeleteError(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("delete error"))

	result, err := service.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_ManageStockItem_AddOperation(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(nil)

	result, err := service.ManageStockItem(ctx, 1, 5, "add")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(15), result.Stock)
	mockRepo.AssertExpectations(t)
}

func TestProductService_ManageStockItem_RemoveOperation(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(nil)

	result, err := service.ManageStockItem(ctx, 1, 5, "remove")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_ManageStockItem_CaseInsensitive(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(nil)

	result, err := service.ManageStockItem(ctx, 1, 5, "ADD")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_ManageStockItem_InvalidProductId(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	result, err := service.ManageStockItem(ctx, 0, 5, "add")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidProductId, err)
	assert.Nil(t, result)
}

func TestProductService_ManageStockItem_InvalidQuantity(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	result, err := service.ManageStockItem(ctx, 1, 0, "add")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidQuantity, err)
	assert.Nil(t, result)
}

func TestProductService_ManageStockItem_InvalidOperation(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	result, err := service.ManageStockItem(ctx, 1, 5, "invalid")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidManageStockOperation, err)
	assert.Nil(t, result)
}

func TestProductService_AddStockItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(nil)

	result, err := service.AddStockItem(ctx, 1, 5)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(15), result.Stock)
	mockRepo.AssertExpectations(t)
}

func TestProductService_AddStockItem_InvalidProductId(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	result, err := service.AddStockItem(ctx, 0, 5)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidProductId, err)
	assert.Nil(t, result)
}

func TestProductService_AddStockItem_InvalidQuantity(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	result, err := service.AddStockItem(ctx, 1, 0)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidQuantity, err)
	assert.Nil(t, result)
}

func TestProductService_AddStockItem_ProductNotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("not found"))

	result, err := service.AddStockItem(ctx, 999, 5)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_AddStockItem_UpdateFails(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(errors.New("update failed"))

	result, err := service.AddStockItem(ctx, 1, 5)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_RemoveStockItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(nil)

	result, err := service.RemoveStockItem(ctx, 1, 5)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_RemoveStockItem_InvalidQuantity(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	result, err := service.RemoveStockItem(ctx, 1, 0)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidQuantity, err)
	assert.Nil(t, result)
}

func TestProductService_RemoveStockItem_ProductNotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("not found"))

	result, err := service.RemoveStockItem(ctx, 999, 5)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_RemoveStockItem_UpdateFails(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()

	productModel := createTestProductModel(1, 10, 10000)
	mockRepo.On("FindByID", ctx, uint(1)).Return(productModel, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*product.Model")).Return(errors.New("update failed"))

	result, err := service.RemoveStockItem(ctx, 1, 5)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
