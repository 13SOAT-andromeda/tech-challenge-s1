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
