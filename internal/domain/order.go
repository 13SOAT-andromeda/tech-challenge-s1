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
	ID                uint           `json:"id"`
	DateIn            time.Time      `json:"date_in"`
	DateOut           *time.Time     `json:"date_out"`
	DateApproved      *time.Time     `json:"date_approved"`
	DateRejected      *time.Time     `json:"date_rejected"`
	Status            OrderStatus    `json:"status"`
	VehicleKilometers int            `json:"vehicle_kilometers"`
	Note              *string        `json:"note"`
	DiagnosticNote    *string        `json:"diagnostic_note"`
	Price             *float64       `json:"price"`
	CustomerVehicleID uint           `json:"customer_vehicle_id"`
	EmployeeID        uint           `json:"employee_id"`
	CompanyID         uint           `json:"company_id"`
	Vehicle           *Vehicle       `json:"vehicle,omitempty"`
	Products          *[]Product     `json:"products,omitempty"`
	Maintenances      *[]Maintenance `json:"maintenances,omitempty"`
}
