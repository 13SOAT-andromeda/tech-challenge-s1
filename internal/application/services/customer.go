package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"gorm.io/gorm"
)

type customerService struct {
	repo ports.CustomerRepository
}

func NewCustomerService(repo ports.CustomerRepository) *customerService {
	return &customerService{repo: repo}
}

func (s *customerService) Create(ctx context.Context, c domain.Customer) (*domain.Customer, error) {

	existentCustomer, err := s.repo.FindByDocument(ctx, c.Document.GetDocumentNumber())

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existentCustomer != nil {
		return nil, errors.New("Customer already exists")
	}

	var model customer.Model
	model.FromDomain(&c)

	response, err := s.repo.Create(ctx, &model)

	if err != nil {
		return nil, err
	}

	result := response.ToDomain()

	return result, nil
}

func (s *customerService) UpdateByID(ctx context.Context, id uint, c domain.Customer) error {

	_, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("Customer with Id %d not found or disabled", id)
	}

	doc := c.Document.GetDocumentNumber()

	if doc != "" {
		other, err := s.repo.FindByDocument(ctx, doc)
		if err == nil && other.ID != id {
			return fmt.Errorf("The customer cannot be updated. Number is invalid or already in use to another customer")
		}
	}

	var model customer.Model

	model.FromDomain(&c)

	if err := s.repo.Update(ctx, &model); err != nil {
		return fmt.Errorf("Failed to update customer: %w", err)
	}

	return nil
}

func (s *customerService) Search(ctx context.Context, customerFilter *filter.CustomerFilter) ([]domain.Customer, error) {

	if customerFilter == nil {
		customerFilter = &filter.CustomerFilter{}
	}

	customerModels, err := s.repo.Search(ctx, *customerFilter)

	if err != nil {
		return nil, err
	}

	domainCustomers := make([]domain.Customer, 0, len(customerModels))

	for _, customerModel := range customerModels {
		domainCustomers = append(domainCustomers, *customerModel.ToDomain())
	}

	return domainCustomers, nil
}

func (s *customerService) GetByID(ctx context.Context, id uint) (*domain.Customer, error) {
	response, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}
	result := response.ToDomain()

	return result, nil
}

func (s *customerService) DeleteByID(ctx context.Context, id uint) (*domain.Customer, error) {

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
