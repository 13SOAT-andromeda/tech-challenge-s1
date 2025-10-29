package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	secret := "test-secret"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour

	service := NewService(secret, accessExpiry, refreshExpiry)

	assert.NotNil(t, service)
	assert.Equal(t, []byte(secret), service.secret)
	assert.Equal(t, accessExpiry, service.accessTokenExpiry)
	assert.Equal(t, refreshExpiry, service.refreshTokenExpiry)
}

func TestService_GenerateAccessToken(t *testing.T) {
	service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	userID := uint(1)
	email := "test@example.com"
	role := "user"
	sessionID := uint(123)

	token, err := service.GenerateAccessToken(userID, email, role, sessionID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, sessionID, claims.SessionID)
	assert.Equal(t, "tech-challenge-s1", claims.Issuer)
	assert.Equal(t, "1", claims.Subject)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Second)))
}

func TestService_GenerateRefreshToken(t *testing.T) {
	service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	userID := uint(1)

	token, err := service.GenerateRefreshToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "tech-challenge-s1", claims.Issuer)
	assert.Equal(t, "1", claims.Subject)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestService_ValidateToken(t *testing.T) {
	service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name: "valid token",
			token: func() string {
				token, _ := service.GenerateAccessToken(1, "test@example.com", "user", 123)
				return token
			}(),
			expectError: false,
		},
		{
			name:        "invalid token",
			token:       "invalid.token.here",
			expectError: true,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "malformed token",
			token:       "not-a-jwt-token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestService_ExtractUserIDFromToken(t *testing.T) {
	service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	userID := uint(123)
	sessionID := uint(456)
	token, err := service.GenerateAccessToken(userID, "test@example.com", "user", sessionID)
	assert.NoError(t, err)

	extractedUserID, err := service.ExtractUserIDFromToken(token)

	assert.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}

func TestService_IsTokenExpired(t *testing.T) {
	service := NewService("test-secret", 1*time.Millisecond, 7*24*time.Hour)

	// Generate a token that expires very quickly
	token, err := service.GenerateAccessToken(1, "test@example.com", "user", 123)
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	assert.True(t, service.IsTokenExpired(token))
}

func TestService_RefreshAccessToken(t *testing.T) {
	service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	userID := uint(1)
	email := "test@example.com"
	role := "user"
	sessionID := uint(789)

	refreshToken, err := service.GenerateRefreshToken(userID)
	assert.NoError(t, err)

	newAccessToken, err := service.RefreshAccessToken(refreshToken, email, role, sessionID)

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)

	// Validate the new access token
	claims, err := service.ValidateToken(newAccessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, sessionID, claims.SessionID)
}

func TestService_EdgeCases(t *testing.T) {
	t.Run("zero user ID", func(t *testing.T) {
		service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

		token, err := service.GenerateAccessToken(0, "test@example.com", "user", 123)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(0), claims.UserID)
	})

	t.Run("empty email and role", func(t *testing.T) {
		service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

		token, err := service.GenerateAccessToken(1, "", "", 123)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "", claims.Email)
		assert.Equal(t, "", claims.Role)
	})

	t.Run("different secrets", func(t *testing.T) {
		service1 := NewService("secret1", 15*time.Minute, 7*24*time.Hour)
		service2 := NewService("secret2", 15*time.Minute, 7*24*time.Hour)

		token, err := service1.GenerateAccessToken(1, "test@example.com", "user", 123)
		assert.NoError(t, err)

		// Token generated with service1 should not be valid with service2
		_, err = service2.ValidateToken(token)
		assert.Error(t, err)
	})
}

func TestClaims_Structure(t *testing.T) {
	service := NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	userID := uint(42)
	email := "user@example.com"
	role := "admin"
	sessionID := uint(999)

	token, err := service.GenerateAccessToken(userID, email, role, sessionID)
	assert.NoError(t, err)

	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)

	// Test all claim fields
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, sessionID, claims.SessionID)
	assert.Equal(t, "tech-challenge-s1", claims.Issuer)
	assert.Equal(t, "42", claims.Subject)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Second)))
	assert.True(t, claims.NotBefore.Before(time.Now().Add(time.Second)))
}
