package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type UserSearch struct {
	Name    string
	Email   string
	Contact string
}

type UserRepository interface {
	Repository[user.Model]
	Search(ctx context.Context, params UserSearch) []user.Model
	GetByEmail(ctx context.Context, email string) (*user.Model, error)
}

type UserService interface {
	Create(ctx context.Context, u domain.User) (*domain.User, error)
	CreateAdminUser(ctx context.Context, email, password string) error
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Search(ctx context.Context, params map[string]interface{}) (*[]domain.User, error)
	Update(ctx context.Context, u domain.User) (*domain.User, error)
	Delete(ctx context.Context, id uint) error
}
