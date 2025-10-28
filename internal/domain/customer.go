package domain

import (
	"fmt"
	"slices"
	"time"
)

type Customer struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Email     string
	Document  *Document
	Type      string
	Contact   string
	Address   *Address
}

func (c *Customer) ValidateCustomerType() error {

	acceptedTypes := []string{"administrator", "attendant", "mechanic"}

	if !slices.Contains(acceptedTypes, c.Type) {
		return fmt.Errorf("customer type '%s' is not valid. Accepted types: %v", c.Type, acceptedTypes)
	}

	return nil
}
