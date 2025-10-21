package domain

type Company struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Address  *Address `json:"address"`
	Document string   `json:"document"`
	Contact  string   `json:"contact"`
}
