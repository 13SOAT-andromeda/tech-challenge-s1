package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// mockSessionService minimally implements ports.SessionService for tests
type mockSessionService struct {
	sessionsByID map[uint]*domain.Session
	err          error
}

func (s *mockSessionService) Create(_ context.Context, _ uint, _ string, _ time.Time) (*domain.Session, error) {
	return nil, nil
}
func (s *mockSessionService) GetByID(_ context.Context, sessionID uint) (*domain.Session, error) {
	if s.err != nil {
		return nil, s.err
	}
	if sess, ok := s.sessionsByID[sessionID]; ok {
		return sess, nil
	}
	return nil, nil
}
func (s *mockSessionService) GetByRefreshToken(_ context.Context, _ string) (*domain.Session, error) {
	return nil, nil
}
func (s *mockSessionService) GetByUserID(_ context.Context, _ uint) ([]*domain.Session, error) {
	return nil, nil
}
func (s *mockSessionService) Update(_ context.Context, _ *domain.Session) (*domain.Session, error) {
	return nil, nil
}
func (s *mockSessionService) Delete(_ context.Context, _ uint) error         { return nil }
func (s *mockSessionService) DeleteByUserID(_ context.Context, _ uint) error { return nil }
func (s *mockSessionService) Validate(_ context.Context, _ string) (*domain.Session, error) {
	return nil, nil
}
func (s *mockSessionService) DeleteByRefreshToken(_ context.Context, _ string) error { return nil }
func (s *mockSessionService) DeleteExpiredSessions(_ context.Context) error          { return nil }

func setupRouter(m *AuthMiddleware, handler gin.HandlerFunc, useRole bool, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if useRole {
		r.Use(m.RoleRequired(role))
	} else {
		r.Use(m.AuthRequired())
	}
	r.GET("/test", handler)
	return r
}

func generateToken(t *testing.T, svc *jwt.Service, userID, sessionID uint, email, role string) string {
	t.Helper()
	token, err := svc.GenerateAccessToken(userID, email, role, sessionID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

func TestAuthRequired_NoAuthorizationHeader(t *testing.T) {
	jwtSvc := jwt.NewService("secret", time.Hour, time.Hour)
	m := &AuthMiddleware{jwtService: jwtSvc, sessionService: &mockSessionService{sessionsByID: map[uint]*domain.Session{}}}

	r := setupRouter(m, func(c *gin.Context) { c.Status(http.StatusOK) }, false, "")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_InvalidHeaderFormat(t *testing.T) {
	jwtSvc := jwt.NewService("secret", time.Hour, time.Hour)
	m := &AuthMiddleware{jwtService: jwtSvc, sessionService: &mockSessionService{sessionsByID: map[uint]*domain.Session{}}}

	r := setupRouter(m, func(c *gin.Context) { c.Status(http.StatusOK) }, false, "")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token abc")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	// Middleware validates with jwtSvcA, but token is signed with jwtSvcB
	jwtSvcA := jwt.NewService("secret-a", time.Hour, time.Hour)
	jwtSvcB := jwt.NewService("secret-b", time.Hour, time.Hour)
	token := generateToken(t, jwtSvcB, 1, 1, "user@example.com", "user")

	m := &AuthMiddleware{jwtService: jwtSvcA, sessionService: &mockSessionService{sessionsByID: map[uint]*domain.Session{}}}
	r := setupRouter(m, func(c *gin.Context) { c.Status(http.StatusOK) }, false, "")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_SessionInvalid(t *testing.T) {
	jwtSvc := jwt.NewService("secret", time.Hour, time.Hour)
	sessionID := uint(42)
	token := generateToken(t, jwtSvc, 7, sessionID, "user@example.com", "user")

	// Expired session
	expired := &domain.Session{ID: sessionID, ExpiresAt: time.Now().Add(-1 * time.Hour)}
	sessStub := &mockSessionService{sessionsByID: map[uint]*domain.Session{sessionID: expired}}
	m := &AuthMiddleware{jwtService: jwtSvc, sessionService: sessStub}

	r := setupRouter(m, func(c *gin.Context) { c.Status(http.StatusOK) }, false, "")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_Success_SetsContextAndCallsNext(t *testing.T) {
	jwtSvc := jwt.NewService("secret", time.Hour, time.Hour)
	sessionID := uint(99)
	userID := uint(123)
	email := "user@example.com"
	role := "user"
	token := generateToken(t, jwtSvc, userID, sessionID, email, role)

	valid := &domain.Session{ID: sessionID, ExpiresAt: time.Now().Add(1 * time.Hour)}
	sessStub := &mockSessionService{sessionsByID: map[uint]*domain.Session{sessionID: valid}}
	m := &AuthMiddleware{jwtService: jwtSvc, sessionService: sessStub}

	r := setupRouter(m, func(c *gin.Context) {
		// Validate context values populated by middleware
		if got := GetUserID(c); got != userID {
			t.Fatalf("user_id not set correctly, got %d", got)
		}
		if got := GetUserEmail(c); got != email {
			t.Fatalf("user_email not set correctly, got %s", got)
		}
		if got := GetUserRole(c); got != role {
			t.Fatalf("user_role not set correctly, got %s", got)
		}
		if claims := GetClaims(c); claims == nil || claims.UserID != userID || claims.SessionID != sessionID {
			t.Fatalf("claims not set correctly")
		}
		c.Status(http.StatusOK)
	}, false, "")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRoleRequired_ForbiddenOnRoleMismatch(t *testing.T) {
	jwtSvc := jwt.NewService("secret", time.Hour, time.Hour)
	sessionID := uint(1)
	token := generateToken(t, jwtSvc, 1, sessionID, "user@example.com", "user")
	valid := &domain.Session{ID: sessionID, ExpiresAt: time.Now().Add(1 * time.Hour)}
	m := &AuthMiddleware{jwtService: jwtSvc, sessionService: &mockSessionService{sessionsByID: map[uint]*domain.Session{sessionID: valid}}}

	r := setupRouter(m, func(c *gin.Context) { c.Status(http.StatusOK) }, true, "admin")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestRoleRequired_SuccessWhenRoleMatches(t *testing.T) {
	jwtSvc := jwt.NewService("secret", time.Hour, time.Hour)
	sessionID := uint(2)
	token := generateToken(t, jwtSvc, 1, sessionID, "admin@example.com", "admin")
	valid := &domain.Session{ID: sessionID, ExpiresAt: time.Now().Add(1 * time.Hour)}
	m := &AuthMiddleware{jwtService: jwtSvc, sessionService: &mockSessionService{sessionsByID: map[uint]*domain.Session{sessionID: valid}}}

	r := setupRouter(m, func(c *gin.Context) { c.Status(http.StatusOK) }, true, "admin")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
