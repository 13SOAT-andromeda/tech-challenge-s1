package repository

import (
	"context"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Save(ctx context.Context, c domain.Company) (*domain.Company, error) {
	m := model.FromDomain(c)

	return nil, r.db.Create(&m).Error
}

func (r *CompanyRepository) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
	var m model.CompanyModel
	if err := r.db.Preload("Address").First(&m, id).Error; err != nil {
		return nil, err
	}

	company := model.ToDomain(m)
	return &company, nil
}

func (r *CompanyRepository) UpdateByID(ctx context.Context, id uint, c domain.Company) (*domain.Company, error) {
	var m model.CompanyModel
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}

	m.Name = c.Name
	m.Email = c.Email
	m.Document = c.Document
	m.Contact = c.Contact

	if m.Address == nil {
		m.Address = &model.AddressModel{}
	}
	m.Address.Address = c.Address.Address
	m.Address.AddressNumber = c.Address.AddressNumber
	m.Address.City = c.Address.City
	m.Address.Neighborhood = c.Address.Neighborhood
	m.Address.Country = c.Address.Country
	m.Address.ZipCode = c.Address.ZipCode

	if err := r.db.Save(&m).Error; err != nil {
		return nil, err
	}

	company := model.ToDomain(m)
	return &company, nil
}

func (r *CompanyRepository) DeleteByID(ctx context.Context, id uint) (*domain.Company, error) {
	var m model.CompanyModel
	if err := r.db.Preload("Address").First(&m, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Delete(&m).Error; err != nil {
		return nil, err
	}
	company := model.ToDomain(m)
	return &company, nil
}
