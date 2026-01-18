package customer

import (
	"context"
	"errors"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
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
		return errors.New("customer not found")
	}

	if customer == nil {
		return errors.New("customer not found")
	}

	vehicle, err := s.vehicleService.GetByID(ctx, vehicleID)

	if err != nil {
		return errors.New("vehicle not found")
	}

	if vehicle == nil {
		return errors.New("vehicle not found")
	}

	existing, err := s.customerVehicleRepo.FindByCustomerAndVehicle(ctx, customerID, vehicleID)
	if err != nil {
		return errors.New("error checking existing association")
	}
	if existing != nil {
		return errors.New("vehicle is already associated with this customer")
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
		return errors.New("error creating customer-vehicle association")
	}

	return nil
}

func (s *UseCase) RemoveVehicleFromCustomer(ctx context.Context, customerID, vehicleID uint) error {
	customer, err := s.customerRepository.FindByID(ctx, customerID)
	if err != nil {
		return errors.New("customer not found")
	}
	if customer == nil {
		return errors.New("customer not found")
	}

	vehicle, err := s.vehicleService.GetByID(ctx, vehicleID)
	if err != nil {
		return errors.New("vehicle not found")
	}
	if vehicle == nil {
		return errors.New("vehicle not found")
	}

	err = s.customerVehicleRepo.DeleteByCustomerAndVehicle(ctx, customerID, vehicleID)
	if err != nil {
		return errors.New("error removing customer-vehicle association")
	}

	return nil
}

func (s *UseCase) GetCustomerVehicles(ctx context.Context, customerID uint) ([]domain.CustomerVehicle, error) {

	customer, err := s.customerRepository.FindByID(ctx, customerID)

	if err != nil {
		return nil, errors.New("customer not found")
	}

	if customer == nil {
		return nil, errors.New("customer not found")
	}

	customerVehicles, err := s.customerVehicleRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, errors.New("error fetching customer vehicles")
	}

	cvDomain := make([]domain.CustomerVehicle, 0, len(customerVehicles))
	for _, cv := range customerVehicles {
		cvDomain = append(cvDomain, *cv.ToDomain())
	}

	return cvDomain, nil
}
