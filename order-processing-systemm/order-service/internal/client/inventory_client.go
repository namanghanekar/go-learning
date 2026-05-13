package client

import (
	"context"

	pb "order-processing-system/proto/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryClient struct {
	Client pb.InventoryServiceClient
}

func NewInventoryClient() *InventoryClient {

	conn, _ := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	client := pb.NewInventoryServiceClient(conn)

	return &InventoryClient{
		Client: client,
	}
}

func (i *InventoryClient) CheckStock(
	productID int,
	quantity int,
) (bool, error) {

	res, err := i.Client.CheckStock(
		context.Background(),
		&pb.InventoryRequest{
			ProductId: int32(productID),
			Quantity:  int32(quantity),
		},
	)

	if err != nil {
		return false, err
	}

	return res.Available, nil
}
