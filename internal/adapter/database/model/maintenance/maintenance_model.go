package maintenance

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance_category"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name         string `gorm:"unique; not null"`
	DefaultPrice *float64
	CategoryId   uint
	Number       string `gorm:"not null"`

	ServiceCategory *maintenance_category.Model `gorm:"foreignKey:CategoryId;references:ID"`
}

func (Model) TableName() string {
	return "Service"
}

func (m *Model) ToDomain() *domain.Maintenance {
	var zeroServiceCategory maintenance_category.Model
	if m == nil {
		return nil
	}
	if m.ServiceCategory == nil {
		m.ServiceCategory = &zeroServiceCategory
	}

	return &domain.Maintenance{
		ID:                  m.ID,
		Name:                m.Name,
		DefaultPrice:        m.DefaultPrice,
		CategoryId:          m.CategoryId,
		Number:              m.Number,
		MaintenanceCategory: *m.ServiceCategory.ToDomain(),
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func (m *Model) FromDomain(d *domain.Maintenance) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	m.Name = d.Name
	m.DefaultPrice = d.DefaultPrice
	m.CategoryId = d.CategoryId
	m.Number = d.Number

	if m.ServiceCategory == nil {
		m.ServiceCategory = &maintenance_category.Model{}
	}

	m.ServiceCategory.FromDomain(&d.MaintenanceCategory)
}
