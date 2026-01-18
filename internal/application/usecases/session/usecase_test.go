package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
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

func TestNewSessionUseCase(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "168h",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	assert.NotNil(t, useCase)
	assert.Implements(t, (*ports.SessionUseCase)(nil), useCase)
}

func TestLogin_Success(t *testing.T) {
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

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("Password123!"), nil)
	mockHasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)

	// Mock user
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Contact:   "123456789",
		Role:      "user",
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
	input := ports.LoginInput{
		Email:    user.Email,
		Password: "Password123!",
	}

	output, err := useCase.Login(context.Background(), input)

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

func TestLogin_UserNotFound(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "168h",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

	// Execute
	input := ports.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	output, err := useCase.Login(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "user not found")

	userService.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
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

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Mock user
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Contact:   "123456789",
		Role:      "user",
		DeletedAt: nil,
	}

	hasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(errors.New("invalid password"))

	user.Password, _ = domain.NewPassword("Correctpassword123>", hasher)

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)

	// Execute
	input := ports.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	output, err := useCase.Login(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "user not found")

	userService.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
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

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Mock inactive user
	deletedAt := time.Now()
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Contact:   "123456789",
		Role:      "user",
		DeletedAt: &deletedAt,
	}
	hasher.On("Generate", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("Password123!"), nil)
	hasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)
	user.Password, _ = domain.NewPassword("Password123!", hasher)
	user.Password.Hash()

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)

	// Execute
	input := ports.LoginInput{
		Email:    "test@example.com",
		Password: "Password123!",
	}

	output, err := useCase.Login(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "user not found")

	userService.AssertExpectations(t)
}

func TestLogin_SessionCreationError(t *testing.T) {
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

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)
	mockHasher.On("Generate", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("Password123!"), nil)
	mockHasher.On("Compare", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)

	// Mock user
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Contact:   "123456789",
		Role:      "user",
		DeletedAt: nil,
	}
	user.Password, _ = domain.NewPassword("Password123!", mockHasher)
	user.Password.Hash()

	// Setup mocks
	userService.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)
	sessionService.On("Create", mock.Anything, uint(1), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return((*domain.Session)(nil), assert.AnError)

	// Execute
	input := ports.LoginInput{
		Email:    "test@example.com",
		Password: "Password123!",
	}

	output, err := useCase.Login(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}

	useCase := NewSessionUseCase(userService, sessionService, nil, nil)

	// Setup mocks: direct delete by refresh token
	sessionService.On("DeleteByRefreshToken", mock.Anything, "refresh-token").Return(nil)

	// Execute
	input := ports.LogoutInput{
		RefreshToken: "refresh-token",
	}

	err := useCase.Logout(context.Background(), input)

	// Assertions
	assert.NoError(t, err)

	sessionService.AssertExpectations(t)
}

func TestLogout_InvalidRefreshToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}

	useCase := NewSessionUseCase(userService, sessionService, nil, nil)

	// Setup mocks - session not found
	sessionService.On("DeleteByRefreshToken", mock.Anything, "invalid-token").Return(errors.New("invalid or expired session"))

	// Execute
	input := ports.LogoutInput{
		RefreshToken: "invalid-token",
	}

	err := useCase.Logout(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid or expired session")

	sessionService.AssertExpectations(t)
}

func TestLogout_SessionNotFound(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}

	useCase := NewSessionUseCase(userService, sessionService, nil, nil)

	// Setup mocks - session not found (retornar error)
	sessionService.On("DeleteByRefreshToken", mock.Anything, "nonexistent-token").Return(errors.New("invalid or expired session"))

	// Execute
	input := ports.LogoutInput{
		RefreshToken: "nonexistent-token",
	}

	err := useCase.Logout(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid or expired session")

	sessionService.AssertExpectations(t)
}

func TestLogout_UpdateError(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}

	useCase := NewSessionUseCase(userService, sessionService, nil, nil)

	// Setup mocks
	sessionService.On("DeleteByRefreshToken", mock.Anything, "refresh-token").Return(assert.AnError)

	// Execute
	input := ports.LogoutInput{
		RefreshToken: "refresh-token",
	}

	err := useCase.Logout(context.Background(), input)

	// Assertions
	assert.Error(t, err)

	sessionService.AssertExpectations(t)
}

func TestLogout_EmptyRefreshToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}

	useCase := NewSessionUseCase(userService, sessionService, nil, nil)

	// Setup mocks - empty token validation fails (retornar um error, não string)
	sessionService.On("DeleteByRefreshToken", mock.Anything, "").Return(errors.New("invalid or expired session"))

	// Execute
	input := ports.LogoutInput{
		RefreshToken: "",
	}

	err := useCase.Logout(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid or expired session")

	sessionService.AssertExpectations(t)
}

func TestRefresh_Success(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Mock session
	session := &domain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	// Mock user
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Contact:   "123456789",
		Role:      "user",
		DeletedAt: nil,
	}

	// Setup mocks
	sessionService.On("Validate", mock.Anything, "refresh-token").Return(session, nil)
	userService.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	// Execute
	input := ports.RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Refresh(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.AccessToken)
	assert.Equal(t, "refresh-token", output.RefreshToken)
	assert.Equal(t, int64(900), output.ExpiresIn) // 15 minutes in seconds

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestRefresh_InvalidRefreshToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails (retorna um error, não uma string)
	sessionService.On("Validate", mock.Anything, "invalid-token").
		Return((*domain.Session)(nil), errors.New("invalid or expired session"))

	// Execute
	input := ports.RefreshInput{
		RefreshToken: "invalid-token",
	}

	output, err := useCase.Refresh(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid or expired session")

	sessionService.AssertExpectations(t)
}

func TestRefresh_ExpiredSession(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails due to expiration
	sessionService.On("Validate", mock.Anything, "refresh-token").Return((*domain.Session)(nil), errors.New("invalid or expired session"))

	// Execute
	input := ports.RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Refresh(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid or expired session")

	sessionService.AssertExpectations(t)
}

func TestRefresh_InactiveSession(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails due to inactive session (retorna error, não string)
	sessionService.On("Validate", mock.Anything, "refresh-token").Return((*domain.Session)(nil), errors.New("invalid or expired session"))

	// Execute
	input := ports.RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Refresh(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid or expired session")

	sessionService.AssertExpectations(t)
}

func TestRefresh_UserNotFound(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Mock session
	session := &domain.Session{
		ID:           1,
		UserID:       999, // Non-existent user
		RefreshToken: stringPtr("refresh-token"),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	// Setup mocks
	sessionService.On("Validate", mock.Anything, "refresh-token").Return(session, nil)
	userService.On("GetByID", mock.Anything, uint(999)).Return((*domain.User)(nil), errors.New("user not found"))

	// Execute
	input := ports.RefreshInput{
		RefreshToken: "refresh-token",
	}

	output, err := useCase.Refresh(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "user not found")

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestRefresh_EmptyRefreshToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)
	config := &config.Config{
		JWT: &config.JWTConfig{
			AccessTokenExpiry:  "15m",
			RefreshTokenExpiry: "7d",
		},
	}

	useCase := NewSessionUseCase(userService, sessionService, jwtService, config)

	// Setup mocks - session validation fails for empty token
	sessionService.On("Validate", mock.Anything, "").Return((*domain.Session)(nil), errors.New("invalid or expired session"))

	// Execute
	input := ports.RefreshInput{
		RefreshToken: "",
	}

	output, err := useCase.Refresh(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid or expired session")

	sessionService.AssertExpectations(t)
}

func TestValidate_ValidToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Generate a valid token with session ID
	userID := uint(1)
	sessionID := uint(123)
	email := "test@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role, sessionID)
	assert.NoError(t, err)

	// Mock user
	user := &domain.User{
		ID:        userID,
		Name:      "Test User",
		Email:     email,
		Contact:   "123456789",
		Role:      role,
		DeletedAt: nil,
	}

	// Mock valid session
	validSession := &domain.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Setup mocks
	userService.On("GetByID", mock.Anything, userID).Return(user, nil)
	sessionService.On("GetByID", mock.Anything, sessionID).Return(validSession, nil)

	// Execute
	input := ports.ValidateInput{
		Token: token,
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.Valid)
	assert.NotNil(t, output.User)
	assert.Equal(t, user.ID, output.User.ID)
	assert.Equal(t, user.Email, output.User.Email)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestValidate_InvalidToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Execute with invalid token
	input := ports.ValidateInput{
		Token: "invalid.token.here",
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestValidate_EmptyToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Execute with empty token
	input := ports.ValidateInput{
		Token: "",
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestValidate_UserNotFound(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Generate a valid token with session ID
	userID := uint(999) // Non-existent user
	sessionID := uint(123)
	email := "nonexistent@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role, sessionID)
	assert.NoError(t, err)

	// Mock valid session
	validSession := &domain.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Setup mocks - user not found
	userService.On("GetByID", mock.Anything, userID).Return((*domain.User)(nil), assert.AnError)
	sessionService.On("GetByID", mock.Anything, sessionID).Return(validSession, nil)

	// Execute
	input := ports.ValidateInput{
		Token: token,
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestValidate_ExpiredToken(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	// Create JWT service with very short expiry
	jwtService := jwt.NewService("test-secret", 1*time.Millisecond, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Generate a token that expires quickly
	userID := uint(1)
	sessionID := uint(123)
	email := "test@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role, sessionID)
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Execute
	input := ports.ValidateInput{
		Token: token,
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestValidate_SessionNotFound(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Generate a valid token with session ID
	userID := uint(1)
	sessionID := uint(999) // Non-existent session
	email := "test@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role, sessionID)
	assert.NoError(t, err)

	// Setup mocks - session not found
	sessionService.On("GetByID", mock.Anything, sessionID).Return((*domain.Session)(nil), assert.AnError)

	// Execute
	input := ports.ValidateInput{
		Token: token,
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid or expired session")

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestValidate_ExpiredSession(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Generate a valid token with session ID
	userID := uint(1)
	sessionID := uint(123)
	email := "test@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role, sessionID)
	assert.NoError(t, err)

	// Mock expired session
	expiredSession := &domain.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(-24 * time.Hour), // Expired
	}

	// Setup mocks - expired session
	sessionService.On("GetByID", mock.Anything, sessionID).Return(expiredSession, nil)

	// Execute
	input := ports.ValidateInput{
		Token: token,
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid or expired session")

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

func TestValidate_SessionServiceError(t *testing.T) {
	userService := &MockUserService{}
	sessionService := &MockSessionService{}
	jwtService := jwt.NewService("test-secret", 15*time.Minute, 7*24*time.Hour)

	useCase := NewSessionUseCase(userService, sessionService, jwtService, nil)

	// Generate a valid token with session ID
	userID := uint(1)
	sessionID := uint(123)
	email := "test@example.com"
	role := "user"
	token, err := jwtService.GenerateAccessToken(userID, email, role, sessionID)
	assert.NoError(t, err)

	// Setup mocks - session service error
	sessionService.On("GetByID", mock.Anything, sessionID).Return((*domain.Session)(nil), assert.AnError)

	// Execute
	input := ports.ValidateInput{
		Token: token,
	}

	output, err := useCase.Validate(context.Background(), input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid or expired session")

	userService.AssertExpectations(t)
	sessionService.AssertExpectations(t)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
