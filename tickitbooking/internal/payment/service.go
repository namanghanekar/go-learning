package payment

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
	"worldtour-tickets/internal/contracts"
	"worldtour-tickets/internal/domain"
	"worldtour-tickets/internal/store"

	"github.com/google/uuid"
)

type queuedPayment struct {
	id      string
	request *contracts.SubmitPaymentRequest
}

type holdTimer struct {
	cancel context.CancelFunc
}

type Service struct {
	repo         *store.PaymentRepository
	inventory    contracts.InventoryServiceClient
	paymentQueue chan queuedPayment
	appCtx       context.Context

	streamMu sync.Mutex
	stream   contracts.InventoryService_OpenSeatStreamClient

	holdTimers   map[string]holdTimer
	holdTimersMu sync.Mutex
}

func NewService(repo *store.PaymentRepository, inventory contracts.InventoryServiceClient) *Service {
	return &Service{
		repo:         repo,
		inventory:    inventory,
		paymentQueue: make(chan queuedPayment, 256),
		holdTimers:   map[string]holdTimer{},
	}
}

func (s *Service) RegisterHold(ctx context.Context, req *contracts.RegisterHoldRequest) (*contracts.PaymentResponse, error) {
	if req.HoldToken == "" {
		return nil, errors.New("hold token is required")
	}

	paymentID := uuid.NewString()
	if err := s.repo.UpsertPending(domain.PaymentRecord{
		ID:          paymentID,
		SeatID:      req.SeatID,
		HoldToken:   req.HoldToken,
		UserID:      req.UserID,
		Status:      "HOLD_REGISTERED",
		AmountCents: req.AmountCents,
	}); err != nil {
		return nil, err
	}

	s.startExpiryTimer(req)

	return &contracts.PaymentResponse{
		Accepted:   true,
		Message:    "hold registered with payment service",
		PaymentRef: paymentID,
		At:         time.Now().UTC(),
	}, nil
}

func (s *Service) SubmitPayment(ctx context.Context, req *contracts.SubmitPaymentRequest) (*contracts.PaymentResponse, error) {
	paymentID := uuid.NewString()
	if err := s.repo.UpsertPending(domain.PaymentRecord{
		ID:          paymentID,
		SeatID:      req.SeatID,
		HoldToken:   req.HoldToken,
		UserID:      req.UserID,
		Status:      "PAYMENT_QUEUED",
		AmountCents: req.AmountCents,
	}); err != nil {
		return nil, err
	}

	select {
	case s.paymentQueue <- queuedPayment{id: paymentID, request: req}:
		return &contracts.PaymentResponse{
			Accepted:   true,
			Message:    "payment queued",
			PaymentRef: paymentID,
			At:         time.Now().UTC(),
		}, nil
	default:
		_ = s.repo.MarkProcessed(paymentID, "REJECTED", "payment queue is full")
		return &contracts.PaymentResponse{
			Accepted:   false,
			Message:    "flash sale queue is full, retry shortly",
			PaymentRef: paymentID,
			At:         time.Now().UTC(),
		}, nil
	}
}

func (s *Service) Start(ctx context.Context) {
	s.appCtx = ctx
	go s.worker(ctx)
	go s.connectStreamLoop(ctx)
}

func (s *Service) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item := <-s.paymentQueue:
			s.processPayment(item)
		}
	}
}

func (s *Service) processPayment(job queuedPayment) {
	s.cancelTimer(job.request.SeatID)
	time.Sleep(250 * time.Millisecond)

	ack, err := s.sendToInventory(&contracts.StreamEnvelope{
		ClientName: "payment-service",
		Confirm: &contracts.ConfirmSeatRequest{
			SeatID:     job.request.SeatID,
			HoldToken:  job.request.HoldToken,
			PaymentRef: job.id,
		},
	})
	if err != nil || ack == nil || !ack.Accepted {
		reason := "inventory confirm failed"
		if ack != nil && ack.Message != "" {
			reason = ack.Message
		}
		if err != nil {
			reason = err.Error()
		}
		_ = s.repo.MarkProcessed(job.id, "FAILED", reason)
		return
	}

	_ = s.repo.MarkProcessed(job.id, "SUCCEEDED", "")
}

func (s *Service) startExpiryTimer(req *contracts.RegisterHoldRequest) {
	s.cancelTimer(req.SeatID)

	timerCtx, cancel := context.WithCancel(context.Background())

	s.holdTimersMu.Lock()
	s.holdTimers[req.SeatID] = holdTimer{cancel: cancel}
	s.holdTimersMu.Unlock()

	go func() {
		wait := time.Until(req.ExpiresAt)
		if wait < 0 {
			wait = 0
		}

		select {
		case <-timerCtx.Done():
			return
		case <-time.After(wait):
			ack, err := s.sendToInventory(&contracts.StreamEnvelope{
				ClientName: "payment-service",
				Release: &contracts.ReleaseSeatRequest{
					SeatID:    req.SeatID,
					HoldToken: req.HoldToken,
					Reason:    "hold expired",
				},
			})
			if err != nil {
				log.Printf("failed to expire seat %s: %v", req.SeatID, err)
				return
			}
			if ack != nil && !ack.Accepted {
				log.Printf("inventory rejected expiry for seat %s: %s", req.SeatID, ack.Message)
			}
		}
	}()
}

func (s *Service) cancelTimer(seatID string) {
	s.holdTimersMu.Lock()
	defer s.holdTimersMu.Unlock()

	if timer, ok := s.holdTimers[seatID]; ok {
		timer.cancel()
		delete(s.holdTimers, seatID)
	}
}

func (s *Service) connectStreamLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		stream, err := s.inventory.OpenSeatStream(ctx)
		if err != nil {
			log.Printf("inventory stream connect failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		s.streamMu.Lock()
		s.stream = stream
		s.streamMu.Unlock()
		return
	}
}

func (s *Service) sendToInventory(msg *contracts.StreamEnvelope) (*contracts.SeatCommandAck, error) {
	for attempt := 0; attempt < 3; attempt++ {
		s.streamMu.Lock()
		stream := s.stream
		s.streamMu.Unlock()

		if stream == nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}

		if err := stream.Send(msg); err != nil {
			log.Printf("inventory stream send failed: %v", err)
			s.streamMu.Lock()
			s.stream = nil
			s.streamMu.Unlock()
			go s.connectStreamLoop(s.appCtx)
			time.Sleep(300 * time.Millisecond)
			continue
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				log.Printf("inventory stream recv failed: %v", err)
				s.streamMu.Lock()
				s.stream = nil
				s.streamMu.Unlock()
				go s.connectStreamLoop(s.appCtx)
				time.Sleep(300 * time.Millisecond)
				break
			}

			if resp.Event != nil {
				log.Printf("inventory event: %s seat=%s status=%s", resp.Event.EventType, resp.Event.Seat.Number, resp.Event.Seat.Status)
				continue
			}

			if resp.Ack != nil {
				return resp.Ack, nil
			}
		}
	}

	return nil, errors.New("inventory stream unavailable")
}
