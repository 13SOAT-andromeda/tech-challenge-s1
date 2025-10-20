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
	Document string        `gorm:"unique"`
	Contact  string        `gorm:"not null"`
	Address  address.Model `gorm:"embedded"`
}

func (Model) TableName() string {
	return "Company"
}

func ToDomain(c Model) domain.Company {
	return domain.Company{
		ID:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		Document: c.Document,
		Contact:  c.Contact,
		Address:  address.ToDomain(c.Address),
	}
}

func FromDomain(d domain.Company) Model {
	return Model{
		Name:     d.Name,
		Email:    d.Email,
		Document: d.Document,
		Contact:  d.Contact,
		Address:  address.FromDomain(d.Address),
	}
}
