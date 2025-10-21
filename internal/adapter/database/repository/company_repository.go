package repository

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type companyRepository struct {
	*BaseRepository[company.Model]
}

func NewCompanyRepository(db *gorm.DB) ports.CompanyRepository {
	return &companyRepository{
		BaseRepository: NewBaseRepository[company.Model](db),
	}
}
