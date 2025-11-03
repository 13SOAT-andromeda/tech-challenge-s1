package order_maintenance

import (
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
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
		Maintenance: maintenance.Model{
			Model:      gorm.Model{ID: 1},
			Name:       "Change oil",
			Price:      1500,
			CategoryId: "padrao",
		},
		Order: order.Model{
			Model:             gorm.Model{ID: 2},
			DateIn:            now,
			DateOut:           nil,
			Status:            string(domain.RECEIVED),
			VehicleKilometers: 100,
			Note:              &note,
			DiagnosticNote:    &diag,
			Price:             &price,
			User:              user.Model{ID: 10, Name: "usr", Email: "e@x.com", Contact: "cont", Password: "hashed", Role: "role", Active: true},
			CustomerVehicle:   customer_vehicle.Model{Model: gorm.Model{ID: 20}, CustomerId: 5, VehicleId: 6},
			Company:           company.Model{Model: gorm.Model{ID: 30}, Name: "Comp", Contact: "123", Document: "doc"},
		},
	}

	d := m.ToDomain()

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
	if d.Order.User.ID != 10 {
		t.Fatalf("user id: expected 10 got %d", d.Order.User.ID)
	}
	if d.Order.CustomerVehicle.ID != 20 {
		t.Fatalf("customer vehicle id: expected 20 got %d", d.Order.CustomerVehicle.ID)
	}
	if d.Order.Company.ID != 30 {
		t.Fatalf("company id: expected 30 got %d", d.Order.Company.ID)
	}
}

func TestModel_FromDomain(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	price := 200.0

	d := &domain.OrderMaintenance{
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
			User:              domain.User{ID: 11},
			CustomerVehicle:   domain.CustomerVehicle{ID: 21, CustomerId: 7, VehicleId: 8},
			Company:           domain.Company{ID: 31},
		},
	}

	m := &Model{}
	m.FromDomain(d)

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

	if m.Order.ID != 4 {
		t.Fatalf("order id: expected 4 got %d", m.Order.ID)
	}
	if !m.Order.DateIn.Equal(now) {
		t.Fatalf("order datein: expected %v got %v", now, m.Order.DateIn)
	}
	if m.Order.Status != string(domain.APPROVED) {
		t.Fatalf("order status: expected %v got %v", domain.APPROVED, m.Order.Status)
	}
	if m.Order.VehicleKilometers != 200 {
		t.Fatalf("vehicle kilometers: expected 200 got %d", m.Order.VehicleKilometers)
	}
	if m.Order.Price == nil || *m.Order.Price != price {
		t.Fatalf("order price: expected %v got %v", price, m.Order.Price)
	}
	if m.Order.UserID != 11 {
		t.Fatalf("user id: expected 11 got %d", m.Order.UserID)
	}
	if m.Order.CustomerVehicle.ID != 21 {
		t.Fatalf("customer vehicle id: expected 21 got %d", m.Order.CustomerVehicle.ID)
	}
	if m.Order.CustomerVehicle.Customer.ID != 7 {
		t.Fatalf("customer id: expected 7 got %d", m.Order.CustomerVehicle.Customer.ID)
	}
	if m.Order.CustomerVehicle.Vehicle.ID != 8 {
		t.Fatalf("vehicle id: expected 8 got %d", m.Order.CustomerVehicle.Vehicle.ID)
	}
	if m.Order.CompanyID != 31 {
		t.Fatalf("company id: expected 31 got %d", m.Order.CompanyID)
	}
}

func TestModel_FromDomain_Nil(t *testing.T) {
	m := &Model{}
	m.FromDomain(nil)
	if m.Maintenance.ID != 0 || m.Order.ID != 0 {
		t.Fatalf("expected zero values after FromDomain(nil), got maintenance id %d order id %d", m.Maintenance.ID, m.Order.ID)
	}
}
