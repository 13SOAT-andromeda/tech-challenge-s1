package model

import "gorm.io/gorm"

type CompanyModel struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string
	Document string       `gorm:"unique"`
	Contact  string       `gorm:"not null"`
	Address  AddressModel `gorm:"embedded"`
}

func (CompanyModel) TableName() string {
	return "Company"
}
