package model

type OrderProductModel struct {
	ProductId uint
	OrderId   uint

	Product ProductModel `gorm:"foreignkey:ProductId;references:Id"`
	Order   OrderModel   `gorm:"foreignkey:OrderId;references:Id"`
}

func (OrderProductModel) TableName() string {
	return "Order_Product"
}
