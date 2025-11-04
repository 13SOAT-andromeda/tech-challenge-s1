package repository

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type OrderProductRepository struct {
	*BaseRepository[order_product.Model]
}

func NewOrderProductRepository(db *gorm.DB) ports.OrderProductRepository {
	return &OrderProductRepository{
		BaseRepository: NewBaseRepository[order_product.Model](db),
	}
}
