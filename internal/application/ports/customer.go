package ports

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
)

type CustomerRepository interface {
	Repository[customer.Model]
	FindByEmail(ctx context.Context, email string) (*customer.Model, error)
	FindByDocument(ctx context.Context, document string) (*customer.Model, error)

	Search(ctx context.Context, filters filter.CustomerFilter) ([]customer.Model, error)
}

type CustomerService interface {
	Create(ctx context.Context, c domain.Customer) (*domain.Customer, error)
	Search(ctx context.Context, filter *filter.CustomerFilter) ([]domain.Customer, error)
	GetByID(ctx context.Context, id uint) (*domain.Customer, error)
	UpdateByID(ctx context.Context, id uint, c domain.Customer) error
	DeleteByID(ctx context.Context, id uint) (*domain.Customer, error)
}

type CustomerUseCase interface {
	AddVehicleToCustomer(ctx context.Context, customerID, vehicleID uint) error
	RemoveVehicleFromCustomer(ctx context.Context, customerID, vehicleID uint) error
	GetCustomerVehicles(ctx context.Context, customerID uint) ([]domain.Vehicle, error)
}
