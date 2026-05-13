package dlq

type FailedEvent struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
	Reason  string `json:"reason"`
}
