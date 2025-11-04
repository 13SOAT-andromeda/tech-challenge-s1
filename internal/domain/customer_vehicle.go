package domain

type CustomerVehicle struct {
	ID         uint     `json:"id"`
	CustomerId uint     `json:"customer_id"`
	VehicleId  uint     `json:"vehicle_id"`
	Vehicle    Vehicle  `json:"vehicle,omitempty"`
	Customer   Customer `json:"customer,omitempty"`
}
