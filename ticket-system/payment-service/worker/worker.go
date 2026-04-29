package worker

import (
	"fmt"

	"ticket-system/shared/logger"
)

type PaymentRequest struct {
	SeatID int
	UserID string
	Result chan error
}

var PaymentQueue = make(chan PaymentRequest, 100)

func StartWorker(workerCount int, process func(PaymentRequest)) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for req := range PaymentQueue {
				logger.Log(fmt.Sprintf("processing payment for seat: %d by %s", req.SeatID, req.UserID))
				process(req)
			}
		}()
	}
}
