package session

import (
	"context"
)

// UseCases interfaces
type LoginUseCase interface {
	Execute(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

type ValidateUseCase interface {
	Execute(ctx context.Context, input ValidateInput) (*ValidateOutput, error)
}

type RefreshUseCase interface {
	Execute(ctx context.Context, input RefreshInput) (*RefreshOutput, error)
}

type LogoutUseCase interface {
	Execute(ctx context.Context, input LogoutInput) error
}
