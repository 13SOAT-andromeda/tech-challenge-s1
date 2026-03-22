package mocks

import (
	"context"

	employeeModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/employee"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
)

type MockEmployeeRepository struct {
	MockGenericRepository[employeeModel.Model]
}

var _ ports.EmployeeRepository = (*MockEmployeeRepository)(nil)

func (m *MockEmployeeRepository) GetByPersonID(ctx context.Context, personID uint) (*employeeModel.Model, error) {
	args := m.Called(ctx, personID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*employeeModel.Model), args.Error(1)
}
