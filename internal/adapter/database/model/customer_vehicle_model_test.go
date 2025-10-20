package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCustomerVehicleModelInitialization(t *testing.T) {
	cv := CustomerVehicleModel{
		CustomerId: 1,
		VehicleId:  1,
	}

	assert.NotNil(t, cv)
	assert.Equal(t, uint(1), cv.CustomerId)
	assert.Equal(t, uint(1), cv.VehicleId)
}

func TestCustomerVehicleModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	modelCustomerVehicle := &CustomerVehicleModel{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		CustomerId: 1,
		VehicleId:  1,
	}

	domainCustomerVehicle := modelCustomerVehicle.ToDomain()

	assert.Equal(t, modelCustomerVehicle.ID, domainCustomerVehicle.ID)
	assert.Equal(t, modelCustomerVehicle.CustomerId, domainCustomerVehicle.CustomerId)
	assert.Equal(t, modelCustomerVehicle.VehicleId, domainCustomerVehicle.VehicleId)
	assert.Equal(t, modelCustomerVehicle.CreatedAt, domainCustomerVehicle.CreatedAt)
	assert.Equal(t, modelCustomerVehicle.UpdatedAt, domainCustomerVehicle.UpdatedAt)
	assert.Equal(t, modelCustomerVehicle.DeletedAt.Time, *domainCustomerVehicle.DeletedAt)

	newModel := FromDomainCustomerVehicle(domainCustomerVehicle)

	assert.Equal(t, modelCustomerVehicle, newModel)
}
