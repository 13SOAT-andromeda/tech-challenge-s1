package model

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type VehicleModel struct {
	ID    uint   `gorm:"primaryKey"`
	Plate string `gorm:"unique; not null"`
	Model string `gorm:"not null"`
	Year  int    `gorm:"not null"`
	Brand string `gorm:"not null"`
	Color string `gorm:"not null"`
}

func (VehicleModel) TableName() string {
	return "Vehicle"
}

func (m *VehicleModel) ToDomain() *domain.VehicleModel {
	if m == nil {
		return nil
	}
	return &domain.VehicleModel{
		ID:    m.ID,
		Plate: m.Plate,
		Model: m.Model,
		Year:  m.Year,
		Brand: m.Brand,
		Color: m.Color,
	}
}

func (m *VehicleModel) FromDomain(d *domain.VehicleModel) {
	if d == nil {
		return
	}
	m.ID = d.ID
	m.Plate = d.Plate
	m.Model = d.Model
	m.Year = d.Year
	m.Brand = d.Brand
	m.Color = d.Color
}
