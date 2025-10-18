package domain

import "time"

type OrderHistory struct {
	OrderId uint64
	Date    time.Time
	Status  string
}
