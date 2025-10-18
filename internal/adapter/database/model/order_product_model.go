package model

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type OrderProductModel struct {
	ProductId uint
	OrderId   uint

	Product ProductModel `gorm:"foreignkey:ProductId;references:Id"`
	Order   OrderModel   `gorm:"foreignkey:OrderId;references:Id"`
}

func (OrderProductModel) TableName() string {
	return "Order_Product"
}

func (m *OrderProductModel) ToDomain() *domain.OrderProduct {
	if m == nil {
		return nil
	}

	return &domain.OrderProduct{
		ProductId: m.ProductId,
		OrderId:   m.OrderId,
		Product:   *(&m.Product).ToDomain(),
		Order:     *(&m.Order).ToDomain(),
	}
}

func FromDomainOrderProduct(d *domain.OrderProduct) *OrderProductModel {
	if d == nil {
		return nil
	}

	return &OrderProductModel{
		ProductId: d.ProductId,
		OrderId:   d.OrderId,
		Product:   *FromDomainProduct(&d.Product),
		Order:     *FromDomainOrder(&d.Order),
	}
}
