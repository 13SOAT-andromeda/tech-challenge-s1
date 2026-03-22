package domain

import (
	"fmt"
	"slices"
	"time"
)

type User struct {
	ID        uint       `json:"id"`
	Password  *Password  `json:"-"`
	Role      string     `json:"role"`
	PersonID  uint       `json:"person_id"`
	Person    *Person    `json:"person,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (c *User) ValidateRole() error {
	acceptedTypes := []string{"customer", "attendant", "mechanic", "administrator"}

	if !slices.Contains(acceptedTypes, c.Role) {
		return fmt.Errorf("User role '%s' is not valid. Accepted types: %v", c.Role, acceptedTypes)
	}

	return nil
}
