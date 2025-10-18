package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type ProductModel struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Quantity uint   `gorm:"not null"`
	Price    uint32 `gorm:"not null"` // em centavos
}

func (ProductModel) TableName() string {
	return "Product"
}

func (m *ProductModel) ToDomain() *domain.Product {
	if m == nil {
		return nil
	}
	return &domain.Product{
		ID:        m.ID,
		Name:      m.Name,
		Quantity:  m.Quantity,
		Price:     m.Price,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromDomainProduct(d *domain.Product) *ProductModel {
	if d == nil {
		return nil
	}
	return &ProductModel{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name:     d.Name,
		Quantity: d.Quantity,
		Price:    d.Price,
	}
}
