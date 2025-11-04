package domain

type OrderMaintenance struct {
	MaintenanceId uint
	OrderId       uint
	Maintenance   Maintenance
	Order         Order
}
