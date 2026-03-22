package services

import (
	"context"
	"errors"
	"testing"

	employeeModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/employee"
	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	userModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestUserService(mockRepo *mocks.MockUserRepository, mockPersonRepo *mocks.MockPersonRepository, mockEmployeeRepo *mocks.MockEmployeeRepository) *userService {
	return NewUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)
}

func makeUserWithPerson(name, email, role string, password *domain.Password) domain.User {
	return domain.User{
		Password: password,
		Role:     role,
		Person: &domain.Person{
			Name:    name,
			Email:   email,
			Contact: "11999999999",
			Address: &domain.Address{
				Address:       "Rua Teste",
				AddressNumber: "123",
				City:          "São Paulo",
				Neighborhood:  "Centro",
				Country:       "Brasil",
				ZipCode:       "01234-567",
			},
		},
	}
}

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	mockHasher := new(encryption.MockHasher)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := makeUserWithPerson("João Silva", "joao@test.com", "administrator", password)

	createdPerson := personModel.Model{}
	createdPerson.ID = 1
	createdPerson.Name = "João Silva"
	createdPerson.Email = "joao@test.com"

	mockUM := userModel.Model{}
	mockUM.Role = "administrator"
	mockUM.PersonID = 1
	mockUM.Person = createdPerson

	mockRepo.On("GetByEmail", ctx, "joao@test.com").Return(nil, nil)
	mockPersonRepo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockEmployeeRepo.On("Create", ctx, mock.AnythingOfType("*employee.Model")).Return(&employeeModel.Model{}, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(&mockUM, nil)

	result, err := service.Create(ctx, inputUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "administrator", result.Role)
	assert.NotNil(t, result.Person)
	assert.Equal(t, "João Silva", result.Person.Name)
	assert.Equal(t, "joao@test.com", result.Person.Email)
	assert.Nil(t, result.DeletedAt)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_Create_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	password, err := domain.NewPassword("TestPass123!", &encryption.MockHasher{})
	assert.NoError(t, err)

	inputUser := makeUserWithPerson("João Silva", "joao@test.com", "administrator", password)

	existingUM := userModel.Model{}
	existingUM.ID = 1

	mockRepo.On("GetByEmail", ctx, "joao@test.com").Return(&existingUM, nil)

	result, err := service.Create(ctx, inputUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "email já existe")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	mockHasher := new(encryption.MockHasher)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := makeUserWithPerson("João Silva", "joao@test.com", "administrator", password)

	expectedError := errors.New("database connection error")

	createdPerson := personModel.Model{}
	createdPerson.ID = 1

	mockRepo.On("GetByEmail", ctx, "joao@test.com").Return(nil, nil)
	mockPersonRepo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockEmployeeRepo.On("Create", ctx, mock.AnythingOfType("*employee.Model")).Return(&employeeModel.Model{}, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(nil, expectedError)

	result, err := service.Create(ctx, inputUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	existingUM := userModel.Model{}
	existingUM.ID = 1
	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(&existingUM, nil)

	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_EmailError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, errors.New("error on consult email"))

	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_PasswordWeak(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, nil)

	err := service.CreateAdminUser(ctx, "admin@admin.com", "weakpass")

	assert.Error(t, err)
	assert.EqualError(t, err, "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	createdPerson := personModel.Model{}
	createdPerson.ID = 1

	createdUM := userModel.Model{}

	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, nil)
	mockPersonRepo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(&createdUM, nil)
	createdEM := employeeModel.Model{}
	mockEmployeeRepo.On("Create", ctx, mock.AnythingOfType("*employee.Model")).Return(&createdEM, nil)

	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
}

func TestUserService_CreateAdminUser_CreateError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	createdPerson := personModel.Model{}
	createdPerson.ID = 1

	mockRepo.On("GetByEmail", ctx, "admin@admin.com").Return(nil, nil)
	mockPersonRepo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(nil, errors.New("error on create admin user"))

	err := service.CreateAdminUser(ctx, "admin@admin.com", "ValidPass123!")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
}

func TestUserService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	userID := uint(1)

	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = "João Silva"
	pm.Email = "joao@test.com"

	expectedUser := &userModel.Model{
		Role:     "administrator",
		PersonID: 1,
		Person:   pm,
	}
	expectedUser.ID = 1

	mockRepo.On("FindByID", ctx, userID).Return(expectedUser, nil)

	result, err := service.GetByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Person)
	assert.Equal(t, "João Silva", result.Person.Name)
	assert.Equal(t, "joao@test.com", result.Person.Email)
	assert.Equal(t, "administrator", result.Role)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	userID := uint(999)

	mockRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found"))

	result, err := service.GetByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "user not found")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Search_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"name":    "João",
		"email":   "joao@test.com",
		"contact": "11999999999",
	}

	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = "João Silva"
	pm.Email = "joao@test.com"

	expectedUsers := []userModel.Model{
		{
			Role:     "administrator",
			PersonID: 1,
			Person:   pm,
		},
	}
	expectedUsers[0].ID = 1

	mockRepo.On("Search", ctx, ports.UserSearch{
		Name:    "João",
		Email:   "joao@test.com",
		Contact: "11999999999",
	}).Return(expectedUsers)

	result, err := service.Search(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	assert.Equal(t, "administrator", (*result)[0].Role)
	assert.NotNil(t, (*result)[0].Person)
	assert.Equal(t, "João Silva", (*result)[0].Person.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Search_WithPartialParams(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	params := map[string]interface{}{
		"name": "João",
	}

	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = "João Silva"
	pm.Email = "joao@test.com"

	expectedUsers := []userModel.Model{
		{
			Role:     "administrator",
			PersonID: 1,
			Person:   pm,
		},
	}
	expectedUsers[0].ID = 1

	mockRepo.On("Search", ctx, ports.UserSearch{
		Name:    "João",
		Email:   "",
		Contact: "",
	}).Return(expectedUsers)

	result, err := service.Search(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	assert.Equal(t, "João Silva", (*result)[0].Person.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	userID := uint(1)

	pm := personModel.Model{}
	pm.ID = 1
	pm.Name = "João Silva"
	pm.Email = "joao@test.com"
	pm.Contact = "11999999999"

	existingUser := &userModel.Model{
		Role:     "administrator",
		PersonID: 1,
		Person:   pm,
	}
	existingUser.ID = 1

	updatedUser := domain.User{
		ID:   userID,
		Role: "administrator",
		Person: &domain.Person{
			Name:    "João Silva Santos",
			Email:   "joao@test.com",
			Contact: "11988888888",
		},
	}

	mockRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	mockPersonRepo.On("Update", ctx, mock.AnythingOfType("*person.Model")).Return(nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*user.Model")).Return(nil)

	result, err := service.Update(ctx, updatedUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	userID := uint(999)

	updatedUser := domain.User{
		ID:   userID,
		Role: "administrator",
	}

	mockRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found"))

	result, err := service.Update(ctx, updatedUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "user not found")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_PasswordUpdateNotAllowed(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	userID := uint(1)

	existingUser := &userModel.Model{
		Role:     "administrator",
		PersonID: 1,
	}
	existingUser.ID = 1

	password, err := domain.NewPassword("NewPass123!", &encryption.MockHasher{})
	assert.NoError(t, err)

	updatedUser := domain.User{
		ID:       userID,
		Role:     "administrator",
		Password: password,
	}

	mockRepo.On("FindByID", ctx, userID).Return(existingUser, nil)

	result, err := service.Update(ctx, updatedUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "senha de usuário não pode ser atualizada")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	userID := uint(1)

	mockRepo.On("Delete", ctx, userID).Return(nil)

	err := service.Delete(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()
	userID := uint(1)
	expectedError := errors.New("database error")

	mockRepo.On("Delete", ctx, userID).Return(expectedError)

	err := service.Delete(ctx, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "ocorreu um erro ao excluir o usuário")
	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_WithNilPerson(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	inputUser := domain.User{
		Role:   "administrator",
		Person: nil,
	}

	result, err := service.Create(ctx, inputUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "person data is needed")
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}

func TestUserService_Create_AsEmployee_WithPosition(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	mockHasher := new(encryption.MockHasher)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := makeUserWithPerson("Carlos Mecânico", "carlos@test.com", "mechanic", password)
	inputUser.Employee = &domain.Employee{Position: "Mecânico Sênior"}

	createdPerson := personModel.Model{}
	createdPerson.ID = 2
	createdPerson.Name = "Carlos Mecânico"
	createdPerson.Email = "carlos@test.com"

	createdEmployee := employeeModel.Model{}
	createdEmployee.ID = 1
	createdEmployee.Position = "Mecânico Sênior"
	createdEmployee.PersonID = 2

	mockUM := userModel.Model{}
	mockUM.Role = "mechanic"
	mockUM.PersonID = 2
	mockUM.Person = createdPerson

	mockRepo.On("GetByEmail", ctx, "carlos@test.com").Return(nil, nil)
	mockPersonRepo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockEmployeeRepo.On("Create", ctx, mock.MatchedBy(func(e *employeeModel.Model) bool {
		return e.Position == "Mecânico Sênior"
	})).Return(&createdEmployee, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(&mockUM, nil)

	result, err := service.Create(ctx, inputUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mechanic", result.Role)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_Create_AsCustomer_DoesNotCreateEmployee(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	mockHasher := new(encryption.MockHasher)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := makeUserWithPerson("Ana Cliente", "ana@test.com", "customer", password)

	createdPerson := personModel.Model{}
	createdPerson.ID = 3
	createdPerson.Name = "Ana Cliente"
	createdPerson.Email = "ana@test.com"

	mockUM := userModel.Model{}
	mockUM.Role = "customer"
	mockUM.PersonID = 3
	mockUM.Person = createdPerson

	mockRepo.On("GetByEmail", ctx, "ana@test.com").Return(nil, nil)
	mockPersonRepo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.Model")).Return(&mockUM, nil)

	result, err := service.Create(ctx, inputUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "customer", result.Role)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_EmployeeRepoError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPersonRepo := new(mocks.MockPersonRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	mockHasher := new(encryption.MockHasher)
	service := newTestUserService(mockRepo, mockPersonRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), 10).Return([]byte("hashed_password"), nil)

	password, err := domain.NewPassword("TestPass123!", mockHasher)
	assert.NoError(t, err)

	inputUser := makeUserWithPerson("Roberto Atendente", "roberto@test.com", "attendant", password)
	inputUser.Employee = &domain.Employee{Position: "Atendente"}

	createdPerson := personModel.Model{}
	createdPerson.ID = 4

	expectedError := errors.New("erro ao criar employee")

	mockRepo.On("GetByEmail", ctx, "roberto@test.com").Return(nil, nil)
	mockPersonRepo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(&createdPerson, nil)
	mockEmployeeRepo.On("Create", ctx, mock.AnythingOfType("*employee.Model")).Return(nil, expectedError)

	result, err := service.Create(ctx, inputUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
	mockPersonRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}
