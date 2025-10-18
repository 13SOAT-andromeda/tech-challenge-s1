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

func (s *CompanyService) Create(ctx context.Context, c domain.Company) (*domain.Company, error) {
	companyModel := model.CompanyFromDomain(c)

	createModel, err := s.repo.Create(ctx, &companyModel)
	if err != nil {
		return nil, err
	}

	createdCompany := model.CompanyToDomain(*createModel)
	return &createdCompany, nil
}

func (s *CompanyService) GetByID(ctx context.Context, id uint) (*domain.Company, error) {
	companyModel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	companyDomain := model.CompanyToDomain(*companyModel)
	return &companyDomain, nil
}

func (s *CompanyService) UpdateByID(ctx context.Context, id uint, c domain.Company) (*domain.Company, error) {
	companyModel := model.CompanyFromDomain(c)

	// Chamada ajustada: Update provavelmente recebe apenas o modelo e retorna um erro
	if err := s.repo.Update(ctx, &companyModel); err != nil {
		return nil, err
	}

	// Converte o modelo atualizado (local) para domínio e retorna
	updatedCompany := model.CompanyToDomain(companyModel)
	return &updatedCompany, nil
}

func (s *CompanyService) DeleteByID(ctx context.Context, id uint) (*domain.Company, error) {
	companyModel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	deletedCompany := model.CompanyToDomain(*companyModel)
	return &deletedCompany, nil
}
