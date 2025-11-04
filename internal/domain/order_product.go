package domain

type OrderProduct struct {
	Quantity  uint
	ProductId uint
	OrderId   uint

	Product Product
	Order   Order
}
