package ports

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
)

type OrderRepository interface {
	Repository[order.Model]
}

type OrderService interface {
}
