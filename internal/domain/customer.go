package domain

import "time"

type Customer struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Email     string
	Document  string
	Type      string
	Contact   string
	Address   *Address
}
