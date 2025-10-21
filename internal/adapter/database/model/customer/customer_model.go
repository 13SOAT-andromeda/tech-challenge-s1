package customer

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string
	Document string         `gorm:"unique"`
	Type     string         `gorm:"not null"`
	Contact  string         `gorm:"not null"`
	Address  *address.Model `gorm:"embedded"`
}

func (*Model) TableName() string {
	return "Customer"
}

func (m *Model) ToDomain() *domain.Customer {
	var addressDomain *domain.Address
	if m.Address != nil {
		addressDomain = m.Address.ToDomain()
	} else {
		addressDomain = nil
	}

	return &domain.Customer{
		ID:       m.ID,
		Name:     m.Name,
		Email:    m.Email,
		Document: m.Document,
		Type:     m.Type,
		Contact:  m.Contact,
		Address:  addressDomain,
	}
}

func (m *Model) FromDomain(d *domain.Customer) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Name = d.Name
	m.Email = d.Email
	m.Document = d.Document
	m.Type = d.Type
	m.Contact = d.Contact

	if m.Address == nil {
		m.Address = &address.Model{}
	}

	m.Address.FromDomain(d.Address)
}
