package domain

import "time"

type Vehicle struct {
	ID        uint       `json:"id,omitempty"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Plate     *Plate     `json:"plate,omitempty"`
	Name      string     `json:"name"`
	Year      int        `json:"year"`
	Brand     string     `json:"brand"`
	Color     string     `json:"color"`
}
