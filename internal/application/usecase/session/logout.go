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

func (uc *logoutUseCase) Execute(ctx context.Context, input LogoutInput) (*LogoutOutput, error) {
	// Delete session by refresh token
	if err := uc.sessionService.DeleteByRefreshToken(ctx, input.RefreshToken); err != nil {
		return nil, services.ErrSessionInvalid
	}

	output := &LogoutOutput{
		Success: true,
		Message: "Logged out successfully",
	}

	return output, nil
}
