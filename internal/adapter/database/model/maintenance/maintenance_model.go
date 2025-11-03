package maintenance

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name       string `gorm:"not null"`
	Price      int64  `gorm:"not null"`
	CategoryID string `gorm:"not null"`
}

func (*Model) TableName() string {
	return "Maintenance"
}

func (m *Model) ToDomain() *domain.Maintenance {

	return &domain.Maintenance{
		ID:         m.ID,
		Name:       m.Name,
		Price:      m.Price,
		CategoryID: domain.MaintenanceCategory(m.CategoryID),
	}
}

func (m *Model) FromDomain(d *domain.Maintenance) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Name = d.Name
	m.Price = d.Price
	m.CategoryID = string(d.CategoryID)
}
