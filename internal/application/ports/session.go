package ports

import (
	"context"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) (*domain.Session, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	FindByUserID(ctx context.Context, userID uint) ([]*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) (*domain.Session, error)
	Delete(ctx context.Context, sessionID uint) error
	DeleteByUserID(ctx context.Context, userID uint) error
}

type SessionService interface {
	Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*domain.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID uint) ([]*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) (*domain.Session, error)
	Delete(ctx context.Context, sessionID uint) error
	DeleteByUserID(ctx context.Context, userID uint) error
	Validate(ctx context.Context, refreshToken string) (*domain.Session, error)
}
