package model

import "time"

type OrderModel struct {
	ID                uint      `gorm:"primaryKey"`
	DateIn            time.Time `gorm:"not null"`
	DateOut           *time.Time
	Number            string `gorm:"not null; unique"`
	Status            string
	VehicleKilometers *float64
	Note              *string
	DiagnosticNote    *string
	Price             float64 `gorm:"not null"`
	UserId            uint
	CustomerVehicleId uint
	CompanyId         uint

	User            UserModel            `gorm:"foreignKey:UserId;references:ID"`
	CustomerVehicle CustomerVehicleModel `gorm:"foreignKey:CustomerVehicleId;references:ID"`
	Company         CompanyModel         `gorm:"foreignKey:CompanyId;references:ID"`
}

func (OrderModel) TableName() string {
	return "Order"
}
