package service

import (
	"time"

	grpcclient "ticket-system/payment-service/grpc"
	"ticket-system/payment-service/timer"
)

type Service struct {
	client *grpcclient.Client
}

func NewService(c *grpcclient.Client) *Service {
	return &Service{client: c}
}

func (s *Service) StartTimer(seatID int) {
	timer.Start(seatID, 5*time.Minute, func() {
		_ = s.client.Unlock(seatID)
	})
}

func (s *Service) Pay(seatID int, userID string) error {
	timer.Stop(seatID)
	return s.client.Confirm(seatID, userID)
}
