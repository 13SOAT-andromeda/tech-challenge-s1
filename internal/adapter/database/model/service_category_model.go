package model

import "gorm.io/gorm"

type ServiceCategoryModel struct {
	gorm.Model
	Name string `gorm:"not null"`
}

func (ServiceCategoryModel) TableName() string {
	return "Service_Category"
}
