package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerInitialization(t *testing.T) {
	c := Customer{
		ID:       1,
		Name:     "Gedan",
		Email:    "gedan@example.com",
		Document: "12345678900",
		Type:     "teste",
		Contact:  "11999999999",
		Address: &Address{
			Address:       "Rua Teste",
			City:          "New York",
			AddressNumber: "317",
			ZipCode:       "1234",
			Neighborhood:  "New York",
			Country:       "Brasil",
		},
	}

	assert.NotNil(t, c)
	assert.Equal(t, "Gedan", c.Name)
	assert.Equal(t, "gedan@example.com", c.Email)
	assert.Equal(t, "12345678900", c.Document)
	assert.Equal(t, "teste", c.Type)
	assert.Equal(t, "11999999999", c.Contact)
	assert.Equal(t, "Rua Teste", c.Address.Address)
	assert.Equal(t, "317", c.Address.AddressNumber)
	assert.Equal(t, "New York", c.Address.Neighborhood)
	assert.Equal(t, "New York", c.Address.City)
	assert.Equal(t, "Brasil", c.Address.Country)
	assert.Equal(t, "1234", c.Address.ZipCode)
}
