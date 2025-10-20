package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderModelInitialization(t *testing.T) {
	now := time.Now()
	vehicleKm := 10000.5
	note := "Test note"
	diagnosticNote := "Test diagnostic note"

	o := OrderModel{
		DateIn:            now,
		DateOut:           &now,
		Number:            "12345",
		Status:            "pending",
		VehicleKilometers: &vehicleKm,
		Note:              &note,
		DiagnosticNote:    &diagnosticNote,
		Price:             150.75,
		UserId:            1,
		CustomerVehicleId: 1,
		CompanyId:         1,
	}

	assert.NotNil(t, o)
	assert.Equal(t, now, o.DateIn)
	assert.Equal(t, &now, o.DateOut)
	assert.Equal(t, "12345", o.Number)
	assert.Equal(t, "pending", o.Status)
	assert.Equal(t, &vehicleKm, o.VehicleKilometers)
	assert.Equal(t, &note, o.Note)
	assert.Equal(t, &diagnosticNote, o.DiagnosticNote)
	assert.Equal(t, 150.75, o.Price)
	assert.Equal(t, uint(1), o.UserId)
	assert.Equal(t, uint(1), o.CustomerVehicleId)
	assert.Equal(t, uint(1), o.CompanyId)
}

func TestOrderModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	vehicleKm := 10000.5
	note := "Test note"
	diagnosticNote := "Test diagnostic note"

	modelOrder := &OrderModel{
		DateIn:            now,
		DateOut:           &now,
		Number:            "12345",
		Status:            "pending",
		VehicleKilometers: &vehicleKm,
		Note:              &note,
		DiagnosticNote:    &diagnosticNote,
		Price:             150.75,
		UserId:            1,
		CustomerVehicleId: 1,
		CompanyId:         1,
		User:              UserModel{Sessions: []SessionModel{}},
		CustomerVehicle:   CustomerVehicleModel{},
		Company:           CompanyModel{},
	}

	domainOrder := modelOrder.ToDomain()

	assert.Equal(t, modelOrder.ID, domainOrder.ID)
	assert.Equal(t, modelOrder.DateIn, domainOrder.DateIn)
	assert.Equal(t, modelOrder.DateOut, domainOrder.DateOut)
	assert.Equal(t, modelOrder.Number, domainOrder.Number)
	assert.Equal(t, modelOrder.Status, domainOrder.Status)
	assert.Equal(t, modelOrder.VehicleKilometers, domainOrder.VehicleKilometers)
	assert.Equal(t, modelOrder.Note, domainOrder.Note)
	assert.Equal(t, modelOrder.DiagnosticNote, domainOrder.DiagnosticNote)
	assert.Equal(t, modelOrder.Price, domainOrder.Price)
	assert.Equal(t, modelOrder.UserId, domainOrder.UserId)
	assert.Equal(t, modelOrder.CustomerVehicleId, domainOrder.CustomerVehicleId)
	assert.Equal(t, modelOrder.CompanyId, domainOrder.CompanyId)

	newModel := FromDomainOrder(domainOrder)

	assert.Equal(t, modelOrder, newModel)
}
