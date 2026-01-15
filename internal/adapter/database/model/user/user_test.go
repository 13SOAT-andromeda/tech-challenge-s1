package user

import (
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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
	user := &Model{}
	assert.Equal(t, "User", user.TableName())
}

func TestNewUserModelFromDomain(t *testing.T) {
	mockHasher := &MockHasher{}
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

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
		ID:        1,
		Name:      "João Silva",
		Email:     "joao@test.com",
		Contact:   "11999999999",
		Role:      "admin",
		Password:  password,
		Address:   address,
		DeletedAt: nil,
	}

	userModel := &Model{}
	userModel.FromDomain(&domainUser)

	assert.Equal(t, domainUser.ID, userModel.ID)
	assert.Equal(t, domainUser.Name, userModel.Name)
	assert.Equal(t, domainUser.Email, userModel.Email)
	assert.Equal(t, domainUser.Contact, userModel.Contact)
	assert.Equal(t, domainUser.Role, userModel.Role)
	assert.Equal(t, domainUser.Password.GetHashed(), userModel.Password)
	assert.False(t, userModel.DeletedAt.Valid)
	assert.Equal(t, domainUser.Address.Address, userModel.Address.Address)
	assert.Equal(t, domainUser.Address.AddressNumber, userModel.Address.AddressNumber)
	assert.Equal(t, domainUser.Address.Neighborhood, userModel.Address.Neighborhood)
	assert.Equal(t, domainUser.Address.City, userModel.Address.City)
	assert.Equal(t, domainUser.Address.Country, userModel.Address.Country)
	assert.Equal(t, domainUser.Address.ZipCode, userModel.Address.ZipCode)

	mockHasher.AssertExpectations(t)
}

func TestUserModel_ToDomain(t *testing.T) {
	userModel := Model{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Password: "hashed_password",
		Role:     "admin",
		Address: &address.Model{
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
	assert.Nil(t, domainUser.DeletedAt)
	assert.Equal(t, userModel.Address.Address, domainUser.Address.Address)
	assert.Equal(t, userModel.Address.AddressNumber, domainUser.Address.AddressNumber)
	assert.Equal(t, userModel.Address.Neighborhood, domainUser.Address.Neighborhood)
	assert.Equal(t, userModel.Address.City, domainUser.Address.City)
	assert.Equal(t, userModel.Address.Country, domainUser.Address.Country)
	assert.Equal(t, userModel.Address.ZipCode, domainUser.Address.ZipCode)
}

func TestNewUserModelFromDomain_WithNilAddress(t *testing.T) {
	mockHasher := &MockHasher{}
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)
	err = password.Hash()
	assert.NoError(t, err)

	domainUser := domain.User{
		ID:        1,
		Name:      "João Silva",
		Email:     "joao@test.com",
		Contact:   "11999999999",
		Role:      "admin",
		Password:  password,
		Address:   nil,
		DeletedAt: nil,
	}

	userModel := &Model{}
	userModel.FromDomain(&domainUser)

	assert.Equal(t, domainUser.ID, userModel.ID)
	assert.Equal(t, domainUser.Name, userModel.Name)
	assert.Equal(t, domainUser.Email, userModel.Email)
	assert.Equal(t, domainUser.Contact, userModel.Contact)
	assert.Equal(t, domainUser.Role, userModel.Role)
	assert.Equal(t, domainUser.Password.GetHashed(), userModel.Password)
	assert.False(t, userModel.DeletedAt.Valid)
	assert.Equal(t, "", userModel.Address.Address)
	assert.Equal(t, "", userModel.Address.AddressNumber)
	assert.Equal(t, "", userModel.Address.Neighborhood)
	assert.Equal(t, "", userModel.Address.City)
	assert.Equal(t, "", userModel.Address.Country)
	assert.Equal(t, "", userModel.Address.ZipCode)

	mockHasher.AssertExpectations(t)
}

func TestUserModel_ToDomain_WithEmptyAddress(t *testing.T) {
	userModel := Model{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Password: "hashed_password",
		Role:     "admin",
		Address:  &address.Model{},
	}

	domainUser := userModel.ToDomain()

	assert.Equal(t, userModel.ID, domainUser.ID)
	assert.Equal(t, userModel.Name, domainUser.Name)
	assert.Equal(t, userModel.Email, domainUser.Email)
	assert.Equal(t, userModel.Contact, domainUser.Contact)
	assert.Equal(t, userModel.Role, domainUser.Role)
	assert.Equal(t, userModel.Password, domainUser.Password.GetHashed())
	assert.Nil(t, domainUser.DeletedAt)
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
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte(""), errors.New("hash error"))

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)
	err = password.Hash()
	assert.Error(t, err)

	domainUser := domain.User{
		ID:        1,
		Name:      "João Silva",
		Email:     "joao@test.com",
		Contact:   "11999999999",
		Role:      "admin",
		Password:  password,
		Address:   nil,
		DeletedAt: nil,
	}

	userModel := &Model{}
	userModel.FromDomain(&domainUser)

	assert.Equal(t, "", userModel.Password)

	mockHasher.AssertExpectations(t)
}
