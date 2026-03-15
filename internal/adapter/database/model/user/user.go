package user

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/address"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/document"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name     string         `gorm:"not null"`
	Email    string         `gorm:"not null"`
	Document document.Model `gorm:"embedded;unique"`
	Contact  string         `gorm:"not null"`
	Address  *address.Model `gorm:"embedded"`
	Password string         `gorm:"not null"`
	Role     string         `gorm:"not null"`
}

func (*Model) TableName() string {
	return "User"
}

func (m *Model) ToDomain() *domain.User {
	pass := domain.NewPasswordFromHash(m.Password, encryption.NewBcryptHasher())
	documentDomain := m.Document.ToDomain()

	var addressDomain *domain.Address

	if m.Address != nil {
		addressDomain = m.Address.ToDomain()
	} else {
		addressDomain = nil
	}

	var deletedAt *time.Time

	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Document:  documentDomain,
		Contact:   m.Contact,
		Role:      m.Role,
		Password:  pass,
		Address:   addressDomain,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *Model) FromDomain(d *domain.User) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Name = d.Name
	m.Email = d.Email
	m.Contact = d.Contact
	m.Role = d.Role
	m.Password = d.Password.GetHashed()
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt

	m.Document.FromDomain(d.Document)

	if d.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
	} else {
		m.DeletedAt = gorm.DeletedAt{Valid: false}
	}

	if m.Address == nil {
		m.Address = &address.Model{}
	}

	m.Address.FromDomain(d.Address)
}
