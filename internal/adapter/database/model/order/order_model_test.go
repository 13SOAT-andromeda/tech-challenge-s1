package order

import (
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/stretchr/testify/assert"
)

func TestOrderModelInitialization(t *testing.T) {
	now := time.Now()
	vehicleKm := 10000
	note := "Test note"
	diagnosticNote := "Test diagnostic note"
	price := 150.75

	o := Model{
		DateIn:            now,
		DateOut:           &now,
		Status:            "pending",
		VehicleKilometers: vehicleKm,
		Note:              &note,
		DiagnosticNote:    &diagnosticNote,
		Price:             &price,
	}

	assert.NotNil(t, o)
	assert.Equal(t, now, o.DateIn)
	assert.Equal(t, &now, o.DateOut)
	assert.Equal(t, "pending", o.Status)
	assert.Equal(t, vehicleKm, o.VehicleKilometers)
	assert.Equal(t, &note, o.Note)
	assert.Equal(t, &diagnosticNote, o.DiagnosticNote)
	assert.Equal(t, &price, o.Price)
}

func TestOrderModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	vehicleKm := 10000
	note := "Test note"
	diagnosticNote := "Test diagnostic note"
	price := 150.75

	modelOrder := Model{
		DateIn:            now,
		DateOut:           &now,
		Status:            "Recebida",
		VehicleKilometers: vehicleKm,
		Note:              &note,
		DiagnosticNote:    &diagnosticNote,
		Price:             &price,
		User:              user.Model{},
		CustomerVehicle:   customer_vehicle.Model{},
		Company:           company.Model{},
	}

	domainOrder := modelOrder.ToDomain()

	assert.Equal(t, modelOrder.ID, domainOrder.ID)
	assert.Equal(t, modelOrder.DateIn, domainOrder.DateIn)
	assert.Equal(t, modelOrder.DateOut, domainOrder.DateOut)

	assert.NotNil(t, domainOrder.Status)
	assert.Equal(t, modelOrder.Status, string(domainOrder.Status))

	assert.Equal(t, modelOrder.VehicleKilometers, domainOrder.VehicleKilometers)
	assert.Equal(t, modelOrder.Note, domainOrder.Note)
	assert.Equal(t, modelOrder.DiagnosticNote, domainOrder.DiagnosticNote)
	assert.Equal(t, modelOrder.Price, domainOrder.Price)
}
