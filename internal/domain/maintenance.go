package domain

type Maintenance struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}
