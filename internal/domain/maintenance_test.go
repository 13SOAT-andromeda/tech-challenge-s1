package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCategoryName_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected MaintenanceCategory
	}{
		{"standard", "standard", MaintenanceCategories.STANDARD},
		{"utility", "utility", MaintenanceCategories.UTILITY},
		{"commercial", "commercial", MaintenanceCategories.COMMERCIAL},
		{"premium", "premium", MaintenanceCategories.PREMIUM},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCategoryName(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseCategoryName_Invalid_Panics(t *testing.T) {
	expr := fmt.Sprintf("Maintenance Category '%s' is not valid.", "unknown")
	assert.PanicsWithValue(t, expr, func() { ParseCategoryName("unknown") })
}

func TestValidateMaintenanceCategory(t *testing.T) {
	tests := []struct {
		name    string
		cat     MaintenanceCategory
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid - standard literal",
			cat:     MaintenanceCategory("standard"),
			wantErr: false,
		},
		{
			name:    "valid - utility literal",
			cat:     MaintenanceCategory("utility"),
			wantErr: false,
		},
		{
			name:    "invalid - constant (padrao)",
			cat:     MaintenanceCategories.STANDARD,
			wantErr: true,
			errMsg:  "Maintenance Category 'padrao' is not valid. Accepted types: [standard utility commercial premium]",
		},
		{
			name:    "invalid - empty",
			cat:     MaintenanceCategory(""),
			wantErr: true,
			errMsg:  "Maintenance Category '' is not valid. Accepted types: [standard utility commercial premium]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Maintenance{CategoryID: tt.cat}
			err := m.ValidateMaintenanceCategory()
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
