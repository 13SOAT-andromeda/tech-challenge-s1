package domain

import "github.com/13SOAT-andromeda/tech-challenge-s1/pkg/errors"

var (
	ErrPasswordTooShort      = &errors.ValidationError{Message: "senha deve ter pelo menos 8 caracteres"}
	ErrPasswordHash          = &errors.ValidationError{Message: "erro ao criar o hash da senha"}
	ErrPasswordInvalidFormat = &errors.ValidationError{Message: "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial"}
	ErrPasswordInvalid       = &errors.ValidationError{Message: "senha inválida"}
)
