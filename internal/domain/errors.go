package domain

import appErrors "github.com/13SOAT-andromeda/tech-challenge-s1/pkg/errors"

var (
	ErrPasswordTooShort      = &appErrors.ValidationError{Message: "senha deve ter pelo menos 8 caracteres"}
	ErrPasswordHash          = &appErrors.ValidationError{Message: "erro ao criar o hash da senha"}
	ErrPasswordInvalidFormat = &appErrors.ValidationError{Message: "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial"}
	ErrPasswordInvalid       = &appErrors.ValidationError{Message: "senha inválida"}
)
