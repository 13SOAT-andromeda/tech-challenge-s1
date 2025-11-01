package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidPlate(t *testing.T) {
	oldPlate := "EML8635"
	mercosurPlate := "LSN4I49"

	pOld, err := NewPlate(oldPlate)
	assert.NoError(t, err)
	assert.Nil(t, err)

	pMercosur, err := NewPlate(mercosurPlate)
	assert.NoError(t, err)
	assert.Nil(t, err)

	assert.Equal(t, oldPlate, pOld.GetPlate())
	assert.Equal(t, mercosurPlate, pMercosur.GetPlate())
}

func TestValidInvalidPlate(t *testing.T) {
	plate := "EM98635"

	p, err := NewPlate(plate)

	assert.ErrorIs(t, err, ErrPlateInvalid)
	assert.Nil(t, p)
}

func TestEmptyPlate(t *testing.T) {
	plate := ""

	p, err := NewPlate(plate)

	assert.ErrorIs(t, err, ErrPlateEmpty)
	assert.Nil(t, p)
}
