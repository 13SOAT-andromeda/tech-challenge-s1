package order_service

import (
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/service"
	"github.com/stretchr/testify/assert"
)

func TestOrderServiceModelInitialization(t *testing.T) {
	os := Model{
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
	modelOrderService := Model{
		ServiceId: 1,
		OrderId:   1,
		Price:     50.0,
		Service:   service.ServiceModel{},
		Order:     order.Model{User: model.Model{Sessions: []model.SessionModel{}}},
	}

	domainOrderService := modelOrderService.ToDomain()

	assert.Equal(t, modelOrderService.ServiceId, domainOrderService.ServiceId)
	assert.Equal(t, modelOrderService.OrderId, domainOrderService.OrderId)
	assert.Equal(t, modelOrderService.Price, domainOrderService.Price)

}
