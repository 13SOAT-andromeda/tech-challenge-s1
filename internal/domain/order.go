package domain

import "time"

type OrderStatus string

const (
	RECEIVED          OrderStatus = "Recebida"
	IN_ANALYSIS       OrderStatus = "Em diagnóstico"
	ANALYSIS_FINISHED OrderStatus = "Diagnóstico finalizado"
	AWAITING_APPROVAL OrderStatus = "Aguardando aprovação"
	APPROVED          OrderStatus = "Aprovado"
	IN_PROGRESS       OrderStatus = "Em execução"
	FINISHED          OrderStatus = "Finalizado"
	DELIVERED         OrderStatus = "Entregue"
)

var OrderStatuses = []OrderStatus{
	RECEIVED,
	IN_ANALYSIS,
	ANALYSIS_FINISHED,
	AWAITING_APPROVAL,
	APPROVED,
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
	VehicleKilometers int
	Note              *string
	DiagnosticNote    *string
	Price             *float64
	User              User
	CustomerVehicle   CustomerVehicle
	Company           Company
}
