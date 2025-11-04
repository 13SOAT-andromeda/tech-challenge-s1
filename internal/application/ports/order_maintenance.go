package ports

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_maintenance"
)

type OrderMaintenanceRepository interface {
	Repository[order_maintenance.Model]
}
