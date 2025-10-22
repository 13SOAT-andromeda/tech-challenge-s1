package service_category

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name string `gorm:"not null"`
}

func (*Model) TableName() string {
	return "Service_Category"
}

func (m *Model) ToDomain() *domain.ServiceCategory {
	if m == nil {
		return nil
	}

	return &domain.ServiceCategory{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *Model) FromDomain(d *domain.ServiceCategory) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	m.Name = d.Name
}
