package customer_vehicle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCustomerVehicleModelInitialization(t *testing.T) {
	cv := Model{
		CustomerID: 1,
		VehicleID:  1,
	}

	assert.NotNil(t, cv)
	assert.Equal(t, uint(1), cv.CustomerID)
	assert.Equal(t, uint(1), cv.VehicleID)
}

func TestCustomerVehicleModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	modelCustomerVehicle := Model{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		CustomerID: 1,
		VehicleID:  1,
	}

	domainCustomerVehicle := modelCustomerVehicle.ToDomain()

	assert.Equal(t, modelCustomerVehicle.ID, domainCustomerVehicle.ID)
	assert.Equal(t, modelCustomerVehicle.CustomerID, domainCustomerVehicle.CustomerId)
	assert.Equal(t, modelCustomerVehicle.VehicleID, domainCustomerVehicle.VehicleId)
}
