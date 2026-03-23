package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/middlewares"
	pkgjwt "github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const roleTestSecret = "role-test-secret"

func makeRoleToken(t *testing.T, role string) string {
	t.Helper()
	claims := pkgjwt.Claims{
		Email: "test@example.com",
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(roleTestSecret))
	require.NoError(t, err)
	return signed
}

func setupRoleRouter(requiredRole string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.AuthRequired(roleTestSecret))
	r.GET("/test", middlewares.RoleRequired(requiredRole), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestRoleRequired_MatchingRole(t *testing.T) {
	r := setupRoleRouter("mechanic")
	tokenStr := makeRoleToken(t, "mechanic")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleRequired_AdministratorWithoutMatchingRoleIsForbidden(t *testing.T) {
	r := setupRoleRouter("mechanic")
	tokenStr := makeRoleToken(t, "administrator")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRoleRequired_InsufficientRole(t *testing.T) {
	r := setupRoleRouter("administrator")
	tokenStr := makeRoleToken(t, "mechanic")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRoleRequired_MissingClaimsInContext(t *testing.T) {
	// No AuthRequired middleware — context will have no claims.
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", middlewares.RoleRequired("mechanic"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestExtractClaims_ReturnsCorrectValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.AuthRequired(roleTestSecret))
	r.GET("/test", func(c *gin.Context) {
		claims, err := middlewares.ExtractClaims(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": claims.ID, "role": claims.Role, "email": claims.Email})
	})

	tokenStr := makeRoleToken(t, "administrator")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, `"role":"administrator"`)
	assert.Contains(t, body, `"email":"test@example.com"`)
}
