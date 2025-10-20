package domain

import (
	"time"
)

type CustomerVehicle struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	CustomerId uint
	VehicleId  uint
}
