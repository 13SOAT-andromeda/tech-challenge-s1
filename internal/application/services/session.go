package services

import (
	"context"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type sessionService struct {
	repo ports.SessionRepository
}

func NewSessionService(repo ports.SessionRepository) *sessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*domain.Session, error) {
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

func (s *sessionService) GetByID(ctx context.Context, sessionID uint) (*domain.Session, error) {
	if sessionID == 0 {
		return nil, ErrSessionIDInvalid
	}

	return s.repo.FindByID(ctx, sessionID)
}

func (s *sessionService) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	if refreshToken == "" {
		return nil, ErrSessionRefreshTokenEmpty
	}

	return s.repo.FindByRefreshToken(ctx, refreshToken)
}

func (s *sessionService) GetByUserID(ctx context.Context, userID uint) ([]*domain.Session, error) {
	if userID == 0 {
		return nil, ErrSessionUserIDInvalid
	}

	return s.repo.FindByUserID(ctx, userID)
}

func (s *sessionService) Update(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	if session == nil {
		return nil, ErrSessionNil
	}

	if session.ID == 0 {
		return nil, ErrSessionIDInvalid
	}

	return s.repo.Update(ctx, session)
}

func (s *sessionService) Delete(ctx context.Context, sessionID uint) error {
	if sessionID == 0 {
		return ErrSessionIDInvalid
	}

	return s.repo.Delete(ctx, sessionID)
}

func (s *sessionService) DeleteByUserID(ctx context.Context, userID uint) error {
	if userID == 0 {
		return ErrSessionUserIDInvalid
	}

	return s.repo.DeleteByUserID(ctx, userID)
}

func (s *sessionService) Validate(ctx context.Context, refreshToken string) (*domain.Session, error) {
	session, err := s.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if !session.IsValid() {
		return nil, ErrSessionInvalid
	}

	return session, nil
}

func (s *sessionService) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return ErrSessionRefreshTokenEmpty
	}
	return s.repo.DeleteByRefreshToken(ctx, refreshToken)
}

func (s *sessionService) DeleteExpiredSessions(ctx context.Context) error {
	return s.repo.DeleteExpiredSessions(ctx)
}
