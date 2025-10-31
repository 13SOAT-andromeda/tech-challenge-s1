package domain

import "time"

type Vehicle struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Plate     *Plate
	Name      string
	Year      int
	Brand     string
	Color     string
}
