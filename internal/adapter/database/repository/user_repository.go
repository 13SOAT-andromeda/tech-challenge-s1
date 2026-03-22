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

func (u *userRepository) FindByID(ctx context.Context, id uint) (*user.Model, error) {
	var entity user.Model
	err := u.db.WithContext(ctx).Joins("Person").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (u *userRepository) Search(ctx context.Context, params ports.UserSearch) []user.Model {
	users := []user.Model{}
	q := u.db.Joins("Person")
	if params.Name != "" {
		q = q.Where(`lower("Person"."name") LIKE ?`, "%"+strings.ToLower(params.Name)+"%")
	}
	if params.Email != "" {
		q = q.Where(`lower("Person"."email") LIKE ?`, "%"+strings.ToLower(params.Email)+"%")
	}
	if params.Contact != "" {
		q = q.Where(`lower("Person"."contact") LIKE ?`, "%"+strings.ToLower(params.Contact)+"%")
	}
	q.Find(&users)

	return users
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*user.Model, error) {
	var entity user.Model
	err := u.db.WithContext(ctx).
		Joins("Person").
		Where(`"Person"."email" = ?`, email).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
