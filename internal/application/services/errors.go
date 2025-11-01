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

// Session errors
var (
	ErrSessionUserIDInvalid     = &appErrors.ValidationError{Message: "ID de usuário inválido"}
	ErrSessionRefreshTokenEmpty = &appErrors.ValidationError{Message: "refresh token não pode estar vazio"}
	ErrSessionExpiresAtPast     = &appErrors.ValidationError{Message: "data de expiração não pode estar no passado"}
	ErrSessionNotFound          = &appErrors.ValidationError{Message: "sessão não encontrada"}
	ErrSessionInvalid           = &appErrors.ValidationError{Message: "sessão inválida ou expirada"}
	ErrSessionIDInvalid         = &appErrors.ValidationError{Message: "ID de sessão inválido"}
	ErrSessionNil               = &appErrors.ValidationError{Message: "sessão não pode ser nula"}
)
