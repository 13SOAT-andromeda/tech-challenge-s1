package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeakPassword(t *testing.T) {
	weakPass := "senhafraca"

	p, err := NewPassword(weakPass)

	assert.ErrorIs(t, ErrPasswordInvalidFormat, err)
	assert.Nil(t, p)
}

func TestShortPassword(t *testing.T) {
	emptyPass := "1234567"

	p, err := NewPassword(emptyPass)

	assert.ErrorIs(t, ErrPasswordTooShort, err)
	assert.Nil(t, p)
}

func TestInvalidPassword(t *testing.T) {
	invalidPass := "12345678"

	p, err := NewPassword(invalidPass)

	assert.ErrorIs(t, ErrPasswordInvalidFormat, err)
	assert.Nil(t, p)
}

func TestValidPassword(t *testing.T) {
	validPass := "P@ss123><!..."

	p, err := NewPassword(validPass)

	assert.NoError(t, err)
	assert.NotEmpty(t, p.value)
}
