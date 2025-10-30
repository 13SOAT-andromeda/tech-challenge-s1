package session

import (
	"context"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
)

type loginUseCase struct {
	userService    ports.UserService
	sessionService ports.SessionService
	jwtService     *jwt.Service
	config         *config.Config
}

func NewLoginUseCase(
	userService ports.UserService,
	sessionService ports.SessionService,
	jwtService *jwt.Service,
	config *config.Config,
) LoginUseCase {
	return &loginUseCase{
		userService:    userService,
		sessionService: sessionService,
		jwtService:     jwtService,
		config:         config,
	}
}

func (uc *loginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.userService.GetByEmail(ctx, input.Email)

	if err != nil || user == nil {
		return nil, services.ErrUserNotFound
	}

	if err := user.Password.Compare(input.Password); err != nil || !user.Active {
		return nil, services.ErrUserNotFound
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshExpiry, _ := time.ParseDuration(uc.config.JWT.RefreshTokenExpiry)
	sessionExpiresAt := time.Now().Add(refreshExpiry)

	session, err := uc.sessionService.Create(ctx, user.ID, refreshToken, sessionExpiresAt)
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, err
	}

	accessExpiry, _ := time.ParseDuration(uc.config.JWT.AccessTokenExpiry)

	userOutput := UserOutput{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Contact: user.Contact,
		Role:    user.Role,
		Active:  user.Active,
	}

	output := &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpiry.Seconds()),
		User:         userOutput,
	}

	return output, nil
}
