package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type customerRepository struct {
	*BaseRepository[model.CustomerModel]
}

func NewCustomerRepository(db *gorm.DB) ports.CustomerRepository {
	return &customerRepository{
		BaseRepository: NewBaseRepository[model.CustomerModel](db),
	}
}

func (r *customerRepository) FindByEmail(ctx context.Context, email string) (*model.CustomerModel, error) {
	var customer model.CustomerModel
	err := r.BaseRepository.db.WithContext(ctx).Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}
