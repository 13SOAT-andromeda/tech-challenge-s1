package domain

type Company struct {
	ID       uint
	Name     string
	Email    string
	Document string
	Contact  string
	Address  Address
}
