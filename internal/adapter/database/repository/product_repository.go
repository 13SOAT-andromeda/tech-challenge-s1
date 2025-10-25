package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	*BaseRepository[product.Model]
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &ProductRepository{
		BaseRepository: NewBaseRepository[product.Model](db),
	}
}

func (r *ProductRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model product.Model

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&model, id).Error; err != nil {
			return err
		}

		model.Stock -= uint(quantity)
		return tx.Save(&model).Error
	})
}

func (r *ProductRepository) FindByIDs(ctx context.Context, productIDs []uint) ([]product.Model, error) {
	var products []product.Model

	err := r.db.WithContext(ctx).Where("id IN ?", productIDs).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
