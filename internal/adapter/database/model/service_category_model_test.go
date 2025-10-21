package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestServiceCategoryModelInitialization(t *testing.T) {
	sc := ServiceCategoryModel{
		Name: "Category A",
	}

	assert.NotNil(t, sc)
	assert.Equal(t, "Category A", sc.Name)
}

func TestServiceCategoryModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)

	modelServiceCategory := &ServiceCategoryModel{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		Name: "Category A",
	}

	domainServiceCategory := modelServiceCategory.ToDomain()

	assert.Equal(t, modelServiceCategory.ID, domainServiceCategory.ID)
	assert.Equal(t, modelServiceCategory.Name, domainServiceCategory.Name)
	assert.Equal(t, modelServiceCategory.CreatedAt, domainServiceCategory.CreatedAt)
	assert.Equal(t, modelServiceCategory.UpdatedAt, domainServiceCategory.UpdatedAt)

}
