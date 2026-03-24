package services

import (
	"context"
	"errors"
	"testing"

	employeeModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/employee"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func newTestEmployeeService(repo *mocks.MockEmployeeRepository) *employeeService {
	return NewEmployeeService(repo).(*employeeService)
}

func makeEmployeeModel(id, personID uint, position string) *employeeModel.Model {
	m := &employeeModel.Model{}
	m.ID = id
	m.PersonID = personID
	m.Position = position
	return m
}

func makeEmployeeDomain(id, personID uint, position string) domain.Employee {
	return domain.Employee{
		ID:       id,
		PersonID: personID,
		Position: position,
	}
}

func TestEmployeeService_Create_Success(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	input := makeEmployeeDomain(0, 1, "mechanic")
	returned := makeEmployeeModel(10, 1, "mechanic")

	repo.On("Create", ctx, mock.AnythingOfType("*employee.Model")).Return(returned, nil)

	result, err := svc.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, "mechanic", result.Position)
	assert.Equal(t, uint(1), result.PersonID)
	repo.AssertExpectations(t)
}

func TestEmployeeService_Create_RepoError(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	input := makeEmployeeDomain(0, 1, "mechanic")
	repoErr := errors.New("db error")

	repo.On("Create", ctx, mock.AnythingOfType("*employee.Model")).Return(nil, repoErr)

	result, err := svc.Create(ctx, input)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestEmployeeService_GetByID_Success(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	returned := makeEmployeeModel(5, 2, "advisor")

	repo.On("FindByID", ctx, uint(5)).Return(returned, nil)

	result, err := svc.GetByID(ctx, 5)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(5), result.ID)
	assert.Equal(t, "advisor", result.Position)
	repo.AssertExpectations(t)
}

func TestEmployeeService_GetByID_NotFound(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	repoErr := errors.New("record not found")

	repo.On("FindByID", ctx, uint(99)).Return(nil, repoErr)

	result, err := svc.GetByID(ctx, 99)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestEmployeeService_GetByPersonID_Success(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	returned := makeEmployeeModel(7, 3, "manager")

	repo.On("GetByPersonID", ctx, uint(3)).Return(returned, nil)

	result, err := svc.GetByPersonID(ctx, 3)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(7), result.ID)
	assert.Equal(t, uint(3), result.PersonID)
	assert.Equal(t, "manager", result.Position)
	repo.AssertExpectations(t)
}

func TestEmployeeService_GetByPersonID_NotFound_ReturnsNil(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	repo.On("GetByPersonID", ctx, uint(99)).Return(nil, gorm.ErrRecordNotFound)

	result, err := svc.GetByPersonID(ctx, 99)

	assert.NoError(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestEmployeeService_GetByPersonID_RepoError(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	repoErr := errors.New("connection error")

	repo.On("GetByPersonID", ctx, uint(3)).Return(nil, repoErr)

	result, err := svc.GetByPersonID(ctx, 3)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestEmployeeService_Delete_Success(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, uint(4)).Return(nil)

	err := svc.Delete(ctx, 4)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestEmployeeService_Delete_RepoError(t *testing.T) {
	repo := new(mocks.MockEmployeeRepository)
	svc := newTestEmployeeService(repo)
	ctx := context.Background()

	repoErr := errors.New("delete failed")

	repo.On("Delete", ctx, uint(4)).Return(repoErr)

	err := svc.Delete(ctx, 4)

	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}
