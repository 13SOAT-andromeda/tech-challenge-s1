package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSessionModelInitialization(t *testing.T) {
	refreshToken := "some-refresh-token"
	s := SessionModel{
		UserID:       1,
		RefreshToken: &refreshToken,
	}

	assert.NotNil(t, s)
	assert.Equal(t, uint(1), s.UserID)
	assert.Equal(t, &refreshToken, s.RefreshToken)
}

func TestSessionModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)
	refreshToken := "some-refresh-token"

	modelSession := &SessionModel{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		UserID:       1,
		RefreshToken: &refreshToken,
	}

	domainSession := modelSession.ToDomain()

	assert.Equal(t, modelSession.ID, domainSession.ID)
	assert.Equal(t, modelSession.UserID, domainSession.UserID)
	assert.Equal(t, modelSession.RefreshToken, domainSession.RefreshToken)
	assert.Equal(t, modelSession.CreatedAt, domainSession.CreatedAt)
	assert.Equal(t, modelSession.UpdatedAt, domainSession.UpdatedAt)

}
