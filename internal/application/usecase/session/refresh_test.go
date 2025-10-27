package session

import (
	"context"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRefreshUseCase(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewRefreshUseCase(userService, sessionService, jwtService, config)

	assert.NotNil(t, useCase)
	assert.Implements(t, (*RefreshUseCase)(nil), useCase)
}

func TestRefreshUseCase_Execute_Success(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewRefreshUseCase(userService, sessionService, jwtService, config)

	// Mock session
	session := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     true,
	}

	// Mock user
	user := &domain.User{
		ID:      1,
		Name:    "Test User",
		Email:   "test@example.com",
		Contact: "123456789",
		Role:    "user",
		Active:  true,
	}

	// Setup mocks
	sessionService.On("Validate", mock.Anything, "refresh-token").Return(session, nil)
	userService.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	// Execute
	input := RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.AccessToken)
	assert.Equal(t, "refresh-token", output.RefreshToken)
	assert.Equal(t, int64(900), output.ExpiresIn) // 15 minutes in seconds
	assert.Equal(t, "Bearer", output.TokenType)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestRefreshUseCase_Execute_InvalidRefreshToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewRefreshUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails
	sessionService.On("Validate", mock.Anything, "invalid-token").Return((*domain.Session)(nil), services.ErrSessionInvalid)

	// Execute
	input := RefreshInput{
		RefreshToken: "invalid-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrSessionInvalid, err)

	sessionService.AssertExpectations(t)
}

func TestRefreshUseCase_Execute_ExpiredSession(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewRefreshUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails due to expiration
	sessionService.On("Validate", mock.Anything, "refresh-token").Return((*domain.Session)(nil), services.ErrSessionInvalid)

	// Execute
	input := RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrSessionInvalid, err)

	sessionService.AssertExpectations(t)
}

func TestRefreshUseCase_Execute_InactiveSession(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewRefreshUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails due to inactive session
	sessionService.On("Validate", mock.Anything, "refresh-token").Return((*domain.Session)(nil), services.ErrSessionInvalid)

	// Execute
	input := RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrSessionInvalid, err)

	sessionService.AssertExpectations(t)
}

func TestRefreshUseCase_Execute_UserNotFound(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewRefreshUseCase(userService, sessionService, jwtService, config)

	// Mock session
	session := &domain.Session{
		ID:           1,
		UserID:       999, // Non-existent user
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     true,
	}

	// Setup mocks
	sessionService.On("Validate", mock.Anything, "refresh-token").Return(session, nil)
	userService.On("GetByID", mock.Anything, uint(999)).Return((*domain.User)(nil), services.ErrUserNotFound)

	// Execute
	input := RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrUserNotFound, err)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestRefreshUseCase_Execute_EmptyRefreshToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewRefreshUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails for empty token
	sessionService.On("Validate", mock.Anything, "").Return((*domain.Session)(nil), services.ErrSessionRefreshTokenEmpty)

	// Execute
	input := RefreshInput{
		RefreshToken: "",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrSessionInvalid, err)

	sessionService.AssertExpectations(t)
}
