package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type ServiceModel struct {
	gorm.Model
	Name         string `gorm:"unique; not null"`
	DefaultPrice *float64
	CategoryId   uint
	Number       string `gorm:"not null"`

	ServiceCategory ServiceCategoryModel `gorm:"foreignKey:CategoryId;references:ID"`
}

func (ServiceModel) TableName() string {
	return "Service"
}

func (m *ServiceModel) ToDomain() *domain.Service {
	if m == nil {
		return nil
	}
	return &domain.Service{
		ID:              m.ID,
		Name:            m.Name,
		DefaultPrice:    m.DefaultPrice,
		CategoryId:      m.CategoryId,
		Number:          m.Number,
		ServiceCategory: *(&m.ServiceCategory).ToDomain(),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (m *ServiceModel) FromDomain(d *domain.Service) {
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
	m.ServiceCategory.FromDomain(&d.ServiceCategory)
}
