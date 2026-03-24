package ports

import (
	"context"

	employeeModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/employee"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type EmployeeRepository interface {
	Repository[employeeModel.Model]
	GetByPersonID(ctx context.Context, personID uint) (*employeeModel.Model, error)
}

type EmployeeService interface {
	Create(ctx context.Context, e domain.Employee) (*domain.Employee, error)
	GetByID(ctx context.Context, id uint) (*domain.Employee, error)
	GetByPersonID(ctx context.Context, personID uint) (*domain.Employee, error)
	Delete(ctx context.Context, id uint) error
}
