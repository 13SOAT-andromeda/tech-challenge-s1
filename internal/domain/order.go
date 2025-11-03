package domain

import (
	"time"

	"errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

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

var OrderStatuses = struct {
	RECEIVED          OrderStatus
	IN_ANALYSIS       OrderStatus
	ANALYSIS_FINISHED OrderStatus
	AWAITING_APPROVAL OrderStatus
	APPROVED          OrderStatus
	IN_PROGRESS       OrderStatus
	FINISHED          OrderStatus
	DELIVERED         OrderStatus
}{
	RECEIVED:          RECEIVED,
	IN_ANALYSIS:       IN_ANALYSIS,
	ANALYSIS_FINISHED: ANALYSIS_FINISHED,
	AWAITING_APPROVAL: AWAITING_APPROVAL,
	APPROVED:          APPROVED,
	IN_PROGRESS:       IN_PROGRESS,
	FINISHED:          FINISHED,
	DELIVERED:         DELIVERED,
}

type Order struct {
	ID                uint
	DateIn            time.Time
	DateOut           *time.Time
	Status            OrderStatus
	VehicleKilometers int
	Note              *string
	DiagnosticNote    *string
	Price             *float64
	User              User
	CustomerVehicle   CustomerVehicle
	Company           Company
}
