package domain

import "time"

type Order struct {
	ID                uint
	DateIn            time.Time
	DateOut           *time.Time
	Number            string
	Status            string
	VehicleKilometers *float64
	Note              *string
	DiagnosticNote    *string
	Price             float64
	UserId            uint
	CustomerVehicleId uint
	CompanyId         uint

	User            User
	CustomerVehicle CustomerVehicle
	Company         Company
}
