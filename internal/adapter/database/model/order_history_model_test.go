package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderHistoryModelInitialization(t *testing.T) {
	now := time.Now()
	oh := OrderHistoryModel{
		OrderId: 1,
		Date:    now,
		Status:  "received",
	}

	assert.NotNil(t, oh)
	assert.Equal(t, uint64(1), oh.OrderId)
	assert.Equal(t, now, oh.Date)
	assert.Equal(t, "received", oh.Status)
}

func TestOrderHistoryModel_ToFromDomain(t *testing.T) {
	now := time.Now()

	modelOrderHistory := &OrderHistoryModel{
		OrderId: 1,
		Date:    now,
		Status:  "received",
	}

	domainOrderHistory := modelOrderHistory.ToDomain()

	assert.Equal(t, modelOrderHistory.OrderId, domainOrderHistory.OrderId)
	assert.Equal(t, modelOrderHistory.Date, domainOrderHistory.Date)
	assert.Equal(t, modelOrderHistory.Status, domainOrderHistory.Status)

	newModel := FromDomainOrderHistory(domainOrderHistory)

	assert.Equal(t, modelOrderHistory, newModel)
}
