package services

import (
	"context"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type SessionService struct {
	repo ports.SessionRepository
}

func NewSessionService(repo ports.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

// Create creates a new session for a user
func (s *SessionService) Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*domain.Session, error) {
	if userID == 0 {
		return nil, ErrSessionUserIDInvalid
	}

	if refreshToken == "" {
		return nil, ErrSessionRefreshTokenEmpty
	}

	if expiresAt.Before(time.Now()) {
		return nil, ErrSessionExpiresAtPast
	}

	session := domain.NewSession(userID, refreshToken, expiresAt)
	return s.repo.Create(ctx, session)
}

// GetByRefreshToken retrieves a session by its refresh token
func (s *SessionService) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	if refreshToken == "" {
		return nil, ErrSessionRefreshTokenEmpty
	}

	return s.repo.FindByRefreshToken(ctx, refreshToken)
}

// GetByUserID retrieves all sessions for a user
func (s *SessionService) GetByUserID(ctx context.Context, userID uint) ([]*domain.Session, error) {
	if userID == 0 {
		return nil, ErrSessionUserIDInvalid
	}

	return s.repo.FindByUserID(ctx, userID)
}

// Update updates an existing session
func (s *SessionService) Update(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	if session == nil {
		return nil, ErrSessionNil
	}

	if session.ID == 0 {
		return nil, ErrSessionIDInvalid
	}

	return s.repo.Update(ctx, session)
}

// Delete deletes a specific session
func (s *SessionService) Delete(ctx context.Context, sessionID uint) error {
	if sessionID == 0 {
		return ErrSessionIDInvalid
	}

	return s.repo.Delete(ctx, sessionID)
}

// DeleteByUserID deletes all sessions for a user
func (s *SessionService) DeleteByUserID(ctx context.Context, userID uint) error {
	if userID == 0 {
		return ErrSessionUserIDInvalid
	}

	return s.repo.DeleteByUserID(ctx, userID)
}

// Validate validates a session by refresh token and checks if it's still valid
func (s *SessionService) Validate(ctx context.Context, refreshToken string) (*domain.Session, error) {
	session, err := s.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if !session.IsValid() {
		return nil, ErrSessionInvalid
	}

	return session, nil
}
