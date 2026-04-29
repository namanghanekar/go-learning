package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	grpcclient "ticket-system/payment-service/grpc"
	"ticket-system/payment-service/service"
	"ticket-system/payment-service/timer"
	"ticket-system/payment-service/worker"

	pb "ticket-system/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
	service *service.Service
}

func NewServer() *Server {
	client := grpcclient.NewClient()
	svc := service.NewService(client)

	return &Server{
		service: svc,
	}
}

func (s *Server) StartTimer(ctx context.Context, req *pb.SeatRequest) (*pb.SeatResponse, error) {
	log.Println("received timer request for seat:", req.SeatId)
	s.service.StartTimer(int(req.SeatId))

	return &pb.SeatResponse{
		Status: "timer started",
	}, nil
}

func (s *Server) Pay(ctx context.Context, req *pb.PaymentRequest) (*pb.SeatResponse, error) {
	log.Println("received payment request for seat:", req.SeatId)

	result := make(chan error, 1)
	select {
	case worker.PaymentQueue <- worker.PaymentRequest{
		SeatID: int(req.SeatId),
		UserID: req.UserId,
		Result: result,
	}:
	default:
		return &pb.SeatResponse{Status: "payment_queue_busy"}, fmt.Errorf("payment queue is full")
	}

	select {
	case err := <-result:
		if err != nil {
			return &pb.SeatResponse{Status: "payment_failed"}, fmt.Errorf("payment failed: %w", err)
		}
	case <-ctx.Done():
		return &pb.SeatResponse{Status: "payment_timeout"}, fmt.Errorf("payment timeout: %w", ctx.Err())
	}

	return &pb.SeatResponse{
		Status: "payment success",
	}, nil
}
func startSeatListener(client pb.InventoryServiceClient, svc *service.Service) {
	for {
		stream, err := client.SeatUpdates(context.Background(), &pb.Empty{})
		if err != nil {
			log.Println("stream error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("listening for seat updates")

		for {
			seat, recvErr := stream.Recv()
			if recvErr != nil {
				log.Println("stream recv error:", recvErr)
				time.Sleep(1 * time.Second)
				break
			}

			log.Println("seat update:", seat.Id, seat.Status)

			switch seat.Status {
			case "locked":
				svc.StartTimer(int(seat.Id))
			case "booked":
				timer.Stop(int(seat.Id))
			case "available":
				timer.Stop(int(seat.Id))
			}
		}
	}
}

func main() {
	client := grpcclient.NewClient()
	svc := service.NewService(client)
	worker.StartWorker(8, func(req worker.PaymentRequest) {
		req.Result <- svc.Pay(req.SeatID, req.UserID)
	})

	go startSeatListener(client.Client, svc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	pb.RegisterPaymentServiceServer(s, &Server{
		service: svc,
	})

	log.Println("payment service running on :50052")

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
