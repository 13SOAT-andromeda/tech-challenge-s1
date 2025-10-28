package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSessionRepository is a mock implementation of SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	args := m.Called(ctx, session)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.Session, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) Update(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	args := m.Called(ctx, session)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) Delete(ctx context.Context, sessionID uint) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNewSessionService(t *testing.T) {
	mockRepo := &MockSessionRepository{}
	service := NewSessionService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestSessionService_Create(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		refreshToken string
		expiresAt    time.Time
		setupMock    func(*MockSessionRepository)
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "successful session creation",
			userID:       1,
			refreshToken: "valid-token",
			expiresAt:    time.Now().Add(24 * time.Hour),
			setupMock: func(m *MockSessionRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Session")).Return(&domain.Session{
					ID:           1,
					UserID:       1,
					RefreshToken: stringPtr("valid-token"),
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				}, nil)
			},
			expectError: false,
		},
		{
			name:         "zero user ID",
			userID:       0,
			refreshToken: "valid-token",
			expiresAt:    time.Now().Add(24 * time.Hour),
			setupMock:    func(m *MockSessionRepository) {},
			expectError:  true,
			errorMsg:     "ID de usuário inválido",
		},
		{
			name:         "empty refresh token",
			userID:       1,
			refreshToken: "",
			expiresAt:    time.Now().Add(24 * time.Hour),
			setupMock:    func(m *MockSessionRepository) {},
			expectError:  true,
			errorMsg:     "refresh token não pode estar vazio",
		},
		{
			name:         "expires at in the past",
			userID:       1,
			refreshToken: "valid-token",
			expiresAt:    time.Now().Add(-time.Hour),
			setupMock:    func(m *MockSessionRepository) {},
			expectError:  true,
			errorMsg:     "data de expiração não pode estar no passado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSessionRepository{}
			tt.setupMock(mockRepo)

			service := NewSessionService(mockRepo)
			session, err := service.Create(context.Background(), tt.userID, tt.refreshToken, tt.expiresAt)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tt.userID, session.UserID)
				assert.Equal(t, tt.refreshToken, *session.RefreshToken)
				assert.WithinDuration(t, tt.expiresAt, session.ExpiresAt, time.Second)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionService_GetByRefreshToken(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
		setupMock    func(*MockSessionRepository)
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "successful retrieval",
			refreshToken: "valid-token",
			setupMock: func(m *MockSessionRepository) {
				m.On("FindByRefreshToken", mock.Anything, "valid-token").Return(&domain.Session{
					ID:           1,
					UserID:       1,
					RefreshToken: stringPtr("valid-token"),
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				}, nil)
			},
			expectError: false,
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			setupMock:    func(m *MockSessionRepository) {},
			expectError:  true,
			errorMsg:     "refresh token não pode estar vazio",
		},
		{
			name:         "session not found",
			refreshToken: "invalid-token",
			setupMock: func(m *MockSessionRepository) {
				m.On("FindByRefreshToken", mock.Anything, "invalid-token").Return(nil, errors.New("session not found"))
			},
			expectError: true,
			errorMsg:    "session not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSessionRepository{}
			tt.setupMock(mockRepo)

			service := NewSessionService(mockRepo)
			session, err := service.GetByRefreshToken(context.Background(), tt.refreshToken)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tt.refreshToken, *session.RefreshToken)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionService_Validate(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
		setupMock    func(*MockSessionRepository)
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "valid session",
			refreshToken: "valid-token",
			setupMock: func(m *MockSessionRepository) {
				m.On("FindByRefreshToken", mock.Anything, "valid-token").Return(&domain.Session{
					ID:           1,
					UserID:       1,
					RefreshToken: stringPtr("valid-token"),
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				}, nil)
			},
			expectError: false,
		},
		{
			name:         "expired session",
			refreshToken: "expired-token",
			setupMock: func(m *MockSessionRepository) {
				m.On("FindByRefreshToken", mock.Anything, "expired-token").Return(&domain.Session{
					ID:           1,
					UserID:       1,
					RefreshToken: stringPtr("expired-token"),
					ExpiresAt:    time.Now().Add(-time.Hour),
				}, nil)
			},
			expectError: true,
			errorMsg:    "sessão inválida ou expirada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSessionRepository{}
			tt.setupMock(mockRepo)

			service := NewSessionService(mockRepo)
			session, err := service.Validate(context.Background(), tt.refreshToken)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.True(t, session.IsValid())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   uint
		setupMock   func(*MockSessionRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:      "successful deletion",
			sessionID: 1,
			setupMock: func(m *MockSessionRepository) {
				m.On("Delete", mock.Anything, uint(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "zero session ID",
			sessionID:   0,
			setupMock:   func(m *MockSessionRepository) {},
			expectError: true,
			errorMsg:    "ID de sessão inválido",
		},
		{
			name:      "session not found",
			sessionID: 999,
			setupMock: func(m *MockSessionRepository) {
				m.On("Delete", mock.Anything, uint(999)).Return(errors.New("session not found"))
			},
			expectError: true,
			errorMsg:    "session not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSessionRepository{}
			tt.setupMock(mockRepo)

			service := NewSessionService(mockRepo)
			err := service.Delete(context.Background(), tt.sessionID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionService_EdgeCases(t *testing.T) {
	t.Run("nil session in Update", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		service := NewSessionService(mockRepo)

		session, err := service.Update(context.Background(), nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sessão não pode ser nula")
		assert.Nil(t, session)
	})

	t.Run("zero user ID in GetByUserID", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		service := NewSessionService(mockRepo)

		sessions, err := service.GetByUserID(context.Background(), 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID de usuário inválido")
		assert.Nil(t, sessions)
	})

	t.Run("zero user ID in DeleteByUserID", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		service := NewSessionService(mockRepo)

		err := service.DeleteByUserID(context.Background(), 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID de usuário inválido")
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
