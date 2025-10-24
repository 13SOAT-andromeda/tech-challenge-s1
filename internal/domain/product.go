package domain

type Product struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Quantity uint   `json:"quantity"`
	Price    int64  `json:"price"`
}

type ProductName struct {
	Name string `json:"name"`
}
