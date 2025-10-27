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
	IsActive     bool
}

// NewSession creates a new session with the given user ID and refresh token
func NewSession(userID uint, refreshToken string, expiresAt time.Time) *Session {
	return &Session{
		UserID:       userID,
		RefreshToken: &refreshToken,
		ExpiresAt:    expiresAt,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Deactivate marks the session as inactive
func (s *Session) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

// IsValid checks if the session is both active and not expired
func (s *Session) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}
