package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var hasher = NewBcryptHasher()

func TestBcryptHasher(t *testing.T) {
	assert.NotNil(t, hasher)
}

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := hasher.Generate([]byte(password), 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashPasswordFailed(t *testing.T) {
	longPassword := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus ac ultrices quam. Sed rutrum est."
	hash, err := hasher.Generate([]byte(longPassword), 1)
	assert.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	assert.Empty(t, hash)
}

func TestComparePassword(t *testing.T) {
	password := "password"
	hash, err := hasher.Generate([]byte(password), 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	err = hasher.Compare(hash, []byte(password))
	assert.NoError(t, err)
}

func TestComparePasswordFailed(t *testing.T) {
	password := "password"
	hash, err := hasher.Generate([]byte(password), 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	err = hasher.Compare(hash, []byte("wrongpassword"))
	assert.Error(t, err)
}
