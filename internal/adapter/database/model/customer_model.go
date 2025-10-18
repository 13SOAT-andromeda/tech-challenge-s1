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

func (m *CustomerModel) ToDomain() *domain.Customer {
	if m == nil {
		return nil
	}
	return &domain.Customer{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Document:  m.Document,
		Type:      m.Type,
		Contact:   m.Contact,
		Address:   m.Address.ToDomain(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromDomainCustomer(d *domain.Customer) *CustomerModel {
	if d == nil {
		return nil
	}
	model := CustomerModel{
		Name:     d.Name,
		Email:    d.Email,
		Document: d.Document,
		Type:     d.Type,
		Contact:  d.Contact,
		Address:  FromDomainAddress(d.Address),
	}
	model.ID = d.ID

	return &model
}
