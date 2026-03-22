package services

import (
	"context"

	employeeModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/employee"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type employeeService struct {
	repo ports.EmployeeRepository
}

func NewEmployeeService(repo ports.EmployeeRepository) ports.EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) Create(ctx context.Context, e domain.Employee) (*domain.Employee, error) {
	m := &employeeModel.Model{}
	m.FromDomain(&e)

	created, err := s.repo.Create(ctx, m)
	if err != nil {
		return nil, err
	}

	return created.ToDomain(), nil
}

func (s *employeeService) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return m.ToDomain(), nil
}

func (s *employeeService) GetByPersonID(ctx context.Context, personID uint) (*domain.Employee, error) {
	m, err := s.repo.GetByPersonID(ctx, personID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (s *employeeService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
