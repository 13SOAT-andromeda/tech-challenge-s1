package errors

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Password errors
var (
	ErrPasswordTooShort      = &ValidationError{Message: "senha deve ter pelo menos 8 caracteres"}
	ErrPasswordHash          = &ValidationError{Message: "erro ao criar o hash da senha"}
	ErrPasswordInvalidFormat = &ValidationError{Message: "senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial"}
	ErrPasswordInvalid       = &ValidationError{Message: "senha inválida"}
)

// User errors
var (
	ErrUserIdInvalid                = &ValidationError{Message: "ID de usuário inválido"}
	ErrUserNotFound                 = &ValidationError{Message: "usuário não encontrado"}
	ErrUserEmailAlreadyExists       = &ValidationError{Message: "email já existe"}
	ErrUserPasswordUpdateNotAllowed = &ValidationError{Message: "senha de usuário não pode ser atualizada"}
	ErrUserDelete                   = &ValidationError{Message: "ocorreu um erro ao excluir o usuário"}
)
