package model

import "time"

type OrderHistoryModel struct {
	OrderId uint64
	Date    time.Time `gorm:"not null"`
	Status  string    `gorm:"not null"`
}

func (OrderHistoryModel) TableName() string {
	return "Order_History"
}
