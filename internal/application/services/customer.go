package services

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
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

	customerModel := model.FromDomainCustomer(&c)

	createdModel, err := s.repo.Create(ctx, customerModel)

	if err != nil {
		return nil, err
	}

	createdCustomer := createdModel.ToDomain()

	return createdCustomer, nil
}

func (s *CustomerService) GetAll(ctx context.Context) ([]domain.Customer, error) {

	customerModels, err := s.repo.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	domainCustomers := make([]domain.Customer, 0, len(customerModels))

	for _, customerModel := range customerModels {
		domainCustomers = append(domainCustomers, *(&customerModel).ToDomain())
	}

	return domainCustomers, nil
}

func (s *CustomerService) GetByID(ctx context.Context, id uint) (*domain.Customer, error) {

	customerModel, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}
	customerDomain := customerModel.ToDomain()

	return customerDomain, nil
}
