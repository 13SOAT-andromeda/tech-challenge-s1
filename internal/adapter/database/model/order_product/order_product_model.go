package order_product

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
	ProductId uint
	OrderId   uint
	Product   product.Model `gorm:"foreignKey:ProductId;references:ID"`
	Order     order.Model   `gorm:"foreignKey:OrderId;references:ID"`
}

func (*Model) TableName() string {
	return "Order_Product"
}

func (m *Model) ToDomain() *domain.OrderProduct {
	return &domain.OrderProduct{
		ProductId: m.ProductId,
		OrderId:   m.OrderId,
		Product:   *m.Product.ToDomain(),
		Order:     *m.Order.ToDomain(),
	}
}

func (m *Model) FromDomain(d *domain.OrderProduct) {
	if d == nil {
		return
	}
	m.ProductId = d.ProductId
	m.OrderId = d.OrderId
	m.Product.FromDomain(&d.Product)
	m.Order.FromDomain(&d.Order)
}
