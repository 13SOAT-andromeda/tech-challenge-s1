package model

import (
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHasher é um mock para o hasher do Password
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Generate(password []byte, cost int) ([]byte, error) {
	args := m.Called(password, cost)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockHasher) Compare(hashedPassword []byte, password []byte) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func TestUserModel_TableName(t *testing.T) {
	user := &UserModel{}
	assert.Equal(t, "Users", user.TableName())
}

func TestNewUserModelFromDomain(t *testing.T) {
	mockHasher := &MockHasher{}
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 15).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)
	err = password.Hash()
	assert.NoError(t, err)

	address := &domain.Address{
		Address:       "Rua Teste",
		AddressNumber: "123",
		Neighborhood:  "Centro",
		City:          "São Paulo",
		Country:       "Brasil",
		ZipCode:       "01234-567",
	}

	domainUser := domain.User{
		ID:       1,
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Role:     "admin",
		Password: password,
		Address:  address,
		Active:   true,
	}

	userModel := NewUserModelFromDomain(domainUser)

	assert.Equal(t, domainUser.ID, userModel.ID)
	assert.Equal(t, domainUser.Name, userModel.Name)
	assert.Equal(t, domainUser.Email, userModel.Email)
	assert.Equal(t, domainUser.Contact, userModel.Contact)
	assert.Equal(t, domainUser.Role, userModel.Role)
	assert.Equal(t, domainUser.Password.GetHashed(), userModel.Password)
	assert.Equal(t, domainUser.Active, userModel.Active)
	assert.Equal(t, domainUser.Address.Address, userModel.Address.Address)
	assert.Equal(t, domainUser.Address.AddressNumber, userModel.Address.AddressNumber)
	assert.Equal(t, domainUser.Address.Neighborhood, userModel.Address.Neighborhood)
	assert.Equal(t, domainUser.Address.City, userModel.Address.City)
	assert.Equal(t, domainUser.Address.Country, userModel.Address.Country)
	assert.Equal(t, domainUser.Address.ZipCode, userModel.Address.ZipCode)

	mockHasher.AssertExpectations(t)
}

func TestUserModel_ToDomain(t *testing.T) {
	userModel := UserModel{
		ID:       1,
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Password: "hashed_password",
		Role:     "admin",
		Active:   true,
		Address: AddressModel{
			Address:       "Rua Teste",
			AddressNumber: "123",
			Neighborhood:  "Centro",
			City:          "São Paulo",
			Country:       "Brasil",
			ZipCode:       "01234-567",
		},
	}

	domainUser := userModel.ToDomain()

	assert.Equal(t, userModel.ID, domainUser.ID)
	assert.Equal(t, userModel.Name, domainUser.Name)
	assert.Equal(t, userModel.Email, domainUser.Email)
	assert.Equal(t, userModel.Contact, domainUser.Contact)
	assert.Equal(t, userModel.Role, domainUser.Role)
	assert.Equal(t, userModel.Password, domainUser.Password.GetHashed())
	assert.Equal(t, userModel.Active, domainUser.Active)
	assert.Equal(t, userModel.Address.Address, domainUser.Address.Address)
	assert.Equal(t, userModel.Address.AddressNumber, domainUser.Address.AddressNumber)
	assert.Equal(t, userModel.Address.Neighborhood, domainUser.Address.Neighborhood)
	assert.Equal(t, userModel.Address.City, domainUser.Address.City)
	assert.Equal(t, userModel.Address.Country, domainUser.Address.Country)
	assert.Equal(t, userModel.Address.ZipCode, domainUser.Address.ZipCode)
}

func TestNewUserModelFromDomain_WithNilAddress(t *testing.T) {
	mockHasher := &MockHasher{}
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 15).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)
	err = password.Hash()
	assert.NoError(t, err)

	domainUser := domain.User{
		ID:       1,
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Role:     "admin",
		Password: password,
		Address:  nil,
		Active:   true,
	}

	userModel := NewUserModelFromDomain(domainUser)

	assert.Equal(t, domainUser.ID, userModel.ID)
	assert.Equal(t, domainUser.Name, userModel.Name)
	assert.Equal(t, domainUser.Email, userModel.Email)
	assert.Equal(t, domainUser.Contact, userModel.Contact)
	assert.Equal(t, domainUser.Role, userModel.Role)
	assert.Equal(t, domainUser.Password.GetHashed(), userModel.Password)
	assert.Equal(t, domainUser.Active, userModel.Active)
	assert.Equal(t, "", userModel.Address.Address)
	assert.Equal(t, "", userModel.Address.AddressNumber)
	assert.Equal(t, "", userModel.Address.Neighborhood)
	assert.Equal(t, "", userModel.Address.City)
	assert.Equal(t, "", userModel.Address.Country)
	assert.Equal(t, "", userModel.Address.ZipCode)

	mockHasher.AssertExpectations(t)
}

func TestUserModel_ToDomain_WithEmptyAddress(t *testing.T) {
	userModel := UserModel{
		ID:       1,
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Password: "hashed_password",
		Role:     "admin",
		Active:   true,
		Address:  AddressModel{},
	}

	domainUser := userModel.ToDomain()

	assert.Equal(t, userModel.ID, domainUser.ID)
	assert.Equal(t, userModel.Name, domainUser.Name)
	assert.Equal(t, userModel.Email, domainUser.Email)
	assert.Equal(t, userModel.Contact, domainUser.Contact)
	assert.Equal(t, userModel.Role, domainUser.Role)
	assert.Equal(t, userModel.Password, domainUser.Password.GetHashed())
	assert.Equal(t, userModel.Active, domainUser.Active)
	assert.NotNil(t, domainUser.Address)
	assert.Equal(t, "", domainUser.Address.Address)
	assert.Equal(t, "", domainUser.Address.AddressNumber)
	assert.Equal(t, "", domainUser.Address.Neighborhood)
	assert.Equal(t, "", domainUser.Address.City)
	assert.Equal(t, "", domainUser.Address.Country)
	assert.Equal(t, "", domainUser.Address.ZipCode)
}

func TestNewUserModelFromDomain_WithPasswordHashError(t *testing.T) {
	mockHasher := &MockHasher{}
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 15).Return([]byte(""), errors.New("hash error"))

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)
	err = password.Hash()
	assert.Error(t, err)

	domainUser := domain.User{
		ID:       1,
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Role:     "admin",
		Password: password,
		Address:  nil,
		Active:   true,
	}

	userModel := NewUserModelFromDomain(domainUser)

	assert.Equal(t, "", userModel.Password)

	mockHasher.AssertExpectations(t)
}
