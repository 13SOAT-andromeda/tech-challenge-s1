package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of ports.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) CreateAdminUser(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) Search(ctx context.Context, params map[string]interface{}) (*[]domain.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*[]domain.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockSessionService is a mock implementation of ports.SessionService
type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*domain.Session, error) {
	args := m.Called(ctx, userID, refreshToken, expiresAt)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionService) GetByID(ctx context.Context, sessionID uint) (*domain.Session, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionService) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionService) GetByUserID(ctx context.Context, userID uint) ([]*domain.Session, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Session), args.Error(1)
}

func (m *MockSessionService) Update(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	args := m.Called(ctx, session)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionService) Delete(ctx context.Context, sessionID uint) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionService) DeleteByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionService) Validate(ctx context.Context, refreshToken string) (*domain.Session, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionService) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockSessionService) DeleteExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNewLoginUseCase(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "168h",
		},
	}

	useCase := NewLoginUseCase(userService, sessionService, jwtService, config)

	assert.NotNil(t, useCase)
	assert.Implements(t, (*LoginUseCase)(nil), useCase)
}

func TestLoginUseCase_Execute_Success(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	mockHasher := &encryption.MockHasher{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewLoginUseCase(userService, sessionService, jwtService, config)
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("Password123!"), nil)
	mockHasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)

	// Mock user
	user := &domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Contact:  "123456789",
		Role:     "user",
		DeletedAt: nil,
	}
	user.Password, _ = domain.NewPassword("Password123!", mockHasher)
	user.Password.Hash()

	// Mock session
	session := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)
	sessionService.On("Create", mock.Anything, uint(1), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(session, nil)

	// Execute
	input := LoginInput{
		Email:    user.Email,
		Password: "Password123!",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.AccessToken)
	assert.NotEmpty(t, output.RefreshToken)
	assert.Equal(t, int64(900), output.ExpiresIn) // 15 minutes in seconds
	assert.Equal(t, user.ID, output.User.ID)
	assert.Equal(t, user.Email, output.User.Email)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestLoginUseCase_Execute_UserNotFound(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "168h",
		},
	}

	useCase := NewLoginUseCase(userService, sessionService, jwtService, config)

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

	// Execute
	input := LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrUserNotFound, err)

	userService.AssertExpectations(t)
}

func TestLoginUseCase_Execute_InvalidPassword(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	hasher := &encryption.MockHasher{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewLoginUseCase(userService, sessionService, jwtService, config)

	// Mock user
	user := &domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Contact:  "123456789",
		Role:     "user",
		DeletedAt: nil,
	}

	hasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(errors.New("invalid password"))

	user.Password, _ = domain.NewPassword("Correctpassword123>", hasher)

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)

	// Execute
	input := LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrUserNotFound, err)

	userService.AssertExpectations(t)
}

func TestLoginUseCase_Execute_InactiveUser(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	hasher := &encryption.MockHasher{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewLoginUseCase(userService, sessionService, jwtService, config)

	// Mock inactive user
	deletedAt := time.Now()
	user := &domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Contact:  "123456789",
		Role:     "user",
		DeletedAt: &deletedAt,
	}
	hasher.On("Generate", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("Password123!"), nil)
	hasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)
	user.Password, _ = domain.NewPassword("Password123!", hasher)
	user.Password.Hash()

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)

	// Execute
	input := LoginInput{
		Email:    "test@example.com",
		Password: "Password123!",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, services.ErrUserNotFound, err)

	userService.AssertExpectations(t)
}

func TestLoginUseCase_Execute_SessionCreationError(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	mockHasher := &encryption.MockHasher{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewLoginUseCase(userService, sessionService, jwtService, config)
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("Password123!"), nil)
	mockHasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)

	// Mock user
	user := &domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Contact:  "123456789",
		Role:     "user",
		DeletedAt: nil,
	}
	user.Password, _ = domain.NewPassword("Password123!", mockHasher)
	user.Password.Hash()

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)
	sessionService.On("Create", mock.Anything, uint(1), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return((*domain.Session)(nil), assert.AnError)

	// Execute
	input := LoginInput{
		Email:    "test@example.com",
		Password: "Password123!",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
