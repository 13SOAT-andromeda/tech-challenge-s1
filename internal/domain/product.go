package domain

import "errors"

var (
	ErrInvalidPrice      = errors.New("preço inválido, deve ser maior ou igual a zero")
	ErrInvalidName       = errors.New("nome do produto não pode ser vazio")
	ErrInsufficientStock = errors.New("quantidade insuficiente em estoque")
)

type Product struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Stock uint   `json:"stock"`
	Price int64  `json:"price"`
}

func (p *Product) Validate() error {
	if p.Price < 0 {
		return ErrInvalidPrice
	}

	if p.Name == "" {
		return ErrInvalidName
	}

	return nil
}

func (p *Product) CanBePurchased(quantity uint) error {
	if p.Stock < quantity {
		return ErrInsufficientStock
	}

	return nil
}
func (p *Product) DecreaseStock(quantity uint) error {
	if err := p.CanBePurchased(quantity); err != nil {
		return err
	}

	p.Stock -= quantity

	return nil
}
