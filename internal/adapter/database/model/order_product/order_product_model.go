package order_product

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
	ProductId uint
	OrderId   uint
	Product   model.ProductModel `gorm:"foreignkey:ProductId;references:Id"`
	Order     order.Model        `gorm:"foreignkey:OrderId;references:Id"`
}

func (Model) TableName() string {
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
