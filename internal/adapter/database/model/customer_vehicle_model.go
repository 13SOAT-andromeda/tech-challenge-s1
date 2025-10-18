package model

import "gorm.io/gorm"

type CustomerVehicleModel struct {
	gorm.Model
	CustomerId uint
	VehicleId  uint

	Vehicle  VehicleModel  `gorm:"foreignKey:VehicleId;references:ID"`
	Customer CustomerModel `gorm:"foreignKey:CustomerId;references:ID"`
}

func (CustomerVehicleModel) TableName() string {
	return "Customer_Vehicle"
}
