package session

import (
	"context"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewValidateUseCase(t *testing.T) {
	userService := &MockUserService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewValidateUseCase(userService, jwtService)

	assert.NotNil(t, useCase)
	assert.Implements(t, (*ValidateUseCase)(nil), useCase)
}

func TestValidateUseCase_Execute_ValidToken(t *testing.T) {
	userService := &MockUserService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewValidateUseCase(userService, jwtService)

	// Generate a valid token
	userID := uint(1)
	email := "test@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role)
	assert.NoError(t, err)

	// Mock user
	user := &domain.User{
		ID:      userID,
		Name:    "Test User",
		Email:   email,
		Contact: "123456789",
		Role:    role,
		Active:  true,
	}

	// Setup mocks
	userService.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Execute
	input := ValidateInput{
		Token: token,
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Valid)
	assert.NotNil(t, output.User)
	assert.Equal(t, user.ID, output.User.ID)
	assert.Equal(t, user.Email, output.User.Email)

	userService.AssertExpectations(t)
}

func TestValidateUseCase_Execute_InvalidToken(t *testing.T) {
	userService := &MockUserService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewValidateUseCase(userService, jwtService)

	// Execute with invalid token
	input := ValidateInput{
		Token: "invalid.token.here",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.False(t, output.Valid)
	assert.Nil(t, output.User)

	userService.AssertExpectations(t)
}

func TestValidateUseCase_Execute_EmptyToken(t *testing.T) {
	userService := &MockUserService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewValidateUseCase(userService, jwtService)

	// Execute with empty token
	input := ValidateInput{
		Token: "",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.False(t, output.Valid)
	assert.Nil(t, output.User)

	userService.AssertExpectations(t)
}

func TestValidateUseCase_Execute_UserNotFound(t *testing.T) {
	userService := &MockUserService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewValidateUseCase(userService, jwtService)

	// Generate a valid token
	userID := uint(999) // Non-existent user
	email := "nonexistent@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role)
	assert.NoError(t, err)

	// Setup mocks - user not found
	userService.On("GetByID", mock.Anything, userID).Return((*domain.User)(nil), assert.AnError)

	// Execute
	input := ValidateInput{
		Token: token,
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.False(t, output.Valid)
	assert.Nil(t, output.User)

	userService.AssertExpectations(t)
}

func TestValidateUseCase_Execute_ExpiredToken(t *testing.T) {
	userService := &MockUserService{}
	// Create JWT service with very short expiry
	jwtService := jwt.NewService("test-secret", 1*time.Millisecond, 7*24*time.Hour)

	useCase := NewValidateUseCase(userService, jwtService)

	// Generate a token that expires quickly
	userID := uint(1)
	email := "test@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role)
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Execute
	input := ValidateInput{
		Token: token,
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.False(t, output.Valid)
	assert.Nil(t, output.User)

	userService.AssertExpectations(t)
}
