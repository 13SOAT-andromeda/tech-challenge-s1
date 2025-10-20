package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderProductModelInitialization(t *testing.T) {
	op := OrderProductModel{
		ProductId: 1,
		OrderId:   1,
	}

	assert.NotNil(t, op)
	assert.Equal(t, uint(1), op.ProductId)
	assert.Equal(t, uint(1), op.OrderId)
}

func TestOrderProductModel_ToFromDomain(t *testing.T) {
	modelOrderProduct := &OrderProductModel{
		ProductId: 1,
		OrderId:   1,
		Product:   ProductModel{},
		Order:     OrderModel{User: UserModel{}},
	}

	domainOrderProduct := modelOrderProduct.ToDomain()

	assert.Equal(t, modelOrderProduct.ProductId, domainOrderProduct.ProductId)
	assert.Equal(t, modelOrderProduct.OrderId, domainOrderProduct.OrderId)

	newModel := FromDomainOrderProduct(domainOrderProduct)

	assert.Equal(t, modelOrderProduct, newModel)
}
