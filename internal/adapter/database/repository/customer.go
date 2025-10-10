package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Save(ctx context.Context, c domain.Customer) (*domain.Customer, error) {
	m := model.FromDomain(c)

	return nil, r.db.Create(&m).Error
}

func (r *CustomerRepository) FindAll(ctx context.Context) ([]domain.Customer, error) {
	var models []model.CustomerModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	customers := make([]domain.Customer, len(models))
	for i, m := range models {
		customers[i] = model.ToDomain(m)
	}
	return customers, nil
}

func (r *CustomerRepository) FindByID(ctx context.Context, id uint) (*domain.Customer, error) {
	var m model.CustomerModel
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}

	customer := model.ToDomain(m)
	return &customer, nil
}
