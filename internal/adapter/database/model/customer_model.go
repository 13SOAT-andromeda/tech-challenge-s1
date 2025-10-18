package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type CustomerModel struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string
	Document string       `gorm:"unique"`
	Type     string       `gorm:"not null"`
	Contact  string       `gorm:"not null"`
	Address  AddressModel `gorm:"embedded"`
}

func (CustomerModel) TableName() string {
	return "Customer"
}

func ToDomain(model CustomerModel) domain.Customer {
	return domain.Customer{
		ID:       model.ID,
		Name:     model.Name,
		Email:    model.Email,
		Document: model.Document,
		Type:     model.Type,
		Contact:  model.Contact,
		Address:  model.Address.ToDomain(),
	}
}

func FromDomain(domain domain.Customer) CustomerModel {
	model := CustomerModel{

		Name:     domain.Name,
		Email:    domain.Email,
		Document: domain.Document,
		Type:     domain.Type,
		Contact:  domain.Contact,
		Address:  FromDomainAddress(domain.Address),
	}
	model.ID = domain.ID

	return model
}
