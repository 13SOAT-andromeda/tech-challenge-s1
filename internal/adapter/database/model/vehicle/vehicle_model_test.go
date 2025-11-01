package vehicle

import (
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestVehicleModel_TableName(t *testing.T) {
	v := &Model{}
	assert.Equal(t, "Vehicle", v.TableName())
}

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

func TestVehicleModel_ToDomain(t *testing.T) {
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

func TestVehicleModel_ToDomainNil(t *testing.T) {
	var modelVehicle *Model

	domainVehicle := modelVehicle.ToDomain()

	assert.Nil(t, domainVehicle)

}

func TestVehicleModel_FromDomain(t *testing.T) {
	plate, _ := domain.NewPlate("ABC1234")
	domainVehicle := &domain.Vehicle{
		ID:    1,
		Plate: plate,
		Name:  "Gol",
		Year:  2023,
		Brand: "Volkswagen",
		Color: "Preto",
	}

	modelVehicle := &Model{}
	modelVehicle.FromDomain(domainVehicle)

	assert.Equal(t, domainVehicle.ID, modelVehicle.ID)
	assert.Equal(t, domainVehicle.Plate.GetPlate(), modelVehicle.Plate)
	assert.Equal(t, domainVehicle.Name, modelVehicle.Name)
	assert.Equal(t, domainVehicle.Year, modelVehicle.Year)
	assert.Equal(t, domainVehicle.Brand, modelVehicle.Brand)
	assert.Equal(t, domainVehicle.Color, modelVehicle.Color)

}
