package services

import (
	"context"
	"time"

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

type orderService struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) *orderService {
	return &orderService{repo: repo}
}

func (s *orderService) Create(ctx context.Context, o domain.Order) (*domain.Order, error) {

	order := &order.Model{}
	order.FromDomain(&o)
	order.Status = string(domain.RECEIVED)
	order.DateIn = time.Now()
	order.DateOut = nil

	_, err := s.repo.Create(ctx, order)

	if err != nil {
		return nil, err
	}

	created := order.ToDomain()

	return created, nil
}

func (s *orderService) GetByID(ctx context.Context, id uint) (*domain.Order, error) {

	order, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}
	o := order.ToDomain()

	return o, nil
}

func (s *orderService) GetAll(ctx context.Context, params map[string]interface{}) (*[]domain.Order, error) {

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

	for _, order := range orders {
		ordersD = append(ordersD, *order.ToDomain())
	}

	return &ordersD, nil
}

func (s *orderService) Delete(ctx context.Context, id uint) error {

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return ErrOrderDelete
	}

	return nil
}
