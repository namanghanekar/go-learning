package grpcclient

import (
	"context"
	"log"
	"time"

	pb "ticket-system/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentClient() *PaymentClient {
	conn, err := grpc.Dial(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &PaymentClient{
		client: pb.NewPaymentServiceClient(conn),
	}
}

// start 5 min timer
func (c *PaymentClient) StartTimer(seatID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.client.StartTimer(ctx, &pb.SeatRequest{
		SeatId: int32(seatID),
	})

	if err != nil {
		log.Println("timer error:", err)
	}
}
