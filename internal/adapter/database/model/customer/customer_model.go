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
	Document string        `gorm:"unique"`
	Type     string        `gorm:"not null"`
	Contact  string        `gorm:"not null"`
	Address  address.Model `gorm:"embedded"`
}

func (Model) TableName() string {
	return "Customer"
}

func ToDomain(c Model) domain.Customer {
	return domain.Customer{
		ID:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		Document: c.Document,
		Type:     c.Type,
		Contact:  c.Contact,
		Address:  address.ToDomain(c.Address),
	}
}

func FromDomain(d domain.Customer) Model {
	return Model{
		Name:     d.Name,
		Email:    d.Email,
		Document: d.Document,
		Type:     d.Type,
		Contact:  d.Contact,
		Address:  address.FromDomain(d.Address),
	}
}
