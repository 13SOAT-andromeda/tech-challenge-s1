package order_maintenance

import (
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

func TestTableName(t *testing.T) {
	var m Model
	if m.TableName() != "Order_Maintenance" {
		t.Fatalf("unexpected table name: %s", m.TableName())
	}
}

func TestModel_ToDomain(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	price := 123.45
	note := "note"
	diag := "diag"

	m := &Model{
		MaintenanceId: 1,
		OrderId:       2,
		Maintenance: maintenance.Model{
			Model:      gorm.Model{ID: 1},
			Name:       "Change oil",
			Price:      1500,
			CategoryId: "padrao",
		},
	}

	orderDomain := &domain.Order{
		ID:                2,
		DateIn:            now,
		DateOut:           nil,
		Status:            domain.RECEIVED,
		VehicleKilometers: 100,
		Note:              &note,
		DiagnosticNote:    &diag,
		Price:             &price,
		UserID:            10,
		CustomerVehicleID: 20,
		CompanyID:         30,
	}

	d := m.ToDomain(orderDomain)

	// Maintenance assertions
	if d.Maintenance.ID != 1 {
		t.Fatalf("maintenance id: expected 1 got %d", d.Maintenance.ID)
	}
	if d.Maintenance.Name != "Change oil" {
		t.Fatalf("maintenance name: expected %s got %s", "Change oil", d.Maintenance.Name)
	}
	if d.Maintenance.Price != int64(1500) {
		t.Fatalf("maintenance price: expected %d got %d", 1500, d.Maintenance.Price)
	}
	if string(d.Maintenance.CategoryId) != "padrao" {
		t.Fatalf("maintenance category: expected %s got %s", "padrao", string(d.Maintenance.CategoryId))
	}

	// Order assertions
	if d.Order.ID != 2 {
		t.Fatalf("order id: expected 2 got %d", d.Order.ID)
	}
	if !d.Order.DateIn.Equal(now) {
		t.Fatalf("order datein: expected %v got %v", now, d.Order.DateIn)
	}
	if d.Order.Status != domain.RECEIVED {
		t.Fatalf("order status: expected %v got %v", domain.RECEIVED, d.Order.Status)
	}
	if d.Order.VehicleKilometers != 100 {
		t.Fatalf("vehicle kilometers: expected 100 got %d", d.Order.VehicleKilometers)
	}
	if d.Order.Note == nil || *d.Order.Note != note {
		t.Fatalf("order note: expected %s got %v", note, d.Order.Note)
	}
	if d.Order.DiagnosticNote == nil || *d.Order.DiagnosticNote != diag {
		t.Fatalf("diagnostic note: expected %s got %v", diag, d.Order.DiagnosticNote)
	}
	if d.Order.Price == nil || *d.Order.Price != price {
		t.Fatalf("order price: expected %v got %v", price, d.Order.Price)
	}
	if d.Order.UserID != 10 {
		t.Fatalf("user id: expected 10 got %d", d.Order.UserID)
	}
	if d.Order.CustomerVehicleID != 20 {
		t.Fatalf("customer vehicle id: expected 20 got %d", d.Order.CustomerVehicleID)
	}
	if d.Order.CompanyID != 30 {
		t.Fatalf("company id: expected 30 got %d", d.Order.CompanyID)
	}
}

func TestModel_FromDomain(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	price := 200.0

	d := &domain.OrderMaintenance{
		MaintenanceId: 3,
		OrderId:       4,
		Maintenance: domain.Maintenance{
			ID:         3,
			Name:       "Filter",
			Price:      2000,
			CategoryId: domain.MaintenanceCategory("utilitario"),
		},
		Order: domain.Order{
			ID:                4,
			DateIn:            now,
			Status:            domain.APPROVED,
			VehicleKilometers: 200,
			Note:              nil,
			DiagnosticNote:    nil,
			Price:             &price,
			UserID:            11,
			CustomerVehicleID: 21,
			CompanyID:         31,
		},
	}

	m := &Model{}
	m.FromDomain(d)

	if m.MaintenanceId != 3 {
		t.Fatalf("maintenance id: expected 3 got %d", m.MaintenanceId)
	}
	if m.OrderId != 4 {
		t.Fatalf("order id: expected 4 got %d", m.OrderId)
	}
	if m.Maintenance.ID != 3 {
		t.Fatalf("maintenance id: expected 3 got %d", m.Maintenance.ID)
	}
	if m.Maintenance.Name != "Filter" {
		t.Fatalf("maintenance name: expected %s got %s", "Filter", m.Maintenance.Name)
	}
	if m.Maintenance.Price != 2000 {
		t.Fatalf("maintenance price: expected %d got %d", 2000, m.Maintenance.Price)
	}
	if m.Maintenance.CategoryId != "utilitario" {
		t.Fatalf("maintenance category: expected %s got %s", "utilitario", m.Maintenance.CategoryId)
	}
}

func TestModel_FromDomain_Nil(t *testing.T) {
	m := &Model{}
	m.FromDomain(nil)
	if m.Maintenance.ID != 0 || m.MaintenanceId != 0 || m.OrderId != 0 {
		t.Fatalf("expected zero values after FromDomain(nil), got maintenance id %d maintenanceId %d orderId %d", m.Maintenance.ID, m.MaintenanceId, m.OrderId)
	}
}
