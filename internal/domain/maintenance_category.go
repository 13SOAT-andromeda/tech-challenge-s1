package domain

import (
	"time"
)

type MaintenanceCategory struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
}
