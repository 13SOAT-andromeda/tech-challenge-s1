package model

import "gorm.io/gorm"

type SessionModel struct {
	gorm.Model
	UserID       uint `gorm:"not null"`
	RefreshToken *string
}

func (SessionModel) TableName() string {
	return "Session"
}
