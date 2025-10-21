package model

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
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

func (m *OrderHistoryModel) FromDomain(d *domain.OrderHistory) {
	if d == nil {
		return
	}
	m.OrderId = d.OrderId
	m.Date = d.Date
	m.Status = d.Status
}
