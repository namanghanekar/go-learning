package client

import (
	"context"

	pb "order-processing-system/proto/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	Client pb.OrderServiceClient
}

func NewOrderClient() *OrderClient {

	conn, _ := grpc.Dial(
		"localhost:50050",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	client := pb.NewOrderServiceClient(conn)

	return &OrderClient{
		Client: client,
	}
}

func (c *OrderClient) UpdateStatus(
	orderID int,
	status string,
) {

	c.Client.UpdateOrderStatus(
		context.Background(),
		&pb.UpdateOrderStatusRequest{
			OrderId: int32(orderID),
			Status:  status,
		},
	)
}
