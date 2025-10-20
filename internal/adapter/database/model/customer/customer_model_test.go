package customer

import (
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAddressModel_ToDomain(t *testing.T) {
	model := address.Model{
		Address:       "Rua Teste",
		AddressNumber: "123",
		Neighborhood:  "Centro",
		City:          "Cidade",
		Country:       "Brasil",
		ZipCode:       "12345-678",
	}

	result := address.ToDomain(model)

	assert.Equal(t, model.Address, result.Address)
	assert.Equal(t, model.AddressNumber, result.AddressNumber)
	assert.Equal(t, model.Neighborhood, result.Neighborhood)
	assert.Equal(t, model.City, result.City)
	assert.Equal(t, model.Country, result.Country)
	assert.Equal(t, model.ZipCode, result.ZipCode)
}

func TestAddressModel_FromDomain(t *testing.T) {
	domainAddr := domain.Address{
		Address:       "Rua Teste",
		AddressNumber: "123",
		Neighborhood:  "Centro",
		City:          "Cidade",
		Country:       "Brasil",
		ZipCode:       "12345-678",
	}

	model := address.FromDomain(domainAddr)

	assert.Equal(t, domainAddr.Address, model.Address)
	assert.Equal(t, domainAddr.AddressNumber, model.AddressNumber)
	assert.Equal(t, domainAddr.Neighborhood, model.Neighborhood)
	assert.Equal(t, domainAddr.City, model.City)
	assert.Equal(t, domainAddr.Country, model.Country)
	assert.Equal(t, domainAddr.ZipCode, model.ZipCode)
}

func TestCustomerModelInitialization(t *testing.T) {
	model := Model{
		Name:     "Gedan",
		Email:    "gedan@example.com",
		Document: "12345678900",
		Type:     "teste",
		Contact:  "11999999999",
		Address: address.Model{
			Address:       "Rua Teste",
			City:          "New York",
			AddressNumber: "317",
			ZipCode:       "1234",
			Neighborhood:  "New York",
			Country:       "Brasil",
		},
	}

	assert.NotNil(t, model)
	assert.Equal(t, "Gedan", model.Name)
	assert.Equal(t, "gedan@example.com", model.Email)
	assert.Equal(t, "12345678900", model.Document)
	assert.Equal(t, "teste", model.Type)
	assert.Equal(t, "11999999999", model.Contact)
	assert.Equal(t, "Rua Teste", model.Address.Address)
	assert.Equal(t, "317", model.Address.AddressNumber)
	assert.Equal(t, "New York", model.Address.Neighborhood)
	assert.Equal(t, "New York", model.Address.City)
	assert.Equal(t, "Brasil", model.Address.Country)
	assert.Equal(t, "1234", model.Address.ZipCode)
}

func TestCustomerModel_ToFromDomain(t *testing.T) {
	model := Model{
		Name:     "Gedan",
		Email:    "gedan@example.com",
		Document: "12345678900",
		Type:     "teste",
		Contact:  "11999999999",
		Address: address.Model{
			Address:       "Rua Teste",
			AddressNumber: "317",
			Neighborhood:  "Centro",
			City:          "New York",
			Country:       "Brasil",
			ZipCode:       "1234",
		},
	}

	result := ToDomain(model)

	assert.Equal(t, model.ID, result.ID)
	assert.Equal(t, model.Name, result.Name)
	assert.Equal(t, model.Email, result.Email)
	assert.Equal(t, model.Document, result.Document)
	assert.Equal(t, model.Type, result.Type)
	assert.Equal(t, model.Contact, result.Contact)
	assert.Equal(t, model.Address.Address, result.Address.Address)

	newModel := FromDomain(result)

	assert.Equal(t, model, newModel)
}
