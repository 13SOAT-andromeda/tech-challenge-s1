package services

import (
	"context"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type ProductService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CheckAvailability(ctx context.Context, id uint, quantity uint) error {
	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	p := result.ToDomain()

	if err = p.CanBePurchased(quantity); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) CheckProductPrice(ctx context.Context, productIDs []uint) (map[uint]float64, error) {
	if len(productIDs) == 0 {
		return make(map[uint]float64), nil
	}

	products, err := s.repo.FindByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by IDs: %w", err)
	}

	priceMap := make(map[uint]float64, len(products))
	for _, p := range products {
		priceMap[p.ID] = float64(p.Price)
	}

	return priceMap, nil
}

//func (s *ProductService) CreateOrderProducts(ctx context.Context, productIDs []uint) (map[uint]float64, error) {
//
//}

func (s *ProductService) ConfirmOrderProducts(ctx context.Context, id uint, quantity int) error {
	return s.repo.UpdateStock(ctx, id, quantity)
}

func (s *ProductService) GetById(ctx context.Context, id uint) (*domain.Product, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]domain.Product, error) {
	records, err := s.repo.FindAll(ctx, false)
	if err != nil {
		return nil, err
	}

	var result []domain.Product
	for _, record := range records {
		result = append(result, *record.ToDomain())
	}

	return result, nil
}

func (s *ProductService) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	var model product.Model
	model.FromDomain(&p)

	response, err := s.repo.Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
