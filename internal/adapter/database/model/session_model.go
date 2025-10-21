package model

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type SessionModel struct {
	gorm.Model
	UserID       uint `gorm:"not null"`
	RefreshToken *string
}

func (SessionModel) TableName() string {
	return "Session"
}

func (m *SessionModel) ToDomain() *domain.Session {
	if m == nil {
		return nil
	}
	return &domain.Session{
		ID:           m.ID,
		UserID:       m.UserID,
		RefreshToken: m.RefreshToken,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (m *SessionModel) FromDomain(d *domain.Session) {
	if d == nil {
		return
	}
	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	m.UserID = d.UserID
	m.RefreshToken = d.RefreshToken
}
