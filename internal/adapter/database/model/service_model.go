package model

import "gorm.io/gorm"

type ServiceModel struct {
	gorm.Model
	Name         string `gorm:"unique; not null"`
	DefaultPrice *float64
	CategoryId   uint
	Number       string `gorm:"not null"`

	ServiceCategory ServiceCategoryModel `gorm:"foreignKey:CategoryId;references:ID"`
}

func (ServiceModel) TableName() string {
	return "Service"
}
