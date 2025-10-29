package session

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
)

type validateUseCase struct {
	userService ports.UserService
	jwtService  *jwt.Service
}

func NewValidateUseCase(
	userService ports.UserService,
	jwtService *jwt.Service,
) ValidateUseCase {
	return &validateUseCase{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (uc *validateUseCase) Execute(ctx context.Context, input ValidateInput) (*ValidateOutput, error) {
	claims, err := uc.jwtService.ValidateToken(input.Token)
	if err != nil {
		return nil, err
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
		Active:  user.Active,
	}

	output := &ValidateOutput{
		Valid: true,
		User:  userOutput,
	}

	return output, nil
}
