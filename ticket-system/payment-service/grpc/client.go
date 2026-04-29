package grpcclient

import (
	"context"
	"fmt"
	"time"

	pb "ticket-system/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Client pb.InventoryServiceClient
}

func NewClient() *Client {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	return &Client{
		Client: pb.NewInventoryServiceClient(conn),
	}
}

func (c *Client) Unlock(seatID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.Client.UnlockSeat(ctx, &pb.SeatRequest{
		SeatId: int32(seatID),
	})
	if err != nil {
		return fmt.Errorf("unlock seat %d: %w", seatID, err)
	}
	return nil
}

func (c *Client) Confirm(seatID int, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.Client.ConfirmSeat(ctx, &pb.SeatConfirmRequest{
		SeatId: int32(seatID),
		UserId: userID,
	})
	if err != nil {
		return fmt.Errorf("confirm seat %d for user %s: %w", seatID, userID, err)
	}
	return nil
}
