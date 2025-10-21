package order

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	customerVehicle "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
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
	User              model.Model           `gorm:"foreignKey:UserId;references:ID"`
	CustomerVehicle   customerVehicle.Model `gorm:"foreignKey:CustomerVehicleId;references:ID"`
	Company           company.Model         `gorm:"foreignKey:CompanyId;references:ID"`
}

func (*Model) TableName() string {
	return "Order"
}

func (m *Model) ToDomain() *domain.Order {
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
		User:              *m.User.ToDomain(),
		CustomerVehicle:   *m.CustomerVehicle.ToDomain(),
		Company:           *m.Company.ToDomain(),
	}
}

func (m *Model) FromDomain(d *domain.Order) {

	m.ID = d.ID
	m.DateIn = d.DateIn
	m.DateOut = d.DateOut
	m.Number = d.Number
	m.Status = d.Status
	m.VehicleKilometers = d.VehicleKilometers
	m.Note = d.Note
	m.DiagnosticNote = d.DiagnosticNote
	m.Price = d.Price
	m.UserId = d.UserId
	m.CustomerVehicleId = d.CustomerVehicleId
	m.CompanyId = d.CompanyId
	m.User.FromDomain(&d.User)
	m.CustomerVehicle.FromDomain(&d.CustomerVehicle)
	m.Company.FromDomain(&d.Company)
}
