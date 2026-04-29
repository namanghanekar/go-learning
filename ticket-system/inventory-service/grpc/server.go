package grpcserver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"ticket-system/inventory-service/service"
	pb "ticket-system/proto"
	"ticket-system/shared/logger"
	"time"
)

type Server struct {
	pb.UnimplementedInventoryServiceServer
	svc *service.Service
}

var subscribers = make(map[int]chan *pb.Seat)
var mu sync.Mutex

func NewServer(s *service.Service) *Server {
	return &Server{
		svc: s,
	}
}

func broadcast(seat *pb.Seat) {
	mu.Lock()
	defer mu.Unlock()

	for id, ch := range subscribers {
		select {
		case ch <- seat:
		default:
			log.Println("slow subscriber:", id)
		}
	}
}

func (s *Server) GetSeats(ctx context.Context, _ *pb.Empty) (*pb.SeatList, error) {
	seats := s.svc.GetSeats()

	var res []*pb.Seat
	for _, seat := range seats {
		res = append(res, &pb.Seat{
			Id:     int32(seat.ID),
			Status: seat.Status,
			UserId: seat.UserID,
		})
	}

	return &pb.SeatList{Seats: res}, nil
}

func (s *Server) LockSeat(ctx context.Context, req *pb.LockRequest) (*pb.SeatResponse, error) {
	err := s.svc.Lock(int(req.SeatId), req.UserId)
	if err != nil {
		return &pb.SeatResponse{Status: "error"}, err
	}

	logger.Log(fmt.Sprintf("seat locked: %d by %s", req.SeatId, req.UserId))

	broadcast(&pb.Seat{
		Id:     req.SeatId,
		Status: "locked",
		UserId: req.UserId,
	})

	return &pb.SeatResponse{Status: "locked"}, nil
}

func (s *Server) UnlockSeat(ctx context.Context, req *pb.SeatRequest) (*pb.SeatResponse, error) {
	if err := s.svc.Unlock(int(req.SeatId)); err != nil {
		return &pb.SeatResponse{Status: "error"}, err
	}
	seat := &pb.Seat{
		Id:     req.SeatId,
		Status: "available",
	}
	logger.Log(fmt.Sprintf("seat unlocked: %d", req.SeatId))
	broadcast(seat)

	return &pb.SeatResponse{Status: "available"}, nil
}
func (s *Server) ConfirmSeat(ctx context.Context, req *pb.SeatConfirmRequest) (*pb.SeatResponse, error) {
	err := s.svc.Confirm(int(req.SeatId), req.UserId)
	if err != nil {
		return &pb.SeatResponse{Status: "error"}, err
	}
	broadcast(&pb.Seat{
		Id:     req.SeatId,
		Status: "booked",
		UserId: req.UserId,
	})
	logger.Log(fmt.Sprintf("seat booked: %d by %s", req.SeatId, req.UserId))

	return &pb.SeatResponse{Status: "booked"}, nil
}

func (s *Server) SeatUpdates(req *pb.Empty, stream pb.InventoryService_SeatUpdatesServer) error {
	ch := make(chan *pb.Seat, 10)
	id := int(time.Now().UnixNano())

	mu.Lock()
	subscribers[id] = ch
	mu.Unlock()

	log.Println("subscriber connected:", id)

	defer func() {
		mu.Lock()
		delete(subscribers, id)
		close(ch)
		mu.Unlock()

		log.Println("subscriber disconnected:", id)
	}()

	for {
		select {
		case seat := <-ch:
			if err := stream.Send(seat); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return nil
		}
	}
}
