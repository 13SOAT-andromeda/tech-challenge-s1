package services

import (
	"context"
	"errors"
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
		return nil, errors.New("invalid user ID")
	}

	if refreshToken == "" {
		return nil, errors.New("refresh token cannot be empty")
	}

	if expiresAt.Before(time.Now()) {
		return nil, errors.New("data de expiração não pode estar no passado")
	}

	session := domain.NewSession(userID, refreshToken, expiresAt)
	return s.repo.Create(ctx, session)
}

func (s *sessionService) GetByID(ctx context.Context, sessionID uint) (*domain.Session, error) {
	if sessionID == 0 {
		return nil, errors.New("invalid session ID")
	}

	return s.repo.FindByID(ctx, sessionID)
}

func (s *sessionService) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token cannot be empty")
	}

	return s.repo.FindByRefreshToken(ctx, refreshToken)
}

func (s *sessionService) GetByUserID(ctx context.Context, userID uint) ([]*domain.Session, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.FindByUserID(ctx, userID)
}

func (s *sessionService) Update(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	if session == nil {
		return nil, errors.New("session cannot be null")
	}

	if session.ID == 0 {
		return nil, errors.New("invalid session ID")
	}

	return s.repo.Update(ctx, session)
}

func (s *sessionService) Delete(ctx context.Context, sessionID uint) error {
	if sessionID == 0 {
		return errors.New("invalid session ID")
	}

	return s.repo.Delete(ctx, sessionID)
}

func (s *sessionService) DeleteByUserID(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	return s.repo.DeleteByUserID(ctx, userID)
}

func (s *sessionService) Validate(ctx context.Context, refreshToken string) (*domain.Session, error) {
	session, err := s.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if !session.IsValid() {
		return nil, errors.New("invalid or expired session")
	}

	return session, nil
}

func (s *sessionService) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token cannot be empty")
	}

	_, err := s.repo.FindByRefreshToken(ctx, refreshToken)

	if err != nil {
		return err
	}

	return s.repo.DeleteByRefreshToken(ctx, refreshToken)
}

func (s *sessionService) DeleteExpiredSessions(ctx context.Context) error {
	return s.repo.DeleteExpiredSessions(ctx)
}
