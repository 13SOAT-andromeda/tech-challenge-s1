package services

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHasher é um mock para o hasher do Password
func TestUserService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockHasher := new(encryption.MockHasher)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	// Configurar o mock do hasher
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 15).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := domain.User{
		Name:     "João Silva",
		Email:    "joao@test.com",
		Contact:  "11999999999",
		Role:     "administrator",
		Password: password,
		Address: &domain.Address{
			Address:       "Rua Teste",
			AddressNumber: "123",
			City:          "São Paulo",
			Neighborhood:  "Centro",
			Country:       "Brasil",
			ZipCode:       "01234-567",
		},
		DeletedAt: nil,
	}

	mockModel := user.Model{}
	mockModel.FromDomain(&inputUser)

	// Configurar o mock
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(&mockModel, nil)
	mockRepo.On("GetByEmail", ctx, inputUser.Email).Return(nil, nil)

	// Act
	result, err := service.Create(ctx, inputUser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "João Silva", result.Name)
	assert.Equal(t, "joao@test.com", result.Email)
	assert.Equal(t, "administrator", result.Role)
	assert.Nil(t, result.DeletedAt)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_Create_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	password, err := domain.NewPassword("TestPass123!", &encryption.MockHasher{})
	assert.NoError(t, err)

	inputUser := domain.User{
		Name:      "João Silva",
		Email:     "joao@test.com",
		Contact:   "11999999999",
		Role:      "administrator",
		Password:  password,
		DeletedAt: nil,
	}

	mockModel := user.Model{}
	mockModel.FromDomain(&inputUser)

	// Configurar o mock para retornar que o email já existe
	mockRepo.On("GetByEmail", ctx, inputUser.Email).Return(&mockModel, nil)

	// Act
	result, err := service.Create(ctx, inputUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrUserEmailAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockHasher := new(encryption.MockHasher)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	// Configurar o mock do hasher
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 15).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := domain.User{
		Name:      "João Silva",
		Email:     "joao@test.com",
		Contact:   "11999999999",
		Role:      "administrator",
		Password:  password,
		DeletedAt: nil,
	}

	expectedError := errors.New("database connection error")

	// Configurar o mock
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(nil, expectedError)
	mockRepo.On("GetByEmail", ctx, inputUser.Email).Return(nil, nil)

	// Act
	result, err := service.Create(ctx, inputUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(&user.Model{Email: "admin@admin.com"}, nil)

	// Act
	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_EmailError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, errors.New("error on consult email"))

	// Act
	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_PasswordWeak(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, nil)

	// Act
	err := service.CreateAdminUser(ctx, "admin@admin.com", "weakpass")

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockModel := user.Model{}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(&mockModel, nil)
	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, nil)

	// Act
	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_CreateError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(nil, errors.New("error on create admin user"))
	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, nil)

	// Act
	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(1)

	expectedUser := &user.Model{
		Name:  "João Silva",
		Email: "joao@test.com",
		Role:  "administrator",
	}
	expectedUser.ID = 1

	mockRepo.On("FindByID", ctx, userID).Return(expectedUser, nil)

	// Act
	result, err := service.GetByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "João Silva", result.Name)
	assert.Equal(t, "joao@test.com", result.Email)
	assert.Equal(t, "administrator", result.Role)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(999)

	mockRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found"))

	// Act
	result, err := service.GetByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Search_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"name":    "João",
		"email":   "joao@test.com",
		"contact": "11999999999",
	}

	expectedUsers := []user.Model{
		{
			Name:  "João Silva",
			Email: "joao@test.com",
			Role:  "administrator",
		},
	}
	expectedUsers[0].ID = 1

	mockRepo.On("Search", ctx, ports.UserSearch{
		Name:    "João",
		Email:   "joao@test.com",
		Contact: "11999999999",
	}).Return(expectedUsers)

	// Act
	result, err := service.Search(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	assert.Equal(t, "João Silva", (*result)[0].Name)
	assert.Equal(t, "joao@test.com", (*result)[0].Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Search_WithPartialParams(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"name": "João",
	}

	expectedUsers := []user.Model{
		{
			Name:  "João Silva",
			Email: "joao@test.com",
			Role:  "administrator",
		},
	}
	expectedUsers[0].ID = 1

	mockRepo.On("Search", ctx, ports.UserSearch{
		Name:    "João",
		Email:   "",
		Contact: "",
	}).Return(expectedUsers)

	// Act
	result, err := service.Search(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	assert.Equal(t, "João Silva", (*result)[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(1)

	existingUser := &user.Model{
		Name:    "João Silva",
		Email:   "joao@test.com",
		Contact: "11999999999",
		Role:    "administrator",
	}
	existingUser.ID = 1

	updatedUser := domain.User{
		ID:        userID,
		Name:      "João Silva Santos",
		Email:     "joao@test.com",
		Contact:   "11988888888",
		Role:      "administrator",
		DeletedAt: nil,
	}

	// Configurar o mock
	mockRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*user.Model")).Return(nil)

	// Act
	result, err := service.Update(ctx, updatedUser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "João Silva Santos", result.Name)
	assert.Equal(t, "11988888888", result.Contact)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(999)

	updatedUser := domain.User{
		ID:   userID,
		Name: "João Silva Santos",
	}

	// Configurar o mock
	mockRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found"))

	// Act
	result, err := service.Update(ctx, updatedUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(2)

	existingUser := &user.Model{
		Name:     "João Silva",
		Email:    "new@test.com",
		Role:     "administrator",
		Password: "hashed_password",
	}

	oldDataUser := &user.Model{
		Name:     "João Silva",
		Email:    "old@test.com",
		Role:     "administrator",
		Password: "hashed_password",
	}
	oldDataUser.ID = userID

	updatedUser := domain.User{
		ID:        userID,
		Name:      "João Silva Santos",
		Email:     "new@test.com",
		Role:      "administrator",
		DeletedAt: nil,
	}

	// Configurar o mock
	mockRepo.On("FindByID", ctx, userID).Return(oldDataUser, nil)
	mockRepo.On("GetByEmail", ctx, updatedUser.Email).Return(existingUser, nil)

	// Act
	result, err := service.Update(ctx, updatedUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrUserEmailAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_PasswordUpdateNotAllowed(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(1)

	existingUser := &user.Model{
		Name:  "João Silva",
		Email: "joao@test.com",
		Role:  "administrator",
	}

	password, err := domain.NewPassword("NewPass123!", &encryption.MockHasher{})
	assert.NoError(t, err)

	updatedUser := domain.User{
		ID:       userID,
		Name:     "João Silva Santos",
		Email:    "joao@test.com",
		Password: password,
		Role:     "administrator",
	}

	// Configurar o mock
	mockRepo.On("FindByID", ctx, userID).Return(existingUser, nil)

	// Act
	result, err := service.Update(ctx, updatedUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrUserPasswordUpdateNotAllowed, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(1)

	// Configurar o mock
	mockRepo.On("Delete", ctx, userID).Return(nil)

	// Act
	err := service.Delete(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(1)
	expectedError := errors.New("database error")

	// Configurar o mock
	mockRepo.On("Delete", ctx, userID).Return(expectedError)

	// Act
	err := service.Delete(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrUserDelete, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_WithNilAddress(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockHasher := new(encryption.MockHasher)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	// Configurar o mock do hasher
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 15).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := domain.User{
		Name:      "João Silva",
		Email:     "joao@test.com",
		Contact:   "11999999999",
		Role:      "administrator",
		Password:  password,
		Address:   nil,
		DeletedAt: nil,
	}

	mockModel := user.Model{}
	mockModel.FromDomain(&inputUser)

	// Configurar o mock
	mockRepo.On("GetByEmail", ctx, "joao@test.com").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(&mockModel, nil)

	// Act
	result, err := service.Create(ctx, inputUser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "João Silva", result.Name)
	assert.Equal(t, "joao@test.com", result.Email)
	assert.NotNil(t, result.Address) // Deve ser inicializado como endereço vazio
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
