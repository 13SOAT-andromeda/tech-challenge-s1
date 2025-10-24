package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
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

func (r *customerRepository) FindByEmail(ctx context.Context, email string) (*customer.Model, error) {
	var data customer.Model

	err := r.BaseRepository.db.WithContext(ctx).Where("email = ?", email).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *customerRepository) FindByDocument(ctx context.Context, document string) (*customer.Model, error) {
	var data customer.Model

	err := r.BaseRepository.db.WithContext(ctx).Unscoped().Where("document = ?", document).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
