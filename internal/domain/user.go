package domain

type User struct {
	ID       uint
	Name     string
	Email    string
	Contact  string
	Address  *Address
	Password *Password
	Role     string
}
