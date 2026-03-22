package services

import (
	"context"
	"errors"
	"testing"

	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestPersonService(repo *mocks.MockPersonRepository) *personService {
	return NewPersonService(repo).(*personService)
}

func makePersonModel(id uint, name, email string) *personModel.Model {
	m := &personModel.Model{}
	m.ID = id
	m.Name = name
	m.Email = email
	m.Contact = "11999999999"
	return m
}

func makePersonDomain(id uint, name, email string) domain.Person {
	return domain.Person{
		ID:      id,
		Name:    name,
		Email:   email,
		Contact: "11999999999",
	}
}

func TestPersonService_Create_Success(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	input := makePersonDomain(0, "Ana Lima", "ana@test.com")
	returned := makePersonModel(1, "Ana Lima", "ana@test.com")

	repo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(returned, nil)

	result, err := svc.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Ana Lima", result.Name)
	assert.Equal(t, "ana@test.com", result.Email)
	repo.AssertExpectations(t)
}

func TestPersonService_Create_RepoError(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	input := makePersonDomain(0, "Ana Lima", "ana@test.com")
	repoErr := errors.New("db error")

	repo.On("Create", ctx, mock.AnythingOfType("*person.Model")).Return(nil, repoErr)

	result, err := svc.Create(ctx, input)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestPersonService_GetByID_Success(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	returned := makePersonModel(5, "Carlos Silva", "carlos@test.com")

	repo.On("FindByID", ctx, uint(5)).Return(returned, nil)

	result, err := svc.GetByID(ctx, 5)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(5), result.ID)
	assert.Equal(t, "Carlos Silva", result.Name)
	repo.AssertExpectations(t)
}

func TestPersonService_GetByID_NotFound(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	repoErr := errors.New("record not found")

	repo.On("FindByID", ctx, uint(99)).Return(nil, repoErr)

	result, err := svc.GetByID(ctx, 99)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestPersonService_Update_Success(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	input := makePersonDomain(3, "Maria Santos", "maria@test.com")

	repo.On("Update", ctx, mock.AnythingOfType("*person.Model")).Return(nil)

	result, err := svc.Update(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(3), result.ID)
	assert.Equal(t, "Maria Santos", result.Name)
	repo.AssertExpectations(t)
}

func TestPersonService_Update_RepoError(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	input := makePersonDomain(3, "Maria Santos", "maria@test.com")
	repoErr := errors.New("update failed")

	repo.On("Update", ctx, mock.AnythingOfType("*person.Model")).Return(repoErr)

	result, err := svc.Update(ctx, input)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestPersonService_Delete_Success(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, uint(7)).Return(nil)

	err := svc.Delete(ctx, 7)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPersonService_Delete_RepoError(t *testing.T) {
	repo := new(mocks.MockPersonRepository)
	svc := newTestPersonService(repo)
	ctx := context.Background()

	repoErr := errors.New("delete failed")

	repo.On("Delete", ctx, uint(7)).Return(repoErr)

	err := svc.Delete(ctx, 7)

	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}
