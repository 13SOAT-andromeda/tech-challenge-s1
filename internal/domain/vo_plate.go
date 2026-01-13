package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type Plate struct {
	Value string
}

func NewPlate(plate string) (*Plate, error) {

	if err := validatePlate(plate); err != nil {
		return nil, err
	}

	return &Plate{Value: normalizePlate(plate)}, nil
}

func normalizePlate(plate string) string {
	return strings.TrimSpace(strings.ToUpper(strings.ReplaceAll(plate, "-", "")))
}

func validatePlate(plate string) error {
	trimmed := strings.TrimSpace(plate)
	if trimmed == "" {
		return fmt.Errorf("placa vazia")
	}

	normalized := normalizePlate(trimmed)

	// padrão antigo: ABC1234
	oldPattern := regexp.MustCompile(`^[A-Z]{3}\d{4}$`)
	// padrão mercosul: ABC1D23
	mercosurPattern := regexp.MustCompile(`^[A-Z]{3}\d[A-Z]\d{2}$`)

	if oldPattern.MatchString(normalized) || mercosurPattern.MatchString(normalized) {
		return nil
	}

	return fmt.Errorf("placa inválida")
}

func (p *Plate) GetPlate() string {
	return p.Value
}
