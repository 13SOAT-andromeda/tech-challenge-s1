package order_service

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
	ServiceId uint
	OrderId   uint
	Price     float64 `gorm:"not null"`

	Service maintenance.Model `gorm:"foreignkey:ServiceId;references:ID"`
	Order   order.Model       `gorm:"foreignkey:OrderId;references:ID"`
}

func (Model) TableName() string {
	return "Order_Service"
}

func (m *Model) ToDomain() *domain.OrderService {
	return &domain.OrderService{
		ServiceId: m.ServiceId,
		OrderId:   m.OrderId,
		Price:     m.Price,
		Service:   *m.Service.ToDomain(),
		Order:     *m.Order.ToDomain(),
	}
}

func (m *Model) FromDomain(d *domain.OrderService) {
	if d == nil {
		return
	}
	m.ServiceId = d.ServiceId
	m.OrderId = d.OrderId
	m.Price = d.Price
	m.Service.FromDomain(&d.Service)
	m.Order.FromDomain(&d.Order)
}
