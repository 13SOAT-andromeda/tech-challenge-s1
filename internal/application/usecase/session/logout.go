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
	// Find session
	session, err := uc.sessionService.GetByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, services.ErrSessionInvalid
	}

	// Deactivate session
	session.Deactivate()
	_, err = uc.sessionService.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	output := &LogoutOutput{
		Success: true,
		Message: "Logged out successfully",
	}

	return output, nil
}
