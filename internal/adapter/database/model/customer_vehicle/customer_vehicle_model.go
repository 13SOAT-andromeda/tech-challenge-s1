package customer_vehicle

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	CustomerId uint
	VehicleId  uint

	Vehicle  model.VehicleModel `gorm:"foreignKey:VehicleId;references:ID"`
	Customer customer.Model     `gorm:"foreignKey:CustomerId;references:ID"`
}

func (Model) TableName() string {
	return "Customer_Vehicle"
}

func ToDomain(m Model) domain.CustomerVehicle {
	var deletedAt *time.Time

	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return domain.CustomerVehicle{
		ID:         m.ID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  deletedAt,
		CustomerId: m.CustomerId,
		VehicleId:  m.VehicleId,
	}
}

func FromDomain(domainEntity domain.CustomerVehicle) Model {
	var deletedAt gorm.DeletedAt

	if domainEntity.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{
			Time:  *domainEntity.DeletedAt,
			Valid: true,
		}
	}

	return Model{
		Model: gorm.Model{
			ID:        domainEntity.ID,
			CreatedAt: domainEntity.CreatedAt,
			UpdatedAt: domainEntity.UpdatedAt,
			DeletedAt: deletedAt,
		},
		CustomerId: domainEntity.CustomerId,
		VehicleId:  domainEntity.VehicleId,
	}
}
