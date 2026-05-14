package model

import "time"

type Notification struct {
	ID uint `gorm:"primaryKey"`

	UserID int

	Message string

	Type string

	CreatedAt time.Time
}
