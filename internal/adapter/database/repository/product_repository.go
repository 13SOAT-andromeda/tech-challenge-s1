package repository

import (
	"context"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	*BaseRepository[product.Model]
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &productRepository{
		BaseRepository: NewBaseRepository[product.Model](db),
	}
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint, quantity int, operation domain.StockOperation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model product.Model

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&model, id).Error; err != nil {
			return err
		}

		if operation == domain.StockOperationAdd {
			model.Stock += uint(quantity)
		} else if operation == domain.StockOperationRemove {
			if model.Stock < uint(quantity) {
				return fmt.Errorf("insufficient stock to remove %d items", quantity)
			}

			model.Stock -= uint(quantity)
		} else {
			return fmt.Errorf("unknown operation %q", operation)
		}

		return tx.Save(&model).Error
	})
}

func (r *productRepository) FindByIDs(ctx context.Context, productIDs []uint) ([]product.Model, error) {
	var products []product.Model

	err := r.db.WithContext(ctx).Where("id IN ?", productIDs).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
