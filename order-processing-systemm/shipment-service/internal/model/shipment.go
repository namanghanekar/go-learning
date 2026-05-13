package model

import "time"

type Shipment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    int       `json:"order_id"`
	TrackingID string    `json:"tracking_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
