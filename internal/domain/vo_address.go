package domain

type Address struct {
	Address       string
	AddressNumber string
	Neighborhood  string
	City          string
	Country       string
	ZipCode       string
}

func (c *Customer) EnsureAddress() {
	if c.Address == nil {
		c.Address = &Address{}
	}
}
