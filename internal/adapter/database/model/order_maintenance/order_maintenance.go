package order_maintenance

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type Model struct {
	MaintenanceId uint `gorm:"primaryKey"`
	OrderId       uint `gorm:"primaryKey"`
	Maintenance   maintenance.Model `gorm:"foreignKey:MaintenanceId;references:ID"`
}

func (*Model) TableName() string {
	return "Order_Maintenance"
}

func (m *Model) ToDomain(order *domain.Order) *domain.OrderMaintenance {
	result := &domain.OrderMaintenance{
		MaintenanceId: m.MaintenanceId,
		OrderId:       m.OrderId,
		Maintenance:   *m.Maintenance.ToDomain(),
	}
	if order != nil {
		result.Order = *order
	}
	return result
}

func (m *Model) FromDomain(d *domain.OrderMaintenance) {
	if d == nil {
		return
	}
	m.MaintenanceId = d.MaintenanceId
	m.OrderId = d.OrderId
	m.Maintenance.FromDomain(&d.Maintenance)
}
