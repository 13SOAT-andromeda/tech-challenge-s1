package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVehicleTableName(t *testing.T) {
	assert.Equal(t, "Vehicle", VehicleModel{}.TableName())
}

func TestNilVehicleToDomain(t *testing.T) {
	assert.Nil(t, (*VehicleModel)(nil).ToDomain())
}

func TestNilVehicleFromDomain(t *testing.T) {
	assert.Nil(t, FromDomainVehicle(nil))
}

func TestVehicleModelInitialization(t *testing.T) {
	v := VehicleModel{
		Plate: "ABC-1234",
		Model: "Gol",
		Year:  2023,
		Brand: "Volkswagen",
		Color: "Preto",
	}

	assert.NotNil(t, v)
	assert.Equal(t, "ABC-1234", v.Plate)
	assert.Equal(t, "Gol", v.Model)
	assert.Equal(t, 2023, v.Year)
	assert.Equal(t, "Volkswagen", v.Brand)
	assert.Equal(t, "Preto", v.Color)
}

func TestVehicleModel_ToFromDomain(t *testing.T) {
	modelVehicle := &VehicleModel{
		ID:    1,
		Plate: "ABC-1234",
		Model: "Gol",
		Year:  2023,
		Brand: "Volkswagen",
		Color: "Preto",
	}

	domainVehicle := modelVehicle.ToDomain()

	assert.Equal(t, modelVehicle.ID, domainVehicle.ID)
	assert.Equal(t, modelVehicle.Plate, domainVehicle.Plate)
	assert.Equal(t, modelVehicle.Model, domainVehicle.Model)
	assert.Equal(t, modelVehicle.Year, domainVehicle.Year)
	assert.Equal(t, modelVehicle.Brand, domainVehicle.Brand)
	assert.Equal(t, modelVehicle.Color, domainVehicle.Color)

	newModel := FromDomainVehicle(domainVehicle)

	assert.Equal(t, modelVehicle, newModel)
}
