package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"time"
)

type OrderHistoryModel struct {
	OrderId uint64
	Date    time.Time `gorm:"not null"`
	Status  string    `gorm:"not null"`
}

func (OrderHistoryModel) TableName() string {
	return "Order_History"
}

func (m *OrderHistoryModel) ToDomain() *domain.OrderHistory {
	if m == nil {
		return nil
	}
	return &domain.OrderHistory{
		OrderId: m.OrderId,
		Date:    m.Date,
		Status:  m.Status,
	}
}

func FromDomainOrderHistory(d *domain.OrderHistory) *OrderHistoryModel {
	if d == nil {
		return nil
	}
	return &OrderHistoryModel{
		OrderId: d.OrderId,
		Date:    d.Date,
		Status:  d.Status,
	}
}
