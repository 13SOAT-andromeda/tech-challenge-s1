package maintenance_category

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name string `gorm:"not null"`
}

func (*Model) TableName() string {
	return "Maintenance_Category"
}

func (m *Model) ToDomain() *domain.MaintenanceCategory {
	if m == nil {
		return nil
	}

	return &domain.MaintenanceCategory{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *Model) FromDomain(d *domain.MaintenanceCategory) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	m.Name = d.Name
}
