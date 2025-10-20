package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModelInitialization(t *testing.T) {
	u := UserModel{
		Name:    "User A",
		Email:   "user.a@example.com",
		Contact: "11987654321",
		Address: AddressModel{
			Address:       "Rua dos Testes, 123",
			AddressNumber: "123",
			Neighborhood:  "Centro",
			City:          "São Paulo",
			Country:       "Brasil",
			ZipCode:       "01234-567",
		},
		Password: "password",
		Role:     "user",
		Active:   true,
	}

	assert.NotNil(t, u)
	assert.Equal(t, "User A", u.Name)
	assert.Equal(t, "user.a@example.com", u.Email)
	assert.Equal(t, "11987654321", u.Contact)
	assert.Equal(t, "Rua dos Testes, 123", u.Address.Address)
	assert.Equal(t, "123", u.Address.AddressNumber)
	assert.Equal(t, "Centro", u.Address.Neighborhood)
	assert.Equal(t, "São Paulo", u.Address.City)
	assert.Equal(t, "Brasil", u.Address.Country)
	assert.Equal(t, "01234-567", u.Address.ZipCode)
	assert.Equal(t, "password", u.Password)
	assert.Equal(t, "user", u.Role)
	assert.Equal(t, true, u.Active)
}

func TestUserModel_ToFromDomain(t *testing.T) {

	modelUser := UserModel{
		ID:      1,
		Name:    "User A",
		Email:   "user.a@example.com",
		Contact: "11987654321",
		Address: AddressModel{
			Address:       "Rua dos Testes, 123",
			AddressNumber: "123",
			Neighborhood:  "Centro",
			City:          "São Paulo",
			Country:       "Brasil",
			ZipCode:       "01234-567",
		},
		Password: "password",
		Role:     "user",
		Active:   true,
	}

	domainUser := modelUser.ToDomain()

	assert.Equal(t, modelUser.ID, domainUser.ID)
	assert.Equal(t, modelUser.Name, domainUser.Name)
	assert.Equal(t, modelUser.Email, domainUser.Email)
	assert.Equal(t, modelUser.Contact, domainUser.Contact)
	assert.Equal(t, modelUser.Address.Address, domainUser.Address.Address)
	assert.Equal(t, modelUser.Address.AddressNumber, domainUser.Address.AddressNumber)
	assert.Equal(t, modelUser.Address.Neighborhood, domainUser.Address.Neighborhood)
	assert.Equal(t, modelUser.Address.City, domainUser.Address.City)
	assert.Equal(t, modelUser.Address.Country, domainUser.Address.Country)
	assert.Equal(t, modelUser.Address.ZipCode, domainUser.Address.ZipCode)
	assert.Equal(t, modelUser.Password, domainUser.Password.GetHashed())
	assert.Equal(t, modelUser.Role, domainUser.Role)
	assert.Equal(t, modelUser.Active, domainUser.Active)

	newModel := NewUserModelFromDomain(domainUser)

	assert.Equal(t, modelUser, newModel)
}
