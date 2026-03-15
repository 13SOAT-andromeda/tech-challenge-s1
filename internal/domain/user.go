package domain

import (
	"fmt"
	"slices"
	"time"
)

type User struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Document  *Document  `json:"document,omitempty"`
	Contact   string     `json:"contact"`
	Address   *Address   `json:"address"`
	Password  *Password  `json:"-"`
	Role      string     `json:"role"`
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
