package inventory

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"time"
	"worldtour-tickets/internal/contracts"
	"worldtour-tickets/internal/domain"
	"worldtour-tickets/internal/pubsub"
	"worldtour-tickets/internal/store"

	"github.com/google/uuid"
)

type Service struct {
	repo *store.SeatRepository
	hub  *pubsub.Hub[contracts.StreamEnvelope]
}

func NewService(repo *store.SeatRepository) *Service {
	return &Service{
		repo: repo,
		hub:  pubsub.NewHub[contracts.StreamEnvelope](),
	}
}

func (s *Service) ListSeats(ctx context.Context, req *contracts.ListSeatsRequest) (*contracts.ListSeatsResponse, error) {
	seats, err := s.repo.List(req.RoomID)
	if err != nil {
		return nil, err
	}
	return &contracts.ListSeatsResponse{Seats: mapSeats(seats)}, nil
}

func (s *Service) HoldSeat(ctx context.Context, req *contracts.HoldSeatRequest) (*contracts.HoldSeatResponse, error) {
	if req.HoldToken == "" {
		req.HoldToken = uuid.NewString()
	}

	ttl := holdDuration(req.TTLSeconds)

	seat, err := s.repo.HoldSeat(
		req.RoomID,
		req.SeatNumber,
		req.UserID,
		req.HoldToken,
		time.Now().UTC().Add(ttl),
	)
	if err != nil {
		return nil, err
	}

	s.broadcast("SEAT_HELD", seat, "inventory")
	return &contracts.HoldSeatResponse{Seat: toSeatDTO(seat)}, nil
}

func (s *Service) OpenSeatStream(stream contracts.InventoryService_OpenSeatStreamServer) error {
	sub := s.hub.Subscribe(64)
	defer s.hub.Unsubscribe(sub)

	done := make(chan struct{})
	var once sync.Once

	closeDone := func() {
		once.Do(func() { close(done) })
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case msg, ok := <-sub:
				if !ok {
					return
				}
				if err := stream.Send(&msg); err != nil {
					closeDone()
					return
				}
			}
		}
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				closeDone()
				return nil
			}
			closeDone()
			return err
		}

		if msg.Confirm != nil {
			seat, confirmErr := s.repo.ConfirmSeat(msg.Confirm.SeatID, msg.Confirm.HoldToken)
			if confirmErr != nil {
				_ = stream.Send(&contracts.StreamEnvelope{
					Ack: &contracts.SeatCommandAck{
						Accepted: false,
						Message:  confirmErr.Error(),
						At:       time.Now().UTC(),
					},
				})
				continue
			}

			s.broadcast("SEAT_SOLD", seat, msg.ClientName)
			_ = stream.Send(&contracts.StreamEnvelope{
				Ack: &contracts.SeatCommandAck{
					Accepted: true,
					Message:  "seat confirmed",
					Seat:     toSeatDTO(seat),
					At:       time.Now().UTC(),
				},
			})
			continue
		}

		if msg.Release != nil {
			seat, releaseErr := s.repo.ReleaseSeat(msg.Release.SeatID, msg.Release.HoldToken)
			if releaseErr != nil {
				_ = stream.Send(&contracts.StreamEnvelope{
					Ack: &contracts.SeatCommandAck{
						Accepted: false,
						Message:  releaseErr.Error(),
						At:       time.Now().UTC(),
					},
				})
				continue
			}

			s.broadcast("SEAT_RELEASED", seat, msg.ClientName)
			_ = stream.Send(&contracts.StreamEnvelope{
				Ack: &contracts.SeatCommandAck{
					Accepted: true,
					Message:  "seat released",
					Seat:     toSeatDTO(seat),
					At:       time.Now().UTC(),
				},
			})
		}
	}
}

func (s *Service) StartReconciler(ctx context.Context, roomID string) {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				seats, err := s.repo.List(roomID)
				if err != nil {
					log.Printf("inventory reconcile list failed: %v", err)
					continue
				}

				now := time.Now().UTC()
				for _, seat := range seats {
					if seat.Status == domain.SeatStatusHeld && seat.ExpiresAt != nil && seat.ExpiresAt.Before(now) {
						released, releaseErr := s.repo.ReleaseSeat(seat.ID, seat.HoldToken)
						if releaseErr != nil {
							continue
						}
						s.broadcast("SEAT_RELEASED", released, "reconciler")
					}
				}
			}
		}
	}()
}

func (s *Service) broadcast(eventType string, seat domain.Seat, source string) {
	s.hub.Publish(contracts.StreamEnvelope{
		Event: &contracts.SeatEvent{
			EventID:   uuid.NewString(),
			EventType: eventType,
			Seat:      toSeatDTO(seat),
			SentAt:    time.Now().UTC(),
			Source:    source,
		},
	})
}

func holdDuration(ttlSeconds int32) time.Duration {
	ttl := time.Duration(ttlSeconds) * time.Second
	if ttl <= 0 {
		return 5 * time.Minute
	}
	return ttl
}

func mapSeats(seats []domain.Seat) []contracts.SeatDTO {
	out := make([]contracts.SeatDTO, 0, len(seats))
	for _, seat := range seats {
		out = append(out, toSeatDTO(seat))
	}
	return out
}

func toSeatDTO(seat domain.Seat) contracts.SeatDTO {
	var expiresAt time.Time
	if seat.ExpiresAt != nil {
		expiresAt = *seat.ExpiresAt
	}
	return contracts.SeatDTO{
		ID:        seat.ID,
		RoomID:    seat.RoomID,
		Number:    seat.Number,
		Status:    seat.Status,
		HeldBy:    seat.HeldBy,
		HoldToken: seat.HoldToken,
		ExpiresAt: expiresAt,
		UpdatedAt: seat.UpdatedAt,
	}
}
