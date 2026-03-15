package document

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type Model struct {
	Document string `gorm:"column:document;not null;default"`
}

func (m *Model) ToDomain() *domain.Document {
	return &domain.Document{
		Number: m.Document,
	}
}

func (m *Model) FromDomain(d *domain.Document) {
	if d == nil {
		return
	}

	m.Document = d.GetDocumentNumber()
}
