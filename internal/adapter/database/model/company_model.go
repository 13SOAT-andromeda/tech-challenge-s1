package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type CompanyModel struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string
	Document string       `gorm:"unique"`
	Contact  string       `gorm:"not null"`
	Address  AddressModel `gorm:"embedded"`
}

func (CompanyModel) TableName() string {
	return "Company"
}

func (m *CompanyModel) ToDomain() *domain.Company {
	if m == nil {
		return nil
	}
	return &domain.Company{
		ID:       m.ID,
		Name:     m.Name,
		Email:    m.Email,
		Document: m.Document,
		Contact:  m.Contact,
		Address:  *m.Address.ToDomain(),
	}
}

func FromDomainCompany(d *domain.Company) *CompanyModel {
	if d == nil {
		return nil
	}
	return &CompanyModel{
		Model: gorm.Model{
			ID: d.ID,
		},
		Name:     d.Name,
		Email:    d.Email,
		Document: d.Document,
		Contact:  d.Contact,
		Address:  FromDomainAddress(&d.Address),
	}
}
