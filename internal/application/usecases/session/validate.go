package session

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
)

type validateUseCase struct {
	userService    ports.UserService
	sessionService ports.SessionService
	jwtService     *jwt.Service
}

func NewValidateUseCase(
	userService ports.UserService,
	sessionService ports.SessionService,
	jwtService *jwt.Service,
) ValidateUseCase {
	return &validateUseCase{
		userService:    userService,
		sessionService: sessionService,
		jwtService:     jwtService,
	}
}

func (uc *validateUseCase) Execute(ctx context.Context, input ValidateInput) (*ValidateOutput, error) {
	claims, err := uc.jwtService.ValidateToken(input.Token)
	if err != nil {
		return nil, err
	}

	session, err := uc.sessionService.GetByID(ctx, claims.SessionID)
	if err != nil {
		return nil, services.ErrSessionInvalid
	}

	if !session.IsValid() {
		return nil, services.ErrSessionInvalid
	}

	user, err := uc.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	userOutput := &UserOutput{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Contact: user.Contact,
		Role:    user.Role,
	}

	output := &ValidateOutput{
		Valid: true,
		User:  userOutput,
	}

	return output, nil
}
