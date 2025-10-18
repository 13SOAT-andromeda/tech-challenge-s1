package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompanyInitialization(t *testing.T) {
	c := Company{
		ID:       1,
		Name:     "Teste Company",
		Email:    "company_test@example.com",
		Document: "12345678900",
		Contact:  "11999999999",
		Address: &Address{
			Address:       "Rua Teste",
			City:          "Sao Paulo",
			AddressNumber: "123",
			ZipCode:       "00012-345",
			Neighborhood:  "Centro",
			Country:       "Brasil",
		},
	}

	assert.NotNil(t, c)
	assert.Equal(t, "Teste Company", c.Name)
	assert.Equal(t, "company_test@example.com", c.Email)
	assert.Equal(t, "12345678900", c.Document)
	assert.Equal(t, "11999999999", c.Contact)
	assert.Equal(t, "Rua Teste", c.Address.Address)
	assert.Equal(t, "Sao Paulo", c.Address.City)
	assert.Equal(t, "123", c.Address.AddressNumber)
	assert.Equal(t, "Centro", c.Address.Neighborhood)
	assert.Equal(t, "00012-345", c.Address.ZipCode)
	assert.Equal(t, "Brasil", c.Address.Country)
}
