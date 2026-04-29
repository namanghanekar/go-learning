package service

import (
	"fmt"
	"time"

	grpcclient "ticket-system/payment-service/grpc"
	"ticket-system/payment-service/timer"
	"ticket-system/shared/logger"
)

type Service struct {
	client *grpcclient.Client
}

func NewService(c *grpcclient.Client) *Service {
	return &Service{client: c}
}

func (s *Service) StartTimer(seatID int) {
	logger.Log(fmt.Sprintf("payment timer started for seat: %d", seatID))
	timer.Start(seatID, 5*time.Minute, func() {
		logger.Log(fmt.Sprintf("payment timer expired for seat: %d", seatID))
		_ = s.client.Unlock(seatID)
	})
}

func (s *Service) Pay(seatID int, userID string) error {
	timer.Stop(seatID)
	logger.Log(fmt.Sprintf("payment confirmed for seat: %d by %s", seatID, userID))
	return s.client.Confirm(seatID, userID)
}
