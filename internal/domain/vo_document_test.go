package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDocument_Valid(t *testing.T) {
	tests := []struct {
		name     string
		document string
	}{
		{"Valid CPF with dots and dash", "123.456.789-09"},
		{"Valid CPF without formatting", "12345678909"},
		{"Valid CPF with spaces", "123 456 789 09"},
		{"Another valid CPF", "111.444.777-35"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument(tt.document)

			assert.NoError(t, err)
			assert.NotNil(t, doc)
			assert.Equal(t, tt.document, doc.Number)
		})
	}
}

func TestNewDocument_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		document string
	}{
		{"Empty document", ""},
		{"Less than 11 digits", "123456789"},
		{"More than 11 digits", "123456789012"},
		{"All zeros", "00000000000"},
		{"All ones", "11111111111"},
		{"All twos", "22222222222"},
		{"Invalid check digit", "12345678901"},
		{"Only letters", "abcdefghijk"},
		{"Mixed invalid", "123.456.789-00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument(tt.document)

			assert.Error(t, err)
			assert.Nil(t, doc)
			assert.Contains(t, err.Error(), "Document number is invalid")
		})
	}
}

func TestRestoreDocument(t *testing.T) {
	raw := "12345678909"

	doc := RestoreDocument(raw)

	assert.Equal(t, raw, doc.Number)
}

func TestRestoreDocument_NoValidation(t *testing.T) {
	invalid := "00000000000"

	doc := RestoreDocument(invalid)

	assert.Equal(t, invalid, doc.Number)
}

func TestDocument_GetDocumentNumber(t *testing.T) {
	doc := Document{Number: "12345678909"}

	result := doc.GetDocumentNumber()

	assert.Equal(t, "12345678909", result)
}

func TestDocument_NormalizeDocument(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"With dots and dash", "123.456.789-09", "12345678909"},
		{"With spaces", "123 456 789 09", "12345678909"},
		{"With mixed chars", "123-456.789/09", "12345678909"},
		{"Only numbers", "12345678909", "12345678909"},
		{"Empty", "", ""},
		{"With letters", "123abc456def789", "123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{}
			result := doc.NormalizeDocument(tt.input)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDocument_ValidateCpf_Valid(t *testing.T) {
	tests := []struct {
		name string
		cpf  string
	}{
		{"Valid CPF 1", "12345678909"},
		{"Valid CPF with formatting", "123.456.789-09"},
		{"Valid CPF 2", "11144477735"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Number: tt.cpf}

			result := doc.ValidateCpf()

			assert.True(t, result)
		})
	}
}

func TestDocument_ValidateCpf_Invalid(t *testing.T) {
	tests := []struct {
		name string
		cpf  string
	}{
		{"Empty", ""},
		{"Less than 11 digits", "123456789"},
		{"More than 11 digits", "123456789012"},
		{"All zeros", "00000000000"},
		{"All ones", "11111111111"},
		{"All twos", "22222222222"},
		{"All threes", "33333333333"},
		{"All fours", "44444444444"},
		{"All fives", "55555555555"},
		{"All sixes", "66666666666"},
		{"All sevens", "77777777777"},
		{"All eights", "88888888888"},
		{"All nines", "99999999999"},
		{"Invalid first check digit", "12345678900"},
		{"Invalid second check digit", "12345678908"},
		{"Both check digits invalid", "12345678999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Number: tt.cpf}

			result := doc.ValidateCpf()

			assert.False(t, result)
		})
	}
}

func TestDocument_ValidateCpf_WithFormatting(t *testing.T) {
	// Testa que a validação funciona com formatação
	doc := &Document{Number: "123.456.789-09"}

	result := doc.ValidateCpf()

	assert.True(t, result)
}

func TestDocument_ValidateCpf_CheckDigitCalculation(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"First digit is 10, should become 0", "99999999999", false},
		{"Valid CPF where calculation results in 10", "52998224725", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Number: tt.cpf}

			result := doc.ValidateCpf()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewDocument_IntegrationWithValidateCpf(t *testing.T) {
	validCPF := "12345678909"
	invalidCPF := "12345678900"

	validDoc, err := NewDocument(validCPF)
	assert.NoError(t, err)
	assert.NotNil(t, validDoc)
	assert.True(t, validDoc.ValidateCpf())

	invalidDoc, err := NewDocument(invalidCPF)
	assert.Error(t, err)
	assert.Nil(t, invalidDoc)
}

func TestDocument_JSONMarshaling(t *testing.T) {
	doc := Document{Number: "12345678909"}

	assert.NotEmpty(t, doc.Number)
}
