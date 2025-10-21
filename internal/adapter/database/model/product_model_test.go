package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProductTableName(t *testing.T) {
	assert.Equal(t, "Product", ProductModel{}.TableName())
}

func TestNilProductToDomain(t *testing.T) {
	assert.Nil(t, (*ProductModel)(nil).ToDomain())
}

func TestNilProductFromDomain(t *testing.T) {
	assert.Nil(t, FromDomainProduct(nil))
}

func TestProductModelInitialization(t *testing.T) {
	p := ProductModel{
		Name:     "Product A",
		Quantity: 10,
		Price:    1000, // R$ 10.00
	}

	assert.NotNil(t, p)
	assert.Equal(t, "Product A", p.Name)
	assert.Equal(t, uint(10), p.Quantity)
	assert.Equal(t, uint32(1000), p.Price)
}

func TestProductModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	modelProduct := &ProductModel{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		Name:     "Product A",
		Quantity: 10,
		Price:    1000,
	}

	domainProduct := modelProduct.ToDomain()

	assert.Equal(t, modelProduct.ID, domainProduct.ID)
	assert.Equal(t, modelProduct.Name, domainProduct.Name)
	assert.Equal(t, modelProduct.Quantity, domainProduct.Quantity)
	assert.Equal(t, modelProduct.Price, domainProduct.Price)
	assert.Equal(t, modelProduct.CreatedAt, domainProduct.CreatedAt)
	assert.Equal(t, modelProduct.UpdatedAt, domainProduct.UpdatedAt)

	newModel := FromDomainProduct(domainProduct)
	newModel.DeletedAt = modelProduct.DeletedAt

	assert.Equal(t, modelProduct, newModel)
}
