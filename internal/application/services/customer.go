package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"gorm.io/gorm"
)

type customerService struct {
	repo                ports.CustomerRepository
	customerVehicleRepo ports.CustomerVehicleRepository
	vehicleService      ports.VehicleService
}

func NewCustomerService(repo ports.CustomerRepository, customerVehicleRepo ports.CustomerVehicleRepository, vehicleService ports.VehicleService) *customerService {
	return &customerService{
		repo:                repo,
		customerVehicleRepo: customerVehicleRepo,
		vehicleService:      vehicleService,
	}
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

	existentCustomer, err := s.repo.FindByID(ctx, id)

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
	model.CreatedAt = existentCustomer.CreatedAt
	model.DeletedAt = existentCustomer.DeletedAt

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

func (s *customerService) AddVehicleToCustomer(ctx context.Context, customerID, vehicleID uint) error {
	customer, err := s.repo.FindByID(ctx, customerID)

	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	if customer == nil {
		return fmt.Errorf("customer with ID %d not found", customerID)
	}

	vehicle, err := s.vehicleService.GetByID(ctx, vehicleID)

	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}

	if vehicle == nil {
		return fmt.Errorf("vehicle with ID %d not found", vehicleID)
	}

	existing, err := s.customerVehicleRepo.FindByCustomerAndVehicle(ctx, customerID, vehicleID)
	if err != nil {
		return fmt.Errorf("error checking existing association: %w", err)
	}
	if existing != nil {
		return errors.New("vehicle is already associated with this customer")
	}

	customerVehicle := &customer_vehicle.Model{
		Model:      gorm.Model{},
		CustomerId: customerID,
		VehicleId:  vehicleID,
	}

	_, err = s.customerVehicleRepo.Create(ctx, customerVehicle)
	if err != nil {
		return fmt.Errorf("error creating customer-vehicle association: %w", err)
	}

	return nil
}

func (s *customerService) RemoveVehicleFromCustomer(ctx context.Context, customerID, vehicleID uint) error {
	customer, err := s.repo.FindByID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	if customer == nil {
		return fmt.Errorf("customer with ID %d not found", customerID)
	}

	vehicle, err := s.vehicleService.GetByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}
	if vehicle == nil {
		return fmt.Errorf("vehicle with ID %d not found", vehicleID)
	}

	err = s.customerVehicleRepo.DeleteByCustomerAndVehicle(ctx, customerID, vehicleID)
	if err != nil {
		return fmt.Errorf("error removing customer-vehicle association: %w", err)
	}

	return nil
}

func (s *customerService) GetCustomerVehicles(ctx context.Context, customerID uint) ([]domain.Vehicle, error) {
	customer, err := s.repo.FindByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	if customer == nil {
		return nil, fmt.Errorf("customer with ID %d not found", customerID)
	}

	customerVehicles, err := s.customerVehicleRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("error fetching customer vehicles: %w", err)
	}

	vehicles := make([]domain.Vehicle, 0, len(customerVehicles))
	for _, cv := range customerVehicles {
		if cv.Vehicle.ID != 0 {
			vehicles = append(vehicles, *cv.Vehicle.ToDomain())
		}
	}

	return vehicles, nil
}
