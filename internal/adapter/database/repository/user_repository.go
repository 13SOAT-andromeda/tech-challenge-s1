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
	if params.Document != "" {
		q = q.Where("document LIKE ?", "%"+params.Document+"%")
	}
	if params.Contact != "" {
		q.Where("lower(contact) LIKE ?", "%"+strings.ToLower(params.Contact)+"%")
	}
	q.Find(&users)

	return users
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*user.Model, error) {
	user := user.Model{}
	err := u.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByDocument(ctx context.Context, document string) (*user.Model, error) {
	var data user.Model

	err := r.BaseRepository.db.WithContext(ctx).Unscoped().Where("document = ?", document).First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &data, nil
}
