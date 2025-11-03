package maintenance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestServiceModelInitialization(t *testing.T) {
	s := Model{
		Name:       "Maintenance A",
		Price:      100,
		CategoryID: "Standard",
	}

	assert.NotNil(t, s)
	assert.Equal(t, "Maintenance A", s.Name)
	assert.Equal(t, int64(100), s.Price)
	assert.Equal(t, "Standard", s.CategoryID)
}

func TestServiceModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	modelService := &Model{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		Name:       "Service A",
		Price:      100,
		CategoryID: "Standard",
	}

	domainService := modelService.ToDomain()

	assert.NotNil(t, domainService.CategoryID)
	assert.Equal(t, modelService.CategoryID, string(domainService.CategoryID))

	assert.Equal(t, modelService.ID, domainService.ID)
	assert.Equal(t, modelService.Name, domainService.Name)
	assert.Equal(t, modelService.Price, domainService.Price)
}
