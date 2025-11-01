package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

var (
	ErrProductNotFound             = errors.New("Product not found")
	ErrInvalidQuantity             = errors.New("Invalid quantity")
	ErrInvalidProductId            = errors.New("Invalid product Id")
	ErrInvalidManageStockOperation = errors.New("Invalid manage stock operation")
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
	return s.updateInternal(ctx, p, false)
}

func (s *productService) UpdateStock(ctx context.Context, p domain.Product) (*domain.Product, error) {
	return s.updateInternal(ctx, p, true)
}

func (s *productService) updateInternal(ctx context.Context, p domain.Product, allowStockChange bool) (*domain.Product, error) {
	existentProduct, err := s.repo.FindByID(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("product with ID %d not found or disabled", p.ID)
	}

	var model product.Model
	model.FromDomain(&p)
	model.CreatedAt = existentProduct.CreatedAt
	model.DeletedAt = existentProduct.DeletedAt

	if !allowStockChange {
		model.Stock = existentProduct.Stock
	}

	if err := s.repo.Update(ctx, &model); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &p, nil
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

func (s *productService) ManageStockItem(ctx context.Context, productID uint, quantity uint, operation string) (*domain.Product, error) {
	if productID == 0 {
		return nil, ErrInvalidProductId
	}

	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	normalizedOperation := strings.ToLower(operation)

	switch normalizedOperation {
	case "add":
		return s.AddStockItem(ctx, productID, quantity)
	case "remove":
		return s.RemoveStockItem(ctx, productID, quantity)
	default:
		return nil, ErrInvalidManageStockOperation
	}
}

func (s *productService) AddStockItem(ctx context.Context, productID uint, quantity uint) (*domain.Product, error) {

	if productID == 0 {
		return nil, ErrInvalidProductId
	}

	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	product, err := s.GetById(ctx, productID)

	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, ErrProductNotFound
	}

	product.Stock += quantity

	updated, err := s.Update(ctx, *product)

	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *productService) RemoveStockItem(ctx context.Context, productID uint, quantity uint) (*domain.Product, error) {

	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	product, err := s.GetById(ctx, productID)

	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, ErrProductNotFound
	}

	if err := product.DecreaseStock(quantity); err != nil {
		return nil, err
	}

	updated, err := s.Update(ctx, *product)

	if err != nil {
		return nil, err
	}

	return updated, nil
}
