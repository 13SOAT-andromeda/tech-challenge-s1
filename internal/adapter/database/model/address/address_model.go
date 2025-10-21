package address

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type Model struct {
	Address       string
	AddressNumber string
	Neighborhood  string
	City          string
	Country       string
	ZipCode       string
}

func (m *Model) ToDomain() *domain.Address {
	return &domain.Address{
		Address:       m.Address,
		AddressNumber: m.AddressNumber,
		Neighborhood:  m.Neighborhood,
		City:          m.City,
		Country:       m.Country,
		ZipCode:       m.ZipCode,
	}
}

func (m *Model) FromDomain(d *domain.Address) {
	if d == nil {
		return
	}

	m.Address = d.Address
	m.AddressNumber = d.AddressNumber
	m.Neighborhood = d.Neighborhood
	m.City = d.City
	m.Country = d.Country
	m.ZipCode = d.ZipCode
}
