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

func FromDomainService(d *domain.Service) *ServiceModel {
	if d == nil {
		return nil
	}
	return &ServiceModel{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name:            d.Name,
		DefaultPrice:    d.DefaultPrice,
		CategoryId:      d.CategoryId,
		Number:          d.Number,
		ServiceCategory: *FromDomainServiceCategory(&d.ServiceCategory),
	}
}
