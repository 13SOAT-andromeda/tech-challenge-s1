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

func (m *ProductModel) FromDomain(d *domain.Product) {
	if d == nil {
		return
	}
	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	m.Name = d.Name
	m.Quantity = d.Quantity
	m.Price = d.Price
}
