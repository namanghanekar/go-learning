package client

import (
	"context"
	"fmt"
	"time"

	pb "ticket-system/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Inventory pb.InventoryServiceClient
	Payment   pb.PaymentServiceClient
}

func NewClient() *Client {
	inventoryConn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	paymentConn, err := grpc.Dial(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	return &Client{
		Inventory: pb.NewInventoryServiceClient(inventoryConn),
		Payment:   pb.NewPaymentServiceClient(paymentConn),
	}
}

func (c *Client) GetSeats() (*pb.SeatList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.Inventory.GetSeats(ctx, &pb.Empty{})
}

func (c *Client) Lock(seatID int, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := c.Inventory.LockSeat(ctx, &pb.LockRequest{
		SeatId: int32(seatID),
		UserId: userID,
	})
	if err != nil {
		return fmt.Errorf("lock seat %d: %w", seatID, err)
	}
	return nil
}

func (c *Client) Pay(seatID int, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := c.Payment.Pay(ctx, &pb.PaymentRequest{
		SeatId: int32(seatID),
		UserId: userID,
	})
	if err != nil {
		return fmt.Errorf("pay seat %d: %w", seatID, err)
	}
	return nil
}
