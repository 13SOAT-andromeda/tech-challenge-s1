package vehicle

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Plate string `gorm:"unique; not null"`
	Name  string `gorm:"not null"`
	Year  int    `gorm:"not null"`
	Brand string `gorm:"not null"`
	Color string `gorm:"not null"`
}

func (Model) TableName() string {
	return "Vehicle"
}

func (m *Model) ToDomain() *domain.Vehicle {
	if m == nil {
		return nil
	}

	p, _ := domain.NewPlate(m.Plate)

	var deletedAt *time.Time

	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &domain.Vehicle{
		ID:        m.ID,
		Plate:     p,
		Name:      m.Name,
		Year:      m.Year,
		Brand:     m.Brand,
		Color:     m.Color,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *Model) FromDomain(d *domain.Vehicle) {
	if d == nil || d.Plate == nil {
		return
	}

	m.ID = d.ID
	m.Plate = d.Plate.GetPlate()
	m.Name = d.Name
	m.Year = d.Year
	m.Brand = d.Brand
	m.Color = d.Color
}
