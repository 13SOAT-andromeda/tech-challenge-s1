package model

import "gorm.io/gorm"

type ProductModel struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Quantity uint   `gorm:"not null"`
	Price    uint32 `gorm:"not null"` // em centavos
}

func (ProductModel) TableName() string {
	return "Product"
}
