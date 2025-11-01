package product

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Stock uint   `gorm:"not null"`
	Price int64  `gorm:"not null"`
}

func (*Model) TableName() string {
	return "Product"
}

func (m *Model) ToDomain() *domain.Product {
	if m == nil {
		return nil
	}
	return &domain.Product{
		ID:    m.ID,
		Name:  m.Name,
		Stock: m.Stock,
		Price: m.Price,
	}
}

func (m *Model) FromDomain(d *domain.Product) {
	if d == nil {
		return
	}
	m.ID = d.ID
	m.Name = d.Name
	m.Stock = d.Stock
	m.Price = d.Price
}
