package contracts

import "time"

type SeatDTO struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"roomId"`
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	HeldBy    string    `json:"heldBy"`
	HoldToken string    `json:"holdToken"`
	ExpiresAt time.Time `json:"expiresAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SeatEvent struct {
	EventID   string    `json:"eventId"`
	EventType string    `json:"eventType"`
	Seat      SeatDTO   `json:"seat"`
	SentAt    time.Time `json:"sentAt"`
	Source    string    `json:"source"`
}
