package customer

import (
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/document"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddressModel_TableName(t *testing.T) {
	customer := &Model{}
	assert.Equal(t, "Customer", customer.TableName())
}

func TestAddressModel_ToDomain(t *testing.T) {
	model := address.Model{
		Address:       "Rua Teste",
		AddressNumber: "123",
		Neighborhood:  "Centro",
		City:          "Cidade",
		Country:       "Brasil",
		ZipCode:       "12345-678",
	}

	result := model.ToDomain()

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

	model := address.Model{}
	model.FromDomain(&domainAddr)

	assert.Equal(t, domainAddr.Address, model.Address)
	assert.Equal(t, domainAddr.AddressNumber, model.AddressNumber)
	assert.Equal(t, domainAddr.Neighborhood, model.Neighborhood)
	assert.Equal(t, domainAddr.City, model.City)
	assert.Equal(t, domainAddr.Country, model.Country)
	assert.Equal(t, domainAddr.ZipCode, model.ZipCode)
}

func TestAddressModel_FromDomain_Nil(t *testing.T) {
	model := address.Model{
		Address: "Existing",
	}

	model.FromDomain(nil)

	assert.Equal(t, "Existing", model.Address)
}

func TestDocumentModel_ToDomain(t *testing.T) {
	model := document.Model{
		Document: "12345678900",
	}

	result := model.ToDomain()

	assert.NotNil(t, result)
	assert.Equal(t, "12345678900", result.Number)
}

func TestDocumentModel_FromDomain(t *testing.T) {
	domainDoc := domain.Document{
		Number: "98765432100",
	}

	model := document.Model{}
	model.FromDomain(&domainDoc)

	assert.Equal(t, "98765432100", model.Document)
}

func TestDocumentModel_FromDomain_Nil(t *testing.T) {
	model := document.Model{
		Document: "12345678900",
	}

	model.FromDomain(nil)

	assert.Equal(t, "12345678900", model.Document)
}

func TestCustomerModelInitialization(t *testing.T) {
	model := Model{
		Name:  "Gedan",
		Email: "gedan@example.com",
		Document: document.Model{
			Document: "12345678900",
		},
		Type:    "administrator",
		Contact: "11999999999",
		Address: &address.Model{
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
	assert.Equal(t, "12345678900", model.Document.Document) // ← CORRIGIDO
	assert.Equal(t, "administrator", model.Type)
	assert.Equal(t, "11999999999", model.Contact)
	assert.NotNil(t, model.Address)
	assert.Equal(t, "Rua Teste", model.Address.Address)
	assert.Equal(t, "317", model.Address.AddressNumber)
	assert.Equal(t, "New York", model.Address.Neighborhood)
	assert.Equal(t, "New York", model.Address.City)
	assert.Equal(t, "Brasil", model.Address.Country)
	assert.Equal(t, "1234", model.Address.ZipCode)
}

func TestCustomerModel_ToDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)
	model := Model{
		Name:  "Gedan",
		Email: "gedan@example.com",
		Document: document.Model{
			Document: "12345678900",
		},
		Type:    "administrator",
		Contact: "11999999999",
		Address: &address.Model{
			Address:       "Rua Teste",
			AddressNumber: "317",
			Neighborhood:  "Centro",
			City:          "New York",
			Country:       "Brasil",
			ZipCode:       "1234",
		},
		Model: gorm.Model{
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
	}
	model.ID = 1
	model.CreatedAt = now
	model.UpdatedAt = now

	result := model.ToDomain()

	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Gedan", result.Name)
	assert.Equal(t, "gedan@example.com", result.Email)
	assert.Equal(t, "12345678900", result.Document.Number)
	assert.Equal(t, "administrator", result.Type)
	assert.Equal(t, "11999999999", result.Contact)
	assert.NotNil(t, result.Address)
	assert.Equal(t, "Rua Teste", result.Address.Address)
	assert.Equal(t, "317", result.Address.AddressNumber)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now, result.UpdatedAt)
	assert.Equal(t, deletedAt, *result.DeletedAt)
}

func TestCustomerModel_ToDomain_WithNilAddress(t *testing.T) {
	model := Model{
		Name:  "Gedan",
		Email: "gedan@example.com",
		Document: document.Model{
			Document: "12345678900",
		},
		Type:    "administrator",
		Contact: "11999999999",
		Address: nil,
	}

	result := model.ToDomain()

	assert.NotNil(t, result)
	assert.Nil(t, result.Address)
}

func TestCustomerModel_FromDomain(t *testing.T) {
	domainCustomer := &domain.Customer{
		ID:    1,
		Name:  "Gedan",
		Email: "gedan@example.com",
		Document: &domain.Document{
			Number: "98765432100",
		},
		Type:    "mechanic",
		Contact: "11988887777",
		Address: &domain.Address{
			Address:       "Rua Nova",
			AddressNumber: "999",
			Neighborhood:  "Bairro",
			City:          "Cidade",
			Country:       "Brasil",
			ZipCode:       "12345-678",
		},
	}

	model := Model{}
	model.FromDomain(domainCustomer)

	assert.Equal(t, uint(1), model.ID)
	assert.Equal(t, "Gedan", model.Name)
	assert.Equal(t, "gedan@example.com", model.Email)
	assert.Equal(t, "98765432100", model.Document.Document)
	assert.Equal(t, "mechanic", model.Type)
	assert.Equal(t, "11988887777", model.Contact)
	assert.NotNil(t, model.Address)
	assert.Equal(t, "Rua Nova", model.Address.Address)
	assert.Equal(t, "999", model.Address.AddressNumber)
}

func TestCustomerModel_FromDomain_Nil(t *testing.T) {
	model := Model{
		Name: "Existing",
	}

	model.FromDomain(nil)

	assert.Equal(t, "Existing", model.Name)
}

func TestCustomerModel_FromDomain_WithNilAddress(t *testing.T) {
	domainCustomer := &domain.Customer{
		ID:    1,
		Name:  "Gedan",
		Email: "gedan@example.com",
		Document: &domain.Document{
			Number: "98765432100",
		},
		Type:    "attendant",
		Contact: "11988887777",
		Address: nil,
	}

	model := Model{}
	model.FromDomain(domainCustomer)

	assert.Equal(t, uint(1), model.ID)
	assert.Equal(t, "Gedan", model.Name)
	assert.NotNil(t, model.Address)
}

func TestCustomerModel_RoundTrip(t *testing.T) {
	original := Model{
		Name:  "Gedan",
		Email: "gedan@example.com",
		Document: document.Model{
			Document: "12345678900",
		},
		Type:    "administrator",
		Contact: "11999999999",
		Address: &address.Model{
			Address:       "Rua Teste",
			AddressNumber: "317",
			City:          "Cidade",
			Country:       "Brasil",
			ZipCode:       "12345",
		},
	}
	original.ID = 1

	domainCustomer := original.ToDomain()

	converted := Model{}
	converted.FromDomain(domainCustomer)

	assert.Equal(t, original.ID, converted.ID)
	assert.Equal(t, original.Name, converted.Name)
	assert.Equal(t, original.Email, converted.Email)
	assert.Equal(t, original.Document.Document, converted.Document.Document)
	assert.Equal(t, original.Type, converted.Type)
	assert.Equal(t, original.Contact, converted.Contact)
}
