package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerInitialization(t *testing.T) {
	c := Customer{
		ID:   1,
		Type: "PF",
		Person: &Person{
			Name:  "Gedan",
			Email: "gedan@example.com",
			Document: &Document{
				Number: "45653421898",
			},
			Contact: "11999999999",
			Address: &Address{
				Address:       "Rua Teste",
				City:          "New York",
				AddressNumber: "317",
				ZipCode:       "1234",
				Neighborhood:  "New York",
				Country:       "Brasil",
			},
		},
	}

	assert.NotNil(t, c)
	assert.NotNil(t, c.Person)
	assert.Equal(t, "Gedan", c.Person.Name)
	assert.Equal(t, "gedan@example.com", c.Person.Email)
	assert.Equal(t, "45653421898", c.Person.Document.GetDocumentNumber())
	assert.Equal(t, "PF", c.Type)
	assert.Equal(t, "11999999999", c.Person.Contact)
	assert.Equal(t, "Rua Teste", c.Person.Address.Address)
	assert.Equal(t, "317", c.Person.Address.AddressNumber)
	assert.Equal(t, "New York", c.Person.Address.Neighborhood)
	assert.Equal(t, "New York", c.Person.Address.City)
	assert.Equal(t, "Brasil", c.Person.Address.Country)
	assert.Equal(t, "1234", c.Person.Address.ZipCode)
}
