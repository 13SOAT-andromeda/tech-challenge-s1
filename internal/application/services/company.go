package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CompanyService struct {
	repo ports.CompanyRepository
}

func NewCompanyService(repo ports.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) CreateCompany(ctx context.Context, c domain.Company) (*domain.Company, error) {
	c.EnsureAddressCompany()

	companyModel := model.CompanyFromDomain(c)

	createModel, err := s.repo.Create(ctx, &companyModel)
	if err != nil {
		return nil, err
	}

	createdCompany := model.CompanyToDomain(*createModel)
	return &createdCompany, nil
}

func (s *CompanyService) GetCompanyById(ctx context.Context, id uint) (*domain.Company, error) {
	companyModel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	companyDomain := model.CompanyToDomain(*companyModel)
	return &companyDomain, nil
}

func (s *CompanyService) UpdateCompanyById(ctx context.Context, id uint, c domain.Company) (*domain.Company, error) {
	c.EnsureAddressCompany()

	companyModel := model.CompanyFromDomain(c)

	updatedModel, err := s.repo.Update(ctx, id, &companyModel)
	if err != nil {
		return nil, err
	}

	updatedCompany := model.CompanyToDomain(*updatedModel)
	return &updatedCompany, nil
}

func (s *CompanyService) DeleteCompanyById(ctx context.Context, id uint) (*domain.Company, error) {
	deletedModel, err := s.repo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	deletedCompany := model.CompanyToDomain(*deletedModel)
	return &deletedCompany, nil
}
