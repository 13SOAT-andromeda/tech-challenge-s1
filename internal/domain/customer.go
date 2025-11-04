package domain

import (
	"time"
)

type Customer struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Document  *Document  `json:"document,omitempty"`
	Type      string     `json:"type"`
	Contact   string     `json:"contact"`
	Address   *Address   `json:"address,omitempty"`
}
