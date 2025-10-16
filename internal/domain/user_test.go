package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserInitialization(t *testing.T) {
	u := &User{
		ID:      1,
		Name:    "Chris Bumstead",
		Email:   "chris@example.com",
		Contact: "11999999999",
		Address: &Address{
			Address:       "123 Main St",
			AddressNumber: "123",
			Neighborhood:  "Manhattan",
			City:          "New York",
			Country:       "USA",
			ZipCode:       "10001",
		},
		Role: "admin",
	}

	assert.NotNil(t, u)
	assert.Equal(t, "Chris Bumstead", u.Name)
	assert.Equal(t, "chris@example.com", u.Email)
	assert.Equal(t, "11999999999", u.Contact)
	assert.Equal(t, "123 Main St", u.Address.Address)
	assert.Equal(t, "123", u.Address.AddressNumber)
	assert.Equal(t, "Manhattan", u.Address.Neighborhood)
	assert.Equal(t, "New York", u.Address.City)
	assert.Equal(t, "USA", u.Address.Country)
	assert.Equal(t, "10001", u.Address.ZipCode)
	assert.Equal(t, "admin", u.Role)
}
