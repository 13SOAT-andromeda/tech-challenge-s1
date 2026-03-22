package customer

import (
	"time"

	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Type     string            `gorm:"not null"`
	PersonID uint              `gorm:"not null"`
	Person   personModel.Model `gorm:"foreignKey:PersonID;references:ID"`
}

func (*Model) TableName() string {
	return "Customer"
}

func (m *Model) ToDomain() *domain.Customer {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	var person *domain.Person
	if m.Person.ID != 0 {
		person = m.Person.ToDomain()
	}

	return &domain.Customer{
		ID:        m.ID,
		Type:      m.Type,
		PersonID:  m.PersonID,
		Person:    person,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *Model) FromDomain(d *domain.Customer) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Type = d.Type
	m.PersonID = d.PersonID

	if d.Person != nil {
		m.Person.FromDomain(d.Person)
	}
}
