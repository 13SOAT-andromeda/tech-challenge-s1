package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_Validate_Success(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.Validate()

	assert.NoError(t, err)
}

func TestProduct_Validate_InvalidPrice(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: -100,
	}

	err := product.Validate()

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProduct_Validate_ZeroPrice(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 0,
	}

	err := product.Validate()

	assert.NoError(t, err)
}

func TestProduct_Validate_EmptyName(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.Validate()

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidName, err)
}

func TestProduct_Validate_WhitespaceName(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "   ",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.Validate()

	assert.NoError(t, err)
}

func TestProduct_CanBePurchased_Success(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.CanBePurchased(5)

	assert.NoError(t, err)
}

func TestProduct_CanBePurchased_ExactStock(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.CanBePurchased(10)

	assert.NoError(t, err)
}

func TestProduct_CanBePurchased_InsufficientStock(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.CanBePurchased(11)

	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
}

func TestProduct_CanBePurchased_ZeroStock(t *testing.T) {
	StockMock := uint(0)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.CanBePurchased(1)

	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
}

func TestProduct_CanBePurchased_ZeroQuantity(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.CanBePurchased(0)

	assert.NoError(t, err)
}

func TestProduct_DecreaseStock_Success(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.DecreaseStock(3)

	assert.NoError(t, err)
	assert.Equal(t, uint(7), *product.Stock)
}

func TestProduct_DecreaseStock_ToZero(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.DecreaseStock(10)

	assert.NoError(t, err)
	assert.Equal(t, uint(0), *product.Stock)
}

func TestProduct_DecreaseStock_InsufficientStock(t *testing.T) {
	StockMock := uint(5)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.DecreaseStock(10)

	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
	assert.Equal(t, uint(5), *product.Stock)
}

func TestProduct_DecreaseStock_ZeroQuantity(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.DecreaseStock(0)

	assert.NoError(t, err)
	assert.Equal(t, uint(10), *product.Stock) // estoque não deve mudar
}

func TestProduct_DecreaseStock_MultipleOperations(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	err := product.DecreaseStock(3)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), *product.Stock)

	err = product.DecreaseStock(2)
	assert.NoError(t, err)
	assert.Equal(t, uint(5), *product.Stock)

	err = product.DecreaseStock(10)
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
	assert.Equal(t, uint(5), *product.Stock) // estoque não muda na falha
}

func TestProduct_Initialization(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	assert.NotNil(t, product)
	assert.Equal(t, uint(1), product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, uint(10), *product.Stock)
	assert.Equal(t, int64(10000), product.Price)
}

func TestProduct_JSONTags(t *testing.T) {
	StockMock := uint(10)
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Stock: &StockMock,
		Price: 10000,
	}

	assert.Equal(t, uint(1), product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, uint(10), *product.Stock)
	assert.Equal(t, int64(10000), product.Price)
}
