package company

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
	Contact  string         `gorm:"not null"`
	Address  *address.Model `gorm:"embedded"`
}

func (*Model) TableName() string {
	return "Company"
}

func (m *Model) ToDomain() *domain.Company {
	var addressDomain *domain.Address
	if m.Address != nil {
		addressDomain = m.Address.ToDomain()
	} else {
		addressDomain = nil
	}

	return &domain.Company{
		ID:       m.ID,
		Name:     m.Name,
		Email:    m.Email,
		Document: m.Document,
		Contact:  m.Contact,
		Address:  addressDomain,
	}
}

func (m *Model) FromDomain(d *domain.Company) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Name = d.Name
	m.Email = d.Email
	m.Document = d.Document
	m.Contact = d.Contact

	if m.Address == nil {
		m.Address = &address.Model{}
	}

	m.Address.FromDomain(d.Address)
}
