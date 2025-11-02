package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
)

var (
	ErrProductNotFound             = errors.New("product not found")
	ErrInvalidQuantity             = errors.New("invalid quantity")
	ErrInvalidProductId            = errors.New("invalid product Id")
	ErrInvalidManageStockOperation = errors.New("invalid manage stock operation")
)

type ProductService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) checkAvailability(ctx context.Context, productID uint, quantity uint) error {

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

func (s *ProductService) ConfirmOrderProducts(ctx context.Context, id uint, quantity int) error {
	return s.repo.UpdateStock(ctx, id, quantity)
}

func (s *ProductService) Update(ctx context.Context, p domain.Product) (*domain.Product, error) {
	return s.updateInternal(ctx, p, false)
}

func (s *ProductService) UpdateStock(ctx context.Context, p domain.Product) (*domain.Product, error) {
	return s.updateInternal(ctx, p, true)
}

func (s *ProductService) updateInternal(ctx context.Context, p domain.Product, allowStockChange bool) (*domain.Product, error) {
	existentProduct, err := s.repo.FindByID(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("product with ID %d not found or disabled", p.ID)
	}

	var model product.Model
	model.FromDomain(&p)
	model.CreatedAt = existentProduct.CreatedAt
	model.DeletedAt = existentProduct.DeletedAt
	model.Stock = existentProduct.Stock

	if allowStockChange {
		model.Stock = p.Stock
	}

	if err := s.repo.Update(ctx, &model); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &p, nil
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

func (s *ProductService) GetAll(ctx context.Context, productFilter *filter.ProductFilter) ([]domain.Product, error) {

	if productFilter == nil {
		productFilter = &filter.ProductFilter{}
	}

	records, err := s.repo.Search(ctx, *productFilter)
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

func (s *ProductService) ManageStockItem(ctx context.Context, productID uint, quantity uint, operation string) (*domain.Product, error) {
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

func (s *ProductService) AddStockItem(ctx context.Context, productID uint, quantity uint) (*domain.Product, error) {

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

	updated, err := s.UpdateStock(ctx, *product)

	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *ProductService) RemoveStockItem(ctx context.Context, productID uint, quantity uint) (*domain.Product, error) {

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

	updated, err := s.UpdateStock(ctx, *product)

	if err != nil {
		return nil, err
	}

	return updated, nil
}
