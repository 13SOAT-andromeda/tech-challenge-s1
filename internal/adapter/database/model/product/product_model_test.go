package product

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProductModelInitialization(t *testing.T) {
	p := Model{
		Name:  "Product A",
		Stock: 10,
		Price: 1000, // R$ 10.00
	}

	assert.NotNil(t, p)
	assert.Equal(t, "Product A", p.Name)
	assert.Equal(t, uint(10), p.Stock)
	assert.Equal(t, int64(1000), p.Price)
}

func TestProductModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	modelProduct := &Model{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		Name:  "Product A",
		Stock: 10,
		Price: 1000,
	}

	domainProduct := modelProduct.ToDomain()

	assert.Equal(t, modelProduct.ID, domainProduct.ID)
	assert.Equal(t, modelProduct.Name, domainProduct.Name)
	assert.Equal(t, modelProduct.Stock, *domainProduct.Stock)
	assert.Equal(t, modelProduct.Price, domainProduct.Price)

}
