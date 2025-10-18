package domain

import (
	"time"
)

type Service struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Name         string
	DefaultPrice *float64
	CategoryId   uint
	Number       string

	ServiceCategory ServiceCategory
}
