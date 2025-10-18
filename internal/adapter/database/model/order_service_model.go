package model

type OrderServiceModel struct {
	ServiceId uint
	OrderId   uint
	Price     float64 `gorm:"not null"`

	Service ServiceModel `gorm:"foreignkey:ServiceId;references:ID"`
	Order   OrderModel   `gorm:"foreignkey:OrderId;references:ID"`
}

func (OrderServiceModel) TableName() string {
	return "Order_Service"
}
