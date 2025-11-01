package middlewares

import (
	"net/http"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService     *jwt.Service
	sessionService ports.SessionService
}

func NewAuthMiddleware(config *config.Config, sessionService ports.SessionService) *AuthMiddleware {
	accessExpiry, _ := time.ParseDuration(config.JWT.AccessTokenExpiry)
	refreshExpiry, _ := time.ParseDuration(config.JWT.RefreshTokenExpiry)
	jwtService := jwt.NewService(config.JWT.Secret, accessExpiry, refreshExpiry)

	return &AuthMiddleware{
		jwtService:     jwtService,
		sessionService: sessionService,
	}
}

func (m *AuthMiddleware) validateSession(c *gin.Context) (*jwt.Claims, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return nil, false
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		c.Abort()
		return nil, false
	}

	token := authHeader[7:]
	claims, err := m.jwtService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return nil, false
	}

	session, err := m.sessionService.GetByID(c.Request.Context(), claims.SessionID)
	if err != nil || !session.IsValid() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired or invalid"})
		c.Abort()
		return nil, false
	}

	return claims, true
}

// AuthRequired middleware that requires authentication
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, valid := m.validateSession(c)
		if !valid {
			return
		}

		// Add user information to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RoleRequired middleware that requires specific role
func (m *AuthMiddleware) RoleRequired(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, valid := m.validateSession(c)
		if !valid {
			return
		}

		// Check role
		if claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Add user information to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// Helper functions to get user info from context
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

func GetUserEmail(c *gin.Context) string {
	if email, exists := c.Get("user_email"); exists {
		if e, ok := email.(string); ok {
			return e
		}
	}
	return ""
}

func GetUserRole(c *gin.Context) string {
	if role, exists := c.Get("user_role"); exists {
		if r, ok := role.(string); ok {
			return r
		}
	}
	return ""
}

func GetClaims(c *gin.Context) *jwt.Claims {
	if claims, exists := c.Get("claims"); exists {
		if c, ok := claims.(*jwt.Claims); ok {
			return c
		}
	}
	return nil
}
