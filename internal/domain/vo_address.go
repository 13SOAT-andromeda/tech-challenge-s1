package domain

type Address struct {
	Address       string `json:"address"`
	AddressNumber string `json:"address_number"`
	Neighborhood  string `json:"neighborhood"`
	City          string `json:"city"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
}
