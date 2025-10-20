package order_service

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
	ServiceId uint
	OrderId   uint
	Price     float64 `gorm:"not null"`

	Service model.ServiceModel `gorm:"foreignkey:ServiceId;references:ID"`
	Order   order.Model        `gorm:"foreignkey:OrderId;references:ID"`
}

func (Model) TableName() string {
	return "Order_Service"
}

func ToDomain(m Model) domain.OrderService {
	return domain.OrderService{
		ServiceId: m.ServiceId,
		OrderId:   m.OrderId,
		Price:     m.Price,
		Service:   *(&m.Service).ToDomain(),
		Order:     order.ToDomain(m.Order),
	}
}

func FromDomain(d domain.OrderService) Model {
	return Model{
		ServiceId: d.ServiceId,
		OrderId:   d.OrderId,
		Price:     d.Price,
		Service:   *model.FromDomainService(&d.Service),
		Order:     order.FromDomain(d.Order),
	}
}
