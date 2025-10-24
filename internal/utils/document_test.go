package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCpf(t *testing.T) {
	tests := []struct {
		name string
		cpf  string
		want bool
	}{
		{"Valid document - CPF", "529.982.247-25", true},
		{"Valid document - CPF - without diacritcs", "52998224725", true},
		{"Invalid Document - CPF - equal numbers", "111.111.111-11", false},
		{"Invalid Document - CPF - invalid range", "123.456.78", false},
		{"Invalid Document - CPF - invalid characters", "529.982.247-2X", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCpf(tt.cpf)

			assert.Equal(t, tt.want, result)
		})
	}
}
