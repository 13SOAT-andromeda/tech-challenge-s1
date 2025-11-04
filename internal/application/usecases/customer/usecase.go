package customer

import (
	"context"
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/errors"
)

var (
	ErrCustomerNotFound               = &errors.ValidationError{Message: "customer not found"}
	ErrVehicleNotFound                = &errors.ValidationError{Message: "vehicle not found"}
	ErrVehicleAlreadyAssociated       = &errors.ValidationError{Message: "vehicle is already associated with this customer"}
	ErrAssociationCheckFailed         = &errors.ValidationError{Message: "error checking existing association"}
	ErrAssociationCreationFailed      = &errors.ValidationError{Message: "error creating customer-vehicle association"}
	ErrAssociationRemovalFailed       = &errors.ValidationError{Message: "error removing customer-vehicle association"}
	ErrFetchingCustomerVehiclesFailed = &errors.ValidationError{Message: "error fetching customer vehicles"}
)

type UseCase struct {
	customerRepository  ports.CustomerRepository
	customerVehicleRepo ports.CustomerVehicleRepository
	vehicleService      ports.VehicleService
}

func NewCustomerUseCase(customerRepository ports.CustomerRepository, customerVehicleRepo ports.CustomerVehicleRepository, vehicleService ports.VehicleService) *UseCase {
	return &UseCase{
		customerRepository:  customerRepository,
		customerVehicleRepo: customerVehicleRepo,
		vehicleService:      vehicleService,
	}
}

func (s *UseCase) AddVehicleToCustomer(ctx context.Context, customerID, vehicleID uint) error {
	customer, err := s.customerRepository.FindByID(ctx, customerID)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrCustomerNotFound, err)
	}

	if customer == nil {
		return fmt.Errorf("%w: customer with ID %d not found", ErrCustomerNotFound, customerID)
	}

	vehicle, err := s.vehicleService.GetByID(ctx, vehicleID)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrVehicleNotFound, err)
	}

	if vehicle == nil {
		return fmt.Errorf("%w: vehicle with ID %d not found", ErrVehicleNotFound, vehicleID)
	}

	existing, err := s.customerVehicleRepo.FindByCustomerAndVehicle(ctx, customerID, vehicleID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAssociationCheckFailed, err)
	}
	if existing != nil {
		return ErrVehicleAlreadyAssociated
	}

	customerVehicleDomain := &domain.CustomerVehicle{
		CustomerId: customerID,
		VehicleId:  vehicleID,
	}

	model := customer_vehicle.Model{}
	model.FromDomain(customerVehicleDomain)

	customerVehicle := model

	_, err = s.customerVehicleRepo.Create(ctx, &customerVehicle)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAssociationCreationFailed, err)
	}

	return nil
}

func (s *UseCase) RemoveVehicleFromCustomer(ctx context.Context, customerID, vehicleID uint) error {
	customer, err := s.customerRepository.FindByID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCustomerNotFound, err)
	}
	if customer == nil {
		return fmt.Errorf("%w: customer with ID %d not found", ErrCustomerNotFound, customerID)
	}

	vehicle, err := s.vehicleService.GetByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrVehicleNotFound, err)
	}
	if vehicle == nil {
		return fmt.Errorf("%w: vehicle with ID %d not found", ErrVehicleNotFound, vehicleID)
	}

	err = s.customerVehicleRepo.DeleteByCustomerAndVehicle(ctx, customerID, vehicleID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAssociationRemovalFailed, err)
	}

	return nil
}

func (s *UseCase) GetCustomerVehicles(ctx context.Context, customerID uint) ([]domain.CustomerVehicle, error) {

	customer, err := s.customerRepository.FindByID(ctx, customerID)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCustomerNotFound, err)
	}

	if customer == nil {
		return nil, fmt.Errorf("%w: customer with ID %d not found", ErrCustomerNotFound, customerID)
	}

	customerVehicles, err := s.customerVehicleRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFetchingCustomerVehiclesFailed, err)
	}

	cvDomain := make([]domain.CustomerVehicle, 0, len(customerVehicles))
	for _, cv := range customerVehicles {
		cvDomain = append(cvDomain, *cv.ToDomain())
	}

	return cvDomain, nil
}
