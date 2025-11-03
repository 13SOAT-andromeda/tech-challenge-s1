package domain

import (
	"fmt"
	"slices"
)

type User struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Contact  string    `json:"contact"`
	Address  *Address  `json:"address"`
	Password *Password `json:"-"`
	Role     string    `json:"role"`
	Active   bool      `json:"active"`
}

func (c *User) ValidateRole() error {

	acceptedTypes := []string{"customer", "attendant", "mechanic", "administrator"}

	if !slices.Contains(acceptedTypes, c.Role) {
		return fmt.Errorf("User role '%s' is not valid. Accepted types: %v", c.Role, acceptedTypes)
	}

	return nil
}
