package order_product

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
	Quantity  uint          `gorm:"not null"`
	ProductId uint          `gorm:"primaryKey"`
	OrderId   uint          `gorm:"primaryKey"`
	Product   product.Model `gorm:"foreignKey:ProductId;references:ID"`
}

func (*Model) TableName() string {
	return "Order_Product"
}

func (m *Model) ToDomain(order *domain.Order) *domain.OrderProduct {
	result := &domain.OrderProduct{
		Quantity:  m.Quantity,
		ProductId: m.ProductId,
		OrderId:   m.OrderId,
		Product:   *m.Product.ToDomain(),
	}
	if order != nil {
		result.Order = *order
	}
	return result
}

func (m *Model) FromDomain(d *domain.OrderProduct) {
	if d == nil {
		return
	}
	m.Quantity = d.Quantity
	m.ProductId = d.ProductId
	m.OrderId = d.OrderId
	if d.Product.ID != 0 {
		m.Product.FromDomain(&d.Product)
	}
}
