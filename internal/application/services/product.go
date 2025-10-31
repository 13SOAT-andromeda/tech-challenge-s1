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
	ErrProductNotFound = errors.New("Product not found")
	ErrInvalidQuantity = errors.New("Invalid Quantity")
)

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) checkAvailability(ctx context.Context, productID uint, quantity uint) error {

	result, err := s.repo.FindByID(ctx, productID)

	if err != nil {
		return err
	}

	p := result.ToDomain()

	if err = p.CanBePurchased(quantity); err != nil {
		return err
	}

	return nil
}

func (s *productService) CheckProductPrice(ctx context.Context, productIDs []uint) (map[uint]float64, error) {
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

func (s *productService) ConfirmOrderProducts(ctx context.Context, id uint, quantity int) error {
	return s.repo.UpdateStock(ctx, id, quantity)
}

func (s *productService) Update(ctx context.Context, p domain.Product) (*domain.Product, error) {

	// todo: implementar update de produto

	return nil, nil
}

func (s *productService) GetById(ctx context.Context, productID uint) (*domain.Product, error) {
	response, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *productService) GetAll(ctx context.Context) ([]domain.Product, error) {
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

func (s *productService) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	var model product.Model
	model.FromDomain(&p)

	response, err := s.repo.Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *productService) Delete(ctx context.Context, productID uint) (*domain.Product, error) {
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

func (s *productService) AddStockItem(ctx context.Context, productID uint, quantity uint) error {
	if quantity == 0 {
		return ErrInvalidQuantity
	}

	product, err := s.GetById(ctx, productID)
	if err != nil {
		return err
	}

	if product == nil {
		return ErrProductNotFound
	}

	product.Stock += quantity

	_, err = s.Update(ctx, *product)
	return err
}

func (s *productService) RemoveStockItem(ctx context.Context, productID uint, quantity uint) error {

	if quantity == 0 {
		return ErrInvalidQuantity
	}

	product, err := s.GetById(ctx, productID)
	if err != nil {
		return err
	}

	if product == nil {
		return ErrProductNotFound
	}

	if err := product.DecreaseStock(quantity); err != nil {
		return err
	}

	_, err = s.Update(ctx, *product)
	return err
}

func (s *productService) GetCurrentStock(ctx context.Context, productID uint) (uint, error) {
	product, err := s.GetById(ctx, productID)

	if err != nil {
		return 0, err
	}

	if product == nil {
		return 0, ErrProductNotFound
	}

	return product.Stock, nil
}

func (s *productService) SetStock(ctx context.Context, productID uint, quantity uint) error {
	product, err := s.GetById(ctx, productID)

	if err != nil {
		return err
	}

	if product == nil {
		return ErrProductNotFound
	}

	product.Stock = quantity

	_, err = s.Update(ctx, *product)

	return err
}
