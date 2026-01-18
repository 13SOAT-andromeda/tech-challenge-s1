package domain

import (
	"errors"
	"unicode"

	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"
)

type Password struct {
	value  string
	hashed string
	hasher encryption.Hasher
}

func NewPassword(value string, hasher encryption.Hasher) (*Password, error) {
	if err := validatePassword(value); err != nil {
		return nil, err
	}

	return &Password{value: value, hasher: hasher}, nil
}

func NewPasswordFromHash(hashed string, hasher encryption.Hasher) *Password {
	return &Password{hashed: hashed, hasher: hasher}
}

func (p *Password) GetValue() string {
	return p.value
}

func (p *Password) GetHashed() string {
	return p.hashed
}

func (p *Password) Hash() error {
	pass := []byte(p.value)
	hash, err := p.hasher.Generate(pass, 10)

	if err != nil {
		return errors.New("erro ao criar o hash da senha")
	}

	p.hashed = string(hash)

	return nil
}

func (p *Password) Compare(password string) error {
	err := p.hasher.Compare([]byte(p.hashed), []byte(password))
	if err != nil {
		return errors.New("senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("senha deve ter pelo menos 8 caracteres")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")

	}

	return nil
}

func (p *Password) ValidateEqual(input string) error {
	if p.hashed != "" {
		if err := p.hasher.Compare([]byte(p.hashed), []byte(input)); err != nil {
			return errors.New("senha incorreta")
		}
		return nil
	}
	if p.value != "" {
		if p.value != input {
			return errors.New("senha incorreta")
		}
		return nil
	}
	return errors.New("nenhuma senha disponível para comparar")
}
