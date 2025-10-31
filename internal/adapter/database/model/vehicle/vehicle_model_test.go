package vehicle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestVehicleModelInitialization(t *testing.T) {
	v := Model{
		Plate: "ABC-1234",
		Name:  "Gol",
		Year:  2023,
		Brand: "Volkswagen",
		Color: "Preto",
	}

	assert.NotNil(t, v)
	assert.Equal(t, "ABC-1234", v.Plate)
	assert.Equal(t, "Gol", v.Name)
	assert.Equal(t, 2023, v.Year)
	assert.Equal(t, "Volkswagen", v.Brand)
	assert.Equal(t, "Preto", v.Color)
}

func TestVehicleModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	modelVehicle := &Model{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Valid: true, Time: now},
		},
		Plate: "ABC1234",
		Name:  "Gol",
		Year:  2023,
		Brand: "Volkswagen",
		Color: "Preto",
	}

	domainVehicle := modelVehicle.ToDomain()

	assert.Equal(t, modelVehicle.ID, domainVehicle.ID)
	assert.Equal(t, modelVehicle.Plate, domainVehicle.Plate.GetPlate())
	assert.Equal(t, modelVehicle.Name, domainVehicle.Name)
	assert.Equal(t, modelVehicle.Year, domainVehicle.Year)
	assert.Equal(t, modelVehicle.Brand, domainVehicle.Brand)
	assert.Equal(t, modelVehicle.Color, domainVehicle.Color)
	assert.Equal(t, modelVehicle.CreatedAt, domainVehicle.CreatedAt)
	assert.Equal(t, modelVehicle.UpdatedAt, domainVehicle.UpdatedAt)
	assert.Equal(t, modelVehicle.DeletedAt.Time, *domainVehicle.DeletedAt)

}
