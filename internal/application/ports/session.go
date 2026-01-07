package ports

import (
	"context"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type UserOutput struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
	Role    string `json:"role"`
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	User         UserOutput `json:"user"`
}

type ValidateInput struct {
	Token string
}

type ValidateOutput struct {
	Valid bool        `json:"valid"`
	User  *UserOutput `json:"user,omitempty"`
}

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LogoutInput struct {
	RefreshToken string
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) (*domain.Session, error)
	FindByID(ctx context.Context, sessionID uint) (*domain.Session, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	FindByUserID(ctx context.Context, userID uint) ([]*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) (*domain.Session, error)
	Delete(ctx context.Context, sessionID uint) error
	DeleteByUserID(ctx context.Context, userID uint) error
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error
	DeleteExpiredSessions(ctx context.Context) error
}

type SessionService interface {
	Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*domain.Session, error)
	GetByID(ctx context.Context, sessionID uint) (*domain.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID uint) ([]*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) (*domain.Session, error)
	Delete(ctx context.Context, sessionID uint) error
	DeleteByUserID(ctx context.Context, userID uint) error
	Validate(ctx context.Context, refreshToken string) (*domain.Session, error)
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error
	DeleteExpiredSessions(ctx context.Context) error
}

type SessionUseCase interface {
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
	Logout(ctx context.Context, input LogoutInput) error
	Refresh(ctx context.Context, input RefreshInput) (*RefreshOutput, error)
	Validate(ctx context.Context, input ValidateInput) (*ValidateOutput, error)
}
