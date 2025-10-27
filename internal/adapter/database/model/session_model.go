package model

import (
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type SessionModel struct {
	gorm.Model
	UserID       uint `gorm:"not null"`
	RefreshToken *string
	ExpiresAt    time.Time `gorm:"not null"`
	IsActive     bool      `gorm:"default:true"`
}

func (SessionModel) TableName() string {
	return "Session"
}

func (m *SessionModel) ToDomain() *domain.Session {
	if m == nil {
		return nil
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &domain.Session{
		ID:           m.ID,
		UserID:       m.UserID,
		RefreshToken: m.RefreshToken,
		ExpiresAt:    m.ExpiresAt,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    deletedAt,
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
	m.ExpiresAt = d.ExpiresAt
	m.IsActive = d.IsActive
	if d.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
	}
}
