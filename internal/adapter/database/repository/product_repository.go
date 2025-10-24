package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type productRepository struct {
	*BaseRepository[product.Model]
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &productRepository{
		BaseRepository: NewBaseRepository[product.Model](db),
	}
}

func (r *productRepository) FindByName(ctx context.Context, name string) (*product.Model, error) {
	var data product.Model

	searchPattern := "%" + name + "%"

	err := r.BaseRepository.db.WithContext(ctx).Where("name LIKE ?", searchPattern).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
