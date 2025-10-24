package order_product

import (
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/stretchr/testify/assert"
)

func TestOrderProductModelInitialization(t *testing.T) {
	op := Model{
		ProductId: 1,
		OrderId:   1,
	}

	assert.NotNil(t, op)
	assert.Equal(t, uint(1), op.ProductId)
	assert.Equal(t, uint(1), op.OrderId)
}

func TestOrderProductModel_ToFromDomain(t *testing.T) {
	modelOrderProduct := Model{
		ProductId: 1,
		OrderId:   1,
		Product:   product.Model{},
		Order:     order.Model{User: user.Model{}},
	}

	domain := modelOrderProduct.ToDomain()

	assert.Equal(t, modelOrderProduct.ProductId, domain.ProductId)
	assert.Equal(t, modelOrderProduct.OrderId, domain.OrderId)

}
