package domain

type CustomerVehicle struct {
	ID         uint
	CustomerId uint
	VehicleId  uint
	Vehicle    *Vehicle `json:"vehicle,omitempty"`
}
