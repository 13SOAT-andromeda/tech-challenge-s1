package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	*BaseRepository[model.UserModel]
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[model.UserModel](db),
	}
}

func (u *userRepository) Search(ctx context.Context, params ports.UserSearch) []model.UserModel {
	users := []model.UserModel{}
	u.db.Where("name LIKE ? OR email LIKE ? OR contact LIKE ?", params.Name, params.Email, params.Contact).Find(&users)
	return users
}
