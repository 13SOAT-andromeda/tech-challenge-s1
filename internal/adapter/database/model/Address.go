package model

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type AddressModel struct {
	Address       string
	AddressNumber string
	Neighborhood  string
	City          string
	Country       string
	ZipCode       string
}

func (modelAddr *AddressModel) ToDomain() *domain.Address {
	return &domain.Address{
		Address:       modelAddr.Address,
		AddressNumber: modelAddr.AddressNumber,
		Neighborhood:  modelAddr.Neighborhood,
		City:          modelAddr.City,
		Country:       modelAddr.Country,
		ZipCode:       modelAddr.ZipCode,
	}
}

func FromDomainAddress(domainAddr *domain.Address) AddressModel {
	if domainAddr == nil {
		return AddressModel{}
	}
	return AddressModel{
		Address:       domainAddr.Address,
		AddressNumber: domainAddr.AddressNumber,
		Neighborhood:  domainAddr.Neighborhood,
		City:          domainAddr.City,
		Country:       domainAddr.Country,
		ZipCode:       domainAddr.ZipCode,
	}
}
