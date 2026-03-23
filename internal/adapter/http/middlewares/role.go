package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserClaims holds the identity information extracted from the validated JWT.
type UserClaims struct {
	ID    string
	Email string
	Role  string
}

// ExtractClaims reads user identity from the Gin context (set by AuthRequired).
// Returns an error if the claims are not present.
func ExtractClaims(c *gin.Context) (*UserClaims, error) {
	val, exists := c.Get(claimsKey)
	if !exists {
		return nil, errors.New("missing user claims in context")
	}

	claims, ok := val.(*UserClaims)
	if !ok || claims == nil {
		return nil, errors.New("invalid user claims type in context")
	}

	return claims, nil
}

// RoleRequired returns a middleware that grants access if the request carries
// the required role (or the "administrator" role, which bypasses all restrictions).
// Returns 401 if claims are missing from context, 403 if the role is insufficient.
func RoleRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ExtractClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing identity claims"})
			c.Abort()
			return
		}

		if claims.Role != role && claims.Role != "administrator" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID reads the user ID from the JWT claims stored in the Gin context.
func GetUserID(c *gin.Context) string {
	claims, err := ExtractClaims(c)
	if err != nil {
		return ""
	}
	return claims.ID
}

// GetUserEmail reads the user email from the JWT claims stored in the Gin context.
func GetUserEmail(c *gin.Context) string {
	claims, err := ExtractClaims(c)
	if err != nil {
		return ""
	}
	return claims.Email
}

// GetUserRole reads the user role from the JWT claims stored in the Gin context.
func GetUserRole(c *gin.Context) string {
	claims, err := ExtractClaims(c)
	if err != nil {
		return ""
	}
	return claims.Role
}
