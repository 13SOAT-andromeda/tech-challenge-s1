package domain

type OrderService struct {
	ServiceId uint
	OrderId   uint
	Price     float64

	Service Service
	Order   Order
}
