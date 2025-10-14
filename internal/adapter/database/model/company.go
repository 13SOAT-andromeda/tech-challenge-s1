package model

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type CompanyModel struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Email    string
	Document string         `gorm:"unique"`
	Contact  string         `gorm:"not null"`
	Address  domain.Address `gorm:"embedded"`
}

func CompanyToDomain(model CompanyModel) domain.Company {
	return domain.Company{
		ID:       model.ID,
		Name:     model.Name,
		Email:    model.Email,
		Document: model.Document,
		Contact:  model.Contact,
		Address:  &model.Address,
	}
}

func CompanyFromDomain(domain domain.Company) CompanyModel {
	return CompanyModel{
		ID:       domain.ID,
		Name:     domain.Name,
		Email:    domain.Email,
		Document: domain.Document,
		Contact:  domain.Contact,
		Address:  *domain.Address,
	}
}
