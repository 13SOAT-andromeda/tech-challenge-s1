package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHasher struct{}

func (b *mockHasher) Generate(password []byte, cost int) ([]byte, error) {
	return []byte("abcdefghijklmnopqrstuvwxyz"), nil
}

func (b *mockHasher) Compare(hashedPassword, password []byte) error {
	return nil
}

type mockHasherFailed struct{}

func (b *mockHasherFailed) Generate(password []byte, cost int) ([]byte, error) {
	return []byte(""), ErrPasswordHash
}

func (b *mockHasherFailed) Compare(hashedPassword, password []byte) error {
	return ErrPasswordInvalid
}

func TestWeakPassword(t *testing.T) {
	weakPass := "senhafraca"

	p, err := NewPassword(weakPass, &mockHasher{})

	assert.ErrorIs(t, ErrPasswordInvalidFormat, err)
	assert.Nil(t, p)
}

func TestShortPassword(t *testing.T) {
	emptyPass := "1234567"

	p, err := NewPassword(emptyPass, &mockHasher{})

	assert.ErrorIs(t, ErrPasswordTooShort, err)
	assert.Nil(t, p)
}

func TestInvalidPassword(t *testing.T) {
	invalidPass := "12345678"

	p, err := NewPassword(invalidPass, &mockHasher{})

	assert.ErrorIs(t, ErrPasswordInvalidFormat, err)
	assert.Nil(t, p)
}

func TestValidPassword(t *testing.T) {
	validPass := "P@ss123><!..."

	p, err := NewPassword(validPass, &mockHasher{})

	assert.NoError(t, err)
	assert.NotEmpty(t, p.GetValue())
}

func TestHashPassword(t *testing.T) {
	validPass := "P@ss123><!..."

	p, err := NewPassword(validPass, &mockHasher{})

	p.Hash()

	if err != nil {
		t.Fatalf("Hash() failed: %v", err)
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, p.GetValue())
	assert.NotEmpty(t, p.GetHashed())
}

func TestCompareValidPassword(t *testing.T) {
	validPass := "P@ss123><!..."

	p, _ := NewPassword(validPass, &mockHasher{})
	p.Hash()

	err := p.Compare(validPass)

	assert.NoError(t, err)
}

func TestCompareInvalidPassword(t *testing.T) {
	validPass := "P@ss123><!..."

	p, _ := NewPassword(validPass, &mockHasherFailed{})
	p.Hash()

	err := p.Compare("invalidpass")

	assert.ErrorIs(t, ErrPasswordInvalid, err)
}

func TestInvalidHashPassword(t *testing.T) {
	hashed := "hashedpassword"

	p := NewPasswordFromHash(hashed, &mockHasherFailed{})

	assert.ErrorIs(t, ErrPasswordHash, p.Hash())
}
