package customer_vehicle

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	VehicleID  uint
	CustomerID uint

	Vehicle  vehicle.Model  `gorm:"foreignKey:VehicleID;references:ID"`
	Customer customer.Model `gorm:"foreignKey:CustomerID;references:ID"`
}

func (*Model) TableName() string {
	return "Customer_Vehicle"
}

func (m *Model) ToDomain() *domain.CustomerVehicle {
	return &domain.CustomerVehicle{
		ID:         m.ID,
		CustomerId: m.CustomerID,
		VehicleId:  m.VehicleID,
		Vehicle:    *m.Vehicle.ToDomain(),
		Customer:   *m.Customer.ToDomain(),
	}
}

func (m *Model) FromDomain(d *domain.CustomerVehicle) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.CustomerID = d.CustomerId
	m.VehicleID = d.VehicleId

	if d.Vehicle.ID != 0 {
		m.Vehicle.FromDomain(&d.Vehicle)
	}
	if d.Customer.ID != 0 {
		m.Customer.FromDomain(&d.Customer)
	}
}
