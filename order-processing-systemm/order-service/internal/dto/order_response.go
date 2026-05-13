package dto

type OrderResponse struct {
	ID      uint    `json:"id"`
	Status  string  `json:"status"`
	Amount  float64 `json:"amount"`
	Message string  `json:"message"`
}
