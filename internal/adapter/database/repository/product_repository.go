package repository

import (
	"context"
	"errors"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
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

func (r *productRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model product.Model
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&model, id).Error; err != nil {
			return err
		}

		if quantity < 0 && model.Stock < uint(-quantity) {
			return errors.New("not enough stock")
		}

		newStock := int(model.Stock) + quantity
		if newStock < 0 {
			return errors.New("stock would become negative")
		}

		model.Stock = uint(newStock)
		return tx.Model(&model).Where("id = ?", id).Select("Stock").Updates(&model).Error
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
