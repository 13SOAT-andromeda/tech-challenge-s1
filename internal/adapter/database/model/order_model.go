package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"time"
)

type OrderModel struct {
	ID                uint      `gorm:"primaryKey"`
	DateIn            time.Time `gorm:"not null"`
	DateOut           *time.Time
	Number            string `gorm:"not null; unique"`
	Status            string
	VehicleKilometers *float64
	Note              *string
	DiagnosticNote    *string
	Price             float64 `gorm:"not null"`
	UserId            uint
	CustomerVehicleId uint
	CompanyId         uint

	User            UserModel            `gorm:"foreignKey:UserId;references:ID"`
	CustomerVehicle CustomerVehicleModel `gorm:"foreignKey:CustomerVehicleId;references:ID"`
	Company         CompanyModel         `gorm:"foreignKey:CompanyId;references:ID"`
}

func (OrderModel) TableName() string {
	return "Order"
}

func (m *OrderModel) ToDomain() *domain.Order {
	if m == nil {
		return nil
	}
	return &domain.Order{
		ID:                m.ID,
		DateIn:            m.DateIn,
		DateOut:           m.DateOut,
		Number:            m.Number,
		Status:            m.Status,
		VehicleKilometers: m.VehicleKilometers,
		Note:              m.Note,
		DiagnosticNote:    m.DiagnosticNote,
		Price:             m.Price,
		UserId:            m.UserId,
		CustomerVehicleId: m.CustomerVehicleId,
		CompanyId:         m.CompanyId,
		User:              *(&m.User).ToDomain(),
		CustomerVehicle:   *(&m.CustomerVehicle).ToDomain(),
		Company:           *(&m.Company).ToDomain(),
	}
}

func FromDomainOrder(d *domain.Order) *OrderModel {
	if d == nil {
		return nil
	}
	return &OrderModel{
		ID:                d.ID,
		DateIn:            d.DateIn,
		DateOut:           d.DateOut,
		Number:            d.Number,
		Status:            d.Status,
		VehicleKilometers: d.VehicleKilometers,
		Note:              d.Note,
		DiagnosticNote:    d.DiagnosticNote,
		Price:             d.Price,
		UserId:            d.UserId,
		CustomerVehicleId: d.CustomerVehicleId,
		CompanyId:         d.CompanyId,
		User:              *FromDomainUser(&d.User),
		CustomerVehicle:   *FromDomainCustomerVehicle(&d.CustomerVehicle),
		Company:           *FromDomainCompany(&d.Company),
	}
}
