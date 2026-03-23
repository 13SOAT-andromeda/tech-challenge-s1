package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgjwt "github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const authTestSecret = "auth-test-secret"

func makeValidToken(t *testing.T) string {
	t.Helper()
	claims := pkgjwt.Claims{
		Email: "user@example.com",
		Role:  "mechanic",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "99",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(authTestSecret))
	require.NoError(t, err)
	return signed
}

func setupAuthRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.AuthRequired(secret))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestAuthRequired_ValidToken(t *testing.T) {
	r := setupAuthRouter(authTestSecret)
	tokenStr := makeValidToken(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthRequired_MissingHeader(t *testing.T) {
	r := setupAuthRouter(authTestSecret)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_MalformedHeader(t *testing.T) {
	r := setupAuthRouter(authTestSecret)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token somevalue") // not "Bearer"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	r := setupAuthRouter(authTestSecret)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_WrongSecret(t *testing.T) {
	r := setupAuthRouter("different-secret")
	tokenStr := makeValidToken(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_ClaimsStoredInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.AuthRequired(authTestSecret))
	r.GET("/test", func(c *gin.Context) {
		id := middlewares.GetUserID(c)
		role := middlewares.GetUserRole(c)
		email := middlewares.GetUserEmail(c)
		c.JSON(http.StatusOK, gin.H{"id": id, "role": role, "email": email})
	})

	tokenStr := makeValidToken(t)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, `"id":"99"`)
	assert.Contains(t, body, `"role":"mechanic"`)
	assert.Contains(t, body, `"email":"user@example.com"`)
}
