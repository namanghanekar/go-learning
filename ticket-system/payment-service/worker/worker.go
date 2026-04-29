package worker

import "log"

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
				log.Println("processing payment for seat:", req.SeatID)
				process(req)
			}
		}()
	}
}
