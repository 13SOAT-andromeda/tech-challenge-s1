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
	model := company.FromDomain(c)

	response, err := s.repo.Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	result := company.ToDomain(*response)
	return &result, nil
}

func (s *CompanyService) GetByID(ctx context.Context, id uint) (*domain.Company, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := company.ToDomain(*response)
	return &result, nil
}

func (s *CompanyService) UpdateByID(ctx context.Context, id uint, c domain.Company) error {

	c.ID = id

	ent := company.FromDomain(c)

	// note: check this doc https://gorm.io/docs/update.html
	findedCompany, err := s.repo.First(ctx, &ent)
	if err != nil {
		return err
	}

	model := company.FromDomain(company.ToDomain(*findedCompany))

	err = s.repo.Update(ctx, &model)
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

	result := company.ToDomain(*response)
	return &result, nil
}
