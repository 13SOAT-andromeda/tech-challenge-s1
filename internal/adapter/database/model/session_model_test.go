package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSessionModelInitialization(t *testing.T) {
	refreshToken := "some-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	s := SessionModel{
		UserID:       1,
		RefreshToken: &refreshToken,
		ExpiresAt:    expiresAt,
	}

	assert.NotNil(t, s)
	assert.Equal(t, uint(1), s.UserID)
	assert.Equal(t, &refreshToken, s.RefreshToken)
	assert.Equal(t, expiresAt, s.ExpiresAt)
}

func TestSessionModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)
	refreshToken := "some-refresh-token"
	expiresAt := now.Add(24 * time.Hour)

	modelSession := &SessionModel{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		UserID:       1,
		RefreshToken: &refreshToken,
		ExpiresAt:    expiresAt,
	}

	domainSession := modelSession.ToDomain()

	assert.Equal(t, modelSession.ID, domainSession.ID)
	assert.Equal(t, modelSession.UserID, domainSession.UserID)
	assert.Equal(t, modelSession.RefreshToken, domainSession.RefreshToken)
	assert.Equal(t, modelSession.ExpiresAt, domainSession.ExpiresAt)
	assert.Equal(t, modelSession.CreatedAt, domainSession.CreatedAt)
	assert.Equal(t, modelSession.UpdatedAt, domainSession.UpdatedAt)
	assert.Equal(t, deletedAt, *domainSession.DeletedAt)

	// Test FromDomain
	newModelSession := &SessionModel{}
	newModelSession.FromDomain(domainSession)

	assert.Equal(t, modelSession.ID, newModelSession.ID)
	assert.Equal(t, modelSession.UserID, newModelSession.UserID)
	assert.Equal(t, modelSession.RefreshToken, newModelSession.RefreshToken)
	assert.Equal(t, modelSession.ExpiresAt, newModelSession.ExpiresAt)
	assert.Equal(t, modelSession.CreatedAt, newModelSession.CreatedAt)
	assert.Equal(t, modelSession.UpdatedAt, newModelSession.UpdatedAt)
	assert.Equal(t, modelSession.DeletedAt.Time, newModelSession.DeletedAt.Time)
}

func TestSessionModel_ToDomain_NilModel(t *testing.T) {
	var modelSession *SessionModel
	domainSession := modelSession.ToDomain()
	assert.Nil(t, domainSession)
}

func TestSessionModel_FromDomain_NilDomain(t *testing.T) {
	modelSession := &SessionModel{
		UserID: 1,
	}

	modelSession.FromDomain(nil)

	// Should not change the model
	assert.Equal(t, uint(1), modelSession.UserID)
}

func TestSessionModel_EdgeCases(t *testing.T) {
	t.Run("nil refresh token", func(t *testing.T) {
		modelSession := &SessionModel{
			UserID:       1,
			RefreshToken: nil,
			ExpiresAt:    time.Now().Add(time.Hour),
		}

		domainSession := modelSession.ToDomain()
		assert.Nil(t, domainSession.RefreshToken)
		assert.Equal(t, uint(1), domainSession.UserID)
	})

	t.Run("empty refresh token", func(t *testing.T) {
		emptyToken := ""
		modelSession := &SessionModel{
			UserID:       1,
			RefreshToken: &emptyToken,
			ExpiresAt:    time.Now().Add(time.Hour),
		}

		domainSession := modelSession.ToDomain()
		assert.Equal(t, "", *domainSession.RefreshToken)
	})

	t.Run("just map basic fields", func(t *testing.T) {
		modelSession := &SessionModel{
			UserID:    1,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		domainSession := modelSession.ToDomain()
		assert.Equal(t, uint(1), domainSession.UserID)
	})

	t.Run("expired session", func(t *testing.T) {
		expiredTime := time.Now().Add(-time.Hour)
		modelSession := &SessionModel{
			UserID:    1,
			ExpiresAt: expiredTime,
		}

		domainSession := modelSession.ToDomain()
		assert.Equal(t, expiredTime, domainSession.ExpiresAt)
		assert.True(t, domainSession.IsExpired())
	})
}

func TestSessionModel_TableName(t *testing.T) {
	model := SessionModel{}
	assert.Equal(t, "Session", model.TableName())
}
