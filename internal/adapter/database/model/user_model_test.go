package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserModelInitialization(t *testing.T) {
	u := Model{
		Name:     "User A",
		Email:    "user.a@example.com",
		Contact:  "11987654321",
		Address:  "Rua dos Testes, 123",
		Password: "password",
		Role:     "user",
	}

	assert.NotNil(t, u)
	assert.Equal(t, "User A", u.Name)
	assert.Equal(t, "user.a@example.com", u.Email)
	assert.Equal(t, "11987654321", u.Contact)
	assert.Equal(t, "Rua dos Testes, 123", u.Address)
	assert.Equal(t, "password", u.Password)
	assert.Equal(t, "user", u.Role)
}

func TestUserModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	modelUser := &Model{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		Name:     "User A",
		Email:    "user.a@example.com",
		Contact:  "11987654321",
		Address:  "Rua dos Testes, 123",
		Password: "password",
		Role:     "user",
		Sessions: []SessionModel{},
	}

	domainUser := modelUser.ToDomain()

	assert.Equal(t, modelUser.ID, domainUser.ID)
	assert.Equal(t, modelUser.Name, domainUser.Name)
	assert.Equal(t, modelUser.Email, domainUser.Email)
	assert.Equal(t, modelUser.Contact, domainUser.Contact)
	assert.Equal(t, modelUser.Address, domainUser.Address)
	assert.Equal(t, modelUser.Password, domainUser.Password)
	assert.Equal(t, modelUser.Role, domainUser.Role)
	assert.Equal(t, modelUser.CreatedAt, domainUser.CreatedAt)
	assert.Equal(t, modelUser.UpdatedAt, domainUser.UpdatedAt)

}
