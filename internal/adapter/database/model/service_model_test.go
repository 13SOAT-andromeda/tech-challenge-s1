package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestServiceTableName(t *testing.T) {
	assert.Equal(t, "Service", ServiceModel{}.TableName())
}

func TestNilServiceToDomain(t *testing.T) {
	assert.Nil(t, (*ServiceModel)(nil).ToDomain())
}

func TestNilServiceFromDomain(t *testing.T) {
	assert.Nil(t, FromDomainService(nil))
}

func TestServiceModelInitialization(t *testing.T) {
	defaultPrice := 100.0
	s := ServiceModel{
		Name:         "Service A",
		DefaultPrice: &defaultPrice,
		CategoryId:   1,
		Number:       "S123",
	}

	assert.NotNil(t, s)
	assert.Equal(t, "Service A", s.Name)
	assert.Equal(t, &defaultPrice, s.DefaultPrice)
	assert.Equal(t, uint(1), s.CategoryId)
	assert.Equal(t, "S123", s.Number)
}

func TestServiceModel_ToFromDomain(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour * 1)
	defaultPrice := 100.0

	modelService := &ServiceModel{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: deletedAt, Valid: true},
		},
		Name:         "Service A",
		DefaultPrice: &defaultPrice,
		CategoryId:   1,
		Number:       "S123",
	}

	domainService := modelService.ToDomain()

	assert.Equal(t, modelService.ID, domainService.ID)
	assert.Equal(t, modelService.Name, domainService.Name)
	assert.Equal(t, modelService.DefaultPrice, domainService.DefaultPrice)
	assert.Equal(t, modelService.CategoryId, domainService.CategoryId)
	assert.Equal(t, modelService.Number, domainService.Number)
	assert.Equal(t, modelService.CreatedAt, domainService.CreatedAt)
	assert.Equal(t, modelService.UpdatedAt, domainService.UpdatedAt)

	newModel := FromDomainService(domainService)
	newModel.DeletedAt = modelService.DeletedAt

	assert.Equal(t, modelService, newModel)
}
