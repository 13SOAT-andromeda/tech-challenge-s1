package domain

import "time"

type Employee struct {
	ID        uint       `json:"id"`
	Position  string     `json:"position"`
	PersonID  uint       `json:"person_id"`
	Person    *Person    `json:"person,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
