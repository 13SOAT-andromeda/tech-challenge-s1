package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type UserSearch struct {
	Name    string
	Email   string
	Contact string
}

type UserRepository interface {
	Repository[model.UserModel]
	Search(ctx context.Context, params UserSearch) []model.UserModel
	Exists(ctx context.Context, id uint, email string) (bool, error)
}

type UserService interface {
	Create(ctx context.Context, u domain.User) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	Search(ctx context.Context, params map[string]interface{}) (*[]domain.User, error)
	Update(ctx context.Context, u domain.User) (*domain.User, error)
	Delete(ctx context.Context, id uint) error
}
