package domain

import (
	"errors"
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
	return nil, errors.New("erro ao criar o hash da senha")
}

func (b *mockHasherFailed) Compare(hashedPassword, password []byte) error {
	return errors.New("senha inválida")
}

func TestWeakPassword(t *testing.T) {
	weakPass := "senhafraca"

	p, err := NewPassword(weakPass, &mockHasher{})

	assert.Error(t, err)
	assert.EqualError(t, err, "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	assert.Nil(t, p)
}

func TestShortPassword(t *testing.T) {
	emptyPass := "1234567"

	p, err := NewPassword(emptyPass, &mockHasher{})

	assert.EqualError(t, err, "senha deve ter pelo menos 8 caracteres")
	assert.Nil(t, p)
}

func TestInvalidPassword(t *testing.T) {
	invalidPass := "12345678"

	p, err := NewPassword(invalidPass, &mockHasher{})

	assert.Error(t, err)
	assert.EqualError(t, err, "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
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

	p, _ := NewPassword(validPass, &mockHasher{})

	provided := "P@ss123><!"

	assert.NotEqual(t, p.GetValue(), provided)

	err := p.ValidateEqual(provided)
	assert.Error(t, err)
	assert.EqualError(t, err, "senha incorreta")
}

func TestNoPasswordAvailable(t *testing.T) {
	p := &Password{}

	err := p.ValidateEqual("P@ss123><!")

	assert.Error(t, err)
	assert.EqualError(t, err, "nenhuma senha disponível para comparar")
}

func TestInvalidHashPassword(t *testing.T) {
	hashed := "hashedpassword"

	p := NewPasswordFromHash(hashed, &mockHasherFailed{})

	err := p.Hash()

	assert.EqualError(t, err, "erro ao criar o hash da senha")
}
