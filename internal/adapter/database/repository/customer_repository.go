package repository

import (
	"context"
	"errors"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"gorm.io/gorm"
)

type customerRepository struct {
	*BaseRepository[customer.Model]
}

func NewCustomerRepository(db *gorm.DB) ports.CustomerRepository {
	return &customerRepository{
		BaseRepository: NewBaseRepository[customer.Model](db),
	}
}

func (r *customerRepository) FindByID(ctx context.Context, id uint) (*customer.Model, error) {
	var entity customer.Model
	err := r.BaseRepository.db.WithContext(ctx).Joins("Person").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *customerRepository) FindByEmail(ctx context.Context, email string) (*customer.Model, error) {
	var data customer.Model

	err := r.BaseRepository.db.WithContext(ctx).
		Joins("Person").
		Where(`"Person"."email" = ?`, email).
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *customerRepository) FindByDocument(ctx context.Context, document string) (*customer.Model, error) {
	var data customer.Model

	err := r.BaseRepository.db.WithContext(ctx).Unscoped().
		Joins("Person").
		Where(`"Person"."document" = ?`, document).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *customerRepository) Search(ctx context.Context, filters filter.CustomerFilter) ([]customer.Model, error) {
	var model []customer.Model

	db := r.db.WithContext(ctx).Joins("Person")

	if !filters.Status {
		db = db.Unscoped()
	}

	if filters.Document != nil {
		db = db.Where(`"Person"."document" = ?`, *filters.Document)
	}

	if filters.Name != nil {
		db = db.Where(`"Person"."name" LIKE ?`, "%"+*filters.Name+"%")
	}

	if filters.Email != nil {
		db = db.Where(`"Person"."email" LIKE ?`, "%"+*filters.Email+"%")
	}

	err := db.Find(&model).Error

	return model, err
}
