package domain

import (
	"time"
)

type Session struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	UserID       uint
	RefreshToken *string
}
