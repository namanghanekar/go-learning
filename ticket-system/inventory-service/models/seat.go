package models

import "time"

type Seat struct {
	ID            int        `gorm:"primaryKey"`
	Status        string     `gorm:"default:available;index"`
	UserID        string
	LockExpiresAt *time.Time `gorm:"index"`
}
