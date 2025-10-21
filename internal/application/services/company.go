package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CompanyService struct {
	repo ports.CompanyRepository
}

func NewCompanyService(repo ports.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) Create(ctx context.Context, c domain.Company) (*domain.Company, error) {
	var model company.Model
	model.FromDomain(&c)

	response, err := s.repo.Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *CompanyService) GetByID(ctx context.Context, id uint) (*domain.Company, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}

func (s *CompanyService) UpdateByID(ctx context.Context, id uint, c domain.Company) error {
	var model company.Model
	model.FromDomain(&c)

	err := s.repo.Update(ctx, &model)
	if err != nil {
		return err
	}

	return nil
}

func (s *CompanyService) DeleteByID(ctx context.Context, id uint) (*domain.Company, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	result := response.ToDomain()
	return result, nil
}
