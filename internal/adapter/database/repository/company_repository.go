package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type companyRepository struct {
	*BaseRepository[model.CompanyModel]
}

func NewCompanyRepository(db *gorm.DB) ports.CompanyRepository {
	return &companyRepository{
		BaseRepository: NewBaseRepository[model.CompanyModel](db),
	}
}

func (r *companyRepository) Save(ctx context.Context, c domain.Company) (*domain.Company, error) {
	m := model.CompanyFromDomain(c)
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	entity := model.CompanyToDomain(m)
	return &entity, nil
}

func (r *companyRepository) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
	var m model.CompanyModel
	if err := r.db.Preload("Address").First(&m, id).Error; err != nil {
		return nil, err
	}
	entity := model.CompanyToDomain(m)
	return &entity, nil
}

func (r *companyRepository) UpdateByID(ctx context.Context, id uint, c domain.Company) (*domain.Company, error) {
	var m model.CompanyModel
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}

	m.Name = c.Name
	m.Email = c.Email
	m.Document = c.Document
	m.Contact = c.Contact

	if c.Address != nil {
		m.Address.Address = c.Address.Address
		m.Address.AddressNumber = c.Address.AddressNumber
		m.Address.City = c.Address.City
		m.Address.Neighborhood = c.Address.Neighborhood
		m.Address.Country = c.Address.Country
		m.Address.ZipCode = c.Address.ZipCode
	} else {
		// domínio não possui endereço -> zera o valor do endereço no model
		m.Address = domain.Address{}
	}

	if err := r.db.Save(&m).Error; err != nil {
		return nil, err
	}

	entity := model.CompanyToDomain(m)
	return &entity, nil
}

func (r *companyRepository) DeleteByID(ctx context.Context, id uint) (*domain.Company, error) {
	var m model.CompanyModel
	if err := r.db.Preload("Address").First(&m, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Delete(&m).Error; err != nil {
		return nil, err
	}
	entity := model.CompanyToDomain(m)
	return &entity, nil
}
