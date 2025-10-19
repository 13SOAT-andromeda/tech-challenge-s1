package domain

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/core/errors"
)

type Password struct {
	value  string
	hashed string
}

func NewPassword(value string) (*Password, error) {
	if err := validatePassword(value); err != nil {
		return nil, err
	}

	return &Password{value: value}, nil
}

func NewPasswordFromHash(hashed string) *Password {
	return &Password{hashed: hashed}
}

func (p *Password) GetValue() string {
	return p.value
}

func (p *Password) GetHashed() string {
	return p.hashed
}

func (p *Password) Hash() error {
	pass := []byte(p.value)
	hash, err := bcrypt.GenerateFromPassword(pass, 15)

	if err != nil {
		return errors.ErrPasswordHash
	}

	p.hashed = string(hash)

	return nil
}

func (p *Password) Compare(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashed), []byte(password))
	if err != nil {
		return errors.ErrPasswordInvalid
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.ErrPasswordTooShort
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
		return errors.ErrPasswordInvalidFormat
	}

	return nil
}
