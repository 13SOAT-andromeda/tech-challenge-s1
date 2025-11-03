package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/errors"
)

var (
	ErrOrderIdInvalid = &errors.ValidationError{Message: "Order Id invalid"}
	ErrOrderNotFound  = &errors.ValidationError{Message: "Order not found"}
	ErrOrderDelete    = &errors.ValidationError{Message: "An error occurred while trying to delete the order"}
)

type OrderService struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, o domain.Order) (*domain.Order, error) {
	model := order.Model{}
	model.FromDomain(&o)

	_, err := s.repo.Create(ctx, &model)

	if err != nil {
		return nil, err
	}

	created := model.ToDomain()

	return created, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (*domain.Order, error) {
	result, err := s.repo.FindOrderByID(ctx, id)

	if err != nil {
		return nil, err
	}
	o := result.ToDomain()

	return o, nil
}

func (s *OrderService) GetAll(ctx context.Context, params map[string]interface{}) (*[]domain.Order, error) {
	oSearch := ports.OrderSearch{Status: "", Enabled: true}

	if params["status"] != nil {
		oSearch.Status = params["status"].(string)
	}

	if params["enabled"] != nil {
		oSearch.Enabled = params["enabled"].(bool)
	}

	orders, err := s.repo.Search(ctx, oSearch)
	if err != nil {
		return nil, err
	}
	ordersD := make([]domain.Order, 0, len(orders))

	for _, item := range orders {
		ordersD = append(ordersD, *item.ToDomain())
	}

	return &ordersD, nil
}

func (s *OrderService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return ErrOrderDelete
	}

	return nil
}

func (s *OrderService) Update(ctx context.Context, o domain.Order) error {
	model := order.Model{}
	model.FromDomain(&o)

	err := s.repo.Update(ctx, &model)
	if err != nil {
		return err
	}

	return nil
}
