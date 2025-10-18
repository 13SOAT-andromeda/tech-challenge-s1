package model

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type CustomerVehicleModel struct {
	gorm.Model
	CustomerId uint
	VehicleId  uint

	Vehicle  VehicleModel  `gorm:"foreignKey:VehicleId;references:ID"`
	Customer CustomerModel `gorm:"foreignKey:CustomerId;references:ID"`
}

func (CustomerVehicleModel) TableName() string {
	return "Customer_Vehicle"
}

func (m *CustomerVehicleModel) ToDomain() *domain.CustomerVehicle {
	if m == nil {
		return nil
	}

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

func FromDomainCustomerVehicle(domainEntity *domain.CustomerVehicle) *CustomerVehicleModel {
	if domainEntity == nil {
		return nil
	}

	var deletedAt gorm.DeletedAt
	if domainEntity.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{
			Time:  *domainEntity.DeletedAt,
			Valid: true,
		}
	}

	return &CustomerVehicleModel{
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
