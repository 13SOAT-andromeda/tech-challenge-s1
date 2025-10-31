package customer_vehicle

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	CustomerId uint
	VehicleId  uint

	Vehicle  vehicle.Model  `gorm:"foreignKey:VehicleId;references:ID"`
	Customer customer.Model `gorm:"foreignKey:CustomerId;references:ID"`
}

func (*Model) TableName() string {
	return "Customer_Vehicle"
}

func (m *Model) ToDomain() *domain.CustomerVehicle {
	var deletedAt *time.Time

	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &domain.CustomerVehicle{
		ID:         m.ID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  deletedAt,
		CustomerId: m.CustomerId,
		VehicleId:  m.VehicleId,
	}
}

func (m *Model) FromDomain(d *domain.CustomerVehicle) {
	var deletedAt gorm.DeletedAt

	if d.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{
			Time:  *d.DeletedAt,
			Valid: true,
		}
	}

	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	m.DeletedAt = deletedAt
	m.CustomerId = d.CustomerId
	m.VehicleId = d.VehicleId
}
