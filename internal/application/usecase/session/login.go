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
	// Find user by email
	users, err := uc.userService.Search(ctx, map[string]interface{}{
		"email": input.Email,
	})
	if err != nil {
		return nil, services.ErrUserNotFound
	}

	if len(*users) == 0 {
		return nil, services.ErrUserNotFound
	}

	user := (*users)[0]

	// Verify password
	if err := user.Password.Compare(input.Password); err != nil {
		return nil, services.ErrUserNotFound
	}

	// Check if user is active
	if !user.Active {
		return nil, services.ErrUserNotFound
	}

	// Generate tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Calculate expiry time
	accessExpiry, _ := time.ParseDuration(uc.config.JWT.AccessTokenExpiry)

	// Create session
	refreshExpiry, _ := time.ParseDuration(uc.config.JWT.RefreshTokenExpiry)
	sessionExpiresAt := time.Now().Add(refreshExpiry)

	_, err = uc.sessionService.Create(ctx, user.ID, refreshToken, sessionExpiresAt)
	if err != nil {
		return nil, err
	}

	// Prepare response
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
		TokenType:    "Bearer",
		User:         userOutput,
	}

	return output, nil
}
