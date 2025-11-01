package session

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
)

type logoutUseCase struct {
	sessionService ports.SessionService
}

func NewLogoutUseCase(
	sessionService ports.SessionService,
) LogoutUseCase {
	return &logoutUseCase{
		sessionService: sessionService,
	}
}

func (uc *logoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	if err := uc.sessionService.DeleteByRefreshToken(ctx, input.RefreshToken); err != nil {
		return services.ErrSessionInvalid
	}
	return nil
}
