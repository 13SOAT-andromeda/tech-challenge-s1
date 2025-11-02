package customer

import (
	"context"
	"errors"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type useCase struct {
	customerRepository  ports.CustomerRepository
	customerVehicleRepo ports.CustomerVehicleRepository
	vehicleService      ports.VehicleService
}

func NewCustomerUseCase(customerRepository ports.CustomerRepository, customerVehicleRepo ports.CustomerVehicleRepository, vehicleService ports.VehicleService) *useCase {
	return &useCase{
		customerRepository:  customerRepository,
		customerVehicleRepo: customerVehicleRepo,
		vehicleService:      vehicleService,
	}
}

func (s *useCase) AddVehicleToCustomer(ctx context.Context, customerID, vehicleID uint) error {
	customer, err := s.customerRepository.FindByID(ctx, customerID)

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

func (s *useCase) RemoveVehicleFromCustomer(ctx context.Context, customerID, vehicleID uint) error {
	customer, err := s.customerRepository.FindByID(ctx, customerID)
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

func (s *useCase) GetCustomerVehicles(ctx context.Context, customerID uint) ([]domain.Vehicle, error) {

	customer, err := s.customerRepository.FindByID(ctx, customerID)

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
