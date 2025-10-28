package services

import (
	"context"

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

//func (s *ProductService) GetByName(ctx context.Context, name string) (*domain.Product, error) {
//	response, err := s.repo.FindByName(ctx, name)
//	if err != nil {
//		return nil, err
//	}
//
//	result := response.ToDomain()
//	return result, nil
//}

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
