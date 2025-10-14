package domain

type Company struct {
	ID       uint
	Name     string
	Email    string
	Address  *Address
	Document string
	Contact  string
}
