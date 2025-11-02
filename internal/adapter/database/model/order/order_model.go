package order

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	customerVehicle "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	DateIn            time.Time `gorm:"not null"`
	DateOut           *time.Time
	Status            string
	VehicleKilometers int
	Note              *string
	DiagnosticNote    *string
	Price             *float64
	UserID            uint
	CustomerVehicleID uint
	CompanyID         uint
	User              user.Model            `gorm:"foreignKey:UserID;references:ID"`
	CustomerVehicle   customerVehicle.Model `gorm:"foreignKey:CustomerVehicleID;references:ID"`
	Company           company.Model         `gorm:"foreignKey:CompanyID;references:ID"`
}

func (*Model) TableName() string {
	return "Orders"
}

func (m *Model) ToDomain() *domain.Order {
	return &domain.Order{
		ID:                m.ID,
		DateIn:            m.DateIn,
		DateOut:           m.DateOut,
		Status:            domain.OrderStatus(m.Status),
		VehicleKilometers: m.VehicleKilometers,
		Note:              m.Note,
		DiagnosticNote:    m.DiagnosticNote,
		Price:             m.Price,
		User:              *m.User.ToDomain(),
		CustomerVehicle:   *m.CustomerVehicle.ToDomain(),
		Company:           *m.Company.ToDomain(),
	}
}

func (m *Model) FromDomain(d *domain.Order) {

	m.ID = d.ID
	m.DateIn = d.DateIn
	m.DateOut = d.DateOut
	m.Status = string(d.Status)
	m.VehicleKilometers = d.VehicleKilometers
	m.Note = d.Note
	m.DiagnosticNote = d.DiagnosticNote
	m.Price = d.Price
	m.UserID = d.User.ID
	m.CustomerVehicleID = d.CustomerVehicle.ID
	m.CompanyID = d.Company.ID
	m.User = user.Model{
		ID: d.User.ID,
	}
	m.CustomerVehicle.FromDomain(&d.CustomerVehicle)
	m.Company.FromDomain(&d.Company)
}
