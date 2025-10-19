package domain

import (
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestWeakPassword(t *testing.T) {
	weakPass := "senhafraca"

	p, err := NewPassword(weakPass)

	assert.ErrorIs(t, errors.ErrPasswordInvalidFormat, err)
	assert.Nil(t, p)
}

func TestShortPassword(t *testing.T) {
	emptyPass := "1234567"

	p, err := NewPassword(emptyPass)

	assert.ErrorIs(t, errors.ErrPasswordTooShort, err)
	assert.Nil(t, p)
}

func TestInvalidPassword(t *testing.T) {
	invalidPass := "12345678"

	p, err := NewPassword(invalidPass)

	assert.ErrorIs(t, errors.ErrPasswordInvalidFormat, err)
	assert.Nil(t, p)
}

func TestValidPassword(t *testing.T) {
	validPass := "P@ss123><!..."

	p, err := NewPassword(validPass)

	assert.NoError(t, err)
	assert.NotEmpty(t, p.GetValue())
}

func TestHashPassword(t *testing.T) {
	validPass := "P@ss123><!..."

	p, err := NewPassword(validPass)

	p.Hash()

	if err != nil {
		t.Fatalf("Hash() failed: %v", err)
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, p.GetValue())
	assert.NotEmpty(t, p.GetHashed())
}
