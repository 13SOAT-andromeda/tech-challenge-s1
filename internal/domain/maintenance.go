package domain

import (
	"time"
)

type Maintenance struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Name         string
	DefaultPrice *float64
	CategoryId   uint
	Number       string

	MaintenanceCategory MaintenanceCategory
}
