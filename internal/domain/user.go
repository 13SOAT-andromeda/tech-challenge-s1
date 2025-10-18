package domain

import (
	"time"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Email     string
	Contact   string
	Address   string
	Password  string
	Role      string

	Sessions []Session
}
