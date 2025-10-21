package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null; unique"`
	Contact  string `gorm:"not null"`
	Address  string `gorm:"not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`

	Sessions []SessionModel `gorm:"foreignKey:UserId;references:ID"`
}

func (*Model) TableName() string {
	return "User"
}

func (m *Model) ToDomain() *domain.User {
	if m == nil {
		return nil
	}

	sessions := make([]domain.Session, len(m.Sessions))
	for i, session := range m.Sessions {
		sessions[i] = *(&session).ToDomain()
	}

	return &domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Contact:   m.Contact,
		Address:   m.Address,
		Password:  m.Password,
		Role:      m.Role,
		Sessions:  sessions,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *Model) FromDomain(d *domain.User) {
	sessions := make([]SessionModel, len(d.Sessions))
	for i, session := range d.Sessions {
		sessions[i] = *FromDomainSession(&session)
	}

	m.ID = d.ID
	m.Name = d.Name
	m.Email = d.Email
	m.Contact = d.Contact
	m.Address = d.Address
	m.Password = d.Password
	m.Role = d.Role
	m.Sessions = sessions
}
