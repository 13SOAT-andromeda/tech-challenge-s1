package services

import appErrors "github.com/13SOAT-andromeda/tech-challenge-s1/pkg/errors"

// User errors
var (
	ErrUserIdInvalid                = &appErrors.ValidationError{Message: "ID de usuário inválido"}
	ErrUserNotFound                 = &appErrors.ValidationError{Message: "usuário não encontrado"}
	ErrUserEmailAlreadyExists       = &appErrors.ValidationError{Message: "email já existe"}
	ErrUserPasswordUpdateNotAllowed = &appErrors.ValidationError{Message: "senha de usuário não pode ser atualizada"}
	ErrUserDelete                   = &appErrors.ValidationError{Message: "ocorreu um erro ao excluir o usuário"}
)
