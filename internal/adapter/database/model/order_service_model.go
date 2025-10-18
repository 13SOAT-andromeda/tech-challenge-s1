package model

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type OrderServiceModel struct {
	ServiceId uint
	OrderId   uint
	Price     float64 `gorm:"not null"`

	Service ServiceModel `gorm:"foreignkey:ServiceId;references:ID"`
	Order   OrderModel   `gorm:"foreignkey:OrderId;references:ID"`
}

func (OrderServiceModel) TableName() string {
	return "Order_Service"
}

func (m *OrderServiceModel) ToDomain() *domain.OrderService {
	if m == nil {
		return nil
	}

	return &domain.OrderService{
		ServiceId: m.ServiceId,
		OrderId:   m.OrderId,
		Price:     m.Price,
		Service:   *(&m.Service).ToDomain(),
		Order:     *(&m.Order).ToDomain(),
	}
}

func FromDomainOrderService(d *domain.OrderService) *OrderServiceModel {
	if d == nil {
		return nil
	}

	return &OrderServiceModel{
		ServiceId: d.ServiceId,
		OrderId:   d.OrderId,
		Price:     d.Price,
		Service:   *FromDomainService(&d.Service),
		Order:     *FromDomainOrder(&d.Order),
	}
}
