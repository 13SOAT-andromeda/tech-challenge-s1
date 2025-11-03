package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

var (
	ErrProductNotFound             = errors.New("product not found")
	ErrInvalidQuantity             = errors.New("invalid quantity")
	ErrInvalidProductId            = errors.New("invalid product Id")
	ErrInvalidManageStockOperation = errors.New("invalid manage stock operation")
	ErrInsufficientItems           = errors.New("insufficient products to perform the operation")
)

type ProductService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetById(ctx context.Context, productID uint) (*domain.Product, error) {
	response, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *ProductService) GetByIds(ctx context.Context, productIDs []uint) ([]domain.Product, error) {
	records, err := s.repo.FindByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	var result []domain.Product
	for _, record := range records {
		result = append(result, *record.ToDomain())
	}

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

func (s *ProductService) Update(ctx context.Context, p domain.Product) (*domain.Product, error) {
	var model product.Model
	model.FromDomain(&p)

	if err := s.repo.Update(ctx, &model); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	updated, err := s.repo.FindByID(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated product: %w", err)
	}

	result := updated.ToDomain()
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

func (s *ProductService) Delete(ctx context.Context, productID uint) (*domain.Product, error) {
	response, err := s.repo.FindByID(ctx, productID)

	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, productID); err != nil {
		return nil, err
	}

	result := response.ToDomain()

	return result, nil
}

func (s *ProductService) CheckAvailability(ctx context.Context, productID uint, quantity uint) (bool, error) {
	result, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return false, err
	}

	p := result.ToDomain()

	if err = p.CanBePurchased(quantity); err != nil {
		return false, err
	}

	return true, nil
}

func (s *ProductService) UpdateStock(ctx context.Context, products []domain.ProductItem, operation domain.StockOperation) error {
	for _, item := range products {
		err := s.repo.UpdateStock(ctx, item.ID, int(item.Quantity), operation)
		if err != nil {
			return err
		}
	}

	return nil
}
