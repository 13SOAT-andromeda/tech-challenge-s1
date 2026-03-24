package person

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/document"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
	Contact  string
	Document document.Model `gorm:"embedded;unique"`
	Address  *address.Model `gorm:"embedded"`
}

func (*Model) TableName() string {
	return "Person"
}

func (m *Model) ToDomain() *domain.Person {
	documentDomain := domain.RestoreDocument(m.Document.Document)

	var addressDomain *domain.Address
	if m.Address != nil {
		addressDomain = m.Address.ToDomain()
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &domain.Person{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Contact:   m.Contact,
		Document:  &documentDomain,
		Address:   addressDomain,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *Model) FromDomain(d *domain.Person) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Name = d.Name
	m.Email = d.Email
	m.Contact = d.Contact

	if d.Document != nil {
		m.Document.FromDomain(d.Document)
	}

	if m.Address == nil {
		m.Address = &address.Model{}
	}
	m.Address.FromDomain(d.Address)
}
