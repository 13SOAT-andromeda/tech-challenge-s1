package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type CustomerService struct {
	repo ports.CustomerRepository
}

func NewCustomerService(repo ports.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(ctx context.Context, c domain.Customer) (*domain.Customer, error) {
	customer := domain.Customer{
		Name:     c.Name,
		Email:    c.Email,
		Document: c.Document,
		Type:     c.Type,
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

	return s.repo.Save(ctx, customer)
}

func (s *CustomerService) GetAll(ctx context.Context) ([]domain.Customer, error) {
	return s.repo.FindAll(ctx)
}

func (s *CustomerService) GetByID(ctx context.Context, id uint) (*domain.Customer, error) {
	return s.repo.FindByID(ctx, id)
}
