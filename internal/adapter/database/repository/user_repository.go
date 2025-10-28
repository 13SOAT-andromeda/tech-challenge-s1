package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	*BaseRepository[user.Model]
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[user.Model](db),
	}
}

func (u *userRepository) Search(ctx context.Context, params ports.UserSearch) []user.Model {
	users := []user.Model{}
	q := u.db.Model(&users)
	if params.Name != "" {
		q.Where("lower(name) LIKE ?", "%"+strings.ToLower(params.Name)+"%")
	}
	if params.Email != "" {
		q.Where("lower(email) LIKE ?", "%"+strings.ToLower(params.Email)+"%")
	}
	if params.Contact != "" {
		q.Where("lower(contact) LIKE ?", "%"+strings.ToLower(params.Contact)+"%")
	}
	q.Find(&users)

	return users
}

func (u *userRepository) Exists(ctx context.Context, id uint, email string) (bool, error) {
	user := user.Model{}
	err := u.db.Where("id <> ? AND email = ?", id, email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u *userRepository) Delete(ctx context.Context, id uint) error {
	return u.db.Model(&user.Model{}).Where("id = ?", id).Update("active", false).Error
}
