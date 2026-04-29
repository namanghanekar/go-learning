package domain

import "time"

const (
	SeatStatusAvailable = "AVAILABLE"
	SeatStatusHeld      = "HELD"
	SeatStatusSold      = "SOLD"
)

type Seat struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	RoomID    string     `gorm:"index;not null" json:"roomId"`
	Number    string     `gorm:"index;not null" json:"number"`
	Status    string     `gorm:"index;not null" json:"status"`
	HoldToken string     `gorm:"index" json:"holdToken"`
	HeldBy    string     `json:"heldBy"`
	ExpiresAt *time.Time `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type PaymentRecord struct {
	ID            string    `gorm:"primaryKey"`
	SeatID        string    `gorm:"index;not null"`
	HoldToken     string    `gorm:"index;not null"`
	UserID        string    `gorm:"index;not null"`
	Status        string    `gorm:"index;not null"`
	AmountCents   int64     `gorm:"not null"`
	RequestedAt   time.Time `gorm:"not null"`
	ProcessedAt   *time.Time
	FailureReason string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
