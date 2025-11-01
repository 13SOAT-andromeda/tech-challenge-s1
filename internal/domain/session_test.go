package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	userID := uint(1)
	refreshToken := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	session := NewSession(userID, refreshToken, expiresAt)

	assert.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, refreshToken, *session.RefreshToken)
	assert.Equal(t, expiresAt, session.ExpiresAt)
	assert.False(t, session.CreatedAt.IsZero())
	assert.False(t, session.UpdatedAt.IsZero())
}

func TestSession_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "session not expired",
			expiresAt: time.Now().Add(time.Hour),
			expected:  false,
		},
		{
			name:      "session expired",
			expiresAt: time.Now().Add(-time.Hour),
			expected:  true,
		},
		{
			name:      "session expires now",
			expiresAt: time.Now().Add(-time.Millisecond),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &Session{
				ExpiresAt: tt.expiresAt,
			}

			result := session.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSession_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "valid session",
			expiresAt: time.Now().Add(time.Hour),
			expected:  true,
		},
		{
			name:      "expired session",
			expiresAt: time.Now().Add(-time.Hour),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &Session{
				ExpiresAt: tt.expiresAt,
			}

			result := session.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSession_EdgeCases(t *testing.T) {
	t.Run("nil refresh token", func(t *testing.T) {
		session := &Session{
			UserID:       1,
			RefreshToken: nil,
			ExpiresAt:    time.Now().Add(time.Hour),
		}

		assert.True(t, session.IsValid())
		assert.False(t, session.IsExpired())
	})

	t.Run("empty refresh token", func(t *testing.T) {
		emptyToken := ""
		session := &Session{
			UserID:       1,
			RefreshToken: &emptyToken,
			ExpiresAt:    time.Now().Add(time.Hour),
		}

		assert.True(t, session.IsValid())
		assert.Equal(t, "", *session.RefreshToken)
	})

	t.Run("zero time values", func(t *testing.T) {
		session := &Session{
			UserID:    1,
			ExpiresAt: time.Time{},
		}

		assert.True(t, session.IsExpired())
		assert.False(t, session.IsValid())
	})
}
