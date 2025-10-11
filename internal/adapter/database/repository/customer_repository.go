package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type customerRepository struct {
	*BaseRepository[domain.Customer]
}

func NewCustomerRepository(db *gorm.DB) ports.CustomerRepository {
	return &customerRepository{
		BaseRepository: NewBaseRepository[domain.Customer](db),
	}
}

func (r *customerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.BaseRepository.db.WithContext(ctx).Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}
