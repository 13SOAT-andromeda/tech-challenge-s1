package middlewares

import (
	"net/http"
	"strings"

	pkgjwt "github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const claimsKey = "userClaims"

// AuthRequired returns a Gin middleware that validates the JWT in the
// Authorization header and stores the parsed claims in the context.
// Returns 401 if the header is missing, malformed, or the token is invalid.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: malformed Authorization header"})
			c.Abort()
			return
		}

		claims, err := pkgjwt.ValidateToken(parts[1], secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token"})
			c.Abort()
			return
		}

		c.Set(claimsKey, &UserClaims{
			ID:    claims.Subject,
			Email: claims.Email,
			Role:  claims.Role,
		})

		c.Next()
	}
}
