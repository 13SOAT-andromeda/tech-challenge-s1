package domain

import "time"

type OrderStatus string

const (
	RECEIVED          OrderStatus = "Recebida"
	IN_ANALYSIS       OrderStatus = "Em diagnóstico"
	AWAITING_APPROVAL OrderStatus = "Aguardando aprovação"
	IN_PROGRESS       OrderStatus = "Em execução"
	FINISHED          OrderStatus = "Finalizado"
	DELIVERED         OrderStatus = "Entregue"
)

var OrderStatuses = []OrderStatus{
	RECEIVED,
	IN_ANALYSIS,
	AWAITING_APPROVAL,
	IN_PROGRESS,
	FINISHED,
	DELIVERED,
}

type Order struct {
	ID                uint
	DateIn            time.Time
	DateOut           *time.Time
	Number            string
	Status            OrderStatus
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
