package services

import (
	"context"

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
	company := domain.Company{
		Name:     c.Name,
		Email:    c.Email,
		Document: c.Document,
		Contact:  c.Contact,
		Address: &domain.Address{
			Address:       c.Address.Address,
			AddressNumber: c.Address.AddressNumber,
			City:          c.Address.City,
			Neighborhood:  c.Address.Neighborhood,
			Country:       c.Address.Country,
			ZipCode:       c.Address.ZipCode,
		},
	}

	return s.repo.Save(ctx, company)
}

func (s *CompanyService) GetCompanyByID(ctx context.Context, id uint) (*domain.Company, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id uint, c domain.Company) (*domain.Company, error) {
	company := domain.Company{
		Name:     c.Name,
		Email:    c.Email,
		Document: c.Document,
		Contact:  c.Contact,
		Address: &domain.Address{
			Address:       c.Address.Address,
			AddressNumber: c.Address.AddressNumber,
			City:          c.Address.City,
			Neighborhood:  c.Address.Neighborhood,
			Country:       c.Address.Country,
			ZipCode:       c.Address.ZipCode,
		},
	}

	return s.repo.UpdateByID(ctx, id, company)
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id uint) (*domain.Company, error) {
	return s.repo.DeleteByID(ctx, id)
}
