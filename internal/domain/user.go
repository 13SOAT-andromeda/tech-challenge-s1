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
	Employee  *Employee  `json:"employee,omitempty"`
	Customer  *Customer  `json:"customer,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (u *User) ValidateRole() error {
	acceptedTypes := []string{"customer", "attendant", "mechanic", "administrator"}

	if !slices.Contains(acceptedTypes, u.Role) {
		return fmt.Errorf("User role '%s' is not valid. Accepted types: %v", u.Role, acceptedTypes)
	}

	return nil
}

func (u *User) IsCustomer() bool {
	return u.Role == "customer"
}

func (u *User) IsEmployee() bool {
	return u.Role == "attendant" || u.Role == "mechanic" || u.Role == "administrator"
}
