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

func FromDomainSession(d *domain.Session) *SessionModel {
	if d == nil {
		return nil
	}
	return &SessionModel{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		UserID:       d.UserID,
		RefreshToken: d.RefreshToken,
	}
}
