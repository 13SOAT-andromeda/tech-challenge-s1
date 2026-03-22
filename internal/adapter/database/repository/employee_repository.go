package repository

import (
	"context"

	employeeModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/employee"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type employeeRepository struct {
	*BaseRepository[employeeModel.Model]
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) ports.EmployeeRepository {
	return &employeeRepository{
		BaseRepository: NewBaseRepository[employeeModel.Model](db),
		db:             db,
	}
}

func (r *employeeRepository) GetByPersonID(ctx context.Context, personID uint) (*employeeModel.Model, error) {
	var entity employeeModel.Model
	err := r.db.WithContext(ctx).Where("person_id = ?", personID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
