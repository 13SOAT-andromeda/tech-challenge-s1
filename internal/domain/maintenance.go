package domain

import (
	"fmt"
	"slices"
	"strings"
)

type MaintenanceCategory string

const (
	STANDARD   MaintenanceCategory = "padrao"
	UTILITY    MaintenanceCategory = "utilitario"
	COMMERCIAL MaintenanceCategory = "comercial"
	PREMIUM    MaintenanceCategory = "premium"
)

var MaintenanceCategories = struct {
	STANDARD   MaintenanceCategory
	UTILITY    MaintenanceCategory
	COMMERCIAL MaintenanceCategory
	PREMIUM    MaintenanceCategory
}{
	STANDARD:   STANDARD,
	UTILITY:    UTILITY,
	COMMERCIAL: COMMERCIAL,
	PREMIUM:    PREMIUM,
}

type Maintenance struct {
	ID         uint                `json:"id"`
	Name       string              `json:"name"`
	Price      int64               `json:"price"`
	CategoryId MaintenanceCategory `json:"category"`
}

func ParseCategoryName(CategoryId string) MaintenanceCategory {
	switch strings.TrimSpace(CategoryId) {
	case "standard":
		return MaintenanceCategories.STANDARD
	case "utility":
		return MaintenanceCategories.UTILITY
	case "commercial":
		return MaintenanceCategories.COMMERCIAL
	case "premium":
		return MaintenanceCategories.PREMIUM
	default:
		panic(fmt.Sprintf("Maintenance Category '%s' is not valid.", CategoryId))
	}
}

func (m *Maintenance) ValidateMaintenanceCategory() error {
	acceptedCategoryId := []MaintenanceCategory{"standard", "utility", "commercial", "premium"}

	if !slices.Contains(acceptedCategoryId, m.CategoryId) {
		var names []string
		for _, c := range acceptedCategoryId {
			names = append(names, string(c))
		}
		return fmt.Errorf("Maintenance Category '%s' is not valid. Accepted types: %v", string(m.CategoryId), names)
	}

	return nil
}
