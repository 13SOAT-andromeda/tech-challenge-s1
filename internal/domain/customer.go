package domain

type Customer struct {
	ID       uint
	Name     string
	Email    string
	Document string
	Type     string
	Contact  string
	Address  *Address
}
