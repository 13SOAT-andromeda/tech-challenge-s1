package customer

import (
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/document"
	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
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

func makePersonModel(name, email, doc, contact string) personModel.Model {
	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = name
	pm.Email = email
	pm.Document = document.Model{Document: doc}
	pm.Contact = contact
	return pm
}

func TestCustomerModelInitialization(t *testing.T) {
	model := Model{
		Type:     "administrator",
		PersonID: 1,
		Person:   makePersonModel("Gedan", "gedan@example.com", "12345678900", "11999999999"),
	}

	assert.NotNil(t, model)
	assert.Equal(t, "administrator", model.Type)
	assert.Equal(t, "Gedan", model.Person.Name)
	assert.Equal(t, "gedan@example.com", model.Person.Email)
	assert.Equal(t, "12345678900", model.Person.Document.Document)
	assert.Equal(t, "11999999999", model.Person.Contact)
}

func TestCustomerModel_ToDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	pm := makePersonModel("Gedan", "gedan@example.com", "12345678900", "11999999999")
	pm.Address = &address.Model{
		Address:       "Rua Teste",
		AddressNumber: "317",
		Neighborhood:  "Centro",
		City:          "New York",
		Country:       "Brasil",
		ZipCode:       "1234",
	}

	model := Model{
		Type:     "administrator",
		PersonID: 1,
		Person:   pm,
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
	assert.Equal(t, "administrator", result.Type)
	assert.NotNil(t, result.Person)
	assert.Equal(t, "Gedan", result.Person.Name)
	assert.Equal(t, "gedan@example.com", result.Person.Email)
	assert.Equal(t, "12345678900", result.Person.Document.GetDocumentNumber())
	assert.Equal(t, "11999999999", result.Person.Contact)
	assert.NotNil(t, result.Person.Address)
	assert.Equal(t, "Rua Teste", result.Person.Address.Address)
	assert.Equal(t, "317", result.Person.Address.AddressNumber)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now, result.UpdatedAt)
	assert.Equal(t, deletedAt, *result.DeletedAt)
}

func TestCustomerModel_ToDomain_WithNilPerson(t *testing.T) {
	model := Model{
		Type:     "administrator",
		PersonID: 0,
	}

	result := model.ToDomain()

	assert.NotNil(t, result)
	assert.Nil(t, result.Person)
}

func TestCustomerModel_FromDomain(t *testing.T) {
	domainCustomer := &domain.Customer{
		ID:   1,
		Type: "mechanic",
		Person: &domain.Person{
			Name:    "Gedan",
			Email:   "gedan@example.com",
			Contact: "11988887777",
			Document: &domain.Document{
				Number: "98765432100",
			},
			Address: &domain.Address{
				Address:       "Rua Nova",
				AddressNumber: "999",
				Neighborhood:  "Bairro",
				City:          "Cidade",
				Country:       "Brasil",
				ZipCode:       "12345-678",
			},
		},
	}

	model := Model{}
	model.FromDomain(domainCustomer)

	assert.Equal(t, uint(1), model.ID)
	assert.Equal(t, "mechanic", model.Type)
	assert.Equal(t, "Gedan", model.Person.Name)
	assert.Equal(t, "gedan@example.com", model.Person.Email)
	assert.Equal(t, "98765432100", model.Person.Document.Document)
	assert.Equal(t, "11988887777", model.Person.Contact)
	assert.NotNil(t, model.Person.Address)
	assert.Equal(t, "Rua Nova", model.Person.Address.Address)
	assert.Equal(t, "999", model.Person.Address.AddressNumber)
}

func TestCustomerModel_FromDomain_Nil(t *testing.T) {
	model := Model{
		Type: "Existing",
	}

	model.FromDomain(nil)

	assert.Equal(t, "Existing", model.Type)
}

func TestCustomerModel_RoundTrip(t *testing.T) {
	pm := makePersonModel("Gedan", "gedan@example.com", "12345678900", "11999999999")
	original := Model{
		Type:     "administrator",
		PersonID: 1,
		Person:   pm,
	}
	original.ID = 1

	domainCustomer := original.ToDomain()

	converted := Model{}
	converted.FromDomain(domainCustomer)

	assert.Equal(t, original.ID, converted.ID)
	assert.Equal(t, original.Type, converted.Type)
	assert.Equal(t, original.Person.Name, converted.Person.Name)
	assert.Equal(t, original.Person.Email, converted.Person.Email)
}
