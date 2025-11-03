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

type ProductItem struct {
	ID       uint
	Quantity uint
}

type MaintenanceItem struct {
	ID uint
}

type Order struct {
	ID                uint               `json:"id"`
	DateIn            time.Time          `json:"date_in"`
	DateOut           *time.Time         `json:"date_out"`
	Status            OrderStatus        `json:"status"`
	VehicleKilometers int                `json:"vehicle_kilometers"`
	Note              *string            `json:"note"`
	DiagnosticNote    *string            `json:"diagnostic_note"`
	Price             *float64           `json:"price"`
	User              User               `json:"user"`
	CustomerVehicle   CustomerVehicle    `json:"customer_vehicle"`
	Company           Company            `json:"company"`
	Products          *[]ProductItem     `json:"products:omitempty"`
	Maintenances      *[]MaintenanceItem `json:"maintenances:omitempty"`
}
