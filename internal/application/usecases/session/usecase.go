package session

import (
	"context"
	"errors"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
)

type UseCase struct {
	userService    ports.UserService
	sessionService ports.SessionService
	jwtService     *jwt.Service
	config         *config.Config
}

func NewSessionUseCase(userService ports.UserService, sessionService ports.SessionService, jwtService *jwt.Service, config *config.Config) *UseCase {
	return &UseCase{
		userService:    userService,
		sessionService: sessionService,
		jwtService:     jwtService,
		config:         config,
	}
}

func (u *UseCase) Login(ctx context.Context, input ports.LoginInput) (*ports.LoginOutput, error) {
	user, err := u.userService.GetByEmail(ctx, input.Email)

	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if err := user.Password.Compare(input.Password); err != nil || user.DeletedAt != nil {
		return nil, errors.New("user not found")
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshExpiry, _ := time.ParseDuration(u.config.JWT.RefreshTokenExpiry)
	sessionExpiresAt := time.Now().Add(refreshExpiry)

	session, err := u.sessionService.Create(ctx, user.ID, refreshToken, sessionExpiresAt)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, err
	}

	accessExpiry, _ := time.ParseDuration(u.config.JWT.AccessTokenExpiry)

	userOutput := ports.UserOutput{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Contact: user.Contact,
		Role:    user.Role,
	}

	output := &ports.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpiry.Seconds()),
		User:         userOutput,
	}

	return output, nil
}

func (u *UseCase) Logout(ctx context.Context, input ports.LogoutInput) error {
	if err := u.sessionService.DeleteByRefreshToken(ctx, input.RefreshToken); err != nil {
		return errors.New("invalid or expired session")
	}
	return nil
}

func (u *UseCase) Refresh(ctx context.Context, input ports.RefreshInput) (*ports.RefreshOutput, error) {
	session, err := u.sessionService.Validate(ctx, input.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired session")
	}

	user, err := u.userService.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := u.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, err
	}

	accessExpiry, _ := time.ParseDuration(u.config.JWT.AccessTokenExpiry)

	output := &ports.RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: input.RefreshToken,
		ExpiresIn:    int64(accessExpiry.Seconds()),
	}

	return output, nil
}

func (u *UseCase) Validate(ctx context.Context, input ports.ValidateInput) (*ports.ValidateOutput, error) {
	claims, err := u.jwtService.ValidateToken(input.Token)
	if err != nil {
		return nil, err
	}

	session, err := u.sessionService.GetByID(ctx, claims.SessionID)
	if err != nil {
		return nil, errors.New("invalid or expired session")
	}

	if !session.IsValid() {
		return nil, errors.New("invalid or expired session")
	}

	user, err := u.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	userOutput := &ports.UserOutput{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Contact: user.Contact,
		Role:    user.Role,
	}

	output := &ports.ValidateOutput{
		Valid: true,
		User:  userOutput,
	}

	return output, nil
}
