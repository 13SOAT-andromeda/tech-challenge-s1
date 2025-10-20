package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderServiceModelInitialization(t *testing.T) {
	os := OrderServiceModel{
		ServiceId: 1,
		OrderId:   1,
		Price:     50.0,
	}

	assert.NotNil(t, os)
	assert.Equal(t, uint(1), os.ServiceId)
	assert.Equal(t, uint(1), os.OrderId)
	assert.Equal(t, 50.0, os.Price)
}

func TestOrderServiceModel_ToFromDomain(t *testing.T) {
	modelOrderService := &OrderServiceModel{
		ServiceId: 1,
		OrderId:   1,
		Price:     50.0,
		Service:   ServiceModel{},
		Order:     OrderModel{User: UserModel{}},
	}

	domainOrderService := modelOrderService.ToDomain()

	assert.Equal(t, modelOrderService.ServiceId, domainOrderService.ServiceId)
	assert.Equal(t, modelOrderService.OrderId, domainOrderService.OrderId)
	assert.Equal(t, modelOrderService.Price, domainOrderService.Price)

	newModel := FromDomainOrderService(domainOrderService)

	assert.Equal(t, modelOrderService, newModel)
}
