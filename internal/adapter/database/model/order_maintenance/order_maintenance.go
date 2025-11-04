package order_maintenance

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
	MaintenanceId uint
	OrderId       uint
	Maintenance   maintenance.Model `gorm:"foreignKey:MaintenanceId;references:ID"`
	Order         order.Model       `gorm:"foreignKey:OrderId;references:ID"`
}

func (*Model) TableName() string {
	return "Order_Maintenance"
}

func (m *Model) ToDomain() *domain.OrderMaintenance {
	return &domain.OrderMaintenance{
		MaintenanceId: m.MaintenanceId,
		OrderId:       m.OrderId,
		Maintenance:   *m.Maintenance.ToDomain(),
		Order:         *m.Order.ToDomain(),
	}
}

func (m *Model) FromDomain(d *domain.OrderMaintenance) {
	if d == nil {
		return
	}
	m.MaintenanceId = d.MaintenanceId
	m.OrderId = d.OrderId
	m.Maintenance.FromDomain(&d.Maintenance)
	m.Order.FromDomain(&d.Order)
}
