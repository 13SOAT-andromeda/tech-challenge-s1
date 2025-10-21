package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompanyTableName(t *testing.T) {
	assert.Equal(t, "Company", CompanyModel{}.TableName())
}

func TestNilCompanyToDomain(t *testing.T) {
	assert.Nil(t, (*CompanyModel)(nil).ToDomain())
}

func TestNilCompanyFromDomain(t *testing.T) {
	assert.Nil(t, FromDomainCompany(nil))
}

func TestCompanyModelInitialization(t *testing.T) {
	c := CompanyModel{
		Name:     "Company Name",
		Email:    "company@example.com",
		Document: "12345678901234",
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
	assert.Equal(t, "Company Name", c.Name)
	assert.Equal(t, "company@example.com", c.Email)
	assert.Equal(t, "12345678901234", c.Document)
	assert.Equal(t, "11999999999", c.Contact)
	assert.Equal(t, "Rua Teste", c.Address.Address)
	assert.Equal(t, "317", c.Address.AddressNumber)
	assert.Equal(t, "New York", c.Address.Neighborhood)
	assert.Equal(t, "New York", c.Address.City)
	assert.Equal(t, "Brasil", c.Address.Country)
	assert.Equal(t, "1234", c.Address.ZipCode)
}

func TestCompanyModel_ToFromDomain(t *testing.T) {
	modelCompany := &CompanyModel{
		Name:     "Company Name",
		Email:    "company@example.com",
		Document: "12345678901234",
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

	domainCompany := modelCompany.ToDomain()

	assert.Equal(t, modelCompany.ID, domainCompany.ID)
	assert.Equal(t, modelCompany.Name, domainCompany.Name)
	assert.Equal(t, modelCompany.Email, domainCompany.Email)
	assert.Equal(t, modelCompany.Document, domainCompany.Document)
	assert.Equal(t, modelCompany.Contact, domainCompany.Contact)
	assert.Equal(t, modelCompany.Address.Address, domainCompany.Address.Address)

	newModel := FromDomainCompany(domainCompany)

	assert.Equal(t, modelCompany, newModel)
}
