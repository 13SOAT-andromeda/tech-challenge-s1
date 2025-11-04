package ports

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_product"
)

type OrderProductRepository interface {
	Repository[order_product.Model]
}
