package order

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
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
	User              model.UserModel                       `gorm:"foreignKey:UserId;references:ID"`
	CustomerVehicle   customer_vehicle.CustomerVehicleModel `gorm:"foreignKey:CustomerVehicleId;references:ID"`
	Company           company.Model                         `gorm:"foreignKey:CompanyId;references:ID"`
}

func (Model) TableName() string {
	return "Order"
}

func ToDomain(m Model) domain.Order {
	return domain.Order{
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
		Company:           company.ToDomain(m.Company),
	}
}

func FromDomain(d domain.Order) Model {
	return Model{
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
		User:              *model.FromDomainUser(&d.User),
		CustomerVehicle:   *customer_vehicle.FromDomainCustomerVehicle(&d.CustomerVehicle),
		Company:           company.FromDomain(d.Company),
	}
}
