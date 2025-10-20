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

func ToDomain(m Model) domain.OrderProduct {
	return domain.OrderProduct{
		ProductId: m.ProductId,
		OrderId:   m.OrderId,
		Product:   *(&m.Product).ToDomain(),
		Order:     order.ToDomain(m.Order),
	}
}

func FromDomain(d domain.OrderProduct) Model {
	return Model{
		ProductId: d.ProductId,
		OrderId:   d.OrderId,
		Product:   *model.FromDomainProduct(&d.Product),
		Order:     order.FromDomain(d.Order),
	}
}
