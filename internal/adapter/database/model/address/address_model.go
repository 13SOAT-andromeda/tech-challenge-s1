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

func ToDomain(a Model) domain.Address {
	return domain.Address{
		Address:       a.Address,
		AddressNumber: a.AddressNumber,
		Neighborhood:  a.Neighborhood,
		City:          a.City,
		Country:       a.Country,
		ZipCode:       a.ZipCode,
	}
}

func FromDomain(a domain.Address) Model {
	return Model{
		Address:       a.Address,
		AddressNumber: a.AddressNumber,
		Neighborhood:  a.Neighborhood,
		City:          a.City,
		Country:       a.Country,
		ZipCode:       a.ZipCode,
	}
}
