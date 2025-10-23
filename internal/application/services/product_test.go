package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_Create_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	inputProduct := domain.Product{
		Name:  "Produto Teste",
		Price: int64(12345), // price as int64 (e.g. cents)
	}

	var mockModel product.Model
	mockModel.FromDomain(&inputProduct)
	// repository should return the created model
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Model")).Return(&mockModel, nil)

	result, err := service.Create(ctx, inputProduct)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, inputProduct.Name, result.Name)
	assert.Equal(t, inputProduct.Price, result.Price)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	inputProduct := domain.Product{
		Name:  "Produto Falha",
		Price: int64(0),
	}

	expectedError := errors.New("database error")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*product.Model")).Return(nil, expectedError)

	result, err := service.Create(ctx, inputProduct)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetById_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	id := uint(1)

	expected := domain.Product{
		Name:  "Produto A",
		Price: int64(10),
	}

	var repoResp product.Model
	repoResp.FromDomain(&expected)

	mockRepo.On("FindByID", ctx, id).Return(&repoResp, nil)

	res, err := service.GetById(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetById_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	id := uint(999)

	mockRepo.On("FindByID", ctx, id).Return(nil, errors.New("not found"))

	res, err := service.GetById(ctx, id)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetAll_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()

	repoList := []product.Model{
		{Name: "P1", Price: int64(11)},
		{Name: "P2", Price: int64(22)},
	}

	mockRepo.On("FindAll", ctx).Return(repoList, nil)

	res, err := service.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 2)
	assert.Equal(t, "P1", res[0].Name)
	assert.Equal(t, "P2", res[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete(t *testing.T) {
	// success
	mockOK := new(MockProductRepository)
	svcOK := NewProductService(mockOK)
	mockOK.On("Delete", mock.Anything, uint(1)).Return(nil)

	if err := svcOK.Delete(context.Background(), 1); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	mockOK.AssertExpectations(t)

	// failure
	mockErr := new(MockProductRepository)
	svcErr := NewProductService(mockErr)
	mockErr.On("Delete", mock.Anything, uint(2)).Return(errors.New("delete failed"))

	if err := svcErr.Delete(context.Background(), 2); err == nil {
		t.Fatalf("expected error on delete, got nil")
	}
	mockErr.AssertExpectations(t)
}
