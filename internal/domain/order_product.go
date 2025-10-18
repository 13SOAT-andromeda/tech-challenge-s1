package domain

type OrderProduct struct {
	ProductId uint
	OrderId   uint

	Product Product
	Order   Order
}
