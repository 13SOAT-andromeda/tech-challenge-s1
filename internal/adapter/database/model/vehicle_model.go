package model

type VehicleModel struct {
	ID    uint   `gorm:"primaryKey"`
	Plate string `gorm:"unique; not null"`
	Model string `gorm:"not null"`
	Year  int    `gorm:"not null"`
	Brand string `gorm:"not null"`
	Color string `gorm:"not null"`
}

func (VehicleModel) TableName() string {
	return "Vehicle"
}
