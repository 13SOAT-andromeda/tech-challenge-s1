package domain

import (
	"time"
)

type Session struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	UserID       uint
	RefreshToken *string
	ExpiresAt    time.Time
}

// NewSession creates a new session with the given user ID and refresh token
func NewSession(userID uint, refreshToken string, expiresAt time.Time) *Session {
	return &Session{
		UserID:       userID,
		RefreshToken: &refreshToken,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsValid() bool {
	return !s.IsExpired()
}
