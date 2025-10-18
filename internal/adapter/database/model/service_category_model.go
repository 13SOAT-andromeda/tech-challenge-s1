package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type ServiceCategoryModel struct {
	gorm.Model
	Name string `gorm:"not null"`
}

func (ServiceCategoryModel) TableName() string {
	return "Service_Category"
}

func (m *ServiceCategoryModel) ToDomain() *domain.ServiceCategory {
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

func FromDomainServiceCategory(d *domain.ServiceCategory) *ServiceCategoryModel {
	if d == nil {
		return nil
	}
	return &ServiceCategoryModel{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name: d.Name,
	}
}
