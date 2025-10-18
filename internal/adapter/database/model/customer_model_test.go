package model

import (
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAddressModel_ToDomain(t *testing.T) {
	modelAddr := &AddressModel{
		Address:       "Rua Teste",
		AddressNumber: "123",
		Neighborhood:  "Centro",
		City:          "Cidade",
		Country:       "Brasil",
		ZipCode:       "12345-678",
	}

	domainAddr := modelAddr.ToDomain()

	assert.Equal(t, modelAddr.Address, domainAddr.Address)
	assert.Equal(t, modelAddr.AddressNumber, domainAddr.AddressNumber)
	assert.Equal(t, modelAddr.Neighborhood, domainAddr.Neighborhood)
	assert.Equal(t, modelAddr.City, domainAddr.City)
	assert.Equal(t, modelAddr.Country, domainAddr.Country)
	assert.Equal(t, modelAddr.ZipCode, domainAddr.ZipCode)
}

func TestFromDomainAddress(t *testing.T) {
	domainAddr := &domain.Address{
		Address:       "Rua Teste",
		AddressNumber: "123",
		Neighborhood:  "Centro",
		City:          "Cidade",
		Country:       "Brasil",
		ZipCode:       "12345-678",
	}

	modelAddr := FromDomainAddress(domainAddr)

	assert.Equal(t, domainAddr.Address, modelAddr.Address)
	assert.Equal(t, domainAddr.AddressNumber, modelAddr.AddressNumber)
	assert.Equal(t, domainAddr.Neighborhood, modelAddr.Neighborhood)
	assert.Equal(t, domainAddr.City, modelAddr.City)
	assert.Equal(t, domainAddr.Country, modelAddr.Country)
	assert.Equal(t, domainAddr.ZipCode, modelAddr.ZipCode)

	nilModel := FromDomainAddress(nil)
	assert.Equal(t, AddressModel{}, nilModel)
}

func TestCustomerModelInitialization(t *testing.T) {
	c := CustomerModel{
		Name:     "Gedan",
		Email:    "gedan@example.com",
		Document: "12345678900",
		Type:     "teste",
		Contact:  "11999999999",
		Address: AddressModel{
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

func Test_EnsureAddress(t *testing.T) {
	c1 := &domain.Customer{}
	c1.EnsureAddress()

	assert.Equal(t, &domain.Address{}, c1.Address)

	existingAddress := &domain.Address{
		Address:       "Rua Teste",
		AddressNumber: "123",
		Neighborhood:  "Centro",
		City:          "Cidade",
		Country:       "Brasil",
		ZipCode:       "12345-678",
	}
	c2 := &domain.Customer{
		Address: existingAddress,
	}
	c2.EnsureAddress()

	assert.Equal(t, existingAddress, c2.Address, "Address existente não deve ser sobrescrito")
}

func TestCustomerModel_ToFromDomain(t *testing.T) {
	modelCustomer := CustomerModel{
		Name:     "Gedan",
		Email:    "gedan@example.com",
		Document: "12345678900",
		Type:     "teste",
		Contact:  "11999999999",
		Address: AddressModel{
			Address:       "Rua Teste",
			AddressNumber: "317",
			Neighborhood:  "Centro",
			City:          "New York",
			Country:       "Brasil",
			ZipCode:       "1234",
		},
	}

	domainCustomer := ToDomain(modelCustomer)

	assert.Equal(t, modelCustomer.ID, domainCustomer.ID)
	assert.Equal(t, modelCustomer.Name, domainCustomer.Name)
	assert.Equal(t, modelCustomer.Email, domainCustomer.Email)
	assert.Equal(t, modelCustomer.Document, domainCustomer.Document)
	assert.Equal(t, modelCustomer.Type, domainCustomer.Type)
	assert.Equal(t, modelCustomer.Contact, domainCustomer.Contact)
	assert.Equal(t, modelCustomer.Address.Address, domainCustomer.Address.Address)

	newModel := FromDomain(domainCustomer)

	assert.Equal(t, modelCustomer, newModel)
}
