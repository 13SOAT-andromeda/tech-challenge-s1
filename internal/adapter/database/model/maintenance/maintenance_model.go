package maintenance

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name         string `gorm:"not null"`
	DefaultPrice *float64
	CategoryId   uint
	Number       string `gorm:"not null"`
}

func (*Model) TableName() string {
	return "Maintenance"
}

func (m *Model) ToDomain() *domain.Maintenance {

	return &domain.Maintenance{
		ID:           m.ID,
		Name:         m.Name,
		DefaultPrice: m.DefaultPrice,
		Number:       m.Number,
	}
}

func (m *Model) FromDomain(d *domain.Maintenance) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Name = d.Name
	m.DefaultPrice = d.DefaultPrice
	m.Number = d.Number
}
