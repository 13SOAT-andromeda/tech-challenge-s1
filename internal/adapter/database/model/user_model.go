package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null; unique"`
	Contact  string `gorm:"not null"`
	Address  string `gorm:"not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`

	Sessions []SessionModel `gorm:"foreignKey:UserId;references:ID"`
}

func (UserModel) TableName() string {
	return "User"
}

func (m *UserModel) ToDomain() *domain.User {
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

func FromDomainUser(d *domain.User) *UserModel {
	if d == nil {
		return nil
	}

	sessions := make([]SessionModel, len(d.Sessions))
	for i, session := range d.Sessions {
		sessions[i] = *FromDomainSession(&session)
	}

	return &UserModel{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name:     d.Name,
		Email:    d.Email,
		Contact:  d.Contact,
		Address:  d.Address,
		Password: d.Password,
		Role:     d.Role,
		Sessions: sessions,
	}
}
