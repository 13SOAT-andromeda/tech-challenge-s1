package user

import (
	"errors"
	"testing"

	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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

	domainUser := domain.User{
		ID:       1,
		Role:     "admin",
		Password: password,
		PersonID: 5,
		Person: &domain.Person{
			Name:     "João Silva",
			Email:    "joao@test.com",
			Contact:  "11999999999",
			Document: &domain.Document{Number: "59597559048"},
		},
		DeletedAt: nil,
	}

	userModel := &Model{}
	userModel.FromDomain(&domainUser)

	assert.Equal(t, domainUser.ID, userModel.ID)
	assert.Equal(t, domainUser.Role, userModel.Role)
	assert.Equal(t, domainUser.Password.GetHashed(), userModel.Password)
	assert.Equal(t, domainUser.PersonID, userModel.PersonID)
	assert.False(t, userModel.DeletedAt.Valid)

	mockHasher.AssertExpectations(t)
}

func TestUserModel_ToDomain(t *testing.T) {
	userModel := Model{
		Model: gorm.Model{
			ID: 1,
		},
		Password: "hashed_password",
		Role:     "admin",
		PersonID: 5,
		Person:   personModel.Model{},
	}
	userModel.Person.ID = 5
	userModel.Person.Name = "João Silva"
	userModel.Person.Email = "joao@test.com"
	userModel.Person.Contact = "11999999999"

	domainUser := userModel.ToDomain()

	assert.Equal(t, userModel.ID, domainUser.ID)
	assert.Equal(t, userModel.Role, domainUser.Role)
	assert.Equal(t, userModel.Password, domainUser.Password.GetHashed())
	assert.Nil(t, domainUser.DeletedAt)
	assert.NotNil(t, domainUser.Person)
	assert.Equal(t, "João Silva", domainUser.Person.Name)
	assert.Equal(t, "joao@test.com", domainUser.Person.Email)
}

func TestNewUserModelFromDomain_WithNilPerson(t *testing.T) {
	mockHasher := &MockHasher{}
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)
	err = password.Hash()
	assert.NoError(t, err)

	domainUser := domain.User{
		ID:        1,
		Role:      "admin",
		Password:  password,
		DeletedAt: nil,
	}

	userModel := &Model{}
	userModel.FromDomain(&domainUser)

	assert.Equal(t, domainUser.ID, userModel.ID)
	assert.Equal(t, domainUser.Role, userModel.Role)
	assert.Equal(t, domainUser.Password.GetHashed(), userModel.Password)
	assert.False(t, userModel.DeletedAt.Valid)

	mockHasher.AssertExpectations(t)
}

func TestUserModel_ToDomain_WithNilPerson(t *testing.T) {
	userModel := Model{
		Model: gorm.Model{
			ID: 1,
		},
		Password: "hashed_password",
		Role:     "admin",
		PersonID: 0,
	}

	domainUser := userModel.ToDomain()

	assert.Equal(t, userModel.ID, domainUser.ID)
	assert.Equal(t, userModel.Role, domainUser.Role)
	assert.Equal(t, userModel.Password, domainUser.Password.GetHashed())
	assert.Nil(t, domainUser.DeletedAt)
	assert.Nil(t, domainUser.Person)
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
		Role:      "admin",
		Password:  password,
		Person:    nil,
		DeletedAt: nil,
	}

	userModel := &Model{}
	userModel.FromDomain(&domainUser)

	assert.Equal(t, "", userModel.Password)

	mockHasher.AssertExpectations(t)
}
