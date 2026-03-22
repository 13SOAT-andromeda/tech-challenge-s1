package services

import (
	"context"
	"errors"
	"fmt"

	customerModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	userModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
)

type customerService struct {
	repo       ports.CustomerRepository
	personRepo ports.PersonRepository
	userRepo   ports.UserRepository
}

func NewCustomerService(repo ports.CustomerRepository, personRepo ports.PersonRepository, userRepo ports.UserRepository) *customerService {
	return &customerService{
		repo:       repo,
		personRepo: personRepo,
		userRepo:   userRepo,
	}
}

func (s *customerService) Create(ctx context.Context, c domain.Customer, password *domain.Password) (*domain.Customer, error) {
	if c.Person == nil {
		c.Person = &domain.Person{}
	}

	docNumber := ""
	if c.Person.Document != nil {
		docNumber = c.Person.Document.GetDocumentNumber()
	}

	if docNumber != "" {
		existing, err := s.repo.FindByDocument(ctx, docNumber)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("customer already exists")
		}
	}

	if c.Person.Email != "" {
		existingUser, err := s.userRepo.GetByEmail(ctx, c.Person.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("email já existe")
		}
	}

	pm := &personModel.Model{}
	pm.FromDomain(c.Person)
	createdPerson, err := s.personRepo.Create(ctx, pm)
	if err != nil {
		return nil, err
	}

	c.PersonID = createdPerson.ID
	c.Person = createdPerson.ToDomain()

	var model customerModel.Model
	model.FromDomain(&c)

	response, err := s.repo.Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	if err := password.Hash(); err != nil {
		return nil, err
	}

	um := &userModel.Model{
		Password: password.GetHashed(),
		Role:     "customer",
		PersonID: createdPerson.ID,
		Person:   *createdPerson,
	}
	if _, err := s.userRepo.Create(ctx, um); err != nil {
		return nil, err
	}

	response.Person = *createdPerson
	return response.ToDomain(), nil
}

func (s *customerService) UpdateByID(ctx context.Context, id uint, c domain.Customer) error {
	existingCustomer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Customer with Id %d not found or disabled", id)
	}

	if c.Person != nil && c.Person.Document != nil {
		doc := c.Person.Document.GetDocumentNumber()
		if doc != "" {
			other, err := s.repo.FindByDocument(ctx, doc)
			if err == nil && other != nil && other.ID != id {
				return fmt.Errorf("The customer cannot be updated. Number is invalid or already in use to another customer")
			}
		}
	}

	if c.Person != nil {
		existingPerson, err := s.personRepo.FindByID(ctx, existingCustomer.PersonID)
		if err != nil {
			return fmt.Errorf("failed to find customer's person: %w", err)
		}

		existingPersonDomain := existingPerson.ToDomain()
		if c.Person.Name != "" {
			existingPersonDomain.Name = c.Person.Name
		}
		if c.Person.Email != "" {
			existingPersonDomain.Email = c.Person.Email
		}
		if c.Person.Contact != "" {
			existingPersonDomain.Contact = c.Person.Contact
		}
		if c.Person.Document != nil {
			existingPersonDomain.Document = c.Person.Document
		}
		if c.Person.Address != nil {
			existingPersonDomain.Address = c.Person.Address
		}

		pm := &personModel.Model{}
		pm.FromDomain(existingPersonDomain)
		if err := s.personRepo.Update(ctx, pm); err != nil {
			return fmt.Errorf("Failed to update person: %w", err)
		}
	}

	var model customerModel.Model
	model.FromDomain(&c)
	model.ID = id
	model.PersonID = existingCustomer.PersonID
	model.CreatedAt = existingCustomer.CreatedAt
	model.DeletedAt = existingCustomer.DeletedAt

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
	for _, cm := range customerModels {
		domainCustomers = append(domainCustomers, *cm.ToDomain())
	}

	return domainCustomers, nil
}

func (s *customerService) GetByID(ctx context.Context, id uint) (*domain.Customer, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return response.ToDomain(), nil
}

func (s *customerService) DeleteByID(ctx context.Context, id uint) (*domain.Customer, error) {
	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return response.ToDomain(), nil
}
