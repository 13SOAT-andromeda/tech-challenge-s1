package domain

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	value string
}

func NewPassword(value string) (*Password, error) {
	if err := validatePassword(value); err != nil {
		return nil, err
	}

	hash, err := hash(value)

	if err != nil {
		return nil, err
	}

	return &Password{value: hash}, nil
}

func (p *Password) Get() string {
	return p.value
}

func hash(value string) (string, error) {
	p := []byte(value)
	hash, err := bcrypt.GenerateFromPassword(p, 15)

	if err != nil {
		return "", ErrPasswordHash
	}

	return string(hash), nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
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
		return ErrPasswordInvalidFormat
	}

	return nil
}

var (
	ErrPasswordTooShort      = &ValidationError{Message: "senha deve ter pelo menos 8 caracteres"}
	ErrPasswordHash          = &ValidationError{Message: "erro ao criar o hash da senha"}
	ErrPasswordInvalidFormat = &ValidationError{Message: "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial"}
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
