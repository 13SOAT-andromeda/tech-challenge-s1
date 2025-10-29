package session

import (
	"context"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
)

type refreshUseCase struct {
	userService    ports.UserService
	sessionService ports.SessionService
	jwtService     *jwt.Service
	config         *config.Config
}

func NewRefreshUseCase(
	userService ports.UserService,
	sessionService ports.SessionService,
	jwtService *jwt.Service,
	config *config.Config,
) RefreshUseCase {
	return &refreshUseCase{
		userService:    userService,
		sessionService: sessionService,
		jwtService:     jwtService,
		config:         config,
	}
}

func (uc *refreshUseCase) Execute(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	session, err := uc.sessionService.Validate(ctx, input.RefreshToken)
	if err != nil {
		return nil, services.ErrSessionInvalid
	}

	user, err := uc.userService.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, services.ErrUserNotFound
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	accessExpiry, _ := time.ParseDuration(uc.config.JWT.AccessTokenExpiry)

	output := &RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: input.RefreshToken,
		ExpiresIn:    int64(accessExpiry.Seconds()),
	}

	return output, nil
}
