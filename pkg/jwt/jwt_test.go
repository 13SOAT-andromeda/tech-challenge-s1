package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	pkgjwt "github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"

func makeToken(claims pkgjwt.Claims, method jwt.SigningMethod, secret interface{}) string {
	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	return signed
}

func TestValidateToken_Valid(t *testing.T) {
	claims := pkgjwt.Claims{
		Email: "user@example.com",
		Role:  "mechanic",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "42",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenStr := makeToken(claims, jwt.SigningMethodHS256, []byte(testSecret))

	got, err := pkgjwt.ValidateToken(tokenStr, testSecret)
	require.NoError(t, err)
	assert.Equal(t, "42", got.Subject)
	assert.Equal(t, "user@example.com", got.Email)
	assert.Equal(t, "mechanic", got.Role)
}

func TestValidateToken_Expired(t *testing.T) {
	claims := pkgjwt.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	tokenStr := makeToken(claims, jwt.SigningMethodHS256, []byte(testSecret))

	_, err := pkgjwt.ValidateToken(tokenStr, testSecret)
	assert.Error(t, err)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	claims := pkgjwt.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenStr := makeToken(claims, jwt.SigningMethodHS256, []byte(testSecret))

	_, err := pkgjwt.ValidateToken(tokenStr, "wrong-secret")
	assert.Error(t, err)
}

func TestValidateToken_WrongAlgorithm(t *testing.T) {
	// Generate a throwaway RSA key for RS256 signing.
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	claims := pkgjwt.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(rsaKey)
	require.NoError(t, err)

	_, err = pkgjwt.ValidateToken(tokenStr, testSecret)
	assert.Error(t, err)
}

func TestValidateToken_Malformed(t *testing.T) {
	_, err := pkgjwt.ValidateToken("not.a.jwt", testSecret)
	assert.Error(t, err)
}

func TestValidateToken_EmptyString(t *testing.T) {
	_, err := pkgjwt.ValidateToken("", testSecret)
	assert.Error(t, err)
}
