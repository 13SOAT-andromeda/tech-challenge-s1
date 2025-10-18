package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null; unique"`
	Contact  string `gorm:"not null"`
	Address  string `gorm:"not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`

	Sessions []SessionModel `gorm:"foreignKey:UserId;references:ID"`
}

func (UserModel) TableName() string {
	return "User"
}
