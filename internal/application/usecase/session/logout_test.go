package session

import (
	"context"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewLogoutUseCase(t *testing.T) {
	sessionService := &MockSessionService{}

	useCase := NewLogoutUseCase(sessionService)

	assert.NotNil(t, useCase)
	assert.Implements(t, (*LogoutUseCase)(nil), useCase)
}

func TestLogoutUseCase_Execute_Success(t *testing.T) {
	sessionService := &MockSessionService{}

	useCase := NewLogoutUseCase(sessionService)

	// Setup mocks: direct delete by refresh token
	sessionService.On("DeleteByRefreshToken", mock.Anything, "refresh-token").Return(nil)

	// Execute
	input := LogoutInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Success)
	assert.Equal(t, "Logged out successfully", output.Message)

	sessionService.AssertExpectations(t)
}

func TestLogoutUseCase_Execute_InvalidRefreshToken(t *testing.T) {
	sessionService := &MockSessionService{}

	useCase := NewLogoutUseCase(sessionService)

	// Setup mocks - session not found
	sessionService.On("DeleteByRefreshToken", mock.Anything, "invalid-token").Return(services.ErrSessionNotFound)

	// Execute
	input := LogoutInput{
		RefreshToken: "invalid-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrSessionInvalid, err)

	sessionService.AssertExpectations(t)
}

func TestLogoutUseCase_Execute_SessionNotFound(t *testing.T) {
	sessionService := &MockSessionService{}

	useCase := NewLogoutUseCase(sessionService)

	// Setup mocks - session not found
	sessionService.On("DeleteByRefreshToken", mock.Anything, "nonexistent-token").Return(services.ErrSessionNotFound)

	// Execute
	input := LogoutInput{
		RefreshToken: "nonexistent-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrSessionInvalid, err)

	sessionService.AssertExpectations(t)
}

func TestLogoutUseCase_Execute_UpdateError(t *testing.T) {
	sessionService := &MockSessionService{}

	useCase := NewLogoutUseCase(sessionService)

	// Setup mocks
	sessionService.On("DeleteByRefreshToken", mock.Anything, "refresh-token").Return(assert.AnError)

	// Execute
	input := LogoutInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)

	sessionService.AssertExpectations(t)
}

func TestLogoutUseCase_Execute_EmptyRefreshToken(t *testing.T) {
	sessionService := &MockSessionService{}

	useCase := NewLogoutUseCase(sessionService)

	// Setup mocks - empty token validation fails
	sessionService.On("DeleteByRefreshToken", mock.Anything, "").Return(services.ErrSessionRefreshTokenEmpty)

	// Execute
	input := LogoutInput{
		RefreshToken: "",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrSessionInvalid, err)

	sessionService.AssertExpectations(t)
}

// AlreadyLoggedOut scenario is no longer applicable with direct delete
