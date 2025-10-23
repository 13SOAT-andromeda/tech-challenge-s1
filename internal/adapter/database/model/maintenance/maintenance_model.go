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

	MaintenanceCategory *maintenance_category.Model `gorm:"foreignKey:CategoryId;references:ID"`
}

func (Model) TableName() string {
	return "Maintenance"
}

func (m *Model) ToDomain() *domain.Maintenance {
	var maintenanceCategoryDomain *domain.MaintenanceCategory
	if m.MaintenanceCategory != nil {
		maintenanceCategoryDomain = m.MaintenanceCategory.ToDomain()
	} else {
		maintenanceCategoryDomain = nil
	}

	return &domain.Maintenance{
		ID:                  m.ID,
		Name:                m.Name,
		DefaultPrice:        m.DefaultPrice,
		Number:              m.Number,
		MaintenanceCategory: maintenanceCategoryDomain,
	}
}

func (m *Model) FromDomain(d *domain.Maintenance) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Name = d.Name
	m.DefaultPrice = d.DefaultPrice
	m.Number = d.Number

	if m.MaintenanceCategory == nil {
		m.MaintenanceCategory = &maintenance_category.Model{}
	}

	m.MaintenanceCategory.FromDomain(d.MaintenanceCategory)
}
