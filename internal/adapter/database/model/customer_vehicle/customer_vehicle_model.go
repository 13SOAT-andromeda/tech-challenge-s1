package customer_vehicle

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	VehicleId  uint
	CustomerId uint

	Vehicle  vehicle.Model  `gorm:"foreignKey:VehicleId;references:ID"`
	Customer customer.Model `gorm:"foreignKey:CustomerId;references:ID"`
}

func (*Model) TableName() string {
	return "Customer_Vehicle"
}

func (m *Model) ToDomain() *domain.CustomerVehicle {
	return &domain.CustomerVehicle{
		ID:         m.ID,
		CustomerId: m.CustomerId,
		VehicleId:  m.VehicleId,
		Vehicle:    *m.Vehicle.ToDomain(),
		Customer:   *m.Customer.ToDomain(),
	}
}

func (m *Model) FromDomain(d *domain.CustomerVehicle) {

	m.ID = d.ID
	m.Customer.ID = d.CustomerId
	m.Vehicle.ID = d.VehicleId
	m.Vehicle.FromDomain(&d.Vehicle)
	m.Customer.FromDomain(&d.Customer)
}
