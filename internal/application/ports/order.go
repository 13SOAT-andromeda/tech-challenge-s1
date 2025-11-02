package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type OrderSearch struct {
	Status  string
	Enabled bool
}

type OrderRepository interface {
	Repository[order.Model]
	Search(ctx context.Context, params OrderSearch) ([]order.Model, error)
}

type OrderService interface {
	Create(ctx context.Context, u domain.Order) (*domain.Order, error)
	GetAll(ctx context.Context, params map[string]interface{}) (*[]domain.Order, error)
	GetByID(ctx context.Context, id uint) (*domain.Order, error)
	Delete(ctx context.Context, id uint) error
}
