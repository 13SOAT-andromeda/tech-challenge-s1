package domain

type Maintenance struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	DefaultPrice *float64 `json:"default_price"`
	Number       string   `json:"number"`
}
