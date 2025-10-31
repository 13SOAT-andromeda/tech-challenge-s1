package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_FieldValues(t *testing.T) {
	p := Product{
		ID:    42,
		Name:  "Sample",
		Stock: 7,
		Price: int64(1999),
	}

	assert.Equal(t, uint(42), p.ID)
	assert.Equal(t, "Sample", p.Name)
	assert.Equal(t, uint(7), p.Stock)
	assert.Equal(t, int64(1999), p.Price)
}

func TestProduct_ZeroValues(t *testing.T) {
	var p Product
	assert.Equal(t, uint(0), p.ID)
	assert.Equal(t, "", p.Name)
	assert.Equal(t, uint(0), p.Stock)
	assert.Equal(t, int64(0), p.Price)
}
