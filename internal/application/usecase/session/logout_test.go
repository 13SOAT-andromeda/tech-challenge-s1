package session

import (
	"context"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
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

	// Mock session
	session := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     true,
	}

	// Mock updated session (after deactivation)
	updatedSession := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     false, // Deactivated
	}

	// Setup mocks
	sessionService.On("GetByRefreshToken", mock.Anything, "refresh-token").Return(session, nil)
	sessionService.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Session) bool {
		return s.ID == 1 && !s.IsActive // Check that session is deactivated
	})).Return(updatedSession, nil)

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
	sessionService.On("GetByRefreshToken", mock.Anything, "invalid-token").Return((*domain.Session)(nil), services.ErrSessionNotFound)

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
	sessionService.On("GetByRefreshToken", mock.Anything, "nonexistent-token").Return((*domain.Session)(nil), services.ErrSessionNotFound)

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

	// Mock session
	session := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     true,
	}

	// Setup mocks
	sessionService.On("GetByRefreshToken", mock.Anything, "refresh-token").Return(session, nil)
	sessionService.On("Update", mock.Anything, mock.Anything).Return((*domain.Session)(nil), assert.AnError)

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
	sessionService.On("GetByRefreshToken", mock.Anything, "").Return((*domain.Session)(nil), services.ErrSessionRefreshTokenEmpty)

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

func TestLogoutUseCase_Execute_AlreadyLoggedOut(t *testing.T) {
	sessionService := &MockSessionService{}

	useCase := NewLogoutUseCase(sessionService)

	// Mock already inactive session
	session := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     false, // Already inactive
	}

	// Mock updated session
	updatedSession := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     false,
	}

	// Setup mocks
	sessionService.On("GetByRefreshToken", mock.Anything, "refresh-token").Return(session, nil)
	sessionService.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Session) bool {
		return s.ID == 1 && !s.IsActive
	})).Return(updatedSession, nil)

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
