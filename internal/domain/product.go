package domain

import "time"

type Product struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Quantity  uint
	Price     uint32
}
